package importers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// WooCommerceImporter gère l'import depuis WooCommerce
type WooCommerceImporter struct {
	ConsumerKey    string
	ConsumerSecret string
	StoreURL       string
	Client         *http.Client
}

// NewWooCommerceImporter crée un nouvel importeur WooCommerce
func NewWooCommerceImporter(consumerKey, consumerSecret, storeURL string) *WooCommerceImporter {
	return &WooCommerceImporter{
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
		StoreURL:       storeURL,
		Client:         &http.Client{Timeout: 30 * time.Second},
	}
}

// ImportProducts importe les produits depuis WooCommerce
func (wi *WooCommerceImporter) ImportProducts(merchantID string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/wp-json/wc/v3/products", wi.StoreURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Authentification Basic Auth pour WooCommerce
	auth := base64.StdEncoding.EncodeToString([]byte(wi.ConsumerKey + ":" + wi.ConsumerSecret))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")

	resp, err := wi.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("erreur API WooCommerce: %s", string(bodyBytes))
	}

	var products []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return nil, err
	}

	// Transformer les produits WooCommerce en format OmniSphere
	transformed := make([]map[string]interface{}, 0)
	for _, wcProduct := range products {
		product := map[string]interface{}{
			"merchant_id": merchantID,
			"currency":    "EUR",
			"images":      extractWooCommerceImages(wcProduct["images"]),
		}

		// Extraire les champs de manière sécurisée
		if name, ok := wcProduct["name"].(string); ok {
			product["name"] = name
		}
		if desc, ok := wcProduct["description"].(string); ok {
			product["description"] = desc
		}
		if sku, ok := wcProduct["sku"].(string); ok {
			product["sku"] = sku
		}
		if price, ok := wcProduct["price"].(string); ok {
			product["price"] = price
		} else if price, ok := wcProduct["price"].(float64); ok {
			product["price"] = price
		}

		transformed = append(transformed, product)
	}

	return transformed, nil
}

func extractWooCommerceImages(images interface{}) []string {
	if images == nil {
		return []string{}
	}

	imgList, ok := images.([]interface{})
	if !ok {
		return []string{}
	}

	result := make([]string, 0)
	for _, img := range imgList {
		if imgMap, ok := img.(map[string]interface{}); ok {
			if src, ok := imgMap["src"].(string); ok {
				result = append(result, src)
			}
		}
	}
	return result
}
