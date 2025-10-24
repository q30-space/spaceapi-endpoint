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
