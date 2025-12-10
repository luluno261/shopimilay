package importers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// ShopifyImporter gère l'import depuis Shopify
type ShopifyImporter struct {
	APIKey    string
	APISecret string
	StoreURL  string
	Client    *http.Client
}

// NewShopifyImporter crée un nouvel importeur Shopify
func NewShopifyImporter(apiKey, apiSecret, storeURL string) *ShopifyImporter {
	return &ShopifyImporter{
		APIKey:    apiKey,
		APISecret: apiSecret,
		StoreURL:  storeURL,
		Client:    &http.Client{Timeout: 30 * time.Second},
	}
}

// ImportProducts importe les produits depuis Shopify
func (si *ShopifyImporter) ImportProducts(merchantID string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("https://%s/admin/api/2023-10/products.json", si.StoreURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(si.APIKey, si.APISecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := si.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("erreur API Shopify: %s", string(bodyBytes))
	}

	var result struct {
		Products []map[string]interface{} `json:"products"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Transformer les produits Shopify en format OmniSphere
	products := make([]map[string]interface{}, 0)
	for _, shopifyProduct := range result.Products {
		// Extraire les variantes de manière sécurisée
		var sku, price interface{}
		if variants, ok := shopifyProduct["variants"].([]interface{}); ok && len(variants) > 0 {
			if variant, ok := variants[0].(map[string]interface{}); ok {
				sku = variant["sku"]
				price = variant["price"]
			}
		}

		product := map[string]interface{}{
			"merchant_id": merchantID,
			"currency":    "EUR",
			"images":      extractImages(shopifyProduct["images"]),
			"tags":        extractTags(shopifyProduct["tags"]),
		}

		// Extraire les champs de manière sécurisée
		if title, ok := shopifyProduct["title"].(string); ok {
			product["name"] = title
		}
		if bodyHTML, ok := shopifyProduct["body_html"].(string); ok {
			product["description"] = bodyHTML
		}
		if sku != nil {
			product["sku"] = sku
		}
		if price != nil {
			product["price"] = price
		}

		products = append(products, product)
	}

	return products, nil
}

func extractImages(images interface{}) []string {
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

func extractTags(tags interface{}) []string {
	if tags == nil {
		return []string{}
	}

	tagStr, ok := tags.(string)
	if !ok {
		return []string{}
	}

	// Les tags Shopify sont séparés par des virgules
	result := make([]string, 0)
	if tagStr != "" {
		// Parser les tags séparés par des virgules
		parts := strings.Split(tagStr, ",")
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed != "" {
				result = append(result, trimmed)
			}
		}
	}
	return result
}
