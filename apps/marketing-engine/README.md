# Marketing Engine

Service d'ingestion d'événements (CDP), segmentation et déclenchement de séquences d'emails pour OmniSphere.

## Fonctionnalités

- Ingestion d'événements depuis Kafka (CDP)
- Segmentation des utilisateurs
- Déclenchement de séquences d'automation (emails)
- Gestion des flows marketing

## Développement

```bash
# Installer les dépendances
npm install

# Lancer en mode développement
npm run start:dev

# Builder
npm run build

# Lancer en production
npm run start:prod
```

## Configuration

Variables d'environnement:
- `PORT` - Port du service (défaut: 8083)
- `KAFKA_BROKERS` - Brokers Kafka (défaut: localhost:9092)

## Architecture

- **EventsConsumer** : Consomme les événements depuis Kafka
- **EventsService** : Traite les événements et enrichit les profils
- **AutomationService** : Gère les séquences d'automation
- **AutomationFlows** : Définit les flows marketing (abandon panier, bienvenue, etc.)

