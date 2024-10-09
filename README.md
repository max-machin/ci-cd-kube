# Projet CI/CD avec Docker et Kubernetes

Ce projet a pour objectif de mettre en place une pipeline d'intégration continue (CI) et de déploiement continu (CD) utilisant **GitHub Actions**.  
L'application est développée avec **Go**, contenue dans une image **Docker**, et déployée sur une infrastructure **Kubernetes**.

## Table des Matières
- [Introduction](#introduction)
- [Prérequis](#prérequis)
- [Structure du Projet](#structure-du-projet)
- [Workflow GitHub Actions](#workflow-github-actions)
- [Ajout de Secrets Docker Hub](#ajout-de-secrets-docker-hub)
- [Construction de l'Image Docker](#construction-de-limage-docker)
- [Tests de l'Application](#tests-de-lapplication)
- [Configuration de la Notification Google Chat](#configuration-de-la-notification-google-chat)
- [Déploiement sur Kubernetes](#déploiement-sur-kubernetes)
- [Conclusion](#conclusion)

---

## Introduction

Ce projet met en œuvre une pipeline CI/CD complète pour une application Go. La pipeline prend en charge les fonctionnalités suivantes :
1. **Récupération du code** : Déclenchée à chaque `push` ou création de `tag` sur la branche `main`.
2. **Exécution des tests** : Les tests unitaires s'exécutent automatiquement. En cas d'échec, la pipeline est arrêtée.
3. **Construction de l'image Docker** : Création d'une image Docker pour l'application.
4. **Push de l'image Docker** : L'image est poussée vers Docker Hub ou GitHub Container Registry.
5. **Envoi de notifications** : Utilisation de Google Chat pour notifier les succès ou échecs de la pipeline.
6. **Déploiement sur Kubernetes** : L'application est déployée automatiquement sur un cluster Kubernetes via GitHub Actions.

---

## Prérequis

Avant de démarrer, vous devez installer :
- **Go** (version 1.23)
- **Docker**
- **Kubernetes** (ou un cluster Kubernetes en nuage)
- **GitHub** pour la gestion des versions et l'intégration CI/CD
- **kubectl** pour interagir avec votre cluster Kubernetes

---

## Structure du Projet

```
.
├── go-app
│   ├── main.go
│   └── go.sum
├── .github
│   └── workflows
│       └── ci-cd.yml
├── .gitignore
├── README.md
└── Dockerfile
```

---

## Workflow GitHub Actions

Le fichier `.github/workflows/ci-cd.yml` automatise l'intégration et le déploiement. Voici un exemple de workflow :

```yaml
name: Go CI

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest

    strategy:
      matrix:
        go-version: [1.23]

    steps:
    - name: Checkout code
      uses: actions/checkout@v2

    - name: Setup Go
      uses: actions/setup-go@v2
      with:
        go-version: ${{ matrix.go-version }}

    # Installer les dépendances
    - name: Install dependencies
      run: |
        cd go-app
        go mod tidy  # Nettoie les dépendances et les met à jour

    # Compiler l'application
    - name: Build Go application
      run: |
        cd go-app
        go build -o myapp ./main.go  # Remplacez ./main.go par le chemin de votre fichier principal

    # Tester l'application Go
    - name: Run Go tests
      run: |
        cd go-app
        go test ./...  # Exécute les tests sur tous les packages

    # Construire une image Docker
    - name: Build Docker image
      run: |
        docker build -t machinmax13/ci-cd:${{ github.sha }} -f go-app/Dockerfile ./go-app

    # Tester le conteneur Docker localement
    - name: Run Docker container for testing
      run: |
        docker run -d --name ci-cd-test -p 8080:3000 machinmax13/ci-cd:${{ github.sha }}
        sleep 5 # Attendre que l'application se lance
        curl -f http://localhost:8080 || exit 1 
        docker stop ci-cd-test
        docker rm ci-cd-test

    # Connexion à Docker Hub
    - name: Log in to Docker Hub
      run: echo "${{ secrets.DOCKER_PASSWORD }}" | docker login -u "${{ secrets.DOCKER_USERNAME }}" --password-stdin

    # Pousser l'image Docker sur Docker Hub
    - name: Push Docker image
      run: |
        docker push machinmax13/ci-cd:${{ github.sha }}


    # # Installer kubectl
    # - name: Set up kubectl
    #   run: |
    #     curl -LO "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl"
    #     chmod +x ./kubectl
    #     sudo mv ./kubectl /usr/local/bin/kubectl

    # # Démarrer Minikube
    # - name: Start Minikube
    #   run: minikube start

    # # Configurer l'environnement Docker de Minikube
    # - name: Set Minikube Docker environment
    #   run: eval $(minikube docker-env)

    # # Construire l'image Docker pour Minikube
    # - name: Build Docker image for Minikube
    #   run: |
    #     docker build -t machinmax13/ci-cd:${{ github.sha }} ./express-app  # Utilisez le même tag

    # # Déploiement sur Kubernetes via Minikube
    # - name: Deploy to Minikube
    #   run: |
    #     kubectl apply -f ./k8s/deployment.yaml

    # Envoyer une notification Google Chat après le succès du déploiement
    - name: Send Success Notification to Google Chat
      if: success()
      run: |
        curl -X POST -H 'Content-Type: application/json' \
        -d '{"text": "✅ Déploiement réussi pour le commit ${{ github.sha }} sur la branche ${{ github.ref_name }}."}' \
        "${{ secrets.GOOGLE_CHAT_WEBHOOK }}"

    # Envoyer une notification Google Chat en cas d'échec de la pipeline
    - name: Send Failure Notification to Google Chat
      if: failure()
      run: |
        curl -X POST -H 'Content-Type: application/json' \
        -d '{"text": "❌ Échec du déploiement pour le commit ${{ github.sha }} sur la branche ${{ github.ref_name }}. Vérifiez les logs pour plus de détails."}' \
        "${{ secrets.GOOGLE_CHAT_WEBHOOK }}"

```

---

## Ajout de Secrets Docker Hub

1. Allez dans **Settings** > **Secrets and variables** > **Actions**.
2. Cliquez sur **New repository secret** et ajoutez les secrets :
   - `DOCKER_USERNAME`
   - `DOCKER_PASSWORD`

---

## Construction de l'Image Docker

Créez un fichier `Dockerfile` à la racine de votre projet :

```dockerfile
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
```

---

## Configuration de la Notification Google Chat

### Étape 1 : Créer un Webhook Google Chat

1. Accédez à **Google Chat** et créez un espace (ou utilisez-en un existant).
2. Cliquez sur le nom de l'espace, puis allez dans **Gérer les Webhooks** et créez-en un nouveau. Copiez l'URL.

### Étape 2 : Ajouter le Webhook à GitHub Secrets

1. Allez dans **Settings** > **Secrets and variables** > **Actions**.
2. Ajoutez un secret nommé `GOOGLE_CHAT_WEBHOOK` avec l'URL du webhook.

### Étape 3 : Intégrer dans le Workflow GitHub Actions

Ajoutez les étapes suivantes pour notifier Google Chat :

```yaml
- name: Send Success Notification to Google Chat
  if: success()
  run: |
    curl -X POST -H 'Content-Type: application/json' \
    -d '{"text": "✅ Déploiement réussi pour le commit ${{ github.sha }} sur la branche ${{ github.ref_name }}."}' \
    "${{ secrets.GOOGLE_CHAT_WEBHOOK }}"

- name: Send Failure Notification to Google Chat
  if: failure()
  run: |
    curl -X POST -H 'Content-Type: application/json' \
    -d '{"text": "❌ Échec du déploiement pour le commit ${{ github.sha }} sur la branche ${{ github.ref_name }}. Vérifiez les logs."}' \
    "${{ secrets.GOOGLE_CHAT_WEBHOOK }}"
```

---

## Déploiement sur Kubernetes

Pour déployer votre application sur Kubernetes, créez un fichier `deployment.yaml` :

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  replicas: 2
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
    spec:
      containers:
        - name: myapp
          image: myapp:latest
          ports:
            - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: myapp-service
spec:
  type: NodePort
  ports:
    - port: 8080
      targetPort: 8080
      nodePort: 30000
  selector:
    app: myapp
```

Appliquez la configuration avec :

```bash
kubectl apply -f deployment.yaml
```

Le déploiement de l'application se fait via le fichier de workflow GitHub Actions lorsque des changements sont poussés vers la branche `main`.

---

## Conclusion

Ce projet présente une configuration complète de CI/CD avec **GitHub Actions**, **Docker**, **Kubernetes**, et des notifications intégrées via **Google Chat**.  
Chaque étape du processus est essentielle pour garantir une livraison continue et une gestion automatisée de vos applications.

---

