package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/q30-space/spaceapi-endpoint/internal/handlers"
	"github.com/q30-space/spaceapi-endpoint/internal/middleware"
	"github.com/q30-space/spaceapi-endpoint/internal/services"
)

func main() {
	// Load initial SpaceAPI data
	spaceAPI := services.LoadSpaceAPIData()

	// Create handlers
	spaceAPIHandler := handlers.NewSpaceAPIHandler(spaceAPI)

	// Create router
	r := mux.NewRouter()

	// Public API routes (no authentication required)
	r.HandleFunc("/api/space", spaceAPIHandler.GetSpaceAPI).Methods("GET")
	r.HandleFunc("/", spaceAPIHandler.GetSpaceAPI).Methods("GET")

	// Protected API routes (authentication required)
	updateRouter := r.PathPrefix("/api/space").Subrouter()
	updateRouter.Use(middleware.AuthMiddleware)
	updateRouter.HandleFunc("/state", spaceAPIHandler.UpdateState).Methods("POST")
	updateRouter.HandleFunc("/people", spaceAPIHandler.UpdatePeopleCount).Methods("POST")
	updateRouter.HandleFunc("/event", spaceAPIHandler.AddEvent).Methods("POST")

	// Health check
	r.HandleFunc("/health", spaceAPIHandler.HealthCheck).Methods("GET")

	// CORS middleware
	r.Use(middleware.CORSMiddleware)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("SpaceAPI server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
