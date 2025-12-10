package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// ElasticsearchClient gère les interactions avec Elasticsearch
type ElasticsearchClient struct {
	BaseURL string
	Client  *http.Client
}

// NewElasticsearchClient crée un nouveau client Elasticsearch
func NewElasticsearchClient() *ElasticsearchClient {
	baseURL := getEnv("ELASTICSEARCH_URL", "http://localhost:9200")
	return &ElasticsearchClient{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

// IndexProduct indexe un produit dans Elasticsearch
func (es *ElasticsearchClient) IndexProduct(product *Product) error {
	indexName := "products"
	url := fmt.Sprintf("%s/%s/_doc/%s", es.BaseURL, indexName, product.ID)
	
	body, err := json.Marshal(product)
	if err != nil {
		return err
	}
	
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := es.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("erreur Elasticsearch: %s", string(bodyBytes))
	}
	
	return nil
}

// DeleteProduct supprime un produit de l'index Elasticsearch
func (es *ElasticsearchClient) DeleteProduct(productID string) error {
	indexName := "products"
	url := fmt.Sprintf("%s/%s/_doc/%s", es.BaseURL, indexName, productID)
	
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	
	resp, err := es.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode >= 400 && resp.StatusCode != 404 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("erreur Elasticsearch: %s", string(bodyBytes))
	}
	
	return nil
}

// SearchProducts recherche des produits dans Elasticsearch
func (es *ElasticsearchClient) SearchProducts(query string, filters map[string]interface{}) ([]Product, error) {
	indexName := "products"
	url := fmt.Sprintf("%s/%s/_search", es.BaseURL, indexName)
	
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"name^2", "description", "tags"},
			},
		},
	}
	
	body, err := json.Marshal(searchQuery)
	if err != nil {
		return nil, err
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := es.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("erreur Elasticsearch: %s", string(bodyBytes))
	}
	
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	
	var products []Product
	hits, ok := result["hits"].(map[string]interface{})
	if !ok {
		return products, nil
	}
	
	hitsList, ok := hits["hits"].([]interface{})
	if !ok {
		return products, nil
	}
	
	for _, hit := range hitsList {
		hitMap, ok := hit.(map[string]interface{})
		if !ok {
			continue
		}
		
		source, ok := hitMap["_source"].(map[string]interface{})
		if !ok {
			continue
		}
		
		productJSON, err := json.Marshal(source)
		if err != nil {
			continue
		}
		
		var product Product
		if err := json.Unmarshal(productJSON, &product); err != nil {
			continue
		}
		
		products = append(products, product)
	}
	
	return products, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

