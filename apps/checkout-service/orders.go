package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/refund"
)

func handleGetOrders(c *gin.Context) {
	merchantID := c.GetHeader("X-Merchant-ID")
	if merchantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Merchant ID manquant"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	rows, err := db.Query(
		"SELECT id, user_id, merchant_id, status, total_amount, currency, payment_intent_id, shipping_address, billing_address, created_at, updated_at FROM orders WHERE merchant_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3",
		merchantID, limit, offset,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des commandes"})
		return
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		err := rows.Scan(&order.ID, &order.UserID, &order.MerchantID, &order.Status, &order.TotalAmount, &order.Currency, &order.PaymentIntentID, &order.ShippingAddress, &order.BillingAddress, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du scan"})
			return
		}
		orders = append(orders, order)
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

func handleGetOrder(c *gin.Context) {
	orderID := c.Param("id")
	merchantID := c.GetHeader("X-Merchant-ID")

	order, err := GetOrder(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération de la commande"})
		return
	}

	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Commande introuvable"})
		return
	}

	// Vérifier que la commande appartient au marchand
	if order.MerchantID != merchantID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Accès non autorisé"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func handleUpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")
	merchantID := c.GetHeader("X-Merchant-ID")

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := GetOrder(orderID)
	if err != nil || order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Commande introuvable"})
		return
	}

	if order.MerchantID != merchantID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Accès non autorisé"})
		return
	}

	err = UpdateOrderStatus(orderID, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Statut mis à jour"})
}

func handleRefundOrder(c *gin.Context) {
	orderID := c.Param("id")
	merchantID := c.GetHeader("X-Merchant-ID")

	var req struct {
		Amount *float64 `json:"amount,omitempty"` // Si nil, remboursement complet
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := GetOrder(orderID)
	if err != nil || order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Commande introuvable"})
		return
	}

	if order.MerchantID != merchantID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Accès non autorisé"})
		return
	}

	if order.PaymentIntentID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Aucun paiement associé"})
		return
	}

	// Implémenter le remboursement Stripe
	amount := order.TotalAmount
	if req.Amount != nil {
		amount = *req.Amount
	}

	amountInCents := int64(amount * 100)
	
	// Utiliser le SDK Stripe pour créer un remboursement
	refundParams := &stripe.RefundParams{
		PaymentIntent: stripe.String(*order.PaymentIntentID),
		Amount:         stripe.Int64(amountInCents),
	}

	refund, err := refund.New(refundParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du remboursement Stripe: " + err.Error()})
		return
	}

	// Mettre à jour le statut de la commande
	newStatus := "refunded"
	if req.Amount != nil && *req.Amount < order.TotalAmount {
		newStatus = "partially_refunded"
	}
	
	err = UpdateOrderStatus(orderID, newStatus)
	if err != nil {
		log.Printf("Erreur lors de la mise à jour du statut: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Remboursement initié",
		"refund_id":    refund.ID,
		"amount":       amount,
		"order_status": newStatus,
	})
}

func handleGetUserOrders(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	orders, err := GetUserOrders(userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

func handleGetAddresses(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	rows, err := db.Query(
		"SELECT id, user_id, merchant_id, type, street, city, state, zip_code, country, is_default, created_at, updated_at FROM addresses WHERE user_id = $1 ORDER BY is_default DESC, created_at DESC",
		userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération"})
		return
	}
	defer rows.Close()

	var addresses []map[string]interface{}
	for rows.Next() {
		var addr struct {
			ID         string
			UserID     string
			MerchantID string
			Type       string
			Street     string
			City       string
			ZipCode    string
			Country    string
			IsDefault  bool
			CreatedAt  time.Time
			UpdatedAt  time.Time
		}
		var state sql.NullString
		err := rows.Scan(&addr.ID, &addr.UserID, &addr.MerchantID, &addr.Type, &addr.Street, &addr.City, &state, &addr.ZipCode, &addr.Country, &addr.IsDefault, &addr.CreatedAt, &addr.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du scan"})
			return
		}
		
		address := map[string]interface{}{
			"id":          addr.ID,
			"user_id":     addr.UserID,
			"merchant_id": addr.MerchantID,
			"type":        addr.Type,
			"street":      addr.Street,
			"city":        addr.City,
			"zip_code":    addr.ZipCode,
			"country":     addr.Country,
			"is_default":  addr.IsDefault,
			"created_at":  addr.CreatedAt,
			"updated_at":  addr.UpdatedAt,
		}
		if state.Valid {
			address["state"] = state.String
		}
		addresses = append(addresses, address)
	}

	c.JSON(http.StatusOK, gin.H{"addresses": addresses})
}

func handleCreateAddress(c *gin.Context) {
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

	var req struct {
		Type      string `json:"type" binding:"required,oneof=shipping billing both"`
		Street    string `json:"street" binding:"required"`
		City      string `json:"city" binding:"required"`
		State     string `json:"state"`
		ZipCode   string `json:"zip_code" binding:"required"`
		Country   string `json:"country" binding:"required"`
		IsDefault bool   `json:"is_default"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Si c'est l'adresse par défaut, désactiver les autres
	if req.IsDefault {
		_, err := db.Exec("UPDATE addresses SET is_default = false WHERE user_id = $1", userID)
		if err != nil {
			log.Printf("Erreur lors de la mise à jour des adresses par défaut: %v", err)
		}
	}

	var state *string
	if req.State != "" {
		state = &req.State
	}

	var addressID string
	err := db.QueryRow(
		"INSERT INTO addresses (user_id, merchant_id, type, street, city, state, zip_code, country, is_default) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id",
		userID, merchantID, req.Type, req.Street, req.City, state, req.ZipCode, req.Country, req.IsDefault,
	).Scan(&addressID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":          addressID,
		"message":     "Adresse créée avec succès",
	})
}

func handleGetDashboardStats(c *gin.Context) {
	merchantID := c.GetHeader("X-Merchant-ID")
	if merchantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Merchant ID manquant"})
		return
	}

	// Calculer les KPIs
	var mrr, gmv, orderCount float64
	var avgCart float64

	// MRR (Monthly Recurring Revenue) - pour l'instant, on utilise le total du mois
	err := db.QueryRow(
		"SELECT COALESCE(SUM(total_amount), 0) FROM orders WHERE merchant_id = $1 AND status = 'paid' AND created_at >= date_trunc('month', CURRENT_DATE)",
		merchantID,
	).Scan(&mrr)
	if err != nil {
		mrr = 0
	}

	// GMV (Gross Merchandise Value) - total des commandes
	err = db.QueryRow(
		"SELECT COALESCE(SUM(total_amount), 0) FROM orders WHERE merchant_id = $1 AND status = 'paid'",
		merchantID,
	).Scan(&gmv)
	if err != nil {
		gmv = 0
	}

	// Nombre de commandes
	err = db.QueryRow(
		"SELECT COUNT(*) FROM orders WHERE merchant_id = $1 AND status = 'paid'",
		merchantID,
	).Scan(&orderCount)
	if err != nil {
		orderCount = 0
	}

	// Panier moyen
	if orderCount > 0 {
		avgCart = gmv / orderCount
	}

	c.JSON(http.StatusOK, gin.H{
		"mrr":         mrr,
		"gmv":         gmv,
		"order_count": orderCount,
		"avg_cart":    avgCart,
	})
}

