package services

import (
	"context"
	"fmt"
	"meeting-planner/backend/internal/db/sqlc"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type CalendarService struct {
	queries         *sqlc.Queries
	passwordManager *PasswordManager
}

func NewCalendarService(queries *sqlc.Queries, passwordManager *PasswordManager) *CalendarService {
	return &CalendarService{
		queries:         queries,
		passwordManager: passwordManager,
	}
}

type CreateCalendarInput struct {
	Title                string
	Description          *string
	Location             *string
	AcceptResponsesUntil *time.Time
	Password             *string
}

func (s *CalendarService) CreateCalendar(ctx context.Context, input CreateCalendarInput) (sqlc.CreateCalendarRow, error) {
	var calendarPassword *string
	var calendarSalt *string

	if input.Password != nil {
		var salt, err = s.passwordManager.GenerateSalt(16)
		if err != nil {
			return sqlc.CreateCalendarRow{}, fmt.Errorf("failed to generate salt: %w", err)
		}

		var passwordHash, hashErr = s.passwordManager.HashPassword(*input.Password, salt)
		if hashErr != nil {
			return sqlc.CreateCalendarRow{}, fmt.Errorf("failed to hash password: %w", hashErr)
		}

		calendarPassword = &passwordHash
		calendarSalt = &salt
	}

	var randomAdminToken, tokenErr = s.passwordManager.GenerateSalt(32)
	if tokenErr != nil {
		return sqlc.CreateCalendarRow{}, fmt.Errorf("failed to generate admin token: %w", tokenErr)
	}

	queryParams := sqlc.CreateCalendarParams{
		Title:       input.Title,
		Description: input.Description,
		Location:    input.Location,
		Password:    calendarPassword,
		Salt:        calendarSalt,
		AdminToken:  randomAdminToken,
	}

	if input.AcceptResponsesUntil != nil {
		queryParams.AcceptResponsesUntil = pgtype.Timestamptz{
			Time:  *input.AcceptResponsesUntil,
			Valid: true,
		}
	}

	calendarRow, creationError := s.queries.CreateCalendar(ctx, queryParams)
	if creationError != nil {
		return sqlc.CreateCalendarRow{}, fmt.Errorf("failed to create calendar: %w", creationError)
	}

	return calendarRow, nil
}

type CreateTimeSlotInput struct {
	StartDate time.Time
	EndDate   time.Time
}

type CreateCalendarTimeSlotsInput struct {
	CalendarID pgtype.UUID
	TimeSlots  []CreateTimeSlotInput
}

func (s *CalendarService) CreateCalendarTimeSlots(ctx context.Context, input CreateCalendarTimeSlotsInput) error {
	for _, slot := range input.TimeSlots {
		queryParams := sqlc.CreateCalendarTimeSlotParams{
			CalendarID: input.CalendarID,
			StartDate: pgtype.Timestamptz{
				Time:  slot.StartDate,
				Valid: true,
			},
			EndDate: pgtype.Timestamptz{
				Time:  slot.EndDate,
				Valid: true,
			},
		}

		_, creationError := s.queries.CreateCalendarTimeSlot(ctx, queryParams)
		if creationError != nil {
			return fmt.Errorf("failed to create calendar time slot: %w", creationError)
		}
	}

	return nil
}

func (s *CalendarService) GetCalendar(ctx context.Context, calendarID pgtype.UUID) (*sqlc.GetCalendarByIDRow, error) {
	calendar, retrievalError := s.queries.GetCalendarByID(ctx, calendarID)
	if retrievalError != nil {
		return nil, fmt.Errorf("failed to retrieve calendar: %w", retrievalError)
	}

	return &calendar, nil
}

func (s *CalendarService) VerifyCalendarAdminToken(ctx context.Context, calendarID pgtype.UUID, adminToken string) error {
	_, retrievalError := s.queries.GetCalendarByIDAndAdminToken(ctx, sqlc.GetCalendarByIDAndAdminTokenParams{
		ID:         calendarID,
		AdminToken: adminToken,
	})
	if retrievalError != nil {
		return fmt.Errorf("failed to verify calendar admin token: %w", retrievalError)
	}

	return nil
}

func (s *CalendarService) GetTimeSlotsByCalendarID(ctx context.Context, calendarID pgtype.UUID) ([]sqlc.CalendarTimeSlot, error) {
	timeSlots, retrievalError := s.queries.GetCalendarTimeSlotsByCalendarID(ctx, calendarID)
	if retrievalError != nil {
		return nil, fmt.Errorf("failed to retrieve time slots: %w", retrievalError)
	}
	return timeSlots, nil
}
