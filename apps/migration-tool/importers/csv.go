package importers

import (
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"strings"
)

// CSVImporter gère l'import depuis un fichier CSV
type CSVImporter struct{}

// NewCSVImporter crée un nouvel importeur CSV
func NewCSVImporter() *CSVImporter {
	return &CSVImporter{}
}

// ImportProducts importe les produits depuis un CSV
func (ci *CSVImporter) ImportProducts(merchantID string, csvData string) ([]map[string]interface{}, error) {
	// Décoder le base64
	data, err := base64.StdEncoding.DecodeString(csvData)
	if err != nil {
		return nil, fmt.Errorf("erreur de décodage base64: %v", err)
	}

	// Parser le CSV
	reader := csv.NewReader(strings.NewReader(string(data)))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("erreur de parsing CSV: %v", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("CSV vide ou sans en-têtes")
	}

	// Première ligne = en-têtes
	headers := records[0]
	
	// Trouver les index des colonnes
	nameIdx := findIndex(headers, "name", "nom", "title", "titre")
	descIdx := findIndex(headers, "description", "desc")
	skuIdx := findIndex(headers, "sku", "reference")
	priceIdx := findIndex(headers, "price", "prix")
	imageIdx := findIndex(headers, "image", "images", "photo")

	if nameIdx == -1 || skuIdx == -1 || priceIdx == -1 {
		return nil, fmt.Errorf("colonnes requises manquantes (name, sku, price)")
	}

	products := make([]map[string]interface{}, 0)
	for i := 1; i < len(records); i++ {
		row := records[i]
		if len(row) < len(headers) {
			continue // Ignorer les lignes incomplètes
		}

		product := map[string]interface{}{
			"merchant_id": merchantID,
			"name":        row[nameIdx],
			"sku":         row[skuIdx],
			"price":       parseFloat(row[priceIdx]),
			"currency":    "EUR",
		}

		if descIdx >= 0 && descIdx < len(row) {
			product["description"] = row[descIdx]
		}

		if imageIdx >= 0 && imageIdx < len(row) {
			images := strings.Split(row[imageIdx], ",")
			product["images"] = images
		}

		products = append(products, product)
	}

	return products, nil
}

func findIndex(headers []string, values ...string) int {
	for i, header := range headers {
		headerLower := strings.ToLower(strings.TrimSpace(header))
		for _, value := range values {
			if headerLower == strings.ToLower(value) {
				return i
			}
		}
	}
	return -1
}

func parseFloat(s string) float64 {
	var result float64
	fmt.Sscanf(s, "%f", &result)
	return result
}

