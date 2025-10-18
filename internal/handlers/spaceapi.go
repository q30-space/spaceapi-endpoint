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
	json.NewEncoder(w).Encode(h.spaceAPI)
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
	json.NewEncoder(w).Encode(h.spaceAPI.State)

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
	json.NewEncoder(w).Encode(h.spaceAPI.Sensors.PeopleNowPresent)

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
	json.NewEncoder(w).Encode(event)

	log.Printf("%s Event added: %+v from %s", time.Unix(event.Timestamp, 0).Format(time.RFC3339), event, r.RemoteAddr)
}

func (h *SpaceAPIHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
