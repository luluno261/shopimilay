package main

import "time"

// Product représente un produit dans le catalogue
type Product struct {
	ID          string    `json:"id" db:"id"`
	MerchantID  string    `json:"merchant_id" db:"merchant_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	SKU         string    `json:"sku" db:"sku"`
	Price       float64   `json:"price" db:"price"`
	Currency    string    `json:"currency" db:"currency"`
	CategoryID  string    `json:"category_id" db:"category_id"`
	Images      []string  `json:"images" db:"images"`
	Tags        []string  `json:"tags" db:"tags"`
	Status      string    `json:"status" db:"status"` // active, inactive, draft
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// ProductVariant représente une variante de produit (taille, couleur, etc.)
type ProductVariant struct {
	ID        string    `json:"id" db:"id"`
	ProductID string    `json:"product_id" db:"product_id"`
	Name      string    `json:"name" db:"name"`
	SKU       string    `json:"sku" db:"sku"`
	Price     float64   `json:"price" db:"price"`
	Stock     int       `json:"stock" db:"stock"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Category représente une catégorie de produits
type Category struct {
	ID          string    `json:"id" db:"id"`
	MerchantID  string    `json:"merchant_id" db:"merchant_id"`
	Name        string    `json:"name" db:"name"`
	Slug        string    `json:"slug" db:"slug"`
	ParentID    *string   `json:"parent_id,omitempty" db:"parent_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CreateProductRequest représente une demande de création de produit
type CreateProductRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	SKU         string   `json:"sku" binding:"required"`
	Price       float64  `json:"price" binding:"required"`
	Currency    string   `json:"currency" binding:"required"`
	CategoryID  string   `json:"category_id"`
	Images      []string `json:"images"`
	Tags        []string `json:"tags"`
}

// UpdateProductRequest représente une demande de mise à jour de produit
type UpdateProductRequest struct {
	Name        *string   `json:"name"`
	Description *string   `json:"description"`
	Price       *float64  `json:"price"`
	CategoryID  *string   `json:"category_id"`
	Images      *[]string `json:"images"`
	Tags        *[]string `json:"tags"`
	Status      *string   `json:"status"`
}

// CreateProduct crée un nouveau produit (utilise CreateProductDB)
func CreateProduct(merchantID string, req *CreateProductRequest) (*Product, error) {
	return CreateProductDB(merchantID, req)
}

// GetProduct récupère un produit par ID (utilise GetProductByID)
func GetProduct(productID string) (*Product, error) {
	return GetProductByID(productID)
}

// UpdateProduct met à jour un produit (utilise UpdateProductDB)
func UpdateProduct(productID string, req *UpdateProductRequest) (*Product, error) {
	return UpdateProductDB(productID, req)
}

// DeleteProduct supprime un produit (utilise DeleteProductDB)
func DeleteProduct(productID string) error {
	return DeleteProductDB(productID)
}

// ListProducts liste les produits avec pagination (utilise ListProductsDB)
func ListProducts(merchantID string, limit, offset int) ([]Product, error) {
	return ListProductsDB(merchantID, limit, offset)
}

