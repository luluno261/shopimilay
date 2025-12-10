package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/stripe/stripe-go/v76/webhook"
)

// StripeConfig contient la configuration Stripe
type StripeConfig struct {
	SecretKey      string
	PublishableKey string
	WebhookSecret  string
}

// InitStripe initialise la configuration Stripe
func InitStripe() *StripeConfig {
	config := &StripeConfig{
		SecretKey:      getEnv("STRIPE_SECRET_KEY", ""),
		PublishableKey: getEnv("STRIPE_PUBLISHABLE_KEY", ""),
		WebhookSecret:  getEnv("STRIPE_WEBHOOK_SECRET", ""),
	}
	
	// Initialiser le client Stripe
	if config.SecretKey != "" {
		stripe.Key = config.SecretKey
	}
	
	return config
}

// CreatePaymentIntent crée un PaymentIntent Stripe
func CreatePaymentIntent(amount int64, currency string, merchantAccountID string) (string, string, error) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount),
		Currency:  stripe.String(currency),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	// Si un compte Stripe Connect est configuré, utiliser on_behalf_of
	if merchantAccountID != "" {
		params.OnBehalfOf = stripe.String(merchantAccountID)
		params.TransferData = &stripe.PaymentIntentTransferDataParams{
			Destination: stripe.String(merchantAccountID),
		}
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return "", "", err
	}

	return pi.ID, pi.ClientSecret, nil
}

// HandleStripeWebhook traite les webhooks Stripe
func HandleStripeWebhook(c *gin.Context) {
	if stripeConfig.WebhookSecret == "" {
		log.Println("Avertissement: STRIPE_WEBHOOK_SECRET non configuré")
		c.JSON(http.StatusOK, gin.H{"received": true})
		return
	}

	// Lire le body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Erreur lors de la lecture du body"})
		return
	}

	// Vérifier la signature
	event, err := webhook.ConstructEvent(body, c.GetHeader("Stripe-Signature"), stripeConfig.WebhookSecret)
	if err != nil {
		log.Printf("Erreur de validation de signature: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Signature invalide"})
		return
	}

	// Traiter l'événement
	switch event.Type {
	case "payment_intent.succeeded":
		var pi stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
			log.Printf("Erreur lors du parsing: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur de parsing"})
			return
		}

		// Mettre à jour le statut de la commande
		_, err = db.Exec(
			"UPDATE orders SET status = 'paid', updated_at = CURRENT_TIMESTAMP WHERE payment_intent_id = $1",
			pi.ID,
		)
		if err != nil {
			log.Printf("Erreur lors de la mise à jour de la commande: %v", err)
		}

		// TODO: Publier un événement Kafka order.paid

	case "payment_intent.payment_failed":
		var pi stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
			log.Printf("Erreur lors du parsing: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur de parsing"})
			return
		}

		// Mettre à jour le statut de la commande
		_, err = db.Exec(
			"UPDATE orders SET status = 'failed', updated_at = CURRENT_TIMESTAMP WHERE payment_intent_id = $1",
			pi.ID,
		)
		if err != nil {
			log.Printf("Erreur lors de la mise à jour de la commande: %v", err)
		}

	default:
		log.Printf("Événement non géré: %s", event.Type)
	}

	c.JSON(http.StatusOK, gin.H{"received": true})
}

// ProcessStripeConnectPayment traite un paiement via Stripe Connect
// Cette fonction est maintenant intégrée dans CreatePaymentIntent
func ProcessStripeConnectPayment(amount int64, currency string, merchantAccountID string, customerID string) error {
	// Le paiement Stripe Connect est géré automatiquement via CreatePaymentIntent
	// avec les paramètres OnBehalfOf et TransferData
	log.Printf("Paiement Stripe Connect configuré: %d %s pour merchant %s", amount, currency, merchantAccountID)
	return nil
}

// GetStripeAccount récupère les informations d'un compte Stripe Connect
func GetStripeAccount(merchantID string) (map[string]interface{}, error) {
	// TODO: Implémenter la récupération du compte Stripe Connect
	// Utiliser stripe.Account.Get() pour récupérer les informations
	return nil, nil
}

