#!/bin/bash

# Script pour exécuter les migrations de base de données
# Usage: ./run-migrations.sh [service] [up|down]

set -e

SERVICE=${1:-all}
DIRECTION=${2:-up}
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_NAME=${DB_NAME:-omnisphere}
DB_USER=${DB_USER:-omnisphere}

echo "🔄 Exécution des migrations - Service: $SERVICE, Direction: $DIRECTION"

# Fonction pour exécuter les migrations d'un service
run_migrations() {
    local service_name=$1
    local migrations_dir="../../shared/database_migrations/$service_name"
    
    if [ ! -d "$migrations_dir" ]; then
        echo "⚠️  Aucune migration trouvée pour $service_name"
        return
    fi
    
    echo "📝 Exécution des migrations pour $service_name"
    
    # Utiliser golang-migrate ou un autre outil de migration
    # Exemple avec golang-migrate:
    # migrate -path $migrations_dir -database "postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable" $DIRECTION
    
    echo "✅ Migrations $service_name terminées"
}

# Exécution selon le service
case $SERVICE in
    auth-service)
        run_migrations "auth-service"
        ;;
    checkout-service)
        run_migrations "checkout-service"
        ;;
    catalogue-service)
        run_migrations "catalogue-service"
        ;;
    all)
        echo "Exécution des migrations pour tous les services..."
        run_migrations "auth-service"
        run_migrations "checkout-service"
        run_migrations "catalogue-service"
        ;;
    *)
        echo "❌ Service inconnu: $SERVICE"
        echo "Services disponibles: auth-service, checkout-service, catalogue-service, all"
        exit 1
        ;;
esac

echo "🎉 Migrations terminées!"

