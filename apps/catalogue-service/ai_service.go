package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// AIService gère les interactions avec l'API OpenAI pour générer des descriptions produits
type AIService struct {
	APIKey  string
	BaseURL string
	Client  *http.Client
}

// NewAIService crée un nouveau service AI
func NewAIService() *AIService {
	apiKey := getEnv("OPENAI_API_KEY", "")
	baseURL := getEnv("OPENAI_BASE_URL", "https://api.openai.com/v1")
	
	return &AIService{
		APIKey:  apiKey,
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

// GenerateProductDescription génère une description SEO optimisée pour un produit
func (ai *AIService) GenerateProductDescription(productName string, productDetails map[string]interface{}) (string, error) {
	if ai.APIKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY non configuré")
	}

	// Construire le prompt
	prompt := fmt.Sprintf(
		"Génère une description SEO optimisée pour le produit suivant:\n\n"+
			"Nom: %s\n"+
			"Détails: %v\n\n"+
			"La description doit être:\n"+
			"- Attrayante et persuasive\n"+
			"- Optimisée pour le SEO (inclure des mots-clés pertinents)\n"+
			"- Entre 150 et 300 mots\n"+
			"- En français\n"+
			"- Formatée en paragraphes courts",
		productName, productDetails,
	)

	// Préparer la requête vers l'API OpenAI
	requestBody := map[string]interface{}{
		"model": "gpt-4",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "Tu es un expert en rédaction e-commerce et SEO. Tu génères des descriptions produits optimisées.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"max_tokens": 500,
		"temperature": 0.7,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", ai.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ai.APIKey)

	resp, err := ai.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("erreur API OpenAI: %s", string(bodyBytes))
	}

	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("aucune réponse de l'API")
	}

	return response.Choices[0].Message.Content, nil
}

// GenerateProductDescriptionRequest représente une demande de génération
type GenerateProductDescriptionRequest struct {
	ProductName string                 `json:"product_name" binding:"required"`
	Details     map[string]interface{} `json:"details,omitempty"`
}

// HandleGenerateDescription gère la requête de génération de description
func HandleGenerateDescription(c *gin.Context) {
	merchantID := c.GetHeader("X-Merchant-ID")
	if merchantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Merchant ID manquant"})
		return
	}

	var req GenerateProductDescriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	aiService := NewAIService()
	description, err := aiService.GenerateProductDescription(req.ProductName, req.Details)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la génération: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"description": description,
	})
}

