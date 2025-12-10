# Schémas d'événements

Définitions des schémas d'événements pour la communication via Kafka.

## Format

Les schémas sont définis en Protobuf pour une sérialisation efficace et un typage fort.

## Événements disponibles

### User Events (`user_events.proto`)
- `UserRegistered` - Inscription d'un nouvel utilisateur
- `UserLoggedIn` - Connexion d'un utilisateur

### Order Events (`order_events.proto`)
- `OrderCreated` - Création d'une commande
- `OrderPaid` - Paiement d'une commande

### Product Events (`product_events.proto`)
- `ProductViewed` - Visualisation d'un produit
- `ProductAddedToCart` - Ajout d'un produit au panier

### Cart Events (`cart_events.proto`)
- `CartAbandoned` - Abandon de panier
- `CartCompleted` - Finalisation d'un panier

## Utilisation

### Génération du code Go

```bash
protoc --go_out=. --go_opt=paths=source_relative *.proto
```

### Génération du code TypeScript

```bash
protoc --ts_out=. *.proto
```

## Topics Kafka

- `user.events` - Événements utilisateur
- `order.events` - Événements de commande
- `product.events` - Événements produit
- `cart.events` - Événements panier

