# Guide de Démarrage Local - OmniSphere

Ce guide vous explique comment tester la plateforme OmniSphere en local.

## Prérequis

Avant de commencer, assurez-vous d'avoir installé :

- **Docker & Docker Compose** (pour l'infrastructure)
- **Go 1.21+** (pour les services backend)
- **Node.js 20+ (LTS)** (pour les applications frontend)
- **npm** ou **yarn** (gestionnaire de paquets Node.js)

## Étape 1 : Démarrer l'infrastructure

L'infrastructure comprend PostgreSQL, Kafka, Redis et Elasticsearch.

```bash
# Depuis la racine du projet
make infra-up
```

Ou manuellement :

```bash
docker-compose up -d
```

Vérifiez que tous les services sont démarrés :

```bash
docker-compose ps
```

## Étape 2 : Exécuter les migrations de base de données

Avant de démarrer les services, vous devez exécuter les migrations SQL :

```bash
# Se connecter à PostgreSQL
docker exec -it omnisphere-postgres psql -U omnisphere -d omnisphere

# Ou utiliser un client PostgreSQL pour exécuter les migrations
# Les fichiers sont dans shared/database_migrations/
```

## Étape 3 : Démarrer les services backend

Ouvrez plusieurs terminaux pour chaque service :

### Terminal 1 : API Gateway (Port 8080)
```bash
cd apps/api-gateway
go mod download  # Si nécessaire
go run main.go
```

### Terminal 2 : Auth Service (Port 8081)
```bash
cd apps/auth-service
go mod download  # Si nécessaire
go run main.go
```

### Terminal 3 : Catalogue Service (Port 8082)
```bash
cd apps/catalogue-service
go mod download  # Si nécessaire
go run main.go
```

### Terminal 4 : Checkout Service (Port 8083)
```bash
cd apps/checkout-service
go mod download  # Si nécessaire
go run main.go
```

### Terminal 5 : Marketing Engine (Port 8084)
```bash
cd apps/marketing-engine
npm install  # Si nécessaire
npm run start:dev
```

### Terminal 6 : Webhook Service (Port 8086)
```bash
cd apps/webhook-service
go mod download  # Si nécessaire
go run main.go
```

### Terminal 7 : Migration Tool (Port 8085)
```bash
cd apps/migration-tool
go mod download  # Si nécessaire
go run main.go
```

## Étape 4 : Démarrer les applications frontend

### Terminal 8 : Admin Dashboard (Port 3000)
```bash
cd web/admin-dashboard
npm install  # Si nécessaire
npm run dev
```

### Terminal 9 : Storefront Template (Port 3001)
```bash
cd web/storefront-template
npm install  # Si nécessaire
npm run dev -- -p 3001
```

## Étape 5 : Accéder aux applications

Une fois tous les services démarrés, vous pouvez accéder à :

- **Admin Dashboard** : http://localhost:3000
- **Storefront Template** : http://localhost:3001
- **API Gateway** : http://localhost:8080
- **API Health Check** : http://localhost:8080/health

## Configuration des variables d'environnement

### Admin Dashboard

Créez un fichier `web/admin-dashboard/.env.local` :

```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

### Storefront Template

Créez un fichier `web/storefront-template/.env.local` :

```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

### Services Backend

Les services Go utilisent des variables d'environnement. Créez des fichiers `.env` dans chaque service ou définissez-les :

```bash
# Exemple pour auth-service
export DATABASE_URL=postgres://omnisphere:omnisphere_dev@localhost:5432/omnisphere?sslmode=disable
export JWT_SECRET=your-secret-key
export PORT=8081
```

## Vérification des services

### Vérifier que l'API Gateway fonctionne

```bash
curl http://localhost:8080/health
```

### Vérifier que PostgreSQL est accessible

```bash
docker exec -it omnisphere-postgres psql -U omnisphere -d omnisphere -c "SELECT version();"
```

### Vérifier que Redis fonctionne

```bash
docker exec -it omnisphere-redis redis-cli ping
```

### Vérifier que Kafka fonctionne

```bash
docker exec -it omnisphere-kafka kafka-broker-api-versions --bootstrap-server localhost:9092
```

### Vérifier qu'Elasticsearch fonctionne

```bash
curl http://localhost:9200/_cluster/health
```

## Script de démarrage automatique (Optionnel)

Vous pouvez créer un script pour démarrer tous les services automatiquement :

### Windows (start-all.bat)

```batch
@echo off
echo Démarrage de l'infrastructure...
docker-compose up -d

echo Attente du démarrage des services...
timeout /t 10

echo Démarrage des services backend...
start "API Gateway" cmd /k "cd apps/api-gateway && go run main.go"
start "Auth Service" cmd /k "cd apps/auth-service && go run main.go"
start "Catalogue Service" cmd /k "cd apps/catalogue-service && go run main.go"
start "Checkout Service" cmd /k "cd apps/checkout-service && go run main.go"
start "Marketing Engine" cmd /k "cd apps/marketing-engine && npm run start:dev"
start "Webhook Service" cmd /k "cd apps/webhook-service && go run main.go"
start "Migration Tool" cmd /k "cd apps/migration-tool && go run main.go"

echo Démarrage des applications frontend...
start "Admin Dashboard" cmd /k "cd web/admin-dashboard && npm run dev"
start "Storefront" cmd /k "cd web/storefront-template && npm run dev -- -p 3001"

echo Tous les services sont en cours de démarrage...
pause
```

### Linux/Mac (start-all.sh)

```bash
#!/bin/bash

echo "Démarrage de l'infrastructure..."
docker-compose up -d

echo "Attente du démarrage des services..."
sleep 10

echo "Démarrage des services backend..."
cd apps/api-gateway && go run main.go &
cd apps/auth-service && go run main.go &
cd apps/catalogue-service && go run main.go &
cd apps/checkout-service && go run main.go &
cd apps/marketing-engine && npm run start:dev &
cd apps/webhook-service && go run main.go &
cd apps/migration-tool && go run main.go &

echo "Démarrage des applications frontend..."
cd web/admin-dashboard && npm run dev &
cd web/storefront-template && npm run dev -- -p 3001 &

echo "Tous les services sont en cours de démarrage..."
wait
```

## Dépannage

### Erreur de connexion à la base de données

Vérifiez que PostgreSQL est bien démarré :
```bash
docker-compose ps postgres
```

### Erreur "port already in use"

Changez le port dans les variables d'environnement ou arrêtez le service qui utilise le port.

### Erreur de dépendances Go

```bash
cd apps/[service-name]
go mod download
go mod tidy
```

### Erreur de dépendances Node.js

```bash
cd web/[app-name]
npm install
```

## Arrêter tous les services

```bash
# Arrêter l'infrastructure
make infra-down

# Ou
docker-compose down

# Arrêter les services manuellement (Ctrl+C dans chaque terminal)
```

## Nettoyer les données

```bash
# Supprimer les volumes Docker (⚠️ supprime toutes les données)
make clean
```

## Prochaines étapes

Une fois tous les services démarrés :

1. Accédez à l'Admin Dashboard : http://localhost:3000
2. Créez un compte marchand
3. Testez le Store Builder : http://localhost:3000/store-builder
4. Testez le Growth Command Center : http://localhost:3000/growth-command-center
5. Visitez le Storefront : http://localhost:3001

