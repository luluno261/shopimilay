package main

import (
	"database/sql"
	"log"
	"os"

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

// User représente un utilisateur
type User struct {
	ID        string `json:"id" db:"id"`
	Email     string `json:"email" db:"email"`
	PasswordHash string `json:"-" db:"password_hash"`
	Role      string `json:"role" db:"role"`
	MerchantID *string `json:"merchant_id,omitempty" db:"merchant_id"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UpdatedAt string `json:"updated_at" db:"updated_at"`
}

// GetUserByEmail récupère un utilisateur par email
func GetUserByEmail(email string) (*User, error) {
	var user User
	err := db.QueryRow(
		"SELECT id, email, password_hash, role, merchant_id, created_at, updated_at FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.MerchantID, &user.CreatedAt, &user.UpdatedAt)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

// GetUserByID récupère un utilisateur par ID
func GetUserByID(userID string) (*User, error) {
	var user User
	err := db.QueryRow(
		"SELECT id, email, password_hash, role, merchant_id, created_at, updated_at FROM users WHERE id = $1",
		userID,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.MerchantID, &user.CreatedAt, &user.UpdatedAt)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

// CreateUser crée un nouvel utilisateur
func CreateUser(email, passwordHash, role string, merchantID *string) (*User, error) {
	var user User
	err := db.QueryRow(
		"INSERT INTO users (email, password_hash, role, merchant_id) VALUES ($1, $2, $3, $4) RETURNING id, email, password_hash, role, merchant_id, created_at, updated_at",
		email, passwordHash, role, merchantID,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Role, &user.MerchantID, &user.CreatedAt, &user.UpdatedAt)
	
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

