# OmniSphere - Plateforme E-commerce Multi-tenant

Structure Monorepo légère pour OmniSphere (Microservices & IaC)

## Architecture

- **Architecture** : Microservices, Événementiel (Kafka)
- **Principaux langages** : Go, Node.js/TypeScript, Next.js
- **Status** : Ready for Development

## Structure du projet

```
.
├── apps/                    # Microservices Backend (API)
│   ├── auth-service/        # Go - Authentification (JWT, Rôles)
│   ├── checkout-service/    # Go - Panier, paiement, Stripe Connect
│   ├── catalogue-service/   # Go - Produits, stocks, Elasticsearch
│   └── marketing-engine/    # Node.js/NestJS - CDP, segmentation, emails
├── web/                     # Applications Frontend (Next.js)
│   ├── admin-dashboard/     # Interface d'administration marchands
│   └── storefront-template/ # Template boutique cliente (SSR/SSG)
├── infra/                   # Infrastructure as Code
│   ├── aws_terraform/       # Terraform - ECS, RDS, MSK, ALB
│   └── devops_scripts/      # Scripts CI/CD, migrations, déploiement
├── shared/                  # Ressources partagées
│   ├── schemas/             # ProtoBuf/JSON Schema pour Kafka
│   ├── libraries/           # Code Go et Node.js commun
│   └── database_migrations/ # Migrations SQL/Go
└── docs/                    # Documentation technique et métier
```

## Services

### Backend Services

- **auth-service** : Gestion de l'authentification (JWT, Rôles), utilisateurs de la plateforme
- **checkout-service** : Gestion critique du panier, tunnel de paiement, intégration Stripe Connect
- **catalogue-service** : Gestion des produits, stocks, synchronisation Elasticsearch
- **marketing-engine** : Ingestion événements (CDP), segmentation, séquences d'emails

### Frontend Applications

- **admin-dashboard** : Interface d'administration pour les marchands (CRM, Commandes, Produits)
- **storefront-template** : Template de base pour la boutique cliente (SSR/SSG)

## Développement local

### Prérequis

- Docker & Docker Compose
- Go 1.21+
- Node.js 20+ (LTS)
- Terraform (pour l'infrastructure)

### Démarrage rapide

```bash
# 1. Démarrer les services d'infrastructure (PostgreSQL, Kafka, Redis, Elasticsearch)
make infra-up

# 2. Installer les dépendances
cd apps/marketing-engine && npm install
cd ../../web/admin-dashboard && npm install
cd ../storefront-template && npm install

# 3. Démarrer les services backend (dans des terminaux séparés)
# Terminal 1: API Gateway (port 8080)
cd apps/api-gateway && go run main.go middleware.go routes.go

# Terminal 2: Auth Service (port 8081)
cd apps/auth-service && go run main.go database.go jwt.go

# Terminal 3: Catalogue Service (port 8082)
cd apps/catalogue-service && go run main.go database.go product_model.go inventory_logic.go elasticsearch_client.go ai_service.go store_builder.go

# Terminal 4: Checkout Service (port 8083)
cd apps/checkout-service && go run main.go database.go checkout_api.go orders.go discounts.go stripe_handlers.go

# Terminal 5: Marketing Engine (port 8084)
cd apps/marketing-engine && npm run start:dev

# 4. Démarrer les applications frontend
# Terminal 6: Admin Dashboard (http://localhost:3000)
cd web/admin-dashboard && npm run dev

# Terminal 7: Storefront (http://localhost:3001)
cd web/storefront-template && npm run dev

# Voir les logs de l'infrastructure
make logs
```

**📖 Guide détaillé** : Voir [docs/LOCAL_SETUP.md](docs/LOCAL_SETUP.md) pour un guide complet avec dépannage.

## Infrastructure

L'infrastructure AWS est gérée via Terraform dans `infra/aws_terraform/`.

Voir la documentation dans `docs/` pour plus de détails.

