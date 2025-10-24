// Copyright (C) 2025  pliski@q30.space
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
