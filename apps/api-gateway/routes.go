package main

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ServiceConfig contient la configuration d'un service backend
type ServiceConfig struct {
	Host string
	Port string
}

// getServiceURL retourne l'URL complète d'un service
func getServiceURL(serviceName string) string {
	// En développement, utiliser localhost avec les ports par défaut
	servicePorts := map[string]string{
		"auth-service":      "8080",
		"checkout-service":  "8081",
		"catalogue-service": "8082",
		"marketing-engine":  "8083",
		"webhook-service":   "8084",
		"migration-tool":    "8085",
	}

	port := getEnv(strings.ToUpper(serviceName)+"_PORT", servicePorts[serviceName])
	host := getEnv(strings.ToUpper(serviceName)+"_HOST", "localhost")

	return "http://" + host + ":" + port
}

// proxyToService crée un handler qui proxy les requêtes vers un service backend
func proxyToService(serviceName, path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceURL := getServiceURL(serviceName)
		targetURL := serviceURL + path

		// Remplacer les paramètres de route
		for _, param := range c.Params {
			targetURL = strings.Replace(targetURL, ":"+param.Key, param.Value, -1)
		}

		// Copier le body de la requête
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Créer la requête vers le service backend
		req, err := http.NewRequest(c.Request.Method, targetURL, bytes.NewBuffer(bodyBytes))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de la requête"})
			return
		}

		// Copier les en-têtes
		for key, values := range c.Request.Header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		// Ajouter les informations utilisateur depuis le contexte (si authentifié)
		if userID, exists := c.Get("user_id"); exists {
			req.Header.Set("X-User-ID", userID.(string))
		}
		if merchantID, exists := c.Get("merchant_id"); exists {
			req.Header.Set("X-Merchant-ID", merchantID.(string))
		}

		// Exécuter la requête
		client := &http.Client{
			Timeout: 30 * http.DefaultClient.Timeout,
		}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "Service backend indisponible"})
			return
		}
		defer resp.Body.Close()

		// Copier les en-têtes de la réponse
		for key, values := range resp.Header {
			for _, value := range values {
				c.Writer.Header().Add(key, value)
			}
		}

		// Copier le status code et le body
		c.Status(resp.StatusCode)
		io.Copy(c.Writer, resp.Body)
	}
}

