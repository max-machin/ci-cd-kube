# Projet CI/CD avec Docker et Kubernetes

Ce projet a pour objectif de mettre en place une pipeline d'intégration continue (CI) et de déploiement continu (CD) utilisant **GitHub Actions**. L'application est développée avec Node.js et Express, et est contenue dans une image Docker qui sera déployée sur une infrastructure Kubernetes.

La pipeline sera déclenchée automatiquement à chaque push sur la branche `main` et lors de la création d'un tag Git, garantissant ainsi une intégration fluide et efficace tout au long du développement du projet.

## Table des Matières
- [Introduction](#introduction)
- [Prérequis](#prérequis)
- [Structure du Projet](#structure-du-projet)
- [Exécution de l'Application Localement](#exécution-de-lapplication-localement)
- [Tests de l'Application](#tests-de-lapplication)
- [Construction de l'Image Docker](#construction-de-limage-docker)
- [Configuration de Docker](#configuration-de-docker)
- [Déploiement sur Kubernetes](#déploiement-sur-kubernetes)
- [Workflow GitHub Actions](#workflow-github-actions)
- [Ajout de Secrets Docker Hub](#ajout-de-secrets-docker-hub)

## Introduction

Dans ce projet, vous allez configurer une pipeline CI/CD qui permettra d'automatiser le cycle de vie d'une application. L'objectif est d'implémenter les étapes suivantes :

1. **Récupération du code** : Chaque push sur la branche `main` ou création d'un tag Git déclenche la pipeline, récupérant le code source.
  
2. **Exécution des tests** : Avant toute construction, des tests seront lancés pour s'assurer du bon fonctionnement de l'application. En cas d'échec d'un test, la pipeline s'arrêtera immédiatement.

3. **Construction de l'image Docker** : Une image Docker sera construite, adaptée à l'environnement de développement ou de production selon la branche ou le tag.

4. **Push de l'image Docker** : L'image Docker sera poussée vers un registre tel que GitHub Container Registry ou Docker Hub, en utilisant des tags pour distinguer les versions de développement et de production.

5. **Envoi de notifications** : À chaque exécution de la pipeline, des notifications seront envoyées via Google Chat, contenant des informations clés comme le commit responsable et le statut de la pipeline.

6. **Déploiement sur l'infrastructure Kubernetes** : Les fichiers de configuration pour déployer l'image Docker seront créés et appliqués sur l'infrastructure Kubernetes fournie, avec une vérification du bon fonctionnement de l'application.

## Prérequis

- Node.js (version 16 ou supérieure)
- Docker
- Kubernetes (ou Minikube pour le développement local)
- GitHub pour la gestion de version et les actions CI/CD

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

## Exécution de l'Application Localement

Pour exécuter l'application localement, suivez ces étapes :

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

4. Accédez à l'application dans votre navigateur à l'adresse [http://localhost:8080](http://localhost:8080).

## Tests de l'Application

Pour exécuter les tests, assurez-vous que les dépendances sont installées, puis exécutez :

```bash
npm test
```

## Construction de l'Image Docker

Pour construire l'image Docker de votre application, créez un fichier `Dockerfile` à la racine de votre projet avec le contenu suivant :

### Dockerfile

```dockerfile
# Utiliser une image de base officielle de Node.js
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

## Configuration de Docker

Pour construire l'image Docker, exécutez les commandes suivantes à la racine de votre projet :

1. **Construire l'image Docker :**
   ```bash
   docker build -t myusername/myapp:latest .
   ```

2. **Exécuter l'image Docker :**
   ```bash
   docker run -p 8080:8080 myusername/myapp:latest
   ```

Accédez à l'application dans votre navigateur à l'adresse [http://localhost:8080](http://localhost:8080) pour vérifier qu'elle fonctionne comme prévu.

## Déploiement sur Kubernetes

Pour déployer votre application sur Kubernetes, vous devrez créer un fichier de configuration Kubernetes (par exemple `deployment.yaml`) qui définit le déploiement et le service. Voici un exemple de fichier `deployment.yaml` :

### deployment.yaml

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
        image: myusername/myapp:latest
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

Pour appliquer cette configuration, exécutez :

```bash
kubectl apply -f deployment.yaml
```

## Workflow GitHub Actions

Le fichier de workflow GitHub Actions (situé dans `.github/workflows/node.js-ci.yml`) est configuré pour automatiser le processus de CI/CD. Voici un exemple de contenu :

### node.js-ci.yml

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

    steps:
    - name: Checkout code
      uses: actions/checkout@v2

    - name: Setup Node.js
      uses: actions/setup-node@v2
      with:
        node-version: 16

    - name: Install dependencies
      working-directory: ./express-app
      run: npm ci

    - name: Run Tests
      working-directory: ./express-app
      run: npm test

    - name: Build Docker image
      run: |
        docker build -t myusername/myapp:${{ github.sha }} ./express-app

    - name: Run Docker container for testing
      run: |
        docker run -d --name myapp-test -p 8080:8080 myusername/myapp:${{ github.sha }}
        sleep 5 # Attendre que l'application se lance
        curl -f http://localhost:8080 || exit 1 # Vérifier si l'application répond
        docker stop myapp-test
        docker rm myapp-test

    - name: Log in to Docker Hub
      run: echo "${{ secrets.DOCKER_PASSWORD }}" | docker login -u "${{ secrets.DOCKER_USERNAME }}" --password-stdin

    - name: Push Docker image
      run: |
        docker push myusername/myapp:${{ github.sha }}

    - name: Deploy to Kubernetes
      run: |
        kubectl apply -f ./k8s/deployment.yaml
```

## Ajout de Secrets Docker Hub

Pour vous connecter à Docker Hub à partir de votre pipeline GitHub Actions, vous devez ajouter des secrets à votre dépôt GitHub. Voici comment faire :

1. Accédez à votre dépôt sur GitHub.
2. Cliquez sur l'onglet **Settings** (paramètres).
3. Dans le menu de gauche, sélectionnez **Secrets and variables** puis cliquez sur **Actions**.
4. Cliquez sur **New repository secret**.
5. Ajoutez les secrets suivants :
   - `DOCKER_USERNAME` : Votre nom d'utilisateur Docker Hub.
   - `DOCKER_PASSWORD` : Votre mot de passe Docker Hub.

Assurez-vous de nommer ces secrets exactement comme mentionné ci-dessus, car ils seront référencés dans le fichier de workflow.

## Conclusion

Ce projet illustre le processus de mise en place d'une pipeline CI/CD avec Docker et Kubernetes. Il est essentiel de s'assurer que chaque étape fonctionne correctement afin d'automatiser le déploiement de l'application.

 du projet.

### Notes Importantes :
- **Remplacez** `myusername/myapp` par votre nom d'utilisateur et le nom de votre application.
- **Personnalisez** ce README selon vos besoins et les détails spécifiques de votre projet.

---

