package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

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
		_, _ = w.Write([]byte("success"))
	}))

	handler.ServeHTTP(w, req)

	suite.Assert().Equal(http.StatusOK, w.Code)
	suite.Assert().Equal("*", w.Header().Get("Access-Control-Allow-Origin"))
	suite.Assert().Equal("GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
	suite.Assert().Equal("Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
	suite.Assert().Equal("success", w.Body.String())
}

func (suite *CORSMiddlewareTestSuite) TestCORSMiddleware_POSTRequest() {
	req := httptest.NewRequest("POST", "/api/space/state", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("success"))
	}))

	handler.ServeHTTP(w, req)

	suite.Assert().Equal(http.StatusOK, w.Code)
	suite.Assert().Equal("*", w.Header().Get("Access-Control-Allow-Origin"))
	suite.Assert().Equal("GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
	suite.Assert().Equal("Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
	suite.Assert().Equal("success", w.Body.String())
}

func (suite *CORSMiddlewareTestSuite) TestCORSMiddleware_OPTIONSRequest() {
	req := httptest.NewRequest("OPTIONS", "/api/space", nil)
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type, Authorization")
	w := httptest.NewRecorder()

	handler := CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This should not be called for OPTIONS requests
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("handler-called"))
	}))

	handler.ServeHTTP(w, req)

	suite.Assert().Equal(http.StatusOK, w.Code)
	suite.Assert().Equal("*", w.Header().Get("Access-Control-Allow-Origin"))
	suite.Assert().Equal("GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
	suite.Assert().Equal("Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
	suite.Assert().Equal("", w.Body.String()) // Handler should not be called
}

func (suite *CORSMiddlewareTestSuite) TestCORSMiddleware_OPTIONSRequestWithOrigin() {
	req := httptest.NewRequest("OPTIONS", "/api/space", nil)
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type, Authorization")
	w := httptest.NewRecorder()

	handler := CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("handler-called"))
	}))

	handler.ServeHTTP(w, req)

	suite.Assert().Equal(http.StatusOK, w.Code)
	suite.Assert().Equal("*", w.Header().Get("Access-Control-Allow-Origin"))
	suite.Assert().Equal("GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
	suite.Assert().Equal("Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
	suite.Assert().Equal("", w.Body.String()) // Handler should not be called
}

func (suite *CORSMiddlewareTestSuite) TestCORSMiddleware_WithAuthorizationHeader() {
	req := httptest.NewRequest("GET", "/api/space", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()

	handler := CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("success"))
	}))

	handler.ServeHTTP(w, req)

	suite.Assert().Equal(http.StatusOK, w.Code)
	suite.Assert().Equal("*", w.Header().Get("Access-Control-Allow-Origin"))
	suite.Assert().Equal("GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
	suite.Assert().Equal("Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
	suite.Assert().Equal("success", w.Body.String())
}

func (suite *CORSMiddlewareTestSuite) TestCORSMiddleware_WithContentTypeHeader() {
	req := httptest.NewRequest("POST", "/api/space/state", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler := CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("success"))
	}))

	handler.ServeHTTP(w, req)

	suite.Assert().Equal(http.StatusOK, w.Code)
	suite.Assert().Equal("*", w.Header().Get("Access-Control-Allow-Origin"))
	suite.Assert().Equal("GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
	suite.Assert().Equal("Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
	suite.Assert().Equal("success", w.Body.String())
}

func (suite *CORSMiddlewareTestSuite) TestCORSMiddleware_ErrorResponse() {
	req := httptest.NewRequest("GET", "/api/space", nil)
	w := httptest.NewRecorder()

	handler := CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("error"))
	}))

	handler.ServeHTTP(w, req)

	suite.Assert().Equal(http.StatusInternalServerError, w.Code)
	suite.Assert().Equal("*", w.Header().Get("Access-Control-Allow-Origin"))
	suite.Assert().Equal("GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
	suite.Assert().Equal("Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
	suite.Assert().Equal("error", w.Body.String())
}

func (suite *CORSMiddlewareTestSuite) TestCORSMiddleware_AllHTTPMethods() {
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD"}

	for _, method := range methods {
		req := httptest.NewRequest(method, "/api/space", nil)
		w := httptest.NewRecorder()

		handler := CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("success"))
		}))

		handler.ServeHTTP(w, req)

		suite.Assert().Equal(http.StatusOK, w.Code)
		suite.Assert().Equal("*", w.Header().Get("Access-Control-Allow-Origin"))
		suite.Assert().Equal("GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
		suite.Assert().Equal("Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
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
		_, _ = w.Write([]byte("success"))
	}))

	handler.ServeHTTP(w, req)

	suite.Assert().Equal(http.StatusOK, w.Code)
	suite.Assert().Equal("*", w.Header().Get("Access-Control-Allow-Origin"))
	suite.Assert().Equal("GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
	suite.Assert().Equal("Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
	suite.Assert().Equal("test-value", w.Header().Get("X-Custom-Header"))
	suite.Assert().Equal("success", w.Body.String())
}
