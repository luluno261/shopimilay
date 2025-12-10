# Checkout Service

Service de gestion du panier et du tunnel de paiement pour OmniSphere.

## Fonctionnalités

- Gestion du panier d'achat
- Tunnel de paiement
- Intégration Stripe Connect (multi-tenant)
- Webhooks Stripe

## Développement

```bash
go mod download
go run main.go
```

## Endpoints

- `GET /health` - Health check
- `GET /ready` - Readiness check
- `GET /api/v1/cart` - Récupérer le panier
- `POST /api/v1/cart/items` - Ajouter au panier
- `DELETE /api/v1/cart/items/:itemId` - Supprimer du panier
- `POST /api/v1/checkout` - Processus de checkout
- `POST /api/v1/checkout/stripe/webhook` - Webhook Stripe

## Configuration Stripe

Variables d'environnement requises:
- `STRIPE_SECRET_KEY`
- `STRIPE_PUBLISHABLE_KEY`
- `STRIPE_WEBHOOK_SECRET`

