package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var esClient *ElasticsearchClient

func main() {
	// Initialisation de la base de données
	initDB()
	defer db.Close()
	
	// Initialisation Elasticsearch
	esClient = NewElasticsearchClient()
	
	port := getEnv("PORT", "8082")
	
	router := gin.Default()
	
	// Routes de santé
	router.GET("/health", healthCheck)
	router.GET("/ready", readinessCheck)
	
	// Routes du catalogue
	api := router.Group("/api/v1")
	{
		api.GET("/products", handleListProducts)
		api.GET("/products/:id", handleGetProduct)
		api.POST("/products", authenticateMiddleware(), handleCreateProduct)
		api.PUT("/products/:id", authenticateMiddleware(), handleUpdateProduct)
		api.DELETE("/products/:id", authenticateMiddleware(), handleDeleteProduct)
		
		api.GET("/inventory/:productId", handleGetInventory)
		api.PUT("/inventory/:productId", authenticateMiddleware(), handleUpdateInventory)
		
		api.POST("/search", handleSearchProducts)
		
		// Route pour génération de description par IA
		api.POST("/products/generate-description", authenticateMiddleware(), HandleGenerateDescription)
		
		// Routes Store Builder (publiques pour GET, protégées pour POST)
		api.GET("/store-builder/config", handleGetStorefrontConfig)
		api.POST("/store-builder/config", authenticateMiddleware(), handleSaveStorefrontConfig)
		api.GET("/store-builder/theme", handleGetTheme)
		api.POST("/store-builder/theme", authenticateMiddleware(), handleSaveTheme)
	}
	
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	go func() {
		log.Printf("Catalogue Service démarré sur le port %s", port)
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
		"service": "catalogue-service",
	})
}

func readinessCheck(c *gin.Context) {
	// Vérifier la connexion à la base de données
	if err := db.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not ready",
			"service": "catalogue-service",
			"error": "Database connection failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
		"service": "catalogue-service",
	})
}

func handleListProducts(c *gin.Context) {
	merchantID := c.Query("merchant_id")
	if merchantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "merchant_id requis"})
		return
	}
	
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	
	products, err := ListProductsDB(merchantID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des produits"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"products": products})
}

func handleGetProduct(c *gin.Context) {
	productID := c.Param("id")
	
	product, err := GetProductByID(productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération du produit"})
		return
	}
	
	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produit introuvable"})
		return
	}
	
	c.JSON(http.StatusOK, product)
}

func handleCreateProduct(c *gin.Context) {
	merchantID := c.GetHeader("X-Merchant-ID")
	if merchantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Merchant ID manquant"})
		return
	}
	
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	product, err := CreateProductDB(merchantID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création du produit"})
		return
	}
	
	// Indexer dans Elasticsearch
	if err := esClient.IndexProduct(product); err != nil {
		log.Printf("Erreur lors de l'indexation Elasticsearch: %v", err)
	}
	
	c.JSON(http.StatusCreated, product)
}

func handleUpdateProduct(c *gin.Context) {
	productID := c.Param("id")
	merchantID := c.GetHeader("X-Merchant-ID")
	
	// Vérifier que le produit appartient au marchand
	product, err := GetProductByID(productID)
	if err != nil || product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produit introuvable"})
		return
	}
	
	if product.MerchantID != merchantID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Accès non autorisé"})
		return
	}
	
	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	updatedProduct, err := UpdateProductDB(productID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour"})
		return
	}
	
	// Mettre à jour dans Elasticsearch
	if err := esClient.IndexProduct(updatedProduct); err != nil {
		log.Printf("Erreur lors de la mise à jour Elasticsearch: %v", err)
	}
	
	c.JSON(http.StatusOK, updatedProduct)
}

func handleDeleteProduct(c *gin.Context) {
	productID := c.Param("id")
	merchantID := c.GetHeader("X-Merchant-ID")
	
	// Vérifier que le produit appartient au marchand
	product, err := GetProductByID(productID)
	if err != nil || product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produit introuvable"})
		return
	}
	
	if product.MerchantID != merchantID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Accès non autorisé"})
		return
	}
	
	err = DeleteProductDB(productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression"})
		return
	}
	
	// Supprimer d'Elasticsearch
	if err := esClient.DeleteProduct(productID); err != nil {
		log.Printf("Erreur lors de la suppression Elasticsearch: %v", err)
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Produit supprimé"})
}

func handleGetInventory(c *gin.Context) {
	productID := c.Param("productId")
	variantID := c.Query("variant_id")
	
	var variantIDPtr *string
	if variantID != "" {
		variantIDPtr = &variantID
	}
	
	inventory, err := GetInventory(productID, variantIDPtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération du stock"})
		return
	}
	
	if inventory == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stock introuvable"})
		return
	}
	
	c.JSON(http.StatusOK, inventory)
}

func handleUpdateInventory(c *gin.Context) {
	productID := c.Param("productId")
	merchantID := c.GetHeader("X-Merchant-ID")
	
	// Vérifier que le produit appartient au marchand
	product, err := GetProductByID(productID)
	if err != nil || product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produit introuvable"})
		return
	}
	
	if product.MerchantID != merchantID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Accès non autorisé"})
		return
	}
	
	var req struct {
		VariantID *string `json:"variant_id,omitempty"`
		Quantity  int     `json:"quantity" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	err = UpdateInventory(productID, req.VariantID, req.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour du stock"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Stock mis à jour"})
}

func handleSearchProducts(c *gin.Context) {
	var req struct {
		Query   string                 `json:"query" binding:"required"`
		Filters map[string]interface{} `json:"filters,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	products, err := esClient.SearchProducts(req.Query, req.Filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la recherche"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"products": products})
}

func authenticateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// L'authentification est gérée par l'API Gateway
		// On vérifie juste que les headers sont présents
		merchantID := c.GetHeader("X-Merchant-ID")
		if merchantID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Merchant ID manquant"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

