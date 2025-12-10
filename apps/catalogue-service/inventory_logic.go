package main

import (
	"database/sql"
	"errors"
	"log"
	"time"
)

// Inventory représente le stock d'un produit
type Inventory struct {
	ProductID   string    `json:"product_id" db:"product_id"`
	VariantID   *string   `json:"variant_id,omitempty" db:"variant_id"`
	Quantity    int       `json:"quantity" db:"quantity"`
	Reserved    int       `json:"reserved" db:"reserved"` // Quantité réservée (dans les paniers)
	Available   int       `json:"available" db:"available"` // Quantité disponible = Quantity - Reserved
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// InventoryMovement représente un mouvement de stock
type InventoryMovement struct {
	ID          string    `json:"id" db:"id"`
	ProductID   string    `json:"product_id" db:"product_id"`
	VariantID   *string   `json:"variant_id,omitempty" db:"variant_id"`
	Type        string    `json:"type" db:"type"` // in, out, adjustment
	Quantity    int       `json:"quantity" db:"quantity"`
	Reason      string    `json:"reason" db:"reason"`
	ReferenceID *string   `json:"reference_id,omitempty" db:"reference_id"` // ID de commande, etc.
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// GetInventory récupère le stock d'un produit
func GetInventory(productID string, variantID *string) (*Inventory, error) {
	var inv Inventory
	var variantIDNull sql.NullString
	
	query := "SELECT product_id, variant_id, quantity, reserved, (quantity - reserved) as available, updated_at FROM inventory WHERE product_id = $1"
	args := []interface{}{productID}
	
	if variantID != nil {
		query += " AND variant_id = $2"
		args = append(args, *variantID)
	} else {
		query += " AND variant_id IS NULL"
	}
	
	err := db.QueryRow(query, args...).Scan(&inv.ProductID, &variantIDNull, &inv.Quantity, &inv.Reserved, &inv.Available, &inv.UpdatedAt)
	
	if err == sql.ErrNoRows {
		// Créer un inventaire vide si il n'existe pas
		inv.ProductID = productID
		inv.VariantID = variantID
		inv.Quantity = 0
		inv.Reserved = 0
		inv.Available = 0
		inv.UpdatedAt = time.Now()
		return &inv, nil
	}
	
	if err != nil {
		return nil, err
	}
	
	if variantIDNull.Valid {
		v := variantIDNull.String
		inv.VariantID = &v
	}
	
	return &inv, nil
}

// UpdateInventory met à jour le stock d'un produit
func UpdateInventory(productID string, variantID *string, quantity int) error {
	// Vérifier si l'inventaire existe
	existing, err := GetInventory(productID, variantID)
	if err != nil {
		return err
	}
	
	if existing == nil || existing.Quantity == 0 {
		// Créer un nouvel inventaire
		_, err = db.Exec(
			"INSERT INTO inventory (product_id, variant_id, quantity, reserved) VALUES ($1, $2, $3, 0) ON CONFLICT (product_id, variant_id) DO UPDATE SET quantity = $3, updated_at = CURRENT_TIMESTAMP",
			productID, variantID, quantity,
		)
	} else {
		// Mettre à jour l'inventaire existant
		_, err = db.Exec(
			"UPDATE inventory SET quantity = $1, updated_at = CURRENT_TIMESTAMP WHERE product_id = $2 AND (variant_id = $3 OR (variant_id IS NULL AND $3 IS NULL))",
			quantity, productID, variantID,
		)
	}
	
	return err
}

// ReserveInventory réserve une quantité de stock (pour un panier)
func ReserveInventory(productID string, variantID *string, quantity int) error {
	available, err := CheckAvailability(productID, variantID, quantity)
	if err != nil || !available {
		return errors.New("stock insuffisant")
	}
	
	_, err = db.Exec(
		"UPDATE inventory SET reserved = reserved + $1, updated_at = CURRENT_TIMESTAMP WHERE product_id = $2 AND (variant_id = $3 OR (variant_id IS NULL AND $3 IS NULL))",
		quantity, productID, variantID,
	)
	
	return err
}

// ReleaseInventory libère une quantité réservée
func ReleaseInventory(productID string, variantID *string, quantity int) error {
	_, err := db.Exec(
		"UPDATE inventory SET reserved = GREATEST(0, reserved - $1), updated_at = CURRENT_TIMESTAMP WHERE product_id = $2 AND (variant_id = $3 OR (variant_id IS NULL AND $3 IS NULL))",
		quantity, productID, variantID,
	)
	return err
}

// DeductInventory déduit une quantité de stock (après commande)
func DeductInventory(productID string, variantID *string, quantity int, orderID string) error {
	available, err := CheckAvailability(productID, variantID, quantity)
	if err != nil || !available {
		return errors.New("stock insuffisant")
	}
	
	// Déduire de la quantité et de la réserve
	_, err = db.Exec(
		"UPDATE inventory SET quantity = quantity - $1, reserved = GREATEST(0, reserved - $1), updated_at = CURRENT_TIMESTAMP WHERE product_id = $2 AND (variant_id = $3 OR (variant_id IS NULL AND $3 IS NULL))",
		quantity, productID, variantID,
	)
	
	return err
}

// CheckAvailability vérifie si une quantité est disponible
func CheckAvailability(productID string, variantID *string, quantity int) (bool, error) {
	inventory, err := GetInventory(productID, variantID)
	if err != nil {
		return false, err
	}
	
	if inventory == nil {
		return false, errors.New("produit non trouvé")
	}
	
	return inventory.Available >= quantity, nil
}

// SyncInventoryToElasticsearch synchronise le stock avec Elasticsearch
func SyncInventoryToElasticsearch(productID string) error {
	// TODO: Récupérer le produit et son stock
	// TODO: Mettre à jour le document Elasticsearch
	return nil
}

