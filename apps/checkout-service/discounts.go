package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Discount représente un code de réduction
type Discount struct {
	ID        string    `json:"id" db:"id"`
	MerchantID string   `json:"merchant_id" db:"merchant_id"`
	Code      string    `json:"code" db:"code"`
	Type      string    `json:"type" db:"type"` // percentage, fixed
	Value     float64   `json:"value" db:"value"`
	ExpiresAt *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func handleGetDiscounts(c *gin.Context) {
	merchantID := c.GetHeader("X-Merchant-ID")
	if merchantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Merchant ID manquant"})
		return
	}

	rows, err := db.Query(
		"SELECT id, merchant_id, code, type, value, expires_at, is_active, created_at, updated_at FROM discounts WHERE merchant_id = $1 ORDER BY created_at DESC",
		merchantID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération"})
		return
	}
	defer rows.Close()

	var discounts []Discount
	for rows.Next() {
		var d Discount
		var expiresAt sql.NullTime
		err := rows.Scan(&d.ID, &d.MerchantID, &d.Code, &d.Type, &d.Value, &expiresAt, &d.IsActive, &d.CreatedAt, &d.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du scan"})
			return
		}
		if expiresAt.Valid {
			d.ExpiresAt = &expiresAt.Time
		}
		discounts = append(discounts, d)
	}

	c.JSON(http.StatusOK, gin.H{"discounts": discounts})
}

func handleCreateDiscount(c *gin.Context) {
	merchantID := c.GetHeader("X-Merchant-ID")
	if merchantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Merchant ID manquant"})
		return
	}

	var req struct {
		Code      string  `json:"code" binding:"required"`
		Type      string  `json:"type" binding:"required,oneof=percentage fixed"`
		Value     float64 `json:"value" binding:"required"`
		ExpiresAt *string `json:"expires_at,omitempty"`
		IsActive  bool    `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var expiresAt *time.Time
	if req.ExpiresAt != nil {
		t, err := time.Parse("2006-01-02", *req.ExpiresAt)
		if err == nil {
			expiresAt = &t
		}
	}

	var discount Discount
	err := db.QueryRow(
		"INSERT INTO discounts (merchant_id, code, type, value, expires_at, is_active) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, merchant_id, code, type, value, expires_at, is_active, created_at, updated_at",
		merchantID, req.Code, req.Type, req.Value, expiresAt, req.IsActive,
	).Scan(&discount.ID, &discount.MerchantID, &discount.Code, &discount.Type, &discount.Value, &expiresAt, &discount.IsActive, &discount.CreatedAt, &discount.UpdatedAt)

	if err != nil {
		log.Printf("Erreur lors de la création: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création"})
		return
	}

	if expiresAt != nil {
		discount.ExpiresAt = expiresAt
	}

	c.JSON(http.StatusCreated, discount)
}

func handleDeleteDiscount(c *gin.Context) {
	discountID := c.Param("id")
	merchantID := c.GetHeader("X-Merchant-ID")

	// Vérifier que le discount appartient au marchand
	var existingMerchantID string
	err := db.QueryRow("SELECT merchant_id FROM discounts WHERE id = $1", discountID).Scan(&existingMerchantID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Code de réduction introuvable"})
		return
	}
	if existingMerchantID != merchantID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Accès non autorisé"})
		return
	}

	_, err = db.Exec("DELETE FROM discounts WHERE id = $1", discountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Code de réduction supprimé"})
}

// ApplyDiscount applique un code de réduction à un montant
func ApplyDiscount(code string, amount float64, merchantID string) (float64, error) {
	var discount Discount
	var expiresAt sql.NullTime
	
	err := db.QueryRow(
		"SELECT id, merchant_id, code, type, value, expires_at, is_active FROM discounts WHERE code = $1 AND merchant_id = $2 AND is_active = true",
		code, merchantID,
	).Scan(&discount.ID, &discount.MerchantID, &discount.Code, &discount.Type, &discount.Value, &expiresAt, &discount.IsActive)

	if err == sql.ErrNoRows {
		return 0, nil // Code non trouvé
	}
	if err != nil {
		return 0, err
	}

	// Vérifier l'expiration
	if expiresAt.Valid && expiresAt.Time.Before(time.Now()) {
		return 0, nil // Code expiré
	}

	// Calculer la réduction
	var discountAmount float64
	if discount.Type == "percentage" {
		discountAmount = amount * (discount.Value / 100)
	} else {
		discountAmount = discount.Value
		if discountAmount > amount {
			discountAmount = amount
		}
	}

	return discountAmount, nil
}

