# Script d'installation et de configuration OmniSphere (PowerShell)
# Ce script installe toutes les dépendances et configure l'environnement

Write-Host "🚀 Configuration d'OmniSphere..." -ForegroundColor Cyan
Write-Host ""

# Fonction pour vérifier si une commande existe
function Test-Command {
    param($Command)
    $null -ne (Get-Command $Command -ErrorAction SilentlyContinue)
}

# Vérification des prérequis
Write-Host "📋 Vérification des prérequis..." -ForegroundColor Yellow

if (-not (Test-Command "docker")) {
    Write-Host "❌ Docker n'est pas installé. Veuillez l'installer depuis https://www.docker.com/" -ForegroundColor Red
    exit 1
}
Write-Host "✅ Docker installé" -ForegroundColor Green

if (-not (Test-Command "docker-compose")) {
    Write-Host "❌ Docker Compose n'est pas installé. Veuillez l'installer." -ForegroundColor Red
    exit 1
}
Write-Host "✅ Docker Compose installé" -ForegroundColor Green

if (-not (Test-Command "node")) {
    Write-Host "❌ Node.js n'est pas installé. Veuillez l'installer depuis https://nodejs.org/" -ForegroundColor Red
    exit 1
}
$nodeVersion = node -v
Write-Host "✅ Node.js installé: $nodeVersion" -ForegroundColor Green

if (-not (Test-Command "npm")) {
    Write-Host "❌ npm n'est pas installé." -ForegroundColor Red
    exit 1
}
$npmVersion = npm -v
Write-Host "✅ npm installé: $npmVersion" -ForegroundColor Green

$goInstalled = Test-Command "go"
if ($goInstalled) {
    $goVersion = go version
    Write-Host "✅ Go installé: $goVersion" -ForegroundColor Green
} else {
    Write-Host "⚠️  Go n'est pas installé. Les services Go ne pourront pas être démarrés." -ForegroundColor Yellow
    Write-Host "   Installez Go depuis https://go.dev/dl/" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "📦 Installation des dépendances Node.js..." -ForegroundColor Yellow

# Marketing Engine
Write-Host "  → Installation des dépendances pour marketing-engine..." -ForegroundColor Cyan
if (Test-Path "apps\marketing-engine\package.json") {
    Push-Location "apps\marketing-engine"
    npm install
    Write-Host "  ✅ marketing-engine" -ForegroundColor Green
    Pop-Location
} else {
    Write-Host "  ⚠️  package.json non trouvé dans marketing-engine" -ForegroundColor Yellow
}

# Admin Dashboard
Write-Host "  → Installation des dépendances pour admin-dashboard..." -ForegroundColor Cyan
if (Test-Path "web\admin-dashboard\package.json") {
    Push-Location "web\admin-dashboard"
    npm install
    Write-Host "  ✅ admin-dashboard" -ForegroundColor Green
    Pop-Location
} else {
    Write-Host "  ⚠️  package.json non trouvé dans admin-dashboard" -ForegroundColor Yellow
}

# Storefront Template
Write-Host "  → Installation des dépendances pour storefront-template..." -ForegroundColor Cyan
if (Test-Path "web\storefront-template\package.json") {
    Push-Location "web\storefront-template"
    npm install
    Write-Host "  ✅ storefront-template" -ForegroundColor Green
    Pop-Location
} else {
    Write-Host "  ⚠️  package.json non trouvé dans storefront-template" -ForegroundColor Yellow
}

Write-Host ""
if ($goInstalled) {
    Write-Host "📦 Téléchargement des dépendances Go..." -ForegroundColor Yellow
    
    # API Gateway
    if (Test-Path "apps\api-gateway\go.mod") {
        Write-Host "  → Téléchargement des dépendances pour api-gateway..." -ForegroundColor Cyan
        Push-Location "apps\api-gateway"
        go mod download
        go mod tidy
        Write-Host "  ✅ api-gateway" -ForegroundColor Green
        Pop-Location
    }
    
    # Auth Service
    if (Test-Path "apps\auth-service\go.mod") {
        Write-Host "  → Téléchargement des dépendances pour auth-service..." -ForegroundColor Cyan
        Push-Location "apps\auth-service"
        go mod download
        go mod tidy
        Write-Host "  ✅ auth-service" -ForegroundColor Green
        Pop-Location
    }
    
    # Catalogue Service
    if (Test-Path "apps\catalogue-service\go.mod") {
        Write-Host "  → Téléchargement des dépendances pour catalogue-service..." -ForegroundColor Cyan
        Push-Location "apps\catalogue-service"
        go mod download
        go mod tidy
        Write-Host "  ✅ catalogue-service" -ForegroundColor Green
        Pop-Location
    }
    
    # Checkout Service
    if (Test-Path "apps\checkout-service\go.mod") {
        Write-Host "  → Téléchargement des dépendances pour checkout-service..." -ForegroundColor Cyan
        Push-Location "apps\checkout-service"
        go mod download
        go mod tidy
        Write-Host "  ✅ checkout-service" -ForegroundColor Green
        Pop-Location
    }
    
    # Migration Tool
    if (Test-Path "apps\migration-tool\go.mod") {
        Write-Host "  → Téléchargement des dépendances pour migration-tool..." -ForegroundColor Cyan
        Push-Location "apps\migration-tool"
        go mod download
        go mod tidy
        Write-Host "  ✅ migration-tool" -ForegroundColor Green
        Pop-Location
    }
    
    # Webhook Service
    if (Test-Path "apps\webhook-service\go.mod") {
        Write-Host "  → Téléchargement des dépendances pour webhook-service..." -ForegroundColor Cyan
        Push-Location "apps\webhook-service"
        go mod download
        go mod tidy
        Write-Host "  ✅ webhook-service" -ForegroundColor Green
        Pop-Location
    }
} else {
    Write-Host "⚠️  Go n'est pas installé. Les dépendances Go ne seront pas téléchargées." -ForegroundColor Yellow
    Write-Host "   Vous pouvez toujours tester les applications frontend." -ForegroundColor Yellow
}

Write-Host ""
Write-Host "⚙️  Configuration des fichiers d'environnement..." -ForegroundColor Yellow

# Créer les fichiers .env.example
Write-Host "  → Création des fichiers .env.example..." -ForegroundColor Cyan

# API Gateway
if (-not (Test-Path "apps\api-gateway\.env.example")) {
    $content = @"
PORT=8080
JWT_SECRET=your-secret-key-change-in-production
AUTH_SERVICE_URL=http://localhost:8081
CATALOGUE_SERVICE_URL=http://localhost:8082
CHECKOUT_SERVICE_URL=http://localhost:8083
MARKETING_ENGINE_URL=http://localhost:8084
WEBHOOK_SERVICE_URL=http://localhost:8085
MIGRATION_TOOL_URL=http://localhost:8086
"@
    $content | Out-File -FilePath "apps\api-gateway\.env.example" -Encoding UTF8
    Write-Host "  ✅ Créé apps\api-gateway\.env.example" -ForegroundColor Green
}

# Marketing Engine
if (-not (Test-Path "apps\marketing-engine\.env.example")) {
    $content = @"
PORT=8084
DATABASE_URL=postgres://omnisphere:omnisphere_dev@localhost:5432/omnisphere?sslmode=disable
KAFKA_BROKER=localhost:9092
NODE_ENV=development
"@
    $content | Out-File -FilePath "apps\marketing-engine\.env.example" -Encoding UTF8
    Write-Host "  ✅ Créé apps\marketing-engine\.env.example" -ForegroundColor Green
}

# Admin Dashboard
if (-not (Test-Path "web\admin-dashboard\.env.example")) {
    "NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1" | Out-File -FilePath "web\admin-dashboard\.env.example" -Encoding UTF8
    Write-Host "  ✅ Créé web\admin-dashboard\.env.example" -ForegroundColor Green
}

# Storefront Template
if (-not (Test-Path "web\storefront-template\.env.example")) {
    "NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1" | Out-File -FilePath "web\storefront-template\.env.example" -Encoding UTF8
    Write-Host "  ✅ Créé web\storefront-template\.env.example" -ForegroundColor Green
}

# Copier les .env.example vers .env.local pour Next.js
if (-not (Test-Path "web\admin-dashboard\.env.local")) {
    if (Test-Path "web\admin-dashboard\.env.example") {
        Copy-Item "web\admin-dashboard\.env.example" "web\admin-dashboard\.env.local"
        Write-Host "  ✅ Créé web\admin-dashboard\.env.local" -ForegroundColor Green
    }
}

if (-not (Test-Path "web\storefront-template\.env.local")) {
    if (Test-Path "web\storefront-template\.env.example") {
        Copy-Item "web\storefront-template\.env.example" "web\storefront-template\.env.local"
        Write-Host "  ✅ Créé web\storefront-template\.env.local" -ForegroundColor Green
    }
}

Write-Host ""
Write-Host "✅ Configuration terminée !" -ForegroundColor Green
Write-Host ""
Write-Host "📝 Prochaines étapes :" -ForegroundColor Cyan
Write-Host "  1. Démarrer l'infrastructure : make infra-up" -ForegroundColor White
Write-Host "  2. Démarrer les services backend (voir docs/LOCAL_SETUP.md)" -ForegroundColor White
Write-Host "  3. Démarrer les applications frontend :" -ForegroundColor White
Write-Host "     - Admin Dashboard: cd web\admin-dashboard ; npm run dev" -ForegroundColor White
Write-Host "     - Storefront: cd web\storefront-template ; npm run dev" -ForegroundColor White
Write-Host ""
Write-Host "📖 Pour plus de détails, consultez docs/LOCAL_SETUP.md" -ForegroundColor Cyan
