package handlers

import (
	"meeting-planner/backend/internal/services"
	"meeting-planner/backend/internal/utils"
	"net/http"
	"time"
)

type CreateCalendarRequest struct {
	Title                string  `json:"title" validate:"required,min=3,max=256"`
	Description          *string `json:"description,omitempty" validate:"omitempty,max=1024"`
	Location             *string `json:"location,omitempty" validate:"omitempty,max=512"`
	AcceptResponsesUntil *string `json:"accept_responses_until,omitempty" validate:"omitempty,rfc3339"`
	Password             *string `json:"password,omitempty" validate:"omitempty,min=3,max=128"`
}

type CreateCalendarResponse struct {
	ID         string `json:"id"`
	AdminToken string `json:"admin_token"`
}

func (h *Handler) CreateCalendarEndpoint(w http.ResponseWriter, r *http.Request) {
	var requestBody CreateCalendarRequest

	if parsingError := ParseRequest(r, RequestOptions{Body: &requestBody}); parsingError != nil {
		RespondError(w, http.StatusBadRequest, parsingError.Error(), nil)
		return
	}

	serviceInput := services.CreateCalendarInput{
		Title:                requestBody.Title,
		Description:          requestBody.Description,
		Location:             requestBody.Location,
		AcceptResponsesUntil: nil,
		Password:             requestBody.Password,
	}

	if requestBody.AcceptResponsesUntil != nil {
		parsedTime, timeParsingError := time.Parse(time.RFC3339, *requestBody.AcceptResponsesUntil)
		if timeParsingError != nil {
			RespondError(w, http.StatusBadRequest, "Invalid time format for accept_responses_until, expected RFC3339", &timeParsingError)
			return
		}
		serviceInput.AcceptResponsesUntil = &parsedTime
	}

	calendarRow, creationError := h.CalendarService.CreateCalendar(r.Context(), serviceInput)
	if creationError != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to create calendar", &creationError)
		return
	}

	RespondJSON(w, http.StatusCreated, CreateCalendarResponse{
		ID:         utils.UUIDToString(calendarRow.ID),
		AdminToken: calendarRow.AdminToken,
	})
}

type CalendarTimeSlotsRequest struct {
	StartDate string `json:"start_date" validate:"required,rfc3339"`
	EndDate   string `json:"end_date" validate:"required,rfc3339"`
}

type CreateCalendarTimeSlotsRequest struct {
	AdminToken string                     `json:"admin_token" validate:"required"`
	TimeSlots  []CalendarTimeSlotsRequest `json:"time_slots" validate:"required,dive,required"`
}

func (h *Handler) CreateCalendarTimeSlotsEndpoint(w http.ResponseWriter, r *http.Request) {
	calendarID := r.PathValue("calendar_id")

	var requestBody CreateCalendarTimeSlotsRequest

	if parsingError := ParseRequest(r, RequestOptions{Body: &requestBody}); parsingError != nil {
		RespondError(w, http.StatusBadRequest, parsingError.Error(), &parsingError)
		return
	}

	calendarUUID, uuidError := utils.StringToUUID(calendarID)
	if uuidError != nil {
		RespondError(w, http.StatusBadRequest, "Invalid calendar ID", &uuidError)
		return
	}

	adminTokenErr := h.CalendarService.VerifyCalendarAdminToken(r.Context(), calendarUUID, requestBody.AdminToken)
	if adminTokenErr != nil {
		RespondError(w, http.StatusUnauthorized, "Invalid admin token", &adminTokenErr)
		return
	}

	var timeSlots []services.CreateTimeSlotInput
	for _, slot := range requestBody.TimeSlots {
		startTime, _ := time.Parse(time.RFC3339, slot.StartDate)
		endTime, _ := time.Parse(time.RFC3339, slot.EndDate)

		if !endTime.After(startTime) {
			RespondError(w, http.StatusBadRequest, "end_date must be after start_date", nil)
			return
		}

		timeSlots = append(timeSlots, services.CreateTimeSlotInput{
			StartDate: startTime,
			EndDate:   endTime,
		})
	}

	serviceInput := services.CreateCalendarTimeSlotsInput{
		CalendarID: calendarUUID,
		TimeSlots:  timeSlots,
	}

	creationError := h.CalendarService.CreateCalendarTimeSlots(r.Context(), serviceInput)
	if creationError != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to create calendar time slots", &creationError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type CalendarResponse struct {
	ID                   string  `json:"id" validate:"required"`
	Title                string  `json:"title"`
	Description          *string `json:"description,omitempty"`
	Location             *string `json:"location,omitempty"`
	AcceptResponsesUntil *string `json:"accept_responses_until,omitempty"`
}

type TimeSlotResponse struct {
	ID        string `json:"id" validate:"required"`
	StartDate string `json:"start_date" validate:"required,rfc3339"`
	EndDate   string `json:"end_date" validate:"required,rfc3339"`
}

type GetCalendarResponse struct {
	Calendar  CalendarResponse   `json:"calendar" validate:"required"`
	TimeSlots []TimeSlotResponse `json:"time_slots" validate:"required,dive,required"`
}

func (h *Handler) GetCalendarEndpoint(w http.ResponseWriter, r *http.Request) {
	calendarID := r.PathValue("calendar_id")

	calendarUUID, uuidError := utils.StringToUUID(calendarID)
	if uuidError != nil {
		RespondError(w, http.StatusBadRequest, "Invalid calendar ID", &uuidError)
		return
	}

	calendar, retrievalError := h.CalendarService.GetCalendar(r.Context(), calendarUUID)
	if retrievalError != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to retrieve calendar", &retrievalError)
		return
	}

	timeSlots, timeSlotsError := h.CalendarService.GetTimeSlotsByCalendarID(r.Context(), calendarUUID)
	if timeSlotsError != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to retrieve time slots", &timeSlotsError)
		return
	}

	var timeSlotResponses []TimeSlotResponse
	for _, slot := range timeSlots {
		startTime := slot.StartDate.Time

		endTime := slot.EndDate.Time

		if !endTime.After(startTime) {
			RespondError(w, http.StatusBadRequest, "end_date must be after start_date", nil)
			return
		}

		timeSlotResponses = append(timeSlotResponses, TimeSlotResponse{
			ID:        slot.ID.String(),
			StartDate: startTime.Format(time.RFC3339),
			EndDate:   endTime.Format(time.RFC3339),
		})
	}

	var response GetCalendarResponse = GetCalendarResponse{
		Calendar: CalendarResponse{
			ID:          utils.UUIDToString(calendar.ID),
			Title:       calendar.Title,
			Description: calendar.Description,
			Location:    calendar.Location,
		},
		TimeSlots: timeSlotResponses,
	}

	RespondJSON(w, http.StatusOK, response)
}
