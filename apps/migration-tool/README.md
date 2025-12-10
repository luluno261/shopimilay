# Migration Tool

Service d'import de données pour OmniSphere.

## Fonctionnalités

- Import depuis Shopify
- Import depuis WooCommerce
- Import depuis CSV/JSON générique

## Développement

```bash
go mod download
go run main.go
```

## Endpoints

- `GET /health` - Health check
- `GET /ready` - Readiness check
- `POST /api/v1/migration/import` - Démarrer un import
- `GET /api/v1/migration/status/:id` - Statut d'un import

