package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/q30-space/spaceapi-endpoint/internal/models"
	"github.com/q30-space/spaceapi-endpoint/internal/testutil"
	"github.com/stretchr/testify/suite"
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

	suite.Assert().Equal(http.StatusOK, w.Code)
	suite.Assert().Equal("application/json", w.Header().Get("Content-Type"))

	var response models.SpaceAPI
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Assert().NoError(err)
	suite.Assert().Equal("Test Space", response.Space)
	suite.Assert().Equal("https://example.com", response.URL)
	suite.Assert().NotNil(response.State)
	suite.Assert().True(*response.State.Open)
}

func (suite *SpaceAPIHandlerTestSuite) TestUpdateState_ValidUpdate() {
	stateUpdate := testutil.NewMockState()
	jsonData, _ := json.Marshal(stateUpdate)

	req := httptest.NewRequest("POST", "/api/space/state", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.handler.UpdateState(w, req)

	suite.Assert().Equal(http.StatusOK, w.Code)
	suite.Assert().Equal("application/json", w.Header().Get("Content-Type"))

	var response models.State
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Assert().NoError(err)
	suite.Assert().False(*response.Open)
	suite.Assert().Equal("Space is closed for testing", response.Message)
	suite.Assert().Equal("Test Admin", response.TriggerPerson)
	suite.Assert().NotZero(response.Lastchange)
}

func (suite *SpaceAPIHandlerTestSuite) TestUpdateState_PartialUpdate() {
	partialUpdate := map[string]interface{}{
		"open":    true,
		"message": "Partial update test",
	}
	jsonData, _ := json.Marshal(partialUpdate)

	req := httptest.NewRequest("POST", "/api/space/state", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.handler.UpdateState(w, req)

	suite.Assert().Equal(http.StatusOK, w.Code)

	var response models.State
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Assert().NoError(err)
	suite.Assert().True(*response.Open)
	suite.Assert().Equal("Partial update test", response.Message)
	// TriggerPerson should remain unchanged from original
	suite.Assert().Equal("Test User", response.TriggerPerson)
}

func (suite *SpaceAPIHandlerTestSuite) TestUpdateState_InvalidJSON() {
	req := httptest.NewRequest("POST", "/api/space/state", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.handler.UpdateState(w, req)

	suite.Assert().Equal(http.StatusBadRequest, w.Code)
	suite.Assert().Equal("Invalid JSON\n", w.Body.String())
}

func (suite *SpaceAPIHandlerTestSuite) TestUpdatePeopleCount_NewLocation() {
	peopleData := testutil.NewMockPeopleCountRequest()
	jsonData, _ := json.Marshal(peopleData)

	req := httptest.NewRequest("POST", "/api/space/people", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.handler.UpdatePeopleCount(w, req)

	suite.Assert().Equal(http.StatusOK, w.Code)
	suite.Assert().Equal("application/json", w.Header().Get("Content-Type"))

	var response []models.SensorValue
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Assert().NoError(err)
	suite.Assert().Len(response, 2) // Original + new

	// Find the new sensor
	var newSensor *models.SensorValue
	for _, sensor := range response {
		if sensor.Location == "Test Location" {
			newSensor = &sensor
			break
		}
	}
	suite.Assert().NotNil(newSensor)
	suite.Assert().Equal(float64(5), newSensor.Value)
	suite.Assert().Equal("People Counter", newSensor.Name)
	suite.Assert().NotZero(newSensor.Lastchange)
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
	suite.Assert().Equal(http.StatusOK, w.Code)

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

	suite.Assert().Equal(http.StatusOK, w.Code)

	var response []models.SensorValue
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Assert().NoError(err)
	suite.Assert().Len(response, 1) // Still only one sensor
	suite.Assert().Equal(float64(7), response[0].Value)
	suite.Assert().Equal("Main Space", response[0].Location)
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

	suite.Assert().Equal(http.StatusOK, w.Code)

	var response []models.SensorValue
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Assert().NoError(err)

	// Should update the existing "Main Space" sensor
	suite.Assert().Len(response, 1)
	suite.Assert().Equal(float64(2), response[0].Value)
	suite.Assert().Equal("Main Space", response[0].Location)
}

func (suite *SpaceAPIHandlerTestSuite) TestUpdatePeopleCount_InvalidJSON() {
	req := httptest.NewRequest("POST", "/api/space/people", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.handler.UpdatePeopleCount(w, req)

	suite.Assert().Equal(http.StatusBadRequest, w.Code)
	suite.Assert().Equal("Invalid JSON\n", w.Body.String())
}

func (suite *SpaceAPIHandlerTestSuite) TestAddEvent_ValidEvent() {
	event := testutil.NewMockEvent()
	jsonData, _ := json.Marshal(event)

	req := httptest.NewRequest("POST", "/api/space/event", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.handler.AddEvent(w, req)

	suite.Assert().Equal(http.StatusOK, w.Code)
	suite.Assert().Equal("application/json", w.Header().Get("Content-Type"))

	var response models.Event
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.Assert().NoError(err)
	suite.Assert().Equal("Test Event", response.Name)
	suite.Assert().Equal("check-in", response.Type)
	suite.Assert().Equal("Test event extra info", response.Extra)
	suite.Assert().NotZero(response.Timestamp)
	suite.Assert().True(response.Timestamp > time.Now().Unix()-5) // Within last 5 seconds
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
		suite.Assert().Equal(http.StatusOK, w.Code)
	}

	// Verify only 10 events remain
	suite.Assert().Len(suite.handler.spaceAPI.Events, 10)
}

func (suite *SpaceAPIHandlerTestSuite) TestAddEvent_InvalidJSON() {
	req := httptest.NewRequest("POST", "/api/space/event", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	suite.handler.AddEvent(w, req)

	suite.Assert().Equal(http.StatusBadRequest, w.Code)
	suite.Assert().Equal("Invalid JSON\n", w.Body.String())
}

func (suite *SpaceAPIHandlerTestSuite) TestHealthCheck() {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	suite.handler.HealthCheck(w, req)

	suite.Assert().Equal(http.StatusOK, w.Code)
	suite.Assert().Equal("OK", w.Body.String())
}
