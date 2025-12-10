# Webhook Service

Service de gestion des webhooks pour OmniSphere.

## Fonctionnalités

- Création, modification et suppression de webhooks
- Configuration d'événements (order.created, order.paid, etc.)
- Envoi de webhooks avec signature HMAC
- Test de webhooks

## Développement

```bash
go mod download
go run main.go
```

## Endpoints

- `GET /health` - Health check
- `GET /ready` - Readiness check
- `GET /api/v1/webhooks` - Liste des webhooks
- `POST /api/v1/webhooks` - Créer un webhook
- `PUT /api/v1/webhooks/:id` - Mettre à jour un webhook
- `DELETE /api/v1/webhooks/:id` - Supprimer un webhook
- `POST /api/v1/webhooks/:id/test` - Tester un webhook

