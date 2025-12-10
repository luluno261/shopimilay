package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

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

// CartItemDB représente un article de panier en base
type CartItemDB struct {
	ID        string  `db:"id"`
	CartID    string  `db:"cart_id"`
	ProductID string  `db:"product_id"`
	VariantID *string `db:"variant_id"`
	Quantity  int     `db:"quantity"`
	Price     float64 `db:"price"`
}

// Order représente une commande
type Order struct {
	ID             string          `json:"id" db:"id"`
	UserID         string          `json:"user_id" db:"user_id"`
	MerchantID     string          `json:"merchant_id" db:"merchant_id"`
	Status         string          `json:"status" db:"status"`
	TotalAmount    float64         `json:"total_amount" db:"total_amount"`
	Currency       string          `json:"currency" db:"currency"`
	PaymentIntentID *string        `json:"payment_intent_id,omitempty" db:"payment_intent_id"`
	ShippingAddress json.RawMessage `json:"shipping_address" db:"shipping_address"`
	BillingAddress  json.RawMessage `json:"billing_address" db:"billing_address"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at" db:"updated_at"`
}

// GetOrCreateCart récupère ou crée un panier pour un utilisateur
func GetOrCreateCart(userID, merchantID string) (*Cart, error) {
	var cart Cart
	var cartID string
	var totalAmount float64
	var currency string
	var createdAt, updatedAt time.Time

	err := db.QueryRow(
		"SELECT id, total_amount, currency, created_at, updated_at FROM carts WHERE user_id = $1 AND merchant_id = $2 AND status = 'active'",
		userID, merchantID,
	).Scan(&cartID, &totalAmount, &currency, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		// Créer un nouveau panier
		err = db.QueryRow(
			"INSERT INTO carts (user_id, merchant_id, status, total_amount, currency) VALUES ($1, $2, 'active', 0, 'EUR') RETURNING id, total_amount, currency, created_at, updated_at",
			userID, merchantID,
		).Scan(&cartID, &totalAmount, &currency, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	cart.ID = cartID
	cart.UserID = userID
	cart.Total = totalAmount
	cart.Currency = currency
	cart.CreatedAt = createdAt.Format(time.RFC3339)
	cart.UpdatedAt = updatedAt.Format(time.RFC3339)

	// Récupérer les articles du panier
	rows, err := db.Query(
		"SELECT id, product_id, variant_id, quantity, price FROM cart_items WHERE cart_id = $1",
		cartID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []CartItem
	for rows.Next() {
		var item CartItemDB
		var variantID sql.NullString
		err := rows.Scan(&item.ID, &item.ProductID, &variantID, &item.Quantity, &item.Price)
		if err != nil {
			return nil, err
		}
		if variantID.Valid {
			v := variantID.String
			item.VariantID = &v
		}

		items = append(items, CartItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
			Subtotal:  item.Price * float64(item.Quantity),
		})
	}

	cart.Items = items
	return &cart, nil
}

// AddItemToCart ajoute un article au panier
func AddItemToCart(cartID, productID string, variantID *string, quantity int, price float64) error {
	// Vérifier si l'article existe déjà
	var existingID string
	var existingQty int
	err := db.QueryRow(
		"SELECT id, quantity FROM cart_items WHERE cart_id = $1 AND product_id = $2 AND (variant_id = $3 OR (variant_id IS NULL AND $3 IS NULL))",
		cartID, productID, variantID,
	).Scan(&existingID, &existingQty)

	if err == sql.ErrNoRows {
		// Créer un nouvel article
		_, err = db.Exec(
			"INSERT INTO cart_items (cart_id, product_id, variant_id, quantity, price) VALUES ($1, $2, $3, $4, $5)",
			cartID, productID, variantID, quantity, price,
		)
	} else if err == nil {
		// Mettre à jour la quantité
		_, err = db.Exec(
			"UPDATE cart_items SET quantity = quantity + $1 WHERE id = $2",
			quantity, existingID,
		)
	}

	if err != nil {
		return err
	}

	// Mettre à jour le total du panier
	return updateCartTotal(cartID)
}

// RemoveItemFromCart supprime un article du panier
func RemoveItemFromCart(cartID, itemID string) error {
	_, err := db.Exec("DELETE FROM cart_items WHERE id = $1 AND cart_id = $2", itemID, cartID)
	if err != nil {
		return err
	}

	return updateCartTotal(cartID)
}

// updateCartTotal met à jour le total du panier
func updateCartTotal(cartID string) error {
	_, err := db.Exec(
		"UPDATE carts SET total_amount = (SELECT COALESCE(SUM(price * quantity), 0) FROM cart_items WHERE cart_id = $1), updated_at = CURRENT_TIMESTAMP WHERE id = $1",
		cartID,
	)
	return err
}

// CreateOrder crée une nouvelle commande
func CreateOrder(userID, merchantID string, totalAmount float64, currency string, shippingAddr, billingAddr Address, paymentIntentID *string) (*Order, error) {
	shippingJSON, _ := json.Marshal(shippingAddr)
	billingJSON, _ := json.Marshal(billingAddr)

	var order Order
	err := db.QueryRow(
		"INSERT INTO orders (user_id, merchant_id, status, total_amount, currency, payment_intent_id, shipping_address, billing_address) VALUES ($1, $2, 'pending', $3, $4, $5, $6, $7) RETURNING id, user_id, merchant_id, status, total_amount, currency, payment_intent_id, shipping_address, billing_address, created_at, updated_at",
		userID, merchantID, totalAmount, currency, paymentIntentID, shippingJSON, billingJSON,
	).Scan(&order.ID, &order.UserID, &order.MerchantID, &order.Status, &order.TotalAmount, &order.Currency, &order.PaymentIntentID, &order.ShippingAddress, &order.BillingAddress, &order.CreatedAt, &order.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &order, nil
}

// GetOrder récupère une commande par ID
func GetOrder(orderID string) (*Order, error) {
	var order Order
	err := db.QueryRow(
		"SELECT id, user_id, merchant_id, status, total_amount, currency, payment_intent_id, shipping_address, billing_address, created_at, updated_at FROM orders WHERE id = $1",
		orderID,
	).Scan(&order.ID, &order.UserID, &order.MerchantID, &order.Status, &order.TotalAmount, &order.Currency, &order.PaymentIntentID, &order.ShippingAddress, &order.BillingAddress, &order.CreatedAt, &order.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &order, nil
}

// UpdateOrderStatus met à jour le statut d'une commande
func UpdateOrderStatus(orderID, status string) error {
	_, err := db.Exec(
		"UPDATE orders SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2",
		status, orderID,
	)
	return err
}

// GetUserOrders récupère les commandes d'un utilisateur
func GetUserOrders(userID string, limit, offset int) ([]Order, error) {
	rows, err := db.Query(
		"SELECT id, user_id, merchant_id, status, total_amount, currency, payment_intent_id, shipping_address, billing_address, created_at, updated_at FROM orders WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3",
		userID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		err := rows.Scan(&order.ID, &order.UserID, &order.MerchantID, &order.Status, &order.TotalAmount, &order.Currency, &order.PaymentIntentID, &order.ShippingAddress, &order.BillingAddress, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

