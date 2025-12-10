package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	// Initialisation de la base de données
	initDB()
	defer db.Close()

	port := getEnv("PORT", "8084")

	router := gin.Default()

	// Routes de santé
	router.GET("/health", healthCheck)
	router.GET("/ready", readinessCheck)

	// Routes API
	api := router.Group("/api/v1")
	{
		api.GET("/webhooks", authenticateMiddleware(), handleListWebhooks)
		api.POST("/webhooks", authenticateMiddleware(), handleCreateWebhook)
		api.PUT("/webhooks/:id", authenticateMiddleware(), handleUpdateWebhook)
		api.DELETE("/webhooks/:id", authenticateMiddleware(), handleDeleteWebhook)
		api.POST("/webhooks/:id/test", authenticateMiddleware(), handleTestWebhook)
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Webhook Service démarré sur le port %s", port)
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

func initDB() {
	connStr := getEnv("DATABASE_URL", "postgres://omnisphere:omnisphere_dev@localhost:5432/omnisphere?sslmode=disable")

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Erreur de connexion à la base de données:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Erreur de ping à la base de données:", err)
	}

	log.Println("Connexion à la base de données établie")
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "webhook-service",
	})
}

func readinessCheck(c *gin.Context) {
	if err := db.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "not ready",
			"service": "webhook-service",
			"error":   "Database connection failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "ready",
		"service": "webhook-service",
	})
}

// Webhook représente une configuration de webhook
type Webhook struct {
	ID          string    `json:"id" db:"id"`
	MerchantID  string    `json:"merchant_id" db:"merchant_id"`
	EventType   string    `json:"event_type" db:"event_type"`
	URL         string    `json:"url" db:"url"`
	Secret      string    `json:"secret,omitempty" db:"secret"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func handleListWebhooks(c *gin.Context) {
	merchantID := c.GetHeader("X-Merchant-ID")
	if merchantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Merchant ID manquant"})
		return
	}

	rows, err := db.Query(
		"SELECT id, merchant_id, event_type, url, is_active, created_at, updated_at FROM webhooks WHERE merchant_id = $1 ORDER BY created_at DESC",
		merchantID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération"})
		return
	}
	defer rows.Close()

	var webhooks []Webhook
	for rows.Next() {
		var w Webhook
		err := rows.Scan(&w.ID, &w.MerchantID, &w.EventType, &w.URL, &w.IsActive, &w.CreatedAt, &w.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du scan"})
			return
		}
		webhooks = append(webhooks, w)
	}

	c.JSON(http.StatusOK, gin.H{"webhooks": webhooks})
}

func handleCreateWebhook(c *gin.Context) {
	merchantID := c.GetHeader("X-Merchant-ID")
	if merchantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Merchant ID manquant"})
		return
	}

	var req struct {
		EventType string `json:"event_type" binding:"required"`
		URL       string `json:"url" binding:"required,url"`
		Secret    string `json:"secret,omitempty"`
		IsActive  bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Générer un secret si non fourni
	secret := req.Secret
	if secret == "" {
		secret = generateSecret()
	}

	var webhook Webhook
	err := db.QueryRow(
		"INSERT INTO webhooks (merchant_id, event_type, url, secret, is_active) VALUES ($1, $2, $3, $4, $5) RETURNING id, merchant_id, event_type, url, is_active, created_at, updated_at",
		merchantID, req.EventType, req.URL, secret, req.IsActive,
	).Scan(&webhook.ID, &webhook.MerchantID, &webhook.EventType, &webhook.URL, &webhook.IsActive, &webhook.CreatedAt, &webhook.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création"})
		return
	}

	c.JSON(http.StatusCreated, webhook)
}

func handleUpdateWebhook(c *gin.Context) {
	webhookID := c.Param("id")
	merchantID := c.GetHeader("X-Merchant-ID")

	var req struct {
		URL      string `json:"url,omitempty"`
		Secret   string `json:"secret,omitempty"`
		IsActive *bool  `json:"is_active,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Vérifier que le webhook appartient au marchand
	var existingMerchantID string
	err := db.QueryRow("SELECT merchant_id FROM webhooks WHERE id = $1", webhookID).Scan(&existingMerchantID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Webhook introuvable"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur serveur"})
		return
	}
	if existingMerchantID != merchantID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Accès non autorisé"})
		return
	}

	// Construire la requête de mise à jour
	updates := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.URL != "" {
		updates = append(updates, "url = $"+strconv.Itoa(argIndex))
		args = append(args, req.URL)
		argIndex++
	}
	if req.Secret != "" {
		updates = append(updates, "secret = $"+strconv.Itoa(argIndex))
		args = append(args, req.Secret)
		argIndex++
	}
	if req.IsActive != nil {
		updates = append(updates, "is_active = $"+strconv.Itoa(argIndex))
		args = append(args, *req.IsActive)
		argIndex++
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Aucune mise à jour fournie"})
		return
	}

	updates = append(updates, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, webhookID)

	query := "UPDATE webhooks SET " + joinStrings(updates, ", ") + " WHERE id = $" + strconv.Itoa(argIndex) + " RETURNING id, merchant_id, event_type, url, is_active, created_at, updated_at"

	var webhook Webhook
	err = db.QueryRow(query, args...).Scan(&webhook.ID, &webhook.MerchantID, &webhook.EventType, &webhook.URL, &webhook.IsActive, &webhook.CreatedAt, &webhook.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour"})
		return
	}

	c.JSON(http.StatusOK, webhook)
}

func handleDeleteWebhook(c *gin.Context) {
	webhookID := c.Param("id")
	merchantID := c.GetHeader("X-Merchant-ID")

	// Vérifier que le webhook appartient au marchand
	var existingMerchantID string
	err := db.QueryRow("SELECT merchant_id FROM webhooks WHERE id = $1", webhookID).Scan(&existingMerchantID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Webhook introuvable"})
		return
	}
	if existingMerchantID != merchantID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Accès non autorisé"})
		return
	}

	_, err = db.Exec("DELETE FROM webhooks WHERE id = $1", webhookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook supprimé"})
}

func handleTestWebhook(c *gin.Context) {
	webhookID := c.Param("id")
	merchantID := c.GetHeader("X-Merchant-ID")

	// Vérifier que le webhook appartient au marchand
	var webhook Webhook
	err := db.QueryRow(
		"SELECT id, merchant_id, event_type, url, secret, is_active FROM webhooks WHERE id = $1",
		webhookID,
	).Scan(&webhook.ID, &webhook.MerchantID, &webhook.EventType, &webhook.URL, &webhook.Secret, &webhook.IsActive)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Webhook introuvable"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur serveur"})
		return
	}
	if webhook.MerchantID != merchantID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Accès non autorisé"})
		return
	}

	// Envoyer un webhook de test
	testPayload := map[string]interface{}{
		"event_type": "test",
		"timestamp":  time.Now().Unix(),
		"data":       map[string]string{"message": "Test webhook"},
	}

	err = sendWebhook(webhook.URL, webhook.Secret, testPayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'envoi: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook de test envoyé avec succès"})
}

func sendWebhook(url, secret string, payload map[string]interface{}) error {
	// TODO: Implémenter l'envoi HTTP POST avec signature HMAC
	// Pour l'instant, on log juste
	log.Printf("Envoi webhook vers %s: %+v", url, payload)
	return nil
}

func generateSecret() string {
	// TODO: Générer un secret aléatoire sécurisé
	return "secret_" + time.Now().Format("20060102150405")
}

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

func authenticateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
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

