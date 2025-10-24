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

package testutil

import (
	"github.com/q30-space/spaceapi-endpoint/internal/models"
	"time"
)

// NewMockSpaceAPI returns a fully populated SpaceAPI struct for testing
func NewMockSpaceAPI() *models.SpaceAPI {
	now := time.Now().Unix()

	return &models.SpaceAPI{
		APICompatibility: []string{"15"},
		Space:            "Test Space",
		Logo:             "https://example.com/logo.png",
		URL:              "https://example.com",
		Location: &models.Location{
			Address:     "123 Test Street, Test City",
			Lat:         40.7128,
			Lon:         -74.0060,
			Timezone:    "America/New_York",
			CountryCode: "US",
		},
		State: &models.State{
			Open:          models.BoolPtr(true),
			Lastchange:    now,
			TriggerPerson: "Test User",
			Message:       "Space is open for testing",
		},
		Events: []models.Event{
			{
				Name:      "Test Event",
				Type:      "check-in",
				Timestamp: now - 3600, // 1 hour ago
				Extra:     "Test event description",
			},
		},
		Contact: models.Contact{
			Email:   "test@example.com",
			IRC:     "#testspace",
			Twitter: "@testspace",
		},
		Sensors: &models.Sensors{
			PeopleNowPresent: []models.SensorValue{
				{
					Value:      3,
					Location:   "Main Space",
					Name:       "People Counter",
					Lastchange: now - 300, // 5 minutes ago
				},
			},
		},
		Feeds: &models.Feeds{
			Blog: &models.Feed{
				Type: "rss",
				URL:  "https://example.com/blog.rss",
			},
		},
		Projects: []string{"Test Project 1", "Test Project 2"},
		Links: []models.Link{
			{
				Name: "Website",
				URL:  "https://example.com",
			},
		},
	}
}

// NewMockState returns a mock State for testing
func NewMockState() models.State {
	return models.State{
		Open:          models.BoolPtr(false),
		Message:       "Space is closed for testing",
		TriggerPerson: "Test Admin",
	}
}

// NewMockEvent returns a mock Event for testing
func NewMockEvent() models.Event {
	return models.Event{
		Name:  "Test Event",
		Type:  "check-in",
		Extra: "Test event extra info",
	}
}

// NewMockPeopleCountRequest returns a mock people count request
func NewMockPeopleCountRequest() map[string]interface{} {
	return map[string]interface{}{
		"value":    5,
		"location": "Test Location",
	}
}
