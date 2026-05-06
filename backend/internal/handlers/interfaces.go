package handlers

import (
	"context"

	"meeting-planner/backend/internal/db/sqlc"
	"meeting-planner/backend/internal/services"

	"github.com/jackc/pgx/v5/pgtype"
)

type CalendarServicer interface {
	CreateCalendar(ctx context.Context, input *services.CreateCalendarInput) (sqlc.CreateCalendarRow, error)
	VerifyCalendarAdminToken(ctx context.Context, calendarID pgtype.UUID, adminToken string) error
	CreateCalendarTimeSlots(ctx context.Context, input *services.CreateCalendarTimeSlotsInput) error
	GetCalendar(ctx context.Context, calendarID pgtype.UUID) (*sqlc.GetCalendarByIDRow, error)
	GetTimeSlotsByCalendarID(ctx context.Context, calendarID pgtype.UUID) ([]sqlc.CalendarTimeSlot, error)
	IsCalendarPasswordProtected(ctx context.Context, calendarID pgtype.UUID) (bool, error)
	VerifyCalendarPassword(ctx context.Context, calendarID pgtype.UUID, password string) (bool, error)
	GetCalendarTimeSlotByID(ctx context.Context, timeSlotID pgtype.UUID) (*sqlc.CalendarTimeSlot, error)
	GetCalendarVotes(ctx context.Context, calendarID pgtype.UUID) ([]sqlc.GetVotesByCalendarIDRow, error)
	Vote(ctx context.Context, calendarID pgtype.UUID, timeSlotID pgtype.UUID, username string) error
}
