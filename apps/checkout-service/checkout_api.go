package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Cart représente un panier d'achat
type Cart struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Items     []CartItem `json:"items"`
	Total     float64   `json:"total"`
	Currency  string    `json:"currency"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

// CartItem représente un article dans le panier
type CartItem struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Subtotal  float64 `json:"subtotal"`
}

// CheckoutRequest représente une demande de checkout
type CheckoutRequest struct {
	CartID         string  `json:"cart_id" binding:"required"`
	PaymentMethod  string  `json:"payment_method" binding:"required"`
	ShippingAddress Address `json:"shipping_address" binding:"required"`
	BillingAddress  Address `json:"billing_address" binding:"required"`
	DiscountCode   *string `json:"discount_code,omitempty"`
}

// Address représente une adresse
type Address struct {
	Street  string `json:"street" binding:"required"`
	City    string `json:"city" binding:"required"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code" binding:"required"`
	Country string `json:"country" binding:"required"`
}

// CheckoutResponse représente la réponse d'un checkout
type CheckoutResponse struct {
	OrderID     string `json:"order_id"`
	Status      string `json:"status"`
	PaymentIntentID string `json:"payment_intent_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
}

// Les fonctions GetCart, AddToCart, RemoveFromCart sont maintenant dans database.go
// ProcessCheckout est géré directement dans main.go

