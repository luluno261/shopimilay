# Guide de démarrage local - OmniSphere

Ce guide vous explique comment tester la plateforme OmniSphere localement sur votre machine.

## Prérequis

Avant de commencer, assurez-vous d'avoir installé :

- **Docker & Docker Compose** (pour l'infrastructure)
- **Go 1.21+** (pour les services backend Go)
- **Node.js 20+ (LTS)** (pour les applications frontend et marketing-engine)
- **npm** ou **yarn** (gestionnaire de paquets Node.js)

## Étape 1 : Démarrer l'infrastructure

L'infrastructure comprend PostgreSQL, Kafka, Redis et Elasticsearch.

```bash
# Depuis la racine du projet
make infra-up

# Ou manuellement avec Docker Compose
docker-compose up -d
```

Vérifiez que les services sont démarrés :

```bash
docker-compose ps
```

Vous devriez voir :
- `omnisphere-postgres` (port 5432)
- `omnisphere-kafka` (port 9092)
- `omnisphere-redis` (port 6379)
- `omnisphere-elasticsearch` (port 9200)
- `omnisphere-zookeeper` (port 2181)

## Étape 2 : Configurer les variables d'environnement

### API Gateway

Créez un fichier `.env` dans `apps/api-gateway/` (optionnel, valeurs par défaut disponibles) :

```env
PORT=8080
JWT_SECRET=your-secret-key-change-in-production
```

### Services Backend Go

Les services Go utilisent des variables d'environnement. Vous pouvez les définir ou utiliser les valeurs par défaut :

**auth-service, checkout-service, catalogue-service :**
```env
DATABASE_URL=postgres://omnisphere:omnisphere_dev@localhost:5432/omnisphere?sslmode=disable
PORT=8081  # auth-service
PORT=8083  # checkout-service
PORT=8082  # catalogue-service
```

### Marketing Engine (NestJS)

Créez un fichier `.env` dans `apps/marketing-engine/` :

```env
DATABASE_URL=postgres://omnisphere:omnisphere_dev@localhost:5432/omnisphere?sslmode=disable
KAFKA_BROKER=localhost:9092
PORT=8084
```

### Frontend Applications

Les applications Next.js utilisent `NEXT_PUBLIC_API_URL` pour pointer vers l'API Gateway.

**admin-dashboard et storefront-template :**
Créez un fichier `.env.local` dans chaque dossier :

```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

## Étape 3 : Installer les dépendances

### Services Backend Go

Les services Go n'ont pas besoin d'installation de dépendances (gérées par `go.mod`).

### Marketing Engine

```bash
cd apps/marketing-engine
npm install
```

### Frontend Applications

```bash
# Admin Dashboard
cd web/admin-dashboard
npm install

# Storefront Template
cd web/storefront-template
npm install
```

## Étape 4 : Exécuter les migrations de base de données

Avant de démarrer les services, exécutez les migrations SQL :

```bash
# Se connecter à PostgreSQL
docker exec -it omnisphere-postgres psql -U omnisphere -d omnisphere

# Ou utiliser un client PostgreSQL pour exécuter les fichiers dans :
# shared/database_migrations/
```

## Étape 5 : Démarrer les services backend

Ouvrez plusieurs terminaux pour démarrer chaque service :

### Terminal 1 : API Gateway
```bash
cd apps/api-gateway
go run main.go middleware.go routes.go
# Le service sera accessible sur http://localhost:8080
```

### Terminal 2 : Auth Service
```bash
cd apps/auth-service
go run main.go database.go jwt.go
# Le service sera accessible sur http://localhost:8081
```

### Terminal 3 : Catalogue Service
```bash
cd apps/catalogue-service
go run main.go database.go product_model.go inventory_logic.go elasticsearch_client.go ai_service.go store_builder.go
# Le service sera accessible sur http://localhost:8082
```

### Terminal 4 : Checkout Service
```bash
cd apps/checkout-service
go run main.go database.go checkout_api.go orders.go discounts.go stripe_handlers.go
# Le service sera accessible sur http://localhost:8083
```

### Terminal 5 : Marketing Engine
```bash
cd apps/marketing-engine
npm run start:dev
# Le service sera accessible sur http://localhost:8084
```

## Étape 6 : Démarrer les applications frontend

### Terminal 6 : Admin Dashboard
```bash
cd web/admin-dashboard
npm run dev
# L'application sera accessible sur http://localhost:3000
```

### Terminal 7 : Storefront Template
```bash
cd web/storefront-template
npm run dev
# L'application sera accessible sur http://localhost:3001
# (ou le prochain port disponible si 3001 est occupé)
```

## Étape 7 : Accéder aux applications

Une fois tous les services démarrés :

- **Admin Dashboard** : http://localhost:3000
  - Dashboard principal
  - Growth Command Center : http://localhost:3000/growth-command-center
  - Store Builder : http://localhost:3000/store-builder
  - Pages : Products, Orders, Customers, Discounts, etc.

- **Storefront Template** : http://localhost:3001
  - Page d'accueil de la boutique
  - Produits : http://localhost:3001/products/[id]
  - Checkout : http://localhost:3001/checkout
  - Compte client : http://localhost:3001/account

- **API Gateway** : http://localhost:8080
  - Health check : http://localhost:8080/health
  - API : http://localhost:8080/api/v1

## Commandes utiles

### Voir les logs de l'infrastructure
```bash
make logs
# Ou
docker-compose logs -f
```

### Arrêter l'infrastructure
```bash
make infra-down
# Ou
docker-compose down
```

### Nettoyer les volumes Docker
```bash
make clean
# Ou
docker-compose down -v
```

### Vérifier les services en cours d'exécution
```bash
# Services Docker
docker-compose ps

# Services Go (dans chaque terminal)
# Les services affichent des logs dans la console

# Services Node.js
# Vérifiez les ports dans les logs
```

## Dépannage

### Problème : Port déjà utilisé

Si un port est déjà utilisé, modifiez le port dans :
- Les variables d'environnement des services
- Les fichiers de configuration Docker Compose
- Les fichiers `.env.local` des applications frontend

### Problème : Base de données non accessible

Vérifiez que PostgreSQL est démarré :
```bash
docker-compose ps
docker-compose logs postgres
```

### Problème : Erreurs de connexion API

Vérifiez que :
1. L'API Gateway est démarré (port 8080)
2. Les services backend sont démarrés
3. La variable `NEXT_PUBLIC_API_URL` est correctement configurée dans les applications frontend

### Problème : Erreurs TypeScript

Les erreurs TypeScript dans certains fichiers sont normales si les dépendances ne sont pas installées. Exécutez :
```bash
npm install
```

## Test rapide

Pour tester rapidement que tout fonctionne :

1. **Vérifier l'API Gateway** :
   ```bash
   curl http://localhost:8080/health
   ```

2. **Vérifier PostgreSQL** :
   ```bash
   docker exec -it omnisphere-postgres psql -U omnisphere -c "SELECT version();"
   ```

3. **Vérifier Redis** :
   ```bash
   docker exec -it omnisphere-redis redis-cli ping
   ```

4. **Vérifier Elasticsearch** :
   ```bash
   curl http://localhost:9200
   ```

5. **Accéder au frontend** :
   - Ouvrez http://localhost:3000 dans votre navigateur
   - Ouvrez http://localhost:3001 dans votre navigateur

## Notes importantes

- Les services Go doivent être démarrés dans l'ordre (API Gateway en dernier)
- Les applications frontend peuvent être démarrées en parallèle
- Pour le développement, utilisez `npm run dev` qui active le hot-reload
- Les migrations de base de données doivent être exécutées avant le premier démarrage

