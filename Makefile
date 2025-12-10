.PHONY: help infra-up infra-down start stop logs clean

help:
	@echo "OmniSphere - Commandes disponibles:"
	@echo "  make infra-up      - Démarrer l'infrastructure (PostgreSQL, Kafka, Redis, Elasticsearch)"
	@echo "  make infra-down    - Arrêter l'infrastructure"
	@echo "  make start         - Démarrer tous les services"
	@echo "  make stop          - Arrêter tous les services"
	@echo "  make logs          - Voir les logs de tous les services"
	@echo "  make clean         - Nettoyer les volumes Docker"

infra-up:
	docker-compose up -d
	@echo "Infrastructure démarrée. Attente de la disponibilité des services..."
	@sleep 5

infra-down:
	docker-compose down

start: infra-up
	@echo "Démarrage des services..."
	@echo "Auth Service: cd apps/auth-service && go run main.go"
	@echo "Checkout Service: cd apps/checkout-service && go run main.go"
	@echo "Catalogue Service: cd apps/catalogue-service && go run main.go"
	@echo "Marketing Engine: cd apps/marketing-engine && npm run start:dev"
	@echo "Admin Dashboard: cd web/admin-dashboard && npm run dev"
	@echo "Storefront: cd web/storefront-template && npm run dev"

stop:
	docker-compose down

logs:
	docker-compose logs -f

clean:
	docker-compose down -v
	@echo "Volumes Docker nettoyés"

