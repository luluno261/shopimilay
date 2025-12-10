package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// StorefrontConfig représente la configuration complète du storefront
type StorefrontConfig struct {
	MerchantID string    `json:"merchant_id"`
	Sections   []Section `json:"sections"`
	Theme      *Theme    `json:"theme,omitempty"`
}

// Section représente une section du storefront
type Section struct {
	ID    string                 `json:"id"`
	Type  string                 `json:"type"`
	Data  map[string]interface{} `json:"data"`
	Order int                    `json:"order"`
}

// Theme représente le thème personnalisé
type Theme struct {
	Colors     ThemeColors     `json:"colors"`
	Typography ThemeTypography  `json:"typography"`
	Layout     ThemeLayout      `json:"layout"`
}

type ThemeColors struct {
	Primary   string `json:"primary"`
	Secondary string `json:"secondary"`
	Background string `json:"background"`
	Text      string `json:"text"`
	Link      string `json:"link"`
	Button    string `json:"button"`
}

type ThemeTypography struct {
	FontFamily   string             `json:"font_family"`
	BaseSize     int                `json:"base_size"`
	HeadingSizes map[string]int     `json:"heading_sizes"`
	LineHeight   float64            `json:"line_height"`
	FontWeight   int                `json:"font_weight"`
}

type ThemeLayout struct {
	ContainerWidth int    `json:"container_width"`
	Padding         int    `json:"padding"`
	Margin          int    `json:"margin"`
	BorderRadius    int    `json:"border_radius"`
	Shadow          string `json:"shadow"`
}

// handleGetStorefrontConfig récupère la configuration du storefront
func handleGetStorefrontConfig(c *gin.Context) {
	merchantID := c.GetHeader("X-Merchant-ID")
	if merchantID == "" {
		// Essayer de récupérer depuis les query params pour le storefront public
		merchantID = c.Query("merchant_id")
	}
	if merchantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Merchant ID manquant"})
		return
	}

	var config StorefrontConfig
	var sectionsJSON sql.NullString
	var themeJSON sql.NullString

	err := db.QueryRow(
		"SELECT sections, theme FROM storefront_configs WHERE merchant_id = $1",
		merchantID,
	).Scan(&sectionsJSON, &themeJSON)

	if err == sql.ErrNoRows {
		// Pas de configuration, retourner une configuration vide
		c.JSON(http.StatusOK, gin.H{
			"sections": []Section{},
			"theme":    nil,
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération"})
		return
	}

	config.MerchantID = merchantID

	if sectionsJSON.Valid {
		json.Unmarshal([]byte(sectionsJSON.String), &config.Sections)
	}

	if themeJSON.Valid {
		json.Unmarshal([]byte(themeJSON.String), &config.Theme)
	}

	c.JSON(http.StatusOK, gin.H{
		"sections": config.Sections,
		"theme":    config.Theme,
	})
}

// handleSaveStorefrontConfig sauvegarde la configuration du storefront
func handleSaveStorefrontConfig(c *gin.Context) {
	merchantID := c.GetHeader("X-Merchant-ID")
	if merchantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Merchant ID manquant"})
		return
	}

	var req struct {
		Sections []Section `json:"sections"`
		Theme    *Theme    `json:"theme,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sectionsJSON, _ := json.Marshal(req.Sections)
	var themeJSON []byte
	if req.Theme != nil {
		themeJSON, _ = json.Marshal(req.Theme)
	}

	_, err := db.Exec(
		`INSERT INTO storefront_configs (merchant_id, sections, theme, updated_at)
		 VALUES ($1, $2, $3, NOW())
		 ON CONFLICT (merchant_id) 
		 DO UPDATE SET sections = $2, theme = $3, updated_at = NOW()`,
		merchantID,
		string(sectionsJSON),
		string(themeJSON),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la sauvegarde"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Configuration sauvegardée"})
}

// handleSaveTheme sauvegarde uniquement le thème
func handleSaveTheme(c *gin.Context) {
	merchantID := c.GetHeader("X-Merchant-ID")
	if merchantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Merchant ID manquant"})
		return
	}

	var theme Theme
	if err := c.ShouldBindJSON(&theme); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	themeJSON, _ := json.Marshal(theme)

	_, err := db.Exec(
		`INSERT INTO storefront_configs (merchant_id, theme, updated_at)
		 VALUES ($1, $2, NOW())
		 ON CONFLICT (merchant_id) 
		 DO UPDATE SET theme = $2, updated_at = NOW()`,
		merchantID,
		string(themeJSON),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la sauvegarde"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Thème sauvegardé"})
}

// handleGetTheme récupère uniquement le thème
func handleGetTheme(c *gin.Context) {
	merchantID := c.GetHeader("X-Merchant-ID")
	if merchantID == "" {
		// Essayer de récupérer depuis les query params pour le storefront public
		merchantID = c.Query("merchant_id")
	}
	if merchantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Merchant ID manquant"})
		return
	}

	var themeJSON sql.NullString
	err := db.QueryRow(
		"SELECT theme FROM storefront_configs WHERE merchant_id = $1",
		merchantID,
	).Scan(&themeJSON)

	if err == sql.ErrNoRows || !themeJSON.Valid {
		c.JSON(http.StatusOK, gin.H{"theme": nil})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération"})
		return
	}

	var theme Theme
	json.Unmarshal([]byte(themeJSON.String), &theme)

	c.JSON(http.StatusOK, gin.H{"theme": theme})
}

