package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	port := getEnv("PORT", "8080")
	
	router := gin.Default()
	
	// Middleware global
	router.Use(corsMiddleware())
	router.Use(rateLimitMiddleware())
	
	// Routes de santé
	router.GET("/health", healthCheck)
	router.GET("/ready", readinessCheck)
	
	// Routes API publiques (sans authentification)
	public := router.Group("/api/v1")
	{
		// Auth routes
		public.POST("/auth/login", proxyToService("auth-service", "/api/v1/auth/login"))
		public.POST("/auth/register", proxyToService("auth-service", "/api/v1/auth/register"))
		public.POST("/auth/refresh", proxyToService("auth-service", "/api/v1/auth/refresh"))
		public.POST("/auth/delete-account", proxyToService("auth-service", "/api/v1/auth/delete-account"))
		
		// Catalogue routes (publiques)
		public.GET("/products", proxyToService("catalogue-service", "/api/v1/products"))
		public.GET("/products/:id", proxyToService("catalogue-service", "/api/v1/products/:id"))
		public.POST("/search", proxyToService("catalogue-service", "/api/v1/search"))
		
		// Store Builder routes (publiques pour le storefront)
		public.GET("/store-builder/config", proxyToService("catalogue-service", "/api/v1/store-builder/config"))
		public.GET("/store-builder/theme", proxyToService("catalogue-service", "/api/v1/store-builder/theme"))
	}
	
	// Routes API protégées (avec authentification)
	protected := router.Group("/api/v1")
	protected.Use(authenticateMiddleware())
	{
		// Auth routes protégées
		protected.GET("/auth/me", proxyToService("auth-service", "/api/v1/auth/me"))
		
		// Catalogue routes protégées
		protected.POST("/products", proxyToService("catalogue-service", "/api/v1/products"))
		protected.PUT("/products/:id", proxyToService("catalogue-service", "/api/v1/products/:id"))
		protected.DELETE("/products/:id", proxyToService("catalogue-service", "/api/v1/products/:id"))
		protected.PUT("/inventory/:productId", proxyToService("catalogue-service", "/api/v1/inventory/:productId"))
		
		// Checkout routes
		protected.GET("/cart", proxyToService("checkout-service", "/api/v1/cart"))
		protected.POST("/cart/items", proxyToService("checkout-service", "/api/v1/cart/items"))
		protected.DELETE("/cart/items/:itemId", proxyToService("checkout-service", "/api/v1/cart/items/:itemId"))
		protected.POST("/checkout", proxyToService("checkout-service", "/api/v1/checkout"))
		
		// Orders routes
		protected.GET("/orders", proxyToService("checkout-service", "/api/v1/orders"))
		protected.GET("/orders/:id", proxyToService("checkout-service", "/api/v1/orders/:id"))
		protected.PUT("/orders/:id/status", proxyToService("checkout-service", "/api/v1/orders/:id/status"))
		protected.POST("/orders/:id/refund", proxyToService("checkout-service", "/api/v1/orders/:id/refund"))
		
		// Account routes
		protected.GET("/account/orders", proxyToService("checkout-service", "/api/v1/account/orders"))
		protected.GET("/account/addresses", proxyToService("checkout-service", "/api/v1/account/addresses"))
		protected.POST("/account/addresses", proxyToService("checkout-service", "/api/v1/account/addresses"))
		
		// Dashboard routes
		protected.GET("/dashboard/stats", proxyToService("checkout-service", "/api/v1/dashboard/stats"))
		
		// Marketing routes
		protected.GET("/automation", proxyToService("marketing-engine", "/api/v1/automation"))
		protected.POST("/automation", proxyToService("marketing-engine", "/api/v1/automation"))
		protected.GET("/segmentation", proxyToService("marketing-engine", "/api/v1/segmentation"))
		protected.POST("/segmentation", proxyToService("marketing-engine", "/api/v1/segmentation"))
		
		// Webhooks routes
		protected.GET("/webhooks", proxyToService("webhook-service", "/api/v1/webhooks"))
		protected.POST("/webhooks", proxyToService("webhook-service", "/api/v1/webhooks"))
		protected.PUT("/webhooks/:id", proxyToService("webhook-service", "/api/v1/webhooks/:id"))
		protected.DELETE("/webhooks/:id", proxyToService("webhook-service", "/api/v1/webhooks/:id"))
		protected.POST("/webhooks/:id/test", proxyToService("webhook-service", "/api/v1/webhooks/:id/test"))
		
		// Discounts routes
		protected.GET("/discounts", proxyToService("checkout-service", "/api/v1/discounts"))
		protected.POST("/discounts", proxyToService("checkout-service", "/api/v1/discounts"))
		protected.DELETE("/discounts/:id", proxyToService("checkout-service", "/api/v1/discounts/:id"))
		
		// Migration routes
				protected.POST("/migration/import", proxyToService("migration-tool", "/api/v1/migration/import"))
				protected.GET("/migration/status/:id", proxyToService("migration-tool", "/api/v1/migration/status/:id"))
				protected.GET("/store-builder/config", proxyToService("catalogue-service", "/api/v1/store-builder/config"))
				protected.POST("/store-builder/config", proxyToService("catalogue-service", "/api/v1/store-builder/config"))
				protected.GET("/store-builder/theme", proxyToService("catalogue-service", "/api/v1/store-builder/theme"))
				protected.POST("/store-builder/theme", proxyToService("catalogue-service", "/api/v1/store-builder/theme"))
			}
	
	// Webhooks (sans authentification mais avec signature)
	router.POST("/api/v1/webhooks/stripe", proxyToService("checkout-service", "/api/v1/checkout/stripe/webhook"))
	
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	go func() {
		log.Printf("API Gateway démarré sur le port %s", port)
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
		"service": "api-gateway",
	})
}

func readinessCheck(c *gin.Context) {
	// TODO: Vérifier la disponibilité des services backend
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
		"service": "api-gateway",
	})
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

