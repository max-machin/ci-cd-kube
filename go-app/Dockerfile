# Étape 1 : Construction du binaire Go
FROM golang:1.23-alpine as builder

# Installer des dépendances pour Go
RUN apk add --no-cache git

# Définir le répertoire de travail dans le container
WORKDIR /app

# Copier les fichiers du projet Go dans le container
COPY . .

# Télécharger les dépendances et compiler l'application
RUN go mod download
RUN go build -o main .

# Étape 2 : Créer l'image finale
FROM alpine:latest

# Définir le répertoire de travail
WORKDIR /root/

# Copier le binaire Go de l'image builder
COPY --from=builder /app/main .

# Exposer le port de l'API
EXPOSE 3000

# Lancer l'application
CMD ["./main"]
