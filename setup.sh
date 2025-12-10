#!/bin/bash

# Script d'installation et de configuration OmniSphere
# Ce script installe toutes les dépendances et configure l'environnement

set -e

echo "🚀 Configuration d'OmniSphere..."
echo ""

# Couleurs pour les messages
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Fonction pour vérifier si une commande existe
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Vérification des prérequis
echo "📋 Vérification des prérequis..."

if ! command_exists docker; then
    echo -e "${RED}❌ Docker n'est pas installé. Veuillez l'installer depuis https://www.docker.com/${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Docker installé${NC}"

if ! command_exists docker-compose; then
    echo -e "${RED}❌ Docker Compose n'est pas installé. Veuillez l'installer.${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Docker Compose installé${NC}"

if ! command_exists node; then
    echo -e "${RED}❌ Node.js n'est pas installé. Veuillez l'installer depuis https://nodejs.org/${NC}"
    exit 1
fi
NODE_VERSION=$(node -v)
echo -e "${GREEN}✅ Node.js installé: $NODE_VERSION${NC}"

if ! command_exists npm; then
    echo -e "${RED}❌ npm n'est pas installé.${NC}"
    exit 1
fi
NPM_VERSION=$(npm -v)
echo -e "${GREEN}✅ npm installé: $NPM_VERSION${NC}"

if command_exists go; then
    GO_VERSION=$(go version | awk '{print $3}')
    echo -e "${GREEN}✅ Go installé: $GO_VERSION${NC}"
else
    echo -e "${YELLOW}⚠️  Go n'est pas installé. Les services Go ne pourront pas être démarrés.${NC}"
    echo -e "${YELLOW}   Installez Go depuis https://go.dev/dl/${NC}"
fi

echo ""
echo "📦 Installation des dépendances Node.js..."

# Marketing Engine
echo "  → Installation des dépendances pour marketing-engine..."
cd apps/marketing-engine
if [ -f "package.json" ]; then
    npm install
    echo -e "${GREEN}  ✅ marketing-engine${NC}"
else
    echo -e "${YELLOW}  ⚠️  package.json non trouvé dans marketing-engine${NC}"
fi
cd ../..

# Admin Dashboard
echo "  → Installation des dépendances pour admin-dashboard..."
cd web/admin-dashboard
if [ -f "package.json" ]; then
    npm install
    echo -e "${GREEN}  ✅ admin-dashboard${NC}"
else
    echo -e "${YELLOW}  ⚠️  package.json non trouvé dans admin-dashboard${NC}"
fi
cd ../..

# Storefront Template
echo "  → Installation des dépendances pour storefront-template..."
cd web/storefront-template
if [ -f "package.json" ]; then
    npm install
    echo -e "${GREEN}  ✅ storefront-template${NC}"
else
    echo -e "${YELLOW}  ⚠️  package.json non trouvé dans storefront-template${NC}"
fi
cd ../..

echo ""
if command_exists go; then
    echo "📦 Téléchargement des dépendances Go..."
    
    # API Gateway
    if [ -f "apps/api-gateway/go.mod" ]; then
        echo "  → Téléchargement des dépendances pour api-gateway..."
        cd apps/api-gateway
        go mod download
        echo -e "${GREEN}  ✅ api-gateway${NC}"
        cd ../..
    fi
    
    # Auth Service
    if [ -f "apps/auth-service/go.mod" ]; then
        echo "  → Téléchargement des dépendances pour auth-service..."
        cd apps/auth-service
        go mod download
        echo -e "${GREEN}  ✅ auth-service${NC}"
        cd ../..
    fi
    
    # Catalogue Service
    if [ -f "apps/catalogue-service/go.mod" ]; then
        echo "  → Téléchargement des dépendances pour catalogue-service..."
        cd apps/catalogue-service
        go mod download
        echo -e "${GREEN}  ✅ catalogue-service${NC}"
        cd ../..
    fi
    
    # Checkout Service
    if [ -f "apps/checkout-service/go.mod" ]; then
        echo "  → Téléchargement des dépendances pour checkout-service..."
        cd apps/checkout-service
        go mod download
        echo -e "${GREEN}  ✅ checkout-service${NC}"
        cd ../..
    fi
    
    # Migration Tool
    if [ -f "apps/migration-tool/go.mod" ]; then
        echo "  → Téléchargement des dépendances pour migration-tool..."
        cd apps/migration-tool
        go mod download
        echo -e "${GREEN}  ✅ migration-tool${NC}"
        cd ../..
    fi
    
    # Webhook Service
    if [ -f "apps/webhook-service/go.mod" ]; then
        echo "  → Téléchargement des dépendances pour webhook-service..."
        cd apps/webhook-service
        go mod download
        echo -e "${GREEN}  ✅ webhook-service${NC}"
        cd ../..
    fi
else
    echo -e "${YELLOW}⚠️  Go n'est pas installé. Les dépendances Go ne seront pas téléchargées.${NC}"
fi

echo ""
echo "⚙️  Configuration des fichiers d'environnement..."

# Créer les fichiers .env si ils n'existent pas
create_env_file() {
    if [ ! -f "$1" ]; then
        cp "$2" "$1"
        echo -e "${GREEN}  ✅ Créé $1${NC}"
    else
        echo -e "${YELLOW}  ⚠️  $1 existe déjà (non modifié)${NC}"
    fi
}

# Créer les fichiers .env.example d'abord
echo "  → Création des fichiers .env.example..."

# API Gateway
if [ ! -f "apps/api-gateway/.env.example" ]; then
    cat > apps/api-gateway/.env.example << EOF
PORT=8080
JWT_SECRET=your-secret-key-change-in-production
AUTH_SERVICE_URL=http://localhost:8081
CATALOGUE_SERVICE_URL=http://localhost:8082
CHECKOUT_SERVICE_URL=http://localhost:8083
MARKETING_ENGINE_URL=http://localhost:8084
WEBHOOK_SERVICE_URL=http://localhost:8085
MIGRATION_TOOL_URL=http://localhost:8086
EOF
fi

# Marketing Engine
if [ ! -f "apps/marketing-engine/.env.example" ]; then
    cat > apps/marketing-engine/.env.example << EOF
PORT=8084
DATABASE_URL=postgres://omnisphere:omnisphere_dev@localhost:5432/omnisphere?sslmode=disable
KAFKA_BROKER=localhost:9092
NODE_ENV=development
EOF
fi

# Admin Dashboard
if [ ! -f "web/admin-dashboard/.env.example" ]; then
    cat > web/admin-dashboard/.env.example << EOF
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
EOF
fi

# Storefront Template
if [ ! -f "web/storefront-template/.env.example" ]; then
    cat > web/storefront-template/.env.example << EOF
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
EOF
fi

# Copier les .env.example vers .env.local pour Next.js
if [ ! -f "web/admin-dashboard/.env.local" ]; then
    cp web/admin-dashboard/.env.example web/admin-dashboard/.env.local
    echo -e "${GREEN}  ✅ Créé web/admin-dashboard/.env.local${NC}"
fi

if [ ! -f "web/storefront-template/.env.local" ]; then
    cp web/storefront-template/.env.example web/storefront-template/.env.local
    echo -e "${GREEN}  ✅ Créé web/storefront-template/.env.local${NC}"
fi

echo ""
echo -e "${GREEN}✅ Configuration terminée !${NC}"
echo ""
echo "📝 Prochaines étapes :"
echo "  1. Démarrer l'infrastructure : make infra-up"
echo "  2. Démarrer les services backend (voir docs/LOCAL_SETUP.md)"
echo "  3. Démarrer les applications frontend :"
echo "     - Admin Dashboard: cd web/admin-dashboard && npm run dev"
echo "     - Storefront: cd web/storefront-template && npm run dev"
echo ""
echo "📖 Pour plus de détails, consultez docs/LOCAL_SETUP.md"

