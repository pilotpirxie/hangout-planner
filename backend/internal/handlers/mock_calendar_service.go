package handlers

import (
	"context"
	"fmt"

	"meeting-planner/backend/internal/db/sqlc"
	"meeting-planner/backend/internal/services"

	"github.com/jackc/pgx/v5/pgtype"
)

var _ CalendarServicer = (*MockCalendarService)(nil)

type Call struct {
	Method string
	Args   []any
}

type MockCalendarService struct {
	Calls []Call

	CreateCalendarFunc              func(ctx context.Context, input *services.CreateCalendarInput) (sqlc.CreateCalendarRow, error)
	VerifyCalendarAdminTokenFunc    func(ctx context.Context, calendarID pgtype.UUID, adminToken string) error
	CreateCalendarTimeSlotsFunc     func(ctx context.Context, input *services.CreateCalendarTimeSlotsInput) error
	GetCalendarFunc                 func(ctx context.Context, calendarID pgtype.UUID) (*sqlc.GetCalendarByIDRow, error)
	GetTimeSlotsByCalendarIDFunc    func(ctx context.Context, calendarID pgtype.UUID) ([]sqlc.CalendarTimeSlot, error)
	IsCalendarPasswordProtectedFunc func(ctx context.Context, calendarID pgtype.UUID) (bool, error)
	VerifyCalendarPasswordFunc      func(ctx context.Context, calendarID pgtype.UUID, password string) (bool, error)
	GetCalendarTimeSlotByIDFunc     func(ctx context.Context, timeSlotID pgtype.UUID) (*sqlc.CalendarTimeSlot, error)
	GetCalendarVotesFunc            func(ctx context.Context, calendarID pgtype.UUID) ([]sqlc.GetVotesByCalendarIDRow, error)
	VoteFunc                        func(ctx context.Context, calendarID pgtype.UUID, timeSlotID pgtype.UUID, username string) error
}

func (m *MockCalendarService) recordCall(method string, args ...any) {
	m.Calls = append(m.Calls, Call{Method: method, Args: args})
}

func (m *MockCalendarService) CreateCalendar(ctx context.Context, input *services.CreateCalendarInput) (sqlc.CreateCalendarRow, error) {
	m.recordCall("CreateCalendar", input)
	if m.CreateCalendarFunc == nil {
		panic(fmt.Sprintf("MockCalendarService.CreateCalendarFunc is nil"))
	}
	return m.CreateCalendarFunc(ctx, input)
}

func (m *MockCalendarService) VerifyCalendarAdminToken(ctx context.Context, calendarID pgtype.UUID, adminToken string) error {
	m.recordCall("VerifyCalendarAdminToken", calendarID, adminToken)
	if m.VerifyCalendarAdminTokenFunc == nil {
		panic(fmt.Sprintf("MockCalendarService.VerifyCalendarAdminTokenFunc is nil"))
	}
	return m.VerifyCalendarAdminTokenFunc(ctx, calendarID, adminToken)
}

func (m *MockCalendarService) CreateCalendarTimeSlots(ctx context.Context, input *services.CreateCalendarTimeSlotsInput) error {
	m.recordCall("CreateCalendarTimeSlots", input)
	if m.CreateCalendarTimeSlotsFunc == nil {
		panic(fmt.Sprintf("MockCalendarService.CreateCalendarTimeSlotsFunc is nil"))
	}
	return m.CreateCalendarTimeSlotsFunc(ctx, input)
}

func (m *MockCalendarService) GetCalendar(ctx context.Context, calendarID pgtype.UUID) (*sqlc.GetCalendarByIDRow, error) {
	m.recordCall("GetCalendar", calendarID)
	if m.GetCalendarFunc == nil {
		panic(fmt.Sprintf("MockCalendarService.GetCalendarFunc is nil"))
	}
	return m.GetCalendarFunc(ctx, calendarID)
}

func (m *MockCalendarService) GetTimeSlotsByCalendarID(ctx context.Context, calendarID pgtype.UUID) ([]sqlc.CalendarTimeSlot, error) {
	m.recordCall("GetTimeSlotsByCalendarID", calendarID)
	if m.GetTimeSlotsByCalendarIDFunc == nil {
		panic(fmt.Sprintf("MockCalendarService.GetTimeSlotsByCalendarIDFunc is nil"))
	}
	return m.GetTimeSlotsByCalendarIDFunc(ctx, calendarID)
}

func (m *MockCalendarService) IsCalendarPasswordProtected(ctx context.Context, calendarID pgtype.UUID) (bool, error) {
	m.recordCall("IsCalendarPasswordProtected", calendarID)
	if m.IsCalendarPasswordProtectedFunc == nil {
		panic(fmt.Sprintf("MockCalendarService.IsCalendarPasswordProtectedFunc is nil"))
	}
	return m.IsCalendarPasswordProtectedFunc(ctx, calendarID)
}

func (m *MockCalendarService) VerifyCalendarPassword(ctx context.Context, calendarID pgtype.UUID, password string) (bool, error) {
	m.recordCall("VerifyCalendarPassword", calendarID, password)
	if m.VerifyCalendarPasswordFunc == nil {
		panic(fmt.Sprintf("MockCalendarService.VerifyCalendarPasswordFunc is nil"))
	}
	return m.VerifyCalendarPasswordFunc(ctx, calendarID, password)
}

func (m *MockCalendarService) GetCalendarTimeSlotByID(ctx context.Context, timeSlotID pgtype.UUID) (*sqlc.CalendarTimeSlot, error) {
	m.recordCall("GetCalendarTimeSlotByID", timeSlotID)
	if m.GetCalendarTimeSlotByIDFunc == nil {
		panic(fmt.Sprintf("MockCalendarService.GetCalendarTimeSlotByIDFunc is nil"))
	}
	return m.GetCalendarTimeSlotByIDFunc(ctx, timeSlotID)
}

func (m *MockCalendarService) GetCalendarVotes(ctx context.Context, calendarID pgtype.UUID) ([]sqlc.GetVotesByCalendarIDRow, error) {
	m.recordCall("GetCalendarVotes", calendarID)
	if m.GetCalendarVotesFunc == nil {
		panic(fmt.Sprintf("MockCalendarService.GetCalendarVotesFunc is nil"))
	}
	return m.GetCalendarVotesFunc(ctx, calendarID)
}

func (m *MockCalendarService) Vote(ctx context.Context, calendarID pgtype.UUID, timeSlotID pgtype.UUID, username string) error {
	m.recordCall("Vote", calendarID, timeSlotID, username)
	if m.VoteFunc == nil {
		panic(fmt.Sprintf("MockCalendarService.VoteFunc is nil"))
	}
	return m.VoteFunc(ctx, calendarID, timeSlotID, username)
}
