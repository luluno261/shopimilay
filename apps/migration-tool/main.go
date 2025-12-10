package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	port := getEnv("PORT", "8085")

	router := gin.Default()

	// Routes de santé
	router.GET("/health", healthCheck)
	router.GET("/ready", readinessCheck)

	// Routes API
	api := router.Group("/api/v1")
	{
		api.POST("/migration/import", authenticateMiddleware(), handleImport)
		api.GET("/migration/status/:id", authenticateMiddleware(), handleGetMigrationStatus)
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Migration Tool démarré sur le port %s", port)
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
		"status":  "healthy",
		"service": "migration-tool",
	})
}

func readinessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ready",
		"service": "migration-tool",
	})
}

// ImportRequest représente une demande d'import
type ImportRequest struct {
	Source    string `json:"source" binding:"required,oneof=shopify woocommerce csv json"`
	FileURL   string `json:"file_url,omitempty"`
	FileData  string `json:"file_data,omitempty"` // Base64 encoded
	APIKey    string `json:"api_key,omitempty"`   // Pour Shopify/WooCommerce
	APISecret string `json:"api_secret,omitempty"`
	StoreURL  string `json:"store_url,omitempty"`
}

func handleImport(c *gin.Context) {
	merchantID := c.GetHeader("X-Merchant-ID")
	if merchantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Merchant ID manquant"})
		return
	}

	var req ImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Créer un job de migration
	migrationID := generateMigrationID()

	// Traiter l'import de manière asynchrone
	go processImport(migrationID, merchantID, req)

	c.JSON(http.StatusAccepted, gin.H{
		"migration_id": migrationID,
		"status":       "processing",
		"message":      "Import démarré",
	})
}

func processImport(migrationID, merchantID string, req ImportRequest) {
	log.Printf("Démarrage import %s pour merchant %s depuis %s", migrationID, merchantID, req.Source)

	switch req.Source {
	case "shopify":
		importFromShopify(migrationID, merchantID, req)
	case "woocommerce":
		importFromWooCommerce(migrationID, merchantID, req)
	case "csv":
		importFromCSV(migrationID, merchantID, req)
	case "json":
		importFromJSON(migrationID, merchantID, req)
	}
}

func importFromShopify(migrationID, merchantID string, req ImportRequest) {
	log.Printf("Import Shopify pour migration %s", migrationID)
	// TODO: Implémenter avec le package importers une fois le module configuré
	log.Printf("Import Shopify - À implémenter avec importers.NewShopifyImporter")
}

func importFromWooCommerce(migrationID, merchantID string, req ImportRequest) {
	log.Printf("Import WooCommerce pour migration %s", migrationID)
	// TODO: Implémenter avec le package importers une fois le module configuré
	log.Printf("Import WooCommerce - À implémenter avec importers.NewWooCommerceImporter")
}

func importFromCSV(migrationID, merchantID string, req ImportRequest) {
	log.Printf("Import CSV pour migration %s", migrationID)

	if req.FileData == "" {
		log.Printf("Aucune donnée CSV fournie")
		return
	}

	// TODO: Implémenter avec le package importers une fois le module configuré
	log.Printf("Import CSV - À implémenter avec importers.NewCSVImporter")
}

func importFromJSON(migrationID, merchantID string, req ImportRequest) {
	log.Printf("Import JSON pour migration %s", migrationID)

	if req.FileData == "" {
		log.Printf("Aucune donnée JSON fournie")
		return
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(req.FileData), &data); err != nil {
		log.Printf("Erreur parsing JSON: %v", err)
		return
	}

	// Extraire les produits du JSON
	products := extractProductsFromJSON(merchantID, data)
	createProducts(merchantID, products, migrationID)
}

// createProducts crée les produits via l'API catalogue-service
func createProducts(merchantID string, products []map[string]interface{}, migrationID string) {
	_ = getEnv("CATALOGUE_SERVICE_URL", "http://localhost:8082")

	for _, product := range products {
		// TODO: Appeler l'API catalogue-service pour créer le produit
		log.Printf("Création produit: %v", product["name"])
	}
}

// extractProductsFromJSON extrait les produits d'un JSON
func extractProductsFromJSON(merchantID string, data map[string]interface{}) []map[string]interface{} {
	products := make([]map[string]interface{}, 0)

	// Format attendu: { "products": [...] }
	if productsList, ok := data["products"].([]interface{}); ok {
		for _, p := range productsList {
			if productMap, ok := p.(map[string]interface{}); ok {
				productMap["merchant_id"] = merchantID
				products = append(products, productMap)
			}
		}
	}

	return products
}

func handleGetMigrationStatus(c *gin.Context) {
	migrationID := c.Param("id")

	// TODO: Récupérer le statut depuis la base de données ou un cache
	c.JSON(http.StatusOK, gin.H{
		"migration_id":   migrationID,
		"status":         "completed",
		"progress":       100,
		"items_imported": 0,
		"items_failed":   0,
	})
}

func generateMigrationID() string {
	return "migration_" + time.Now().Format("20060102150405")
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
