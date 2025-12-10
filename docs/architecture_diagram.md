# Architecture OmniSphere

## Vue d'ensemble

OmniSphere est une plateforme e-commerce multi-tenant construite avec une architecture microservices.

## Composants principaux

### Frontend

- **Admin Dashboard** (Next.js) - Interface d'administration pour les marchands
- **Storefront Template** (Next.js) - Template de boutique cliente

### Backend Services

- **auth-service** (Go) - Authentification et gestion des utilisateurs
- **checkout-service** (Go) - Panier et paiement (Stripe Connect)
- **catalogue-service** (Go) - Produits et stocks (Elasticsearch)
- **marketing-engine** (NestJS) - CDP, segmentation, automations

### Infrastructure

- **PostgreSQL** (RDS) - Base de données principale
- **Kafka** (MSK) - Bus d'événements
- **Elasticsearch** - Recherche de produits
- **Redis** - Cache et sessions
- **ECS** - Orchestration des conteneurs
- **ALB** - Load balancer

## Flux de données

### Authentification
1. Client → Admin Dashboard / Storefront
2. Admin Dashboard / Storefront → auth-service
3. auth-service → PostgreSQL
4. auth-service → JWT token → Client

### Checkout
1. Client → Storefront
2. Storefront → checkout-service
3. checkout-service → catalogue-service (vérification stock)
4. checkout-service → Stripe Connect (paiement)
5. checkout-service → Kafka (événement order.created)

### Catalogue
1. Client → Storefront
2. Storefront → catalogue-service
3. catalogue-service → PostgreSQL (données)
4. catalogue-service → Elasticsearch (recherche)

### Marketing
1. Services → Kafka (événements)
2. marketing-engine → Kafka (consommation)
3. marketing-engine → Segmentation et automations
4. marketing-engine → Envoi d'emails

## Communication inter-services

- **HTTP/REST** : Communication synchrone entre services
- **Kafka** : Communication asynchrone via événements
- **gRPC** (optionnel) : Pour les appels haute performance

## Sécurité

- JWT pour l'authentification
- HTTPS/TLS pour toutes les communications
- Secrets managés via AWS Secrets Manager
- VPC avec sous-réseaux privés pour les services backend

## Scalabilité

- Services déployés sur ECS Fargate (auto-scaling)
- Base de données RDS avec Multi-AZ en production
- Kafka MSK avec plusieurs brokers
- Cache Redis pour améliorer les performances

## Diagramme de déploiement

```
                    [Internet]
                         |
                    [ALB]
                         |
        +----------------+----------------+
        |                |                |
   [Admin Dashboard] [Storefront]   [API Gateway]
        |                |                |
        +----------------+----------------+
                         |
        +----------------+----------------+
        |                |                |
   [auth-service]  [checkout-service] [catalogue-service] [marketing-engine]
        |                |                |                    |
        +----------------+----------------+--------------------+
                         |                |                    |
                    [PostgreSQL]      [Elasticsearch]      [Kafka]
                                                          |
                                                    [Consumers]
```

