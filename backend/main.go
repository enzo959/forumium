package main

import (
	"fmt"
	"forumium/routes"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	routes.SetupRoutes(mux)

	fmt.Println("Serveur démarré sur http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}