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

package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/q30-space/spaceapi-endpoint/internal/models"
)

type SpaceAPIHandler struct {
	spaceAPI *models.SpaceAPI
}

func NewSpaceAPIHandler(spaceAPI *models.SpaceAPI) *SpaceAPIHandler {
	return &SpaceAPIHandler{
		spaceAPI: spaceAPI,
	}
}

func (h *SpaceAPIHandler) GetSpaceAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(h.spaceAPI); err != nil {
		log.Printf("Error encoding SpaceAPI response: %v", err)
	}
}

func (h *SpaceAPIHandler) UpdateState(w http.ResponseWriter, r *http.Request) {
	var newState models.State
	if err := json.NewDecoder(r.Body).Decode(&newState); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if h.spaceAPI.State == nil {
		h.spaceAPI.State = &models.State{}
	}

	if newState.Open != nil {
		h.spaceAPI.State.Open = newState.Open
	}
	if newState.Message != "" {
		h.spaceAPI.State.Message = newState.Message
	}
	if newState.TriggerPerson != "" {
		h.spaceAPI.State.TriggerPerson = newState.TriggerPerson
	}

	h.spaceAPI.State.Lastchange = time.Now().Unix()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(h.spaceAPI.State); err != nil {
		log.Printf("Error encoding State response: %v", err)
	}

	var ip_address string = r.RemoteAddr
	log.Printf("%s State updated: %+v from %s", time.Unix(h.spaceAPI.State.Lastchange, 0).Format(time.RFC3339), h.spaceAPI.State, ip_address)
}

func (h *SpaceAPIHandler) UpdatePeopleCount(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Value    int    `json:"value"`
		Location string `json:"location,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if h.spaceAPI.Sensors == nil {
		h.spaceAPI.Sensors = &models.Sensors{}
	}

	// Update or add people count sensor
	found := false
	for i, sensor := range h.spaceAPI.Sensors.PeopleNowPresent {
		if sensor.Location == request.Location || (request.Location == "" && sensor.Location == "Main Space") {
			h.spaceAPI.Sensors.PeopleNowPresent[i].Value = request.Value
			h.spaceAPI.Sensors.PeopleNowPresent[i].Lastchange = time.Now().Unix()
			found = true
			break
		}
	}

	if !found {
		h.spaceAPI.Sensors.PeopleNowPresent = append(h.spaceAPI.Sensors.PeopleNowPresent, models.SensorValue{
			Value:      request.Value,
			Location:   request.Location,
			Name:       "People Counter",
			Lastchange: time.Now().Unix(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(h.spaceAPI.Sensors.PeopleNowPresent); err != nil {
		log.Printf("Error encoding PeopleNowPresent response: %v", err)
	}

	log.Printf("%s People count updated: %+v from %s", time.Unix(h.spaceAPI.Sensors.PeopleNowPresent[0].Lastchange, 0).Format(time.RFC3339), h.spaceAPI.Sensors.PeopleNowPresent[0], r.RemoteAddr)
}

func (h *SpaceAPIHandler) AddEvent(w http.ResponseWriter, r *http.Request) {
	var event models.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	event.Timestamp = time.Now().Unix()
	h.spaceAPI.Events = append(h.spaceAPI.Events, event)

	// Keep only last 10 events
	if len(h.spaceAPI.Events) > 10 {
		h.spaceAPI.Events = h.spaceAPI.Events[len(h.spaceAPI.Events)-10:]
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(event); err != nil {
		log.Printf("Error encoding Event response: %v", err)
	}

	log.Printf("%s Event added: %+v from %s", time.Unix(event.Timestamp, 0).Format(time.RFC3339), event, r.RemoteAddr)
}

func (h *SpaceAPIHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("OK")); err != nil {
		log.Printf("Error writing health check response: %v", err)
	}
}
