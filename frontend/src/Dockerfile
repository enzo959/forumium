# On part d'une image Node.js officielle
FROM node:20-alpine

# On définit le dossier de travail dans le conteneur
WORKDIR /app

# On copie les fichiers de dépendances en premier
COPY package.json package-lock.json ./

# On installe les dépendances
RUN npm install

# On copie tout le reste du code source
COPY . .

# port sur lequel Vite écoute
EXPOSE 5173

# Commande pour lancer Vite en mode dev avec hot reload
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]