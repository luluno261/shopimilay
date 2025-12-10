package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

var stripeConfig *StripeConfig

func main() {
	// Initialisation de la base de données
	initDB()
	defer db.Close()
	
	// Initialisation Stripe
	stripeConfig = InitStripe()
	
	port := getEnv("PORT", "8081")
	
	router := gin.Default()
	
	// Routes de santé
	router.GET("/health", healthCheck)
	router.GET("/ready", readinessCheck)
	
	// Routes de checkout
	api := router.Group("/api/v1")
	{
		api.GET("/cart", authenticateMiddleware(), handleGetCart)
		api.POST("/cart/items", authenticateMiddleware(), handleAddToCart)
		api.DELETE("/cart/items/:itemId", authenticateMiddleware(), handleRemoveFromCart)
		api.POST("/checkout", authenticateMiddleware(), handleCheckout)
		api.POST("/checkout/stripe/webhook", handleStripeWebhook)
		
		// Routes de commandes
		api.GET("/orders", authenticateMiddleware(), handleGetOrders)
		api.GET("/orders/:id", authenticateMiddleware(), handleGetOrder)
		api.PUT("/orders/:id/status", authenticateMiddleware(), handleUpdateOrderStatus)
		api.POST("/orders/:id/refund", authenticateMiddleware(), handleRefundOrder)
		
		// Routes compte client
		api.GET("/account/orders", authenticateMiddleware(), handleGetUserOrders)
		api.GET("/account/addresses", authenticateMiddleware(), handleGetAddresses)
		api.POST("/account/addresses", authenticateMiddleware(), handleCreateAddress)
		
		// Routes dashboard
		api.GET("/dashboard/stats", authenticateMiddleware(), handleGetDashboardStats)
		
		// Routes discounts
		api.GET("/discounts", authenticateMiddleware(), handleGetDiscounts)
		api.POST("/discounts", authenticateMiddleware(), handleCreateDiscount)
		api.DELETE("/discounts/:id", authenticateMiddleware(), handleDeleteDiscount)
	}
	
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	go func() {
		log.Printf("Checkout Service démarré sur le port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erreur serveur: %v", err)
		}
	}()
	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	
	log.Println("Arrêt du serveur...")
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Erreur lors de l'arrêt du serveur:", err)
	}
	
	log.Println("Serveur arrêté")
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"service": "checkout-service",
	})
}

func readinessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
		"service": "checkout-service",
	})
}

func handleGetCart(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	merchantID := c.GetHeader("X-Merchant-ID")
	
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}
	
	if merchantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Merchant ID manquant"})
		return
	}

	cart, err := GetOrCreateCart(userID, merchantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération du panier"})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// AddToCartRequest représente une demande d'ajout au panier
type AddToCartRequest struct {
	ProductID string  `json:"product_id" binding:"required"`
	VariantID *string `json:"variant_id,omitempty"`
	Quantity  int     `json:"quantity" binding:"required,min=1"`
}

func handleAddToCart(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	merchantID := c.GetHeader("X-Merchant-ID")
	
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}
	
	if merchantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Merchant ID manquant"})
		return
	}

	var req AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Récupérer le prix depuis le catalogue-service
	price, err := getProductPrice(req.ProductID, req.VariantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Produit introuvable ou erreur lors de la récupération du prix"})
		return
	}

	cart, err := GetOrCreateCart(userID, merchantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération du panier"})
		return
	}

	err = AddItemToCart(cart.ID, req.ProductID, req.VariantID, req.Quantity, price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'ajout au panier"})
		return
	}

	// Récupérer le panier mis à jour
	updatedCart, err := GetOrCreateCart(userID, merchantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération du panier"})
		return
	}

	c.JSON(http.StatusOK, updatedCart)
}

func handleRemoveFromCart(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	merchantID := c.GetHeader("X-Merchant-ID")
	itemID := c.Param("itemId")
	
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}
	
	if merchantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Merchant ID manquant"})
		return
	}

	cart, err := GetOrCreateCart(userID, merchantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération du panier"})
		return
	}

	err = RemoveItemFromCart(cart.ID, itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression"})
		return
	}

	// Récupérer le panier mis à jour
	updatedCart, err := GetOrCreateCart(userID, merchantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération du panier"})
		return
	}

	c.JSON(http.StatusOK, updatedCart)
}

func handleCheckout(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	merchantID := c.GetHeader("X-Merchant-ID")
	
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}
	
	if merchantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Merchant ID manquant"})
		return
	}

	var req CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Récupérer le panier
	cart, err := GetOrCreateCart(userID, merchantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération du panier"})
		return
	}

	if len(cart.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Le panier est vide"})
		return
	}

	// Appliquer le code de réduction si fourni
	finalTotal := cart.Total
	if req.DiscountCode != nil && *req.DiscountCode != "" {
		discountAmount, err := ApplyDiscount(*req.DiscountCode, cart.Total, merchantID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Erreur lors de l'application du code de réduction"})
			return
		}
		if discountAmount > 0 {
			finalTotal = cart.Total - discountAmount
			if finalTotal < 0 {
				finalTotal = 0
			}
		}
	}

	// Créer la commande
	order, err := CreateOrder(userID, merchantID, finalTotal, cart.Currency, req.ShippingAddress, req.BillingAddress, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de la commande"})
		return
	}

	// Créer le PaymentIntent Stripe
	amount := int64(finalTotal * 100) // Convertir en centimes
	paymentIntentID, clientSecret, err := CreatePaymentIntent(amount, cart.Currency, merchantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création du paiement"})
		return
	}

	// Mettre à jour la commande avec le PaymentIntent ID
	_, err = db.Exec("UPDATE orders SET payment_intent_id = $1 WHERE id = $2", paymentIntentID, order.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour de la commande"})
		return
	}

	// Vider le panier (changer le statut)
	_, err = db.Exec("UPDATE carts SET status = 'completed' WHERE id = $1", cart.ID)
	if err != nil {
		log.Printf("Erreur lors de la mise à jour du panier: %v", err)
	}

	c.JSON(http.StatusOK, CheckoutResponse{
		OrderID:        order.ID,
		Status:          order.Status,
		PaymentIntentID: paymentIntentID,
		ClientSecret:    clientSecret,
	})
}

func handleStripeWebhook(c *gin.Context) {
	HandleStripeWebhook(c)
}

func authenticateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// L'authentification est gérée par l'API Gateway
		// On vérifie juste que les headers sont présents
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// getProductPrice récupère le prix d'un produit depuis le catalogue-service
func getProductPrice(productID string, variantID *string) (float64, error) {
	catalogueURL := getEnv("CATALOGUE_SERVICE_URL", "http://localhost:8082")
	
	url := catalogueURL + "/api/v1/products/" + productID
	if variantID != nil {
		url += "?variant_id=" + *variantID
	}
	
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("produit non trouvé")
	}
	
	var product struct {
		Price float64 `json:"price"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return 0, err
	}
	
	return product.Price, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

