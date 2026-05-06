package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRespondJSON(t *testing.T) {
	rr := httptest.NewRecorder()

	RespondJSON(rr, http.StatusOK, map[string]string{"key": "value"})

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}

	var body map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if body["key"] != "value" {
		t.Errorf("expected key=value, got %s", body["key"])
	}
}

func TestRespondError(t *testing.T) {
	rr := httptest.NewRecorder()

	RespondError(rr, http.StatusBadRequest, "something went wrong", nil)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	var body map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if body["error"] != "something went wrong" {
		t.Errorf("expected error 'something went wrong', got %v", body["error"])
	}
}

func TestToJSON(t *testing.T) {
	result := ToJSON(map[string]string{"name": "test"})
	expected := `{"name":"test"}`
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestToJSONPretty(t *testing.T) {
	result := ToJSONPretty(map[string]string{"name": "test"})
	if !strings.Contains(result, `"name"`) || !strings.Contains(result, `"test"`) {
		t.Errorf("expected pretty JSON with name and test, got %s", result)
	}
	if !strings.Contains(result, "\n") {
		t.Error("expected pretty JSON to contain newlines")
	}
}

func TestParseRequest_BodyValidation(t *testing.T) {
	type testBody struct {
		Name string `json:"name" validate:"required"`
	}

	t.Run("valid body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"test"}`))
		rr := httptest.NewRecorder()

		var body testBody
		err := ParseRequest(rr, req, RequestOptions{Body: &body})

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if body.Name != "test" {
			t.Errorf("expected name 'test', got %s", body.Name)
		}
	})

	t.Run("missing required field", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{}`))
		rr := httptest.NewRecorder()

		var body testBody
		err := ParseRequest(rr, req, RequestOptions{Body: &body})

		if err == nil {
			t.Error("expected error for missing required field, got nil")
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`not json`))
		rr := httptest.NewRecorder()

		var body testBody
		err := ParseRequest(rr, req, RequestOptions{Body: &body})

		if err == nil {
			t.Error("expected error for invalid JSON, got nil")
		}
	})
}

func TestParseRequest_QueryParams(t *testing.T) {
	type testQuery struct {
		Name string `query:"name" validate:"required"`
		Age  int    `query:"age" validate:"gte=18"`
	}

	t.Run("valid query params", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?name=John&age=25", nil)
		rr := httptest.NewRecorder()

		var query testQuery
		err := ParseRequest(rr, req, RequestOptions{Query: &query})

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if query.Name != "John" {
			t.Errorf("expected name 'John', got %s", query.Name)
		}
		if query.Age != 25 {
			t.Errorf("expected age 25, got %d", query.Age)
		}
	})

	t.Run("missing required query param", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?age=25", nil)
		rr := httptest.NewRecorder()

		var query testQuery
		err := ParseRequest(rr, req, RequestOptions{Query: &query})

		if err == nil {
			t.Error("expected error for missing required query param, got nil")
		}
	})
}

func TestParseRequest_Headers(t *testing.T) {
	type testHeaders struct {
		AuthToken string `header:"Authorization" validate:"required"`
	}

	t.Run("valid headers", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer token123")
		rr := httptest.NewRecorder()

		var headers testHeaders
		err := ParseRequest(rr, req, RequestOptions{Headers: &headers})

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if headers.AuthToken != "Bearer token123" {
			t.Errorf("expected 'Bearer token123', got %s", headers.AuthToken)
		}
	})

	t.Run("missing required header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		var headers testHeaders
		err := ParseRequest(rr, req, RequestOptions{Headers: &headers})

		if err == nil {
			t.Error("expected error for missing required header, got nil")
		}
	})
}
