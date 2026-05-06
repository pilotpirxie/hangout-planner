package handlers
package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestHandler(mockService CalendarServicer) *Handler {
	return &Handler{
		CalendarService: mockService,
	}
}

func TestHealthcheckEndpoint_Success(t *testing.T) {
	handler := newTestHandler(&MockCalendarService{})

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rr := httptest.NewRecorder()

	handler.HealthcheckEndpoint(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}

	var body map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body["status"] != "ok" {
		t.Errorf("expected status ok, got %v", body["status"])
	}
}

func TestEchoEndpoint_Success(t *testing.T) {
	handler := newTestHandler(&MockCalendarService{})

	body := `{"message": "hello"}`
	req := httptest.NewRequest(http.MethodPost, "/api/echo/test-id?name=John&age=25", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer token123")
	req.SetPathValue("id", "test-id")

	rr := httptest.NewRecorder()

	handler.EchoEndpoint(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestEchoEndpoint_MissingBody(t *testing.T) {
	handler := newTestHandler(&MockCalendarService{})

	req := httptest.NewRequest(http.MethodPost, "/api/echo/test-id?name=John&age=25", nil)
	req.Header.Set("Authorization", "Bearer token123")
	req.SetPathValue("id", "test-id")

	rr := httptest.NewRecorder()

	handler.EchoEndpoint(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestEchoEndpoint_MissingQuery(t *testing.T) {
	handler := newTestHandler(&MockCalendarService{})

	body := `{"message": "hello"}`
	req := httptest.NewRequest(http.MethodPost, "/api/echo/test-id?age=25", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer token123")
	req.SetPathValue("id", "test-id")

	rr := httptest.NewRecorder()

	handler.EchoEndpoint(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestEchoEndpoint_MissingParams(t *testing.T) {
	handler := newTestHandler(&MockCalendarService{})

	body := `{"message": "hello"}`
	req := httptest.NewRequest(http.MethodPost, "/api/echo/?name=John&age=25", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer token123")

	rr := httptest.NewRecorder()

	handler.EchoEndpoint(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestEchoEndpoint_MissingHeaders(t *testing.T) {
	handler := newTestHandler(&MockCalendarService{})

	body := `{"message": "hello"}`
	req := httptest.NewRequest(http.MethodPost, "/api/echo/test-id?name=John&age=25", strings.NewReader(body))
	req.SetPathValue("id", "test-id")

	rr := httptest.NewRecorder()

	handler.EchoEndpoint(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}
