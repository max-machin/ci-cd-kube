# Projet CI/CD avec Docker et Kubernetes

Ce projet a pour objectif de mettre en place une pipeline d'intégration continue (CI) et de déploiement continu (CD) utilisant **GitHub Actions**.  
L'application est développée avec **Node.js** et **Express**, contenue dans une image **Docker**, et déployée sur une infrastructure **Kubernetes**.

## Table des Matières
- [Introduction](#introduction)
- [Prérequis](#prérequis)
- [Structure du Projet](#structure-du-projet)
- [Workflow GitHub Actions](#workflow-github-actions)
- [Ajout de Secrets Docker Hub](#ajout-de-secrets-docker-hub)
- [Construction de l'Image Docker](#construction-de-limage-docker)
- [Exécution de l'Application Localement](#exécution-de-lapplication-localement)
- [Tests de l'Application](#tests-de-lapplication)
- [Configuration de la Notification Google Chat](#configuration-de-la-notification-google-chat)
- [Déploiement sur Kubernetes](#déploiement-sur-kubernetes)
- [Déploiement Local avec Minikube via GitHub Actions](#déploiement-local-avec-minikube-via-github-actions)
- [Conclusion](#conclusion)

---

## Introduction

Ce projet met en œuvre une pipeline CI/CD complète pour une application Node.js. La pipeline prend en charge les fonctionnalités suivantes :
1. **Récupération du code** : Déclenchée à chaque `push` ou création de `tag` sur la branche `main`.
2. **Exécution des tests** : Les tests unitaires s'exécutent automatiquement. En cas d'échec, la pipeline est arrêtée.
3. **Construction de l'image Docker** : Création d'une image Docker pour l'application.
4. **Push de l'image Docker** : L'image est poussée vers Docker Hub ou GitHub Container Registry.
5. **Envoi de notifications** : Utilisation de Google Chat pour notifier les succès ou échecs de la pipeline.
6. **Déploiement sur Kubernetes** : L'application est déployée automatiquement sur un cluster Kubernetes via GitHub Actions.

---

## Prérequis

Avant de démarrer, vous devez installer :
- **Node.js** (version 16 ou supérieure)
- **Docker**
- **Kubernetes** (ou un cluster Kubernetes en nuage)
- **GitHub** pour la gestion des versions et l'intégration CI/CD
- **kubectl** pour interagir avec votre cluster Kubernetes

---

## Structure du Projet

```
.
├── express-app
│   ├── index.js
│   ├── package.json
│   ├── tests
│   │   └── test.js
│   └── package-lock.json
├── .github
│   └── workflows
│       └── node.js-ci.yml
├── .gitignore
├── README.md
└── Dockerfile
```

---

## Workflow GitHub Actions

Le fichier `.github/workflows/node.js-ci.yml` automatise l'intégration et le déploiement. Voici un exemple de workflow :

```yaml
name: Node.js CI

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
        node-version: [16]  # Choisir la version de Node.js que tu veux tester

    steps:
    - name: Checkout code
      uses: actions/checkout@v2

    - name: Setup Node.js
      uses: actions/setup-node@v2
      with:
        node-version: ${{ matrix.node-version }}

    - name: Install dependencies
      working-directory: ./express-app
      run: npm ci

    - name: Run Tests
      working-directory: ./express-app
      run: npm test

    - name: Build Docker image
      run: |
        docker build -t machinmax13/ci-cd:${{ github.sha }} ./express-app

    - name: Run Docker container for testing
      run: |
        docker run -d --name ci-cd-test -p 8080:8080 machinmax13/ci-cd:${{ github.sha }}
        sleep 5 # Attendre que l'application se lance
        curl -f http://localhost:8080 || exit 1 # Vérifier si l'application répond
        docker stop ci-cd-test
        docker rm ci-cd-test

    - name: Log in to Docker Hub
      run: echo "${{ secrets.DOCKER_PASSWORD }}" | docker login -u "${{ secrets.DOCKER_USERNAME }}" --password-stdin

    - name: Push Docker image
      run: |
        docker push machinmax13/ci-cd:${{ github.sha }}

    # Démarrer Minikube
    - name: Start Minikube
      run: |
        minikube start
        eval $(minikube docker-env)  # Configure Docker to use Minikube's Docker daemon

    # Construire l'image Docker pour Minikube
    - name: Build Docker image for Minikube
      run: |
        cd express-app
        docker build -t myapp:latest .

    # Déploiement sur Kubernetes via Minikube
    - name: Deploy to Kubernetes
      run: |
        kubectl apply -f ./deployment.yaml

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
# Utiliser une image officielle de Node.js
FROM node:16

# Créer un répertoire de travail pour l'application
WORKDIR /usr/src/app

# Copier les fichiers de dépendances
COPY package*.json ./ 

# Installer les dépendances
RUN npm install

# Copier le reste de l'application
COPY . .

# Exposer le port sur lequel l'application s'exécute
EXPOSE 8080

# Commande pour démarrer l'application
CMD ["node", "index.js"]
```

---

## Exécution de l'Application Localement

1. Clonez le dépôt :
   ```bash
   git clone https://github.com/prenom-nom/ci-cd-kube.git
   cd ci-cd-kube/express-app
   ```

2. Installez les dépendances :
   ```bash
   npm install
   ```

3. Exécutez l'application :
   ```bash
   node index.js
   ```

4. Accédez à l'application sur [http://localhost:8080](http://localhost:8080).

---

## Tests de l'Application 

Pour exécuter les tests localement :

```bash
npm test
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

## Déploiement Local avec Minikube via GitHub Actions

### Prérequis

Avant de déployer localement, assurez-vous d'avoir installé :

- **Minikube** : Pour créer et gérer un cluster Kubernetes local.
- **kubectl** : L'outil en ligne de commande pour interagir avec votre cluster Kubernetes.

### Étapes d'installation de Minikube

1. **Installer Minikube** :
  

 Suivez les instructions officielles pour installer Minikube : [Guide d'installation de Minikube](https://minikube.sigs.k8s.io/docs/start/).

2. **Démarrer Minikube** :
   ```bash
   minikube start
   ```

3. **Vérifier le statut du cluster** :
   ```bash
   minikube status
   ```

### Déploiement via GitHub Actions

Lorsque vous souhaitez déployer votre application sur votre cluster local Minikube via GitHub Actions, assurez-vous que votre fichier `.github/workflows/node.js-ci.yml` contient la configuration suivante :

```yaml
    - name: Set up kubectl
      run: |
        curl -LO "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl"
        chmod +x ./kubectl
        sudo mv ./kubectl /usr/local/bin/kubectl

    - name: Start Minikube
      run: minikube start

    - name: Set Minikube Docker environment
      run: eval $(minikube docker-env)

    - name: Build Docker image for Minikube
      run: |
        docker build -t myapp:latest ./express-app

    - name: Deploy to Minikube
      run: kubectl apply -f ./k8s/deployment.yaml
```

---

## Conclusion

Ce projet présente une configuration complète de CI/CD avec **GitHub Actions**, **Docker**, **Kubernetes**, et des notifications intégrées via **Google Chat**.  
Chaque étape du processus est essentielle pour garantir une livraison continue et une gestion automatisée de vos applications.

---

