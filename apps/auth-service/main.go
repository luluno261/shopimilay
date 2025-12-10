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
	// Initialisation de la base de données
	initDB()
	defer db.Close()
	
	// Configuration du serveur
	port := getEnv("PORT", "8080")
	
	// Initialisation du routeur Gin
	router := gin.Default()
	
	// Routes de santé
	router.GET("/health", healthCheck)
	router.GET("/ready", readinessCheck)
	
	// Routes d'authentification
	api := router.Group("/api/v1")
	{
		api.POST("/auth/login", handleLogin)
		api.POST("/auth/register", handleRegister)
		api.POST("/auth/refresh", handleRefreshToken)
		api.GET("/auth/me", authenticateMiddleware(), handleGetMe)
		api.POST("/auth/delete-account", handleDeleteAccount)
	}
	
	// Configuration du serveur HTTP
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	// Démarrage du serveur en goroutine
	go func() {
		log.Printf("Auth Service démarré sur le port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erreur serveur: %v", err)
		}
	}()
	
	// Attente d'un signal d'interruption pour arrêt gracieux
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

// healthCheck vérifie l'état de santé du service
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"service": "auth-service",
	})
}

// readinessCheck vérifie si le service est prêt à recevoir du trafic
func readinessCheck(c *gin.Context) {
	if err := db.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not ready",
			"service": "auth-service",
			"error": "Database connection failed",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
		"service": "auth-service",
	})
}

// LoginRequest représente une demande de connexion
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// handleLogin gère la connexion d'un utilisateur
func handleLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur serveur"})
		return
	}

	if user == nil || !CheckPassword(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email ou mot de passe incorrect"})
		return
	}

	token, err := GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la génération du token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"role":       user.Role,
			"merchant_id": user.MerchantID,
		},
	})
}

// RegisterRequest représente une demande d'inscription
type RegisterRequest struct {
	Email     string  `json:"email" binding:"required,email"`
	Password  string  `json:"password" binding:"required,min=8"`
	Role      string  `json:"role"`
	MerchantID *string `json:"merchant_id,omitempty"`
}

// handleRegister gère l'inscription d'un nouvel utilisateur
func handleRegister(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Vérifier si l'utilisateur existe déjà
	existingUser, err := GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur serveur"})
		return
	}
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Cet email est déjà utilisé"})
		return
	}

	// Hash le mot de passe
	passwordHash, err := HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du hashage du mot de passe"})
		return
	}

	// Définir le rôle par défaut
	role := req.Role
	if role == "" {
		role = "user"
	}

	// Créer l'utilisateur
	user, err := CreateUser(req.Email, passwordHash, role, req.MerchantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de l'utilisateur"})
		return
	}

	token, err := GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la génération du token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"token": token,
		"user": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"role":       user.Role,
			"merchant_id": user.MerchantID,
		},
	})
}

// RefreshTokenRequest représente une demande de rafraîchissement de token
type RefreshTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

// handleRefreshToken gère le rafraîchissement du token JWT
func handleRefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, err := ValidateToken(req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalide"})
		return
	}

	// Récupérer l'utilisateur
	user, err := GetUserByID(claims.UserID)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur introuvable"})
		return
	}

	// Générer un nouveau token
	newToken, err := GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la génération du token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": newToken,
	})
}

// handleGetMe retourne les informations de l'utilisateur connecté
func handleGetMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	user, err := GetUserByID(userID.(string))
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur introuvable"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"email":      user.Email,
		"role":       user.Role,
		"merchant_id": user.MerchantID,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}

// authenticateMiddleware vérifie et valide le token JWT
func authenticateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header manquant"})
			c.Abort()
			return
		}

		// Extraire le token (format: "Bearer <token>")
		tokenString := authHeader
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		}

		claims, err := ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalide"})
			c.Abort()
			return
		}

		// Ajouter les informations utilisateur au contexte
		c.Set("user_id", claims.UserID)
		c.Set("merchant_id", claims.MerchantID)
		c.Set("role", claims.Role)
		c.Set("email", claims.Email)

		c.Next()
	}
}

// DeleteAccountRequest représente une demande de suppression de compte
type DeleteAccountRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// handleDeleteAccount gère la demande de suppression de compte (RGPD)
func handleDeleteAccount(c *gin.Context) {
	var req DeleteAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur serveur"})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur introuvable"})
		return
	}

	// Supprimer l'utilisateur et toutes ses données associées
	// TODO: Supprimer aussi les données dans les autres services (commandes, paniers, etc.)
	_, err = db.Exec("DELETE FROM users WHERE id = $1", user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la suppression"})
		return
	}

	// Supprimer les refresh tokens
	_, err = db.Exec("DELETE FROM refresh_tokens WHERE user_id = $1", user.ID)
	if err != nil {
		log.Printf("Erreur lors de la suppression des tokens: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Compte supprimé avec succès"})
}

// getEnv récupère une variable d'environnement ou retourne une valeur par défaut
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

