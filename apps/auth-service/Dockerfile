# Stage 1: Build
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Installer les dépendances système si nécessaire
RUN apk add --no-cache git

# Copier les fichiers de dépendances
COPY go.mod go.sum ./
RUN go mod download

# Copier le code source
COPY . .

# Builder l'application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o auth-service .

# Stage 2: Runtime
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Copier le binaire depuis le stage de build
COPY --from=builder /app/auth-service .

# Exposer le port
EXPOSE 8080

# Commande par défaut
CMD ["./auth-service"]

