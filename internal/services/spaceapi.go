package services

import (
	"encoding/json"
	"github.com/q30-space/spaceapi-endpoint/internal/models"
	"log"
	"os"
)

// LoadSpaceAPIData loads the SpaceAPI configuration from spaceapi.json
func LoadSpaceAPIData() *models.SpaceAPI {
	data, err := os.ReadFile("spaceapi.json")
	if err != nil {
		log.Fatalf("Fatal error: Could not load spaceapi.json: %v", err)
	}

	var spaceAPI models.SpaceAPI
	if err := json.Unmarshal(data, &spaceAPI); err != nil {
		log.Fatalf("Fatal error: Could not parse spaceapi.json: %v", err)
	}

	return &spaceAPI
}
