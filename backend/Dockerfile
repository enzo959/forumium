##---- ÉTAPE 1 : Compilation ----

FROM golang:1.26-alpine

#DÉFINIR un dossier de travail dans le conteneur (ex: /app)
WORKDIR /app
##COPIER les fichiers go.mod et go.sum en premier
    #→ go.mod = liste des dépendances de ton projet
COPY go.mod && go.sum ./

#TÉLÉCHARGER les dépendances Go
    #→ commande : go mod download
RUN go mod download

# COPIER tout le reste du code source
COPY . .
# COMPILER le code en un binaire exécutable
#     → commande : go build -o nom_du_binaire .
RUN go build -o app

# // ---- ÉTAPE 2 : Image finale ----

# UTILISER une image Linux minimaliste comme base
#     → cherche "alpine docker hub", c'est une image Linux de 5MB
FROM alpine:latest
# DÉFINIR le même dossier de travail
WORKDIR /app
# COPIER uniquement le binaire compilé depuis l'étape 1
#     → syntaxe spéciale : COPY --from=nom_etape_1
COPY --from=builder /app/app .

# CRÉER le dossier uploads/ pour les images uploadées
RUN mkdir -p uploads
# # EXPOSER le port 8080
EXPOSE 8080
# DÉFINIR la commande de démarrage
#     → lancer le binaire compilé
CMD ["./app"]