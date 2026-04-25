package handlers

import (
	"fmt"
	"meeting-planner/backend/internal/services"
	"meeting-planner/backend/internal/utils"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type CreateCalendarRequest struct {
	Title       string  `json:"title" validate:"required,min=3,max=256"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1024"`
	Password    *string `json:"password,omitempty" validate:"omitempty,min=3,max=128"`
}

type CreateCalendarResponse struct {
	ID         string `json:"id"`
	AdminToken string `json:"admin_token"`
}

func (h *Handler) CreateCalendarEndpoint(w http.ResponseWriter, r *http.Request) {
	var requestBody CreateCalendarRequest

	if parsingError := ParseRequest(w, r, RequestOptions{Body: &requestBody}); parsingError != nil {
		RespondError(w, http.StatusBadRequest, parsingError.Error(), nil)
		return
	}

	serviceInput := services.CreateCalendarInput{
		Title:       requestBody.Title,
		Description: requestBody.Description,
		Password:    requestBody.Password,
	}

	calendarRow, creationError := h.CalendarService.CreateCalendar(r.Context(), &serviceInput)
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
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" validate:"required"`
}

type CreateCalendarTimeSlotsRequest struct {
	AdminToken string                     `json:"admin_token" validate:"required"`
	TimeSlots  []CalendarTimeSlotsRequest `json:"time_slots" validate:"required,dive,required"`
}

func (h *Handler) CreateCalendarTimeSlotsEndpoint(w http.ResponseWriter, r *http.Request) {
	calendarID := r.PathValue("calendar_id")

	var requestBody CreateCalendarTimeSlotsRequest

	if parsingError := ParseRequest(w, r, RequestOptions{Body: &requestBody}); parsingError != nil {
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

	timeSlots := make([]services.CreateTimeSlotInput, 0, len(requestBody.TimeSlots))
	for _, slot := range requestBody.TimeSlots {
		if !slot.EndDate.After(slot.StartDate) {
			RespondError(w, http.StatusBadRequest, "end_date must be after start_date", nil)
			return
		}

		timeSlots = append(timeSlots, services.CreateTimeSlotInput{
			StartDate: slot.StartDate,
			EndDate:   slot.EndDate,
		})
	}

	serviceInput := services.CreateCalendarTimeSlotsInput{
		CalendarID: calendarUUID,
		TimeSlots:  timeSlots,
	}

	creationError := h.CalendarService.CreateCalendarTimeSlots(r.Context(), &serviceInput)
	if creationError != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to create calendar time slots", &creationError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type CalendarResponse struct {
	ID          string  `json:"id" validate:"required"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
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

func (h *Handler) authorizeCalendarAccess(w http.ResponseWriter, r *http.Request, calendarID string, password string) (pgtype.UUID, bool) {
	calendarUUID, uuidError := utils.StringToUUID(calendarID)
	if uuidError != nil {
		RespondError(w, http.StatusBadRequest, "Invalid calendar ID", &uuidError)
		return pgtype.UUID{}, false
	}

	isProtected, checkError := h.CalendarService.IsCalendarPasswordProtected(r.Context(), calendarUUID)
	if checkError != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to check calendar password protection", &checkError)
		return pgtype.UUID{}, false
	}

	if isProtected {
		if password == "" {
			RespondError(w, http.StatusUnauthorized, "Password is required to access this calendar", nil)
			return pgtype.UUID{}, false
		}

		isValid, passwordErr := h.CalendarService.VerifyCalendarPassword(r.Context(), calendarUUID, password)
		if passwordErr != nil {
			RespondError(w, http.StatusInternalServerError, "Failed to verify calendar password", &passwordErr)
			return pgtype.UUID{}, false
		}

		if !isValid {
			RespondError(w, http.StatusUnauthorized, "Invalid password for this calendar", nil)
			return pgtype.UUID{}, false
		}
	}

	return calendarUUID, true
}

func (h *Handler) GetCalendarEndpoint(w http.ResponseWriter, r *http.Request) {
	calendarID := r.PathValue("calendar_id")
	calendarPassword := r.URL.Query().Get("password")

	calendarUUID, ok := h.authorizeCalendarAccess(w, r, calendarID, calendarPassword)
	if !ok {
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

	timeSlotResponses := make([]TimeSlotResponse, 0, len(timeSlots))
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

	response := GetCalendarResponse{
		Calendar: CalendarResponse{
			ID:          utils.UUIDToString(calendar.ID),
			Title:       calendar.Title,
			Description: calendar.Description,
		},
		TimeSlots: timeSlotResponses,
	}

	RespondJSON(w, http.StatusOK, response)
}

type CheckCalendarPasswordProtectionResponse struct {
	IsPasswordProtected bool `json:"is_password_protected"`
}

func (h *Handler) CheckIfCalendarPasswordProtectedEndpoint(w http.ResponseWriter, r *http.Request) {
	calendarID := r.PathValue("calendar_id")

	calendarUUID, uuidError := utils.StringToUUID(calendarID)
	if uuidError != nil {
		RespondError(w, http.StatusBadRequest, "Invalid calendar ID", &uuidError)
		return
	}

	isProtected, checkError := h.CalendarService.IsCalendarPasswordProtected(r.Context(), calendarUUID)
	if checkError != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to check calendar password protection", &checkError)
		return
	}

	RespondJSON(w, http.StatusOK, CheckCalendarPasswordProtectionResponse{
		IsPasswordProtected: isProtected,
	})
}

type GetCalendarVotesResponse struct {
	ID         string           `json:"id"`
	CalendarID string           `json:"calendar_id"`
	Username   string           `json:"username"`
	TimeSlot   TimeSlotResponse `json:"time_slot"`
}

func (h *Handler) GetCalendarVotesEndpoint(w http.ResponseWriter, r *http.Request) {
	calendarID := r.PathValue("calendar_id")
	password := r.URL.Query().Get("password")

	calendarUUID, ok := h.authorizeCalendarAccess(w, r, calendarID, password)
	if !ok {
		return
	}

	votes, retrievalError := h.CalendarService.GetCalendarVotes(r.Context(), calendarUUID)
	if retrievalError != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to retrieve votes for calendar", &retrievalError)
		return
	}

	responseVotes := make([]GetCalendarVotesResponse, 0, len(votes))
	for _, vote := range votes {
		responseVotes = append(responseVotes, GetCalendarVotesResponse{
			ID:         vote.ID.String(),
			CalendarID: vote.CalendarID.String(),
			Username:   vote.Username,
			TimeSlot: TimeSlotResponse{
				ID:        vote.CalendarTimeSlotID.String(),
				StartDate: vote.StartDate.Time.Format(time.RFC3339),
				EndDate:   vote.EndDate.Time.Format(time.RFC3339),
			},
		})
	}

	RespondJSON(w, http.StatusOK, responseVotes)
}

type VoteOnTimeSlotBodyRequest struct {
	Username string  `json:"username" validate:"required,min=1,max=256"`
	Password *string `json:"password,omitempty" validate:"omitempty,min=3,max=128"`
}

func (h *Handler) VoteOnTimeSlotEndpoint(w http.ResponseWriter, r *http.Request) {
	var requestBody VoteOnTimeSlotBodyRequest
	calendarID := r.PathValue("calendar_id")
	timeSlotID := r.PathValue("time_slot_id")

	if parsingError := ParseRequest(w, r, RequestOptions{Body: &requestBody}); parsingError != nil {
		RespondError(w, http.StatusBadRequest, parsingError.Error(), &parsingError)
		return
	}

	passwordStr := ""
	if requestBody.Password != nil {
		passwordStr = *requestBody.Password
	}

	calendarUUID, ok := h.authorizeCalendarAccess(w, r, calendarID, passwordStr)
	if !ok {
		return
	}

	timeSlotUUID, timeSlotUUIDErr := utils.StringToUUID(timeSlotID)
	if timeSlotUUIDErr != nil {
		RespondError(w, http.StatusBadRequest, "Invalid time slot ID", &timeSlotUUIDErr)
		return
	}

	timeSlot, retrievalError := h.CalendarService.GetCalendarTimeSlotByID(r.Context(), timeSlotUUID)
	if retrievalError != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to retrieve calendar time slot", &retrievalError)
		return
	}

	if timeSlot == nil {
		RespondError(w, http.StatusNotFound, "Time slot not found", nil)
		return
	}

	voteError := h.CalendarService.Vote(r.Context(), calendarUUID, timeSlotUUID, requestBody.Username)
	fmt.Println("Vote error:", voteError)
	if voteError != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to vote on time slot", &voteError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
