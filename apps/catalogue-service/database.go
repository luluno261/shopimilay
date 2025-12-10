package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

var db *sql.DB

// initDB initialise la connexion à la base de données
func initDB() {
	connStr := getEnv("DATABASE_URL", "postgres://omnisphere:omnisphere_dev@localhost:5432/omnisphere?sslmode=disable")
	
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Erreur de connexion à la base de données:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Erreur de ping à la base de données:", err)
	}

	log.Println("Connexion à la base de données établie")
}

// ProductDB représente un produit en base de données
type ProductDB struct {
	ID          string    `db:"id"`
	MerchantID  string    `db:"merchant_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	SKU         string    `db:"sku"`
	Price       float64   `db:"price"`
	Currency    string    `db:"currency"`
	CategoryID  *string   `db:"category_id"`
	Images      []string  `db:"images"`
	Tags        []string  `db:"tags"`
	Status      string    `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// GetProductByID récupère un produit par ID
func GetProductByID(productID string) (*Product, error) {
	var p ProductDB
	var imagesArray pq.StringArray
	var tagsArray pq.StringArray
	var categoryID sql.NullString

	err := db.QueryRow(
		"SELECT id, merchant_id, name, description, sku, price, currency, category_id, images, tags, status, created_at, updated_at FROM products WHERE id = $1",
		productID,
	).Scan(&p.ID, &p.MerchantID, &p.Name, &p.Description, &p.SKU, &p.Price, &p.Currency, &categoryID, &imagesArray, &tagsArray, &p.Status, &p.CreatedAt, &p.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	product := &Product{
		ID:          p.ID,
		MerchantID:  p.MerchantID,
		Name:        p.Name,
		Description: p.Description,
		SKU:         p.SKU,
		Price:       p.Price,
		Currency:    p.Currency,
		Status:      p.Status,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		Images:      []string(imagesArray),
		Tags:        []string(tagsArray),
	}

	if categoryID.Valid {
		product.CategoryID = categoryID.String
	}

	return product, nil
}

// GetProductBySlug récupère un produit par slug (généré à partir du nom)
func GetProductBySlug(slug string) (*Product, error) {
	// Pour l'instant, on utilise l'ID. Le slug sera implémenté plus tard
	return GetProductByID(slug)
}

// CreateProductDB crée un produit en base de données
func CreateProductDB(merchantID string, req *CreateProductRequest) (*Product, error) {
	imagesArray := pq.StringArray(req.Images)
	tagsArray := pq.StringArray(req.Tags)

	var categoryID *string
	if req.CategoryID != "" {
		categoryID = &req.CategoryID
	}

	var p ProductDB
	var imagesResult pq.StringArray
	var tagsResult pq.StringArray
	var categoryIDResult sql.NullString

	err := db.QueryRow(
		"INSERT INTO products (merchant_id, name, description, sku, price, currency, category_id, images, tags, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 'active') RETURNING id, merchant_id, name, description, sku, price, currency, category_id, images, tags, status, created_at, updated_at",
		merchantID, req.Name, req.Description, req.SKU, req.Price, req.Currency, categoryID, imagesArray, tagsArray,
	).Scan(&p.ID, &p.MerchantID, &p.Name, &p.Description, &p.SKU, &p.Price, &p.Currency, &categoryIDResult, &imagesResult, &tagsResult, &p.Status, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		return nil, err
	}

	product := &Product{
		ID:          p.ID,
		MerchantID:  p.MerchantID,
		Name:        p.Name,
		Description: p.Description,
		SKU:         p.SKU,
		Price:       p.Price,
		Currency:    p.Currency,
		Status:      p.Status,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		Images:      []string(imagesResult),
		Tags:        []string(tagsResult),
	}

	if categoryIDResult.Valid {
		product.CategoryID = categoryIDResult.String
	}

	return product, nil
}

// UpdateProductDB met à jour un produit
func UpdateProductDB(productID string, req *UpdateProductRequest) (*Product, error) {
	// Construire la requête dynamiquement
	updates := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Name != nil {
		updates = append(updates, "name = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Name)
		argIndex++
	}
	if req.Description != nil {
		updates = append(updates, "description = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Description)
		argIndex++
	}
	if req.Price != nil {
		updates = append(updates, "price = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Price)
		argIndex++
	}
	if req.CategoryID != nil {
		updates = append(updates, "category_id = $"+strconv.Itoa(argIndex))
		args = append(args, *req.CategoryID)
		argIndex++
	}
	if req.Images != nil {
		imagesArray := pq.StringArray(*req.Images)
		updates = append(updates, "images = $"+strconv.Itoa(argIndex))
		args = append(args, imagesArray)
		argIndex++
	}
	if req.Tags != nil {
		tagsArray := pq.StringArray(*req.Tags)
		updates = append(updates, "tags = $"+strconv.Itoa(argIndex))
		args = append(args, tagsArray)
		argIndex++
	}
	if req.Status != nil {
		updates = append(updates, "status = $"+strconv.Itoa(argIndex))
		args = append(args, *req.Status)
		argIndex++
	}

	if len(updates) == 0 {
		return GetProductByID(productID)
	}

	updates = append(updates, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, productID)

	query := "UPDATE products SET " + joinStrings(updates, ", ") + " WHERE id = $" + strconv.Itoa(argIndex) + " RETURNING id, merchant_id, name, description, sku, price, currency, category_id, images, tags, status, created_at, updated_at"

	var p ProductDB
	var imagesArray pq.StringArray
	var tagsArray pq.StringArray
	var categoryID sql.NullString

	err := db.QueryRow(query, args...).Scan(&p.ID, &p.MerchantID, &p.Name, &p.Description, &p.SKU, &p.Price, &p.Currency, &categoryID, &imagesArray, &tagsArray, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}

	product := &Product{
		ID:          p.ID,
		MerchantID:  p.MerchantID,
		Name:        p.Name,
		Description: p.Description,
		SKU:         p.SKU,
		Price:       p.Price,
		Currency:    p.Currency,
		Status:      p.Status,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
		Images:      []string(imagesArray),
		Tags:        []string(tagsArray),
	}

	if categoryID.Valid {
		product.CategoryID = categoryID.String
	}

	return product, nil
}

// DeleteProductDB supprime un produit
func DeleteProductDB(productID string) error {
	_, err := db.Exec("DELETE FROM products WHERE id = $1", productID)
	return err
}

// ListProductsDB liste les produits avec pagination
func ListProductsDB(merchantID string, limit, offset int) ([]Product, error) {
	rows, err := db.Query(
		"SELECT id, merchant_id, name, description, sku, price, currency, category_id, images, tags, status, created_at, updated_at FROM products WHERE merchant_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3",
		merchantID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p ProductDB
		var imagesArray pq.StringArray
		var tagsArray pq.StringArray
		var categoryID sql.NullString

		err := rows.Scan(&p.ID, &p.MerchantID, &p.Name, &p.Description, &p.SKU, &p.Price, &p.Currency, &categoryID, &imagesArray, &tagsArray, &p.Status, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}

		product := Product{
			ID:          p.ID,
			MerchantID:  p.MerchantID,
			Name:        p.Name,
			Description: p.Description,
			SKU:         p.SKU,
			Price:       p.Price,
			Currency:    p.Currency,
			Status:      p.Status,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
			Images:      []string(imagesArray),
			Tags:        []string(tagsArray),
		}

		if categoryID.Valid {
			product.CategoryID = categoryID.String
		}

		products = append(products, product)
	}

	return products, nil
}

// Helper functions
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

