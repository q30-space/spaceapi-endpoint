package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/q30-space/spaceapi-endpoint/internal/models"
	"github.com/q30-space/spaceapi-endpoint/internal/testutil"
)

type SpaceAPIHandlerTestSuite struct {
	suite.Suite
	handler *SpaceAPIHandler
}

func (suite *SpaceAPIHandlerTestSuite) SetupTest() {
	mockSpaceAPI := testutil.NewMockSpaceAPI()
	suite.handler = NewSpaceAPIHandler(mockSpaceAPI)
}

func TestSpaceAPIHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(SpaceAPIHandlerTestSuite))
}

func (suite *SpaceAPIHandlerTestSuite) TestGetSpaceAPI() {
	req := httptest.NewRequest("GET", "/api/space", nil)
	w := httptest.NewRecorder()

	suite.handler.GetSpaceAPI(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "application/json", w.Header().Get("Content-Type"))

	var response models.SpaceAPI
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Test Space", response.Space)
	assert.Equal(suite.T(), "https://example.com", response.URL)
	assert.NotNil(suite.T(), response.State)
	assert.True(suite.T(), *response.State.Open)
}

func (suite *SpaceAPIHandlerTestSuite) TestUpdateState_ValidUpdate() {
	stateUpdate := testutil.NewMockState()
	jsonData, _ := json.Marshal(stateUpdate)
	
	req := httptest.NewRequest("POST", "/api/space/state", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.handler.UpdateState(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "application/json", w.Header().Get("Content-Type"))

	var response models.State
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), *response.Open)
	assert.Equal(suite.T(), "Space is closed for testing", response.Message)
	assert.Equal(suite.T(), "Test Admin", response.TriggerPerson)
	assert.NotZero(suite.T(), response.Lastchange)
}

func (suite *SpaceAPIHandlerTestSuite) TestUpdateState_PartialUpdate() {
	partialUpdate := map[string]interface{}{
		"open": true,
		"message": "Partial update test",
	}
	jsonData, _ := json.Marshal(partialUpdate)
	
	req := httptest.NewRequest("POST", "/api/space/state", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.handler.UpdateState(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response models.State
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), *response.Open)
	assert.Equal(suite.T(), "Partial update test", response.Message)
	// TriggerPerson should remain unchanged from original
	assert.Equal(suite.T(), "Test User", response.TriggerPerson)
}

func (suite *SpaceAPIHandlerTestSuite) TestUpdateState_InvalidJSON() {
	req := httptest.NewRequest("POST", "/api/space/state", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.handler.UpdateState(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	assert.Equal(suite.T(), "Invalid JSON\n", w.Body.String())
}

func (suite *SpaceAPIHandlerTestSuite) TestUpdatePeopleCount_NewLocation() {
	peopleData := testutil.NewMockPeopleCountRequest()
	jsonData, _ := json.Marshal(peopleData)
	
	req := httptest.NewRequest("POST", "/api/space/people", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.handler.UpdatePeopleCount(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "application/json", w.Header().Get("Content-Type"))

	var response []models.SensorValue
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), response, 2) // Original + new

	// Find the new sensor
	var newSensor *models.SensorValue
	for _, sensor := range response {
		if sensor.Location == "Test Location" {
			newSensor = &sensor
			break
		}
	}
	assert.NotNil(suite.T(), newSensor)
	assert.Equal(suite.T(), float64(5), newSensor.Value)
	assert.Equal(suite.T(), "People Counter", newSensor.Name)
	assert.NotZero(suite.T(), newSensor.Lastchange)
}

func (suite *SpaceAPIHandlerTestSuite) TestUpdatePeopleCount_UpdateExisting() {
	// First, add a sensor
	peopleData := map[string]interface{}{
		"value":    3,
		"location": "Main Space",
	}
	jsonData, _ := json.Marshal(peopleData)
	
	req := httptest.NewRequest("POST", "/api/space/people", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.handler.UpdatePeopleCount(w, req)
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// Now update the same location
	updateData := map[string]interface{}{
		"value":    7,
		"location": "Main Space",
	}
	jsonData, _ = json.Marshal(updateData)
	
	req = httptest.NewRequest("POST", "/api/space/people", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	suite.handler.UpdatePeopleCount(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response []models.SensorValue
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), response, 1) // Still only one sensor
	assert.Equal(suite.T(), float64(7), response[0].Value)
	assert.Equal(suite.T(), "Main Space", response[0].Location)
}

func (suite *SpaceAPIHandlerTestSuite) TestUpdatePeopleCount_DefaultLocation() {
	peopleData := map[string]interface{}{
		"value": 2,
	}
	jsonData, _ := json.Marshal(peopleData)
	
	req := httptest.NewRequest("POST", "/api/space/people", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.handler.UpdatePeopleCount(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var response []models.SensorValue
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	// Should update the existing "Main Space" sensor
	assert.Len(suite.T(), response, 1)
	assert.Equal(suite.T(), float64(2), response[0].Value)
	assert.Equal(suite.T(), "Main Space", response[0].Location)
}

func (suite *SpaceAPIHandlerTestSuite) TestUpdatePeopleCount_InvalidJSON() {
	req := httptest.NewRequest("POST", "/api/space/people", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.handler.UpdatePeopleCount(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	assert.Equal(suite.T(), "Invalid JSON\n", w.Body.String())
}

func (suite *SpaceAPIHandlerTestSuite) TestAddEvent_ValidEvent() {
	event := testutil.NewMockEvent()
	jsonData, _ := json.Marshal(event)
	
	req := httptest.NewRequest("POST", "/api/space/event", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.handler.AddEvent(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "application/json", w.Header().Get("Content-Type"))

	var response models.Event
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Test Event", response.Name)
	assert.Equal(suite.T(), "check-in", response.Type)
	assert.Equal(suite.T(), "Test event extra info", response.Extra)
	assert.NotZero(suite.T(), response.Timestamp)
	assert.True(suite.T(), response.Timestamp > time.Now().Unix()-5) // Within last 5 seconds
}

func (suite *SpaceAPIHandlerTestSuite) TestAddEvent_EventListLimit() {
	// Add 12 events (original has 1, so we'll have 13 total, should keep last 10)
	for i := 0; i < 12; i++ {
		event := models.Event{
			Name:  "Event",
			Type:  "test",
			Extra: "Test event",
		}
		jsonData, _ := json.Marshal(event)
		
		req := httptest.NewRequest("POST", "/api/space/event", bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		suite.handler.AddEvent(w, req)
		assert.Equal(suite.T(), http.StatusOK, w.Code)
	}

	// Verify only 10 events remain
	assert.Len(suite.T(), suite.handler.spaceAPI.Events, 10)
}

func (suite *SpaceAPIHandlerTestSuite) TestAddEvent_InvalidJSON() {
	req := httptest.NewRequest("POST", "/api/space/event", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.handler.AddEvent(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	assert.Equal(suite.T(), "Invalid JSON\n", w.Body.String())
}

func (suite *SpaceAPIHandlerTestSuite) TestHealthCheck() {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	suite.handler.HealthCheck(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "OK", w.Body.String())
}