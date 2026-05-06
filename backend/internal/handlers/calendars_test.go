package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"meeting-planner/backend/internal/db/sqlc"
	"meeting-planner/backend/internal/services"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func makeUUID() pgtype.UUID {
	id := uuid.New()
	return pgtype.UUID{Bytes: [16]byte(id), Valid: true}
}

func makeJSONBody(t *testing.T, body any) *bytes.Reader {
	t.Helper()
	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("failed to marshal body: %v", err)
	}
	return bytes.NewReader(data)
}

func decodeResponse(t *testing.T, rr *httptest.ResponseRecorder, target any) {
	t.Helper()
	if err := json.NewDecoder(rr.Body).Decode(target); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
}

func mustMakeTimestamptz(tm time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: tm, Valid: true}
}

func TestCreateCalendarEndpoint_Success(t *testing.T) {
	calendarUUID := makeUUID()
	adminToken := "test-admin-token"

	mock := &MockCalendarService{
		CreateCalendarFunc: func(_ context.Context, _ *services.CreateCalendarInput) (sqlc.CreateCalendarRow, error) {
			return sqlc.CreateCalendarRow{
				ID:         calendarUUID,
				AdminToken: adminToken,
			}, nil
		},
	}

	handler := newTestHandler(mock)

	body := map[string]any{
		"title":       "Test Calendar",
		"description": "A test calendar",
	}

	req := httptest.NewRequest(http.MethodPost, "/api/calendars", makeJSONBody(t, body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.CreateCalendarEndpoint(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d; body: %s", http.StatusCreated, rr.Code, rr.Body.String())
	}

	var resp CreateCalendarResponse
	decodeResponse(t, rr, &resp)

	if resp.ID != uuid.UUID(calendarUUID.Bytes).String() {
		t.Errorf("expected ID %s, got %s", uuid.UUID(calendarUUID.Bytes).String(), resp.ID)
	}
	if resp.AdminToken != adminToken {
		t.Errorf("expected admin token %s, got %s", adminToken, resp.AdminToken)
	}
}

func TestCreateCalendarEndpoint_MissingTitle(t *testing.T) {
	mock := &MockCalendarService{}
	handler := newTestHandler(mock)

	body := map[string]any{
		"description": "A test calendar",
	}

	req := httptest.NewRequest(http.MethodPost, "/api/calendars", makeJSONBody(t, body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.CreateCalendarEndpoint(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d; body: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
}

func TestCreateCalendarEndpoint_TitleTooShort(t *testing.T) {
	mock := &MockCalendarService{}
	handler := newTestHandler(mock)

	body := map[string]any{
		"title": "ab",
	}

	req := httptest.NewRequest(http.MethodPost, "/api/calendars", makeJSONBody(t, body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.CreateCalendarEndpoint(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d; body: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
}

func TestCreateCalendarEndpoint_ServiceError(t *testing.T) {
	mock := &MockCalendarService{
		CreateCalendarFunc: func(_ context.Context, _ *services.CreateCalendarInput) (sqlc.CreateCalendarRow, error) {
			return sqlc.CreateCalendarRow{}, errors.New("database error")
		},
	}

	handler := newTestHandler(mock)

	body := map[string]any{
		"title": "Test Calendar",
	}

	req := httptest.NewRequest(http.MethodPost, "/api/calendars", makeJSONBody(t, body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.CreateCalendarEndpoint(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d; body: %s", http.StatusInternalServerError, rr.Code, rr.Body.String())
	}
}

func TestCreateCalendarTimeSlotsEndpoint_Success(t *testing.T) {
	calendarUUID := makeUUID()

	mock := &MockCalendarService{
		VerifyCalendarAdminTokenFunc: func(_ context.Context, _ pgtype.UUID, _ string) error {
			return nil
		},
		CreateCalendarTimeSlotsFunc: func(_ context.Context, _ *services.CreateCalendarTimeSlotsInput) error {
			return nil
		},
	}

	handler := newTestHandler(mock)

	startTime := time.Date(2026, 5, 10, 10, 0, 0, 0, time.UTC)
	endTime := time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC)

	body := map[string]any{
		"admin_token": "test-token",
		"time_slots": []map[string]any{
			{"start_date": startTime.Format(time.RFC3339), "end_date": endTime.Format(time.RFC3339)},
		},
	}

	req := httptest.NewRequest(http.MethodPost, "/api/calendars/"+uuid.UUID(calendarUUID.Bytes).String()+"/time-slots", makeJSONBody(t, body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("calendar_id", uuid.UUID(calendarUUID.Bytes).String())
	rr := httptest.NewRecorder()

	handler.CreateCalendarTimeSlotsEndpoint(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d; body: %s", http.StatusCreated, rr.Code, rr.Body.String())
	}
}

func TestCreateCalendarTimeSlotsEndpoint_InvalidCalendarUUID(t *testing.T) {
	mock := &MockCalendarService{}
	handler := newTestHandler(mock)

	body := map[string]any{
		"admin_token": "test-token",
		"time_slots":  []map[string]any{},
	}

	req := httptest.NewRequest(http.MethodPost, "/api/calendars/not-a-uuid/time-slots", makeJSONBody(t, body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("calendar_id", "not-a-uuid")
	rr := httptest.NewRecorder()

	handler.CreateCalendarTimeSlotsEndpoint(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d; body: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
}

func TestCreateCalendarTimeSlotsEndpoint_InvalidAdminToken(t *testing.T) {
	calendarUUID := makeUUID()

	mock := &MockCalendarService{
		VerifyCalendarAdminTokenFunc: func(_ context.Context, _ pgtype.UUID, _ string) error {
			return errors.New("invalid admin token")
		},
	}

	handler := newTestHandler(mock)

	startTime := time.Date(2026, 5, 10, 10, 0, 0, 0, time.UTC)
	endTime := time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC)

	body := map[string]any{
		"admin_token": "wrong-token",
		"time_slots": []map[string]any{
			{"start_date": startTime.Format(time.RFC3339), "end_date": endTime.Format(time.RFC3339)},
		},
	}

	req := httptest.NewRequest(http.MethodPost, "/api/calendars/"+uuid.UUID(calendarUUID.Bytes).String()+"/time-slots", makeJSONBody(t, body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("calendar_id", uuid.UUID(calendarUUID.Bytes).String())
	rr := httptest.NewRecorder()

	handler.CreateCalendarTimeSlotsEndpoint(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d; body: %s", http.StatusUnauthorized, rr.Code, rr.Body.String())
	}
}

func TestCreateCalendarTimeSlotsEndpoint_EndDateBeforeStartDate(t *testing.T) {
	calendarUUID := makeUUID()

	mock := &MockCalendarService{
		VerifyCalendarAdminTokenFunc: func(_ context.Context, _ pgtype.UUID, _ string) error {
			return nil
		},
	}

	handler := newTestHandler(mock)

	startTime := time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC)
	endTime := time.Date(2026, 5, 10, 10, 0, 0, 0, time.UTC)

	body := map[string]any{
		"admin_token": "test-token",
		"time_slots": []map[string]any{
			{"start_date": startTime.Format(time.RFC3339), "end_date": endTime.Format(time.RFC3339)},
		},
	}

	req := httptest.NewRequest(http.MethodPost, "/api/calendars/"+uuid.UUID(calendarUUID.Bytes).String()+"/time-slots", makeJSONBody(t, body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("calendar_id", uuid.UUID(calendarUUID.Bytes).String())
	rr := httptest.NewRecorder()

	handler.CreateCalendarTimeSlotsEndpoint(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d; body: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
}

func TestCreateCalendarTimeSlotsEndpoint_MissingBody(t *testing.T) {
	calendarUUID := makeUUID()

	mock := &MockCalendarService{}
	handler := newTestHandler(mock)

	req := httptest.NewRequest(http.MethodPost, "/api/calendars/"+uuid.UUID(calendarUUID.Bytes).String()+"/time-slots", nil)
	req.SetPathValue("calendar_id", uuid.UUID(calendarUUID.Bytes).String())
	rr := httptest.NewRecorder()

	handler.CreateCalendarTimeSlotsEndpoint(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d; body: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
}

func TestCreateCalendarTimeSlotsEndpoint_ServiceError(t *testing.T) {
	calendarUUID := makeUUID()

	mock := &MockCalendarService{
		VerifyCalendarAdminTokenFunc: func(_ context.Context, _ pgtype.UUID, _ string) error {
			return nil
		},
		CreateCalendarTimeSlotsFunc: func(_ context.Context, _ *services.CreateCalendarTimeSlotsInput) error {
			return errors.New("database error")
		},
	}

	handler := newTestHandler(mock)

	startTime := time.Date(2026, 5, 10, 10, 0, 0, 0, time.UTC)
	endTime := time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC)

	body := map[string]any{
		"admin_token": "test-token",
		"time_slots": []map[string]any{
			{"start_date": startTime.Format(time.RFC3339), "end_date": endTime.Format(time.RFC3339)},
		},
	}

	req := httptest.NewRequest(http.MethodPost, "/api/calendars/"+uuid.UUID(calendarUUID.Bytes).String()+"/time-slots", makeJSONBody(t, body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("calendar_id", uuid.UUID(calendarUUID.Bytes).String())
	rr := httptest.NewRecorder()

	handler.CreateCalendarTimeSlotsEndpoint(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d; body: %s", http.StatusInternalServerError, rr.Code, rr.Body.String())
	}
}

func TestGetCalendarEndpoint_Success_NoPassword(t *testing.T) {
	calendarUUID := makeUUID()
	description := "A test calendar"

	mock := &MockCalendarService{
		IsCalendarPasswordProtectedFunc: func(_ context.Context, _ pgtype.UUID) (bool, error) {
			return false, nil
		},
		GetCalendarFunc: func(_ context.Context, _ pgtype.UUID) (*sqlc.GetCalendarByIDRow, error) {
			return &sqlc.GetCalendarByIDRow{
				ID:          calendarUUID,
				Title:       "Test Calendar",
				Description: &description,
			}, nil
		},
		GetTimeSlotsByCalendarIDFunc: func(_ context.Context, _ pgtype.UUID) ([]sqlc.CalendarTimeSlot, error) {
			return []sqlc.CalendarTimeSlot{}, nil
		},
	}

	handler := newTestHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/api/calendars/"+uuid.UUID(calendarUUID.Bytes).String(), nil)
	req.SetPathValue("calendar_id", uuid.UUID(calendarUUID.Bytes).String())
	rr := httptest.NewRecorder()

	handler.GetCalendarEndpoint(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d; body: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	var resp GetCalendarResponse
	decodeResponse(t, rr, &resp)

	if resp.Calendar.Title != "Test Calendar" {
		t.Errorf("expected title 'Test Calendar', got %s", resp.Calendar.Title)
	}
}

func TestGetCalendarEndpoint_Success_WithPassword(t *testing.T) {
	calendarUUID := makeUUID()

	mock := &MockCalendarService{
		IsCalendarPasswordProtectedFunc: func(_ context.Context, _ pgtype.UUID) (bool, error) {
			return true, nil
		},
		VerifyCalendarPasswordFunc: func(_ context.Context, _ pgtype.UUID, _ string) (bool, error) {
			return true, nil
		},
		GetCalendarFunc: func(_ context.Context, _ pgtype.UUID) (*sqlc.GetCalendarByIDRow, error) {
			return &sqlc.GetCalendarByIDRow{
				ID:    calendarUUID,
				Title: "Protected Calendar",
			}, nil
		},
		GetTimeSlotsByCalendarIDFunc: func(_ context.Context, _ pgtype.UUID) ([]sqlc.CalendarTimeSlot, error) {
			return []sqlc.CalendarTimeSlot{}, nil
		},
	}

	handler := newTestHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/api/calendars/"+uuid.UUID(calendarUUID.Bytes).String()+"?password=secret123", nil)
	req.SetPathValue("calendar_id", uuid.UUID(calendarUUID.Bytes).String())
	rr := httptest.NewRecorder()

	handler.GetCalendarEndpoint(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d; body: %s", http.StatusOK, rr.Code, rr.Body.String())
	}
}

func TestGetCalendarEndpoint_InvalidUUID(t *testing.T) {
	mock := &MockCalendarService{}
	handler := newTestHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/api/calendars/not-a-uuid", nil)
	req.SetPathValue("calendar_id", "not-a-uuid")
	rr := httptest.NewRecorder()

	handler.GetCalendarEndpoint(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d; body: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
}

func TestGetCalendarEndpoint_PasswordRequiredButNotProvided(t *testing.T) {
	calendarUUID := makeUUID()

	mock := &MockCalendarService{
		IsCalendarPasswordProtectedFunc: func(_ context.Context, _ pgtype.UUID) (bool, error) {
			return true, nil
		},
	}

	handler := newTestHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/api/calendars/"+uuid.UUID(calendarUUID.Bytes).String(), nil)
	req.SetPathValue("calendar_id", uuid.UUID(calendarUUID.Bytes).String())
	rr := httptest.NewRecorder()

	handler.GetCalendarEndpoint(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d; body: %s", http.StatusUnauthorized, rr.Code, rr.Body.String())
	}
}

func TestGetCalendarEndpoint_WrongPassword(t *testing.T) {
	calendarUUID := makeUUID()

	mock := &MockCalendarService{
		IsCalendarPasswordProtectedFunc: func(_ context.Context, _ pgtype.UUID) (bool, error) {
			return true, nil
		},
		VerifyCalendarPasswordFunc: func(_ context.Context, _ pgtype.UUID, _ string) (bool, error) {
			return false, nil
		},
	}

	handler := newTestHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/api/calendars/"+uuid.UUID(calendarUUID.Bytes).String()+"?password=wrong", nil)
	req.SetPathValue("calendar_id", uuid.UUID(calendarUUID.Bytes).String())
	rr := httptest.NewRecorder()

	handler.GetCalendarEndpoint(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d; body: %s", http.StatusUnauthorized, rr.Code, rr.Body.String())
	}
}

func TestGetCalendarEndpoint_CalendarRetrievalError(t *testing.T) {
	calendarUUID := makeUUID()

	mock := &MockCalendarService{
		IsCalendarPasswordProtectedFunc: func(_ context.Context, _ pgtype.UUID) (bool, error) {
			return false, nil
		},
		GetCalendarFunc: func(_ context.Context, _ pgtype.UUID) (*sqlc.GetCalendarByIDRow, error) {
			return nil, errors.New("database error")
		},
	}

	handler := newTestHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/api/calendars/"+uuid.UUID(calendarUUID.Bytes).String(), nil)
	req.SetPathValue("calendar_id", uuid.UUID(calendarUUID.Bytes).String())
	rr := httptest.NewRecorder()

	handler.GetCalendarEndpoint(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d; body: %s", http.StatusInternalServerError, rr.Code, rr.Body.String())
	}
}

func TestGetCalendarEndpoint_TimeSlotsRetrievalError(t *testing.T) {
	calendarUUID := makeUUID()

	mock := &MockCalendarService{
		IsCalendarPasswordProtectedFunc: func(_ context.Context, _ pgtype.UUID) (bool, error) {
			return false, nil
		},
		GetCalendarFunc: func(_ context.Context, _ pgtype.UUID) (*sqlc.GetCalendarByIDRow, error) {
			return &sqlc.GetCalendarByIDRow{ID: calendarUUID, Title: "Test"}, nil
		},
		GetTimeSlotsByCalendarIDFunc: func(_ context.Context, _ pgtype.UUID) ([]sqlc.CalendarTimeSlot, error) {
			return nil, errors.New("database error")
		},
	}

	handler := newTestHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/api/calendars/"+uuid.UUID(calendarUUID.Bytes).String(), nil)
	req.SetPathValue("calendar_id", uuid.UUID(calendarUUID.Bytes).String())
	rr := httptest.NewRecorder()

	handler.GetCalendarEndpoint(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d; body: %s", http.StatusInternalServerError, rr.Code, rr.Body.String())
	}
}

func TestCheckIfCalendarPasswordProtectedEndpoint_Protected(t *testing.T) {
	calendarUUID := makeUUID()

	mock := &MockCalendarService{
		IsCalendarPasswordProtectedFunc: func(_ context.Context, _ pgtype.UUID) (bool, error) {
			return true, nil
		},
	}

	handler := newTestHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/api/calendars/"+uuid.UUID(calendarUUID.Bytes).String()+"/password-protected", nil)
	req.SetPathValue("calendar_id", uuid.UUID(calendarUUID.Bytes).String())
	rr := httptest.NewRecorder()

	handler.CheckIfCalendarPasswordProtectedEndpoint(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d; body: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	var resp CheckCalendarPasswordProtectionResponse
	decodeResponse(t, rr, &resp)

	if !resp.IsPasswordProtected {
		t.Error("expected IsPasswordProtected to be true")
	}
}

func TestCheckIfCalendarPasswordProtectedEndpoint_NotProtected(t *testing.T) {
	calendarUUID := makeUUID()

	mock := &MockCalendarService{
		IsCalendarPasswordProtectedFunc: func(_ context.Context, _ pgtype.UUID) (bool, error) {
			return false, nil
		},
	}

	handler := newTestHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/api/calendars/"+uuid.UUID(calendarUUID.Bytes).String()+"/password-protected", nil)
	req.SetPathValue("calendar_id", uuid.UUID(calendarUUID.Bytes).String())
	rr := httptest.NewRecorder()

	handler.CheckIfCalendarPasswordProtectedEndpoint(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d; body: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	var resp CheckCalendarPasswordProtectionResponse
	decodeResponse(t, rr, &resp)

	if resp.IsPasswordProtected {
		t.Error("expected IsPasswordProtected to be false")
	}
}

func TestCheckIfCalendarPasswordProtectedEndpoint_InvalidUUID(t *testing.T) {
	mock := &MockCalendarService{}
	handler := newTestHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/api/calendars/not-a-uuid/password-protected", nil)
	req.SetPathValue("calendar_id", "not-a-uuid")
	rr := httptest.NewRecorder()

	handler.CheckIfCalendarPasswordProtectedEndpoint(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d; body: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
}

func TestCheckIfCalendarPasswordProtectedEndpoint_ServiceError(t *testing.T) {
	calendarUUID := makeUUID()

	mock := &MockCalendarService{
		IsCalendarPasswordProtectedFunc: func(_ context.Context, _ pgtype.UUID) (bool, error) {
			return false, errors.New("database error")
		},
	}

	handler := newTestHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/api/calendars/"+uuid.UUID(calendarUUID.Bytes).String()+"/password-protected", nil)
	req.SetPathValue("calendar_id", uuid.UUID(calendarUUID.Bytes).String())
	rr := httptest.NewRecorder()

	handler.CheckIfCalendarPasswordProtectedEndpoint(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d; body: %s", http.StatusInternalServerError, rr.Code, rr.Body.String())
	}
}

func TestGetCalendarVotesEndpoint_Success(t *testing.T) {
	calendarUUID := makeUUID()
	timeSlotUUID := makeUUID()
	voteUUID := makeUUID()

	mock := &MockCalendarService{
		IsCalendarPasswordProtectedFunc: func(_ context.Context, _ pgtype.UUID) (bool, error) {
			return false, nil
		},
		GetCalendarVotesFunc: func(_ context.Context, _ pgtype.UUID) ([]sqlc.GetVotesByCalendarIDRow, error) {
			return []sqlc.GetVotesByCalendarIDRow{
				{
					ID:                 voteUUID,
					CalendarID:         calendarUUID,
					CalendarTimeSlotID: timeSlotUUID,
					Username:           "alice",
					StartDate:          mustMakeTimestamptz(time.Date(2026, 5, 10, 10, 0, 0, 0, time.UTC)),
					EndDate:            mustMakeTimestamptz(time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC)),
				},
			}, nil
		},
	}

	handler := newTestHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/api/calendars/"+uuid.UUID(calendarUUID.Bytes).String()+"/votes", nil)
	req.SetPathValue("calendar_id", uuid.UUID(calendarUUID.Bytes).String())
	rr := httptest.NewRecorder()

	handler.GetCalendarVotesEndpoint(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d; body: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	var resp []GetCalendarVotesResponse
	decodeResponse(t, rr, &resp)

	if len(resp) != 1 {
		t.Fatalf("expected 1 vote, got %d", len(resp))
	}
	if resp[0].Username != "alice" {
		t.Errorf("expected username 'alice', got %s", resp[0].Username)
	}
}

func TestGetCalendarVotesEndpoint_InvalidUUID(t *testing.T) {
	mock := &MockCalendarService{}
	handler := newTestHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/api/calendars/not-a-uuid/votes", nil)
	req.SetPathValue("calendar_id", "not-a-uuid")
	rr := httptest.NewRecorder()

	handler.GetCalendarVotesEndpoint(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d; body: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
}

func TestGetCalendarVotesEndpoint_PasswordRequired(t *testing.T) {
	calendarUUID := makeUUID()

	mock := &MockCalendarService{
		IsCalendarPasswordProtectedFunc: func(_ context.Context, _ pgtype.UUID) (bool, error) {
			return true, nil
		},
	}

	handler := newTestHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/api/calendars/"+uuid.UUID(calendarUUID.Bytes).String()+"/votes", nil)
	req.SetPathValue("calendar_id", uuid.UUID(calendarUUID.Bytes).String())
	rr := httptest.NewRecorder()

	handler.GetCalendarVotesEndpoint(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d; body: %s", http.StatusUnauthorized, rr.Code, rr.Body.String())
	}
}

func TestGetCalendarVotesEndpoint_WrongPassword(t *testing.T) {
	calendarUUID := makeUUID()

	mock := &MockCalendarService{
		IsCalendarPasswordProtectedFunc: func(_ context.Context, _ pgtype.UUID) (bool, error) {
			return true, nil
		},
		VerifyCalendarPasswordFunc: func(_ context.Context, _ pgtype.UUID, _ string) (bool, error) {
			return false, nil
		},
	}

	handler := newTestHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/api/calendars/"+uuid.UUID(calendarUUID.Bytes).String()+"/votes?password=wrong", nil)
	req.SetPathValue("calendar_id", uuid.UUID(calendarUUID.Bytes).String())
	rr := httptest.NewRecorder()

	handler.GetCalendarVotesEndpoint(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d; body: %s", http.StatusUnauthorized, rr.Code, rr.Body.String())
	}
}

func TestGetCalendarVotesEndpoint_ServiceError(t *testing.T) {
	calendarUUID := makeUUID()

	mock := &MockCalendarService{
		IsCalendarPasswordProtectedFunc: func(_ context.Context, _ pgtype.UUID) (bool, error) {
			return false, nil
		},
		GetCalendarVotesFunc: func(_ context.Context, _ pgtype.UUID) ([]sqlc.GetVotesByCalendarIDRow, error) {
			return nil, errors.New("database error")
		},
	}

	handler := newTestHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/api/calendars/"+uuid.UUID(calendarUUID.Bytes).String()+"/votes", nil)
	req.SetPathValue("calendar_id", uuid.UUID(calendarUUID.Bytes).String())
	rr := httptest.NewRecorder()

	handler.GetCalendarVotesEndpoint(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d; body: %s", http.StatusInternalServerError, rr.Code, rr.Body.String())
	}
}

func TestVoteOnTimeSlotEndpoint_Success(t *testing.T) {
	calendarUUID := makeUUID()
	timeSlotUUID := makeUUID()

	mock := &MockCalendarService{
		IsCalendarPasswordProtectedFunc: func(_ context.Context, _ pgtype.UUID) (bool, error) {
			return false, nil
		},
		GetCalendarTimeSlotByIDFunc: func(_ context.Context, _ pgtype.UUID) (*sqlc.CalendarTimeSlot, error) {
			return &sqlc.CalendarTimeSlot{ID: timeSlotUUID}, nil
		},
		VoteFunc: func(_ context.Context, _ pgtype.UUID, _ pgtype.UUID, _ string) error {
			return nil
		},
	}

	handler := newTestHandler(mock)

	body := map[string]any{
		"username": "alice",
	}

	req := httptest.NewRequest(http.MethodPost, "/api/calendars/"+uuid.UUID(calendarUUID.Bytes).String()+"/time-slots/"+uuid.UUID(timeSlotUUID.Bytes).String()+"/votes", makeJSONBody(t, body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("calendar_id", uuid.UUID(calendarUUID.Bytes).String())
	req.SetPathValue("time_slot_id", uuid.UUID(timeSlotUUID.Bytes).String())
	rr := httptest.NewRecorder()

	handler.VoteOnTimeSlotEndpoint(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d; body: %s", http.StatusCreated, rr.Code, rr.Body.String())
	}
}

func TestVoteOnTimeSlotEndpoint_InvalidCalendarUUID(t *testing.T) {
	timeSlotUUID := makeUUID()

	mock := &MockCalendarService{}
	handler := newTestHandler(mock)

	body := map[string]any{
		"username": "alice",
	}

	req := httptest.NewRequest(http.MethodPost, "/api/calendars/not-a-uuid/time-slots/"+uuid.UUID(timeSlotUUID.Bytes).String()+"/votes", makeJSONBody(t, body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("calendar_id", "not-a-uuid")
	req.SetPathValue("time_slot_id", uuid.UUID(timeSlotUUID.Bytes).String())
	rr := httptest.NewRecorder()

	handler.VoteOnTimeSlotEndpoint(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d; body: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
}

func TestVoteOnTimeSlotEndpoint_InvalidTimeSlotUUID(t *testing.T) {
	calendarUUID := makeUUID()

	mock := &MockCalendarService{
		IsCalendarPasswordProtectedFunc: func(_ context.Context, _ pgtype.UUID) (bool, error) {
			return false, nil
		},
	}

	handler := newTestHandler(mock)

	body := map[string]any{
		"username": "alice",
	}

	req := httptest.NewRequest(http.MethodPost, "/api/calendars/"+uuid.UUID(calendarUUID.Bytes).String()+"/time-slots/not-a-uuid/votes", makeJSONBody(t, body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("calendar_id", uuid.UUID(calendarUUID.Bytes).String())
	req.SetPathValue("time_slot_id", "not-a-uuid")
	rr := httptest.NewRecorder()

	handler.VoteOnTimeSlotEndpoint(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d; body: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
}

func TestVoteOnTimeSlotEndpoint_PasswordRequired(t *testing.T) {
	calendarUUID := makeUUID()
	timeSlotUUID := makeUUID()

	mock := &MockCalendarService{
		IsCalendarPasswordProtectedFunc: func(_ context.Context, _ pgtype.UUID) (bool, error) {
			return true, nil
		},
	}

	handler := newTestHandler(mock)

	body := map[string]any{
		"username": "alice",
	}

	req := httptest.NewRequest(http.MethodPost, "/api/calendars/"+uuid.UUID(calendarUUID.Bytes).String()+"/time-slots/"+uuid.UUID(timeSlotUUID.Bytes).String()+"/votes", makeJSONBody(t, body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("calendar_id", uuid.UUID(calendarUUID.Bytes).String())
	req.SetPathValue("time_slot_id", uuid.UUID(timeSlotUUID.Bytes).String())
	rr := httptest.NewRecorder()

	handler.VoteOnTimeSlotEndpoint(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d; body: %s", http.StatusUnauthorized, rr.Code, rr.Body.String())
	}
}

func TestVoteOnTimeSlotEndpoint_TimeSlotNotFound(t *testing.T) {
	calendarUUID := makeUUID()
	timeSlotUUID := makeUUID()

	mock := &MockCalendarService{
		IsCalendarPasswordProtectedFunc: func(_ context.Context, _ pgtype.UUID) (bool, error) {
			return false, nil
		},
		GetCalendarTimeSlotByIDFunc: func(_ context.Context, _ pgtype.UUID) (*sqlc.CalendarTimeSlot, error) {
			return nil, errors.New("not found")
		},
	}

	handler := newTestHandler(mock)

	body := map[string]any{
		"username": "alice",
	}

	req := httptest.NewRequest(http.MethodPost, "/api/calendars/"+uuid.UUID(calendarUUID.Bytes).String()+"/time-slots/"+uuid.UUID(timeSlotUUID.Bytes).String()+"/votes", makeJSONBody(t, body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("calendar_id", uuid.UUID(calendarUUID.Bytes).String())
	req.SetPathValue("time_slot_id", uuid.UUID(timeSlotUUID.Bytes).String())
	rr := httptest.NewRecorder()

	handler.VoteOnTimeSlotEndpoint(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d; body: %s", http.StatusInternalServerError, rr.Code, rr.Body.String())
	}
}

func TestVoteOnTimeSlotEndpoint_MissingUsername(t *testing.T) {
	calendarUUID := makeUUID()
	timeSlotUUID := makeUUID()

	mock := &MockCalendarService{}
	handler := newTestHandler(mock)

	body := map[string]any{}

	req := httptest.NewRequest(http.MethodPost, "/api/calendars/"+uuid.UUID(calendarUUID.Bytes).String()+"/time-slots/"+uuid.UUID(timeSlotUUID.Bytes).String()+"/votes", makeJSONBody(t, body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("calendar_id", uuid.UUID(calendarUUID.Bytes).String())
	req.SetPathValue("time_slot_id", uuid.UUID(timeSlotUUID.Bytes).String())
	rr := httptest.NewRecorder()

	handler.VoteOnTimeSlotEndpoint(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d; body: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
	}
}

func TestVoteOnTimeSlotEndpoint_ServiceError(t *testing.T) {
	calendarUUID := makeUUID()
	timeSlotUUID := makeUUID()

	mock := &MockCalendarService{
		IsCalendarPasswordProtectedFunc: func(_ context.Context, _ pgtype.UUID) (bool, error) {
			return false, nil
		},
		GetCalendarTimeSlotByIDFunc: func(_ context.Context, _ pgtype.UUID) (*sqlc.CalendarTimeSlot, error) {
			return &sqlc.CalendarTimeSlot{ID: timeSlotUUID}, nil
		},
		VoteFunc: func(_ context.Context, _ pgtype.UUID, _ pgtype.UUID, _ string) error {
			return errors.New("database error")
		},
	}

	handler := newTestHandler(mock)

	body := map[string]any{
		"username": "alice",
	}

	req := httptest.NewRequest(http.MethodPost, "/api/calendars/"+uuid.UUID(calendarUUID.Bytes).String()+"/time-slots/"+uuid.UUID(timeSlotUUID.Bytes).String()+"/votes", makeJSONBody(t, body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("calendar_id", uuid.UUID(calendarUUID.Bytes).String())
	req.SetPathValue("time_slot_id", uuid.UUID(timeSlotUUID.Bytes).String())
	rr := httptest.NewRecorder()

	handler.VoteOnTimeSlotEndpoint(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d; body: %s", http.StatusInternalServerError, rr.Code, rr.Body.String())
	}
}
