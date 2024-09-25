# Utiliser une image Node.js officielle
FROM node:14-alpine

# Définir le répertoire de travail dans le conteneur
WORKDIR /usr/src/app

# Copier les fichiers package.json et package-lock.json pour installer les dépendances
COPY app/package*.json ./

# Installer les dépendances de l'application
RUN npm install --production

# Copier tout le code dans le conteneur
COPY app/ .

# Exposer le port sur lequel l'application va tourner
EXPOSE 3000

# Commande pour démarrer l'application
CMD ["npm", "start"]
