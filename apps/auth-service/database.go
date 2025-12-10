package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	connStr := getEnv("DATABASE_URL", "postgres://omnisphere:omnisphere_dev@localhost:5432/omnisphere?sslmode=disable")

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Database connection error:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Database ping error:", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Database connection established")
}

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	FullName     string    `json:"full_name"`
	Phone        string    `json:"phone"`
	AvatarURL    string    `json:"avatar_url"`
	Role         string    `json:"role"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	LastLoginAt  *time.Time `json:"last_login_at"`
}

type MerchantAccount struct {
	ID                  string    `json:"id"`
	UserID              string    `json:"user_id"`
	StoreName           string    `json:"store_name"`
	StoreSlug           string    `json:"store_slug"`
	Description         string    `json:"description"`
	LogoURL             string    `json:"logo_url"`
	BannerURL           string    `json:"banner_url"`
	Country             string    `json:"country"`
	Currency            string    `json:"currency"`
	Timezone            string    `json:"timezone"`
	StripeAccountID     string    `json:"stripe_account_id"`
	StripeConnectedAt   *time.Time `json:"stripe_connected_at"`
	SubscriptionTier    string    `json:"subscription_tier"`
	SubscriptionStatus  string    `json:"subscription_status"`
	TrialEndsAt         *time.Time `json:"trial_ends_at"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func GetUserByEmail(email string) (*User, error) {
	var user User
	err := db.QueryRow(`
		SELECT id, email, password_hash, full_name, phone, avatar_url, role, status, created_at, updated_at, last_login_at
		FROM users WHERE email = $1
	`, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FullName, &user.Phone, &user.AvatarURL,
		&user.Role, &user.Status, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByID(id string) (*User, error) {
	var user User
	err := db.QueryRow(`
		SELECT id, email, full_name, phone, avatar_url, role, status, created_at, updated_at, last_login_at
		FROM users WHERE id = $1
	`, id).Scan(
		&user.ID, &user.Email, &user.FullName, &user.Phone, &user.AvatarURL,
		&user.Role, &user.Status, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func CreateUser(email, passwordHash, fullName string) (*User, error) {
	var userID string
	err := db.QueryRow(`
		INSERT INTO users (email, password_hash, full_name, role, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`, email, passwordHash, fullName, "merchant", "active").Scan(&userID)

	if err != nil {
		return nil, err
	}

	return GetUserByID(userID)
}

func DeleteUser(id string) error {
	_, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}

func UpdateLastLogin(id string) error {
	_, err := db.Exec("UPDATE users SET last_login_at = NOW() WHERE id = $1", id)
	return err
}

func GetMerchantAccount(userID string) (*MerchantAccount, error) {
	var account MerchantAccount
	err := db.QueryRow(`
		SELECT id, user_id, store_name, store_slug, description, logo_url, banner_url, country,
		       currency, timezone, stripe_account_id, stripe_connected_at, subscription_tier,
		       subscription_status, trial_ends_at, created_at, updated_at
		FROM merchant_accounts WHERE user_id = $1
	`, userID).Scan(
		&account.ID, &account.UserID, &account.StoreName, &account.StoreSlug, &account.Description,
		&account.LogoURL, &account.BannerURL, &account.Country, &account.Currency, &account.Timezone,
		&account.StripeAccountID, &account.StripeConnectedAt, &account.SubscriptionTier,
		&account.SubscriptionStatus, &account.TrialEndsAt, &account.CreatedAt, &account.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func CreateMerchantAccount(userID, storeName, storeSlug string) (*MerchantAccount, error) {
	var accountID string
	err := db.QueryRow(`
		INSERT INTO merchant_accounts (user_id, store_name, store_slug, currency, timezone, subscription_tier, subscription_status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, userID, storeName, storeSlug, "USD", "UTC", "basic", "active").Scan(&accountID)

	if err != nil {
		return nil, err
	}

	return &MerchantAccount{
		ID:        accountID,
		UserID:    userID,
		StoreName: storeName,
		StoreSlug: storeSlug,
	}, nil
}

func GetAPIKeyByPrefix(prefix string) (string, error) {
	var merchantID string
	err := db.QueryRow(`
		SELECT merchant_id FROM api_keys WHERE prefix = $1 AND status = 'active'
	`, prefix).Scan(&merchantID)

	if err == sql.ErrNoRows {
		return "", nil
	}

	return merchantID, err
}

func LogAuditEvent(userID, merchantID, action, resourceType, resourceID string, changes interface{}) error {
	_, err := db.Exec(`
		INSERT INTO audit_logs (user_id, merchant_id, action, resource_type, resource_id, status)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, userID, merchantID, action, resourceType, resourceID, "success")

	return err
}
