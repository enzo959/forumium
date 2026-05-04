package routes

import (
	"forumium/handlers"
	"net/http"
)

func SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/health", handlers.HealthCheck)
}
