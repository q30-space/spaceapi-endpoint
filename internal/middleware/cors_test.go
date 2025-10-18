package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CORSMiddlewareTestSuite struct {
	suite.Suite
}

func TestCORSMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(CORSMiddlewareTestSuite))
}

func (suite *CORSMiddlewareTestSuite) TestCORSMiddleware_RegularRequest() {
	req := httptest.NewRequest("GET", "/api/space", nil)
	w := httptest.NewRecorder()

	handler := CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}))

	handler.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(suite.T(), "GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(suite.T(), "Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(suite.T(), "success", w.Body.String())
}

func (suite *CORSMiddlewareTestSuite) TestCORSMiddleware_POSTRequest() {
	req := httptest.NewRequest("POST", "/api/space/state", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}))

	handler.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(suite.T(), "GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(suite.T(), "Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(suite.T(), "success", w.Body.String())
}

func (suite *CORSMiddlewareTestSuite) TestCORSMiddleware_OPTIONSRequest() {
	req := httptest.NewRequest("OPTIONS", "/api/space", nil)
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type, Authorization")
	w := httptest.NewRecorder()

	handler := CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This should not be called for OPTIONS requests
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("handler-called"))
	}))

	handler.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(suite.T(), "GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(suite.T(), "Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(suite.T(), "", w.Body.String()) // Handler should not be called
}

func (suite *CORSMiddlewareTestSuite) TestCORSMiddleware_OPTIONSRequestWithOrigin() {
	req := httptest.NewRequest("OPTIONS", "/api/space", nil)
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type, Authorization")
	w := httptest.NewRecorder()

	handler := CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("handler-called"))
	}))

	handler.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(suite.T(), "GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(suite.T(), "Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(suite.T(), "", w.Body.String()) // Handler should not be called
}

func (suite *CORSMiddlewareTestSuite) TestCORSMiddleware_WithAuthorizationHeader() {
	req := httptest.NewRequest("GET", "/api/space", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()

	handler := CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}))

	handler.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(suite.T(), "GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(suite.T(), "Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(suite.T(), "success", w.Body.String())
}

func (suite *CORSMiddlewareTestSuite) TestCORSMiddleware_WithContentTypeHeader() {
	req := httptest.NewRequest("POST", "/api/space/state", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}))

	handler.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(suite.T(), "GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(suite.T(), "Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(suite.T(), "success", w.Body.String())
}

func (suite *CORSMiddlewareTestSuite) TestCORSMiddleware_ErrorResponse() {
	req := httptest.NewRequest("GET", "/api/space", nil)
	w := httptest.NewRecorder()

	handler := CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
	}))

	handler.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, w.Code)
	assert.Equal(suite.T(), "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(suite.T(), "GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(suite.T(), "Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(suite.T(), "error", w.Body.String())
}

func (suite *CORSMiddlewareTestSuite) TestCORSMiddleware_AllHTTPMethods() {
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD"}
	
	for _, method := range methods {
		req := httptest.NewRequest(method, "/api/space", nil)
		w := httptest.NewRecorder()

		handler := CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		}))

		handler.ServeHTTP(w, req)

		assert.Equal(suite.T(), http.StatusOK, w.Code)
		assert.Equal(suite.T(), "*", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(suite.T(), "GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
		assert.Equal(suite.T(), "Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
	}
}

func (suite *CORSMiddlewareTestSuite) TestCORSMiddleware_ChainedMiddleware() {
	req := httptest.NewRequest("GET", "/api/space", nil)
	w := httptest.NewRecorder()

	// Chain CORS with another middleware
	handler := CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add custom header
		w.Header().Set("X-Custom-Header", "test-value")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}))

	handler.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(suite.T(), "GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(suite.T(), "Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(suite.T(), "test-value", w.Header().Get("X-Custom-Header"))
	assert.Equal(suite.T(), "success", w.Body.String())
}