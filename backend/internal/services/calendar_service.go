package services

import (
	"context"
	"crypto/subtle"
	"fmt"
	"meeting-planner/backend/internal/db"
	"meeting-planner/backend/internal/db/sqlc"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type CalendarService struct {
	db              *db.DB
	queries         *sqlc.Queries
	passwordManager *PasswordManager
}

func NewCalendarService(database *db.DB, passwordManager *PasswordManager) *CalendarService {
	return &CalendarService{
		db:              database,
		queries:         database.Queries,
		passwordManager: passwordManager,
	}
}

type CreateCalendarInput struct {
	Title                string
	Description          *string
	AcceptResponsesUntil *time.Time
	Password             *string
}

func (s *CalendarService) CreateCalendar(ctx context.Context, input *CreateCalendarInput) (sqlc.CreateCalendarRow, error) {
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
		Password:    calendarPassword,
		Salt:        calendarSalt,
		AdminToken:  randomAdminToken,
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

func (s *CalendarService) CreateCalendarTimeSlots(ctx context.Context, input *CreateCalendarTimeSlotsInput) error {
	tx, err := s.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

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

		_, creationError := qtx.CreateCalendarTimeSlot(ctx, queryParams)
		if creationError != nil {
			return fmt.Errorf("failed to create calendar time slot: %w", creationError)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
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

func (s *CalendarService) IsCalendarPasswordProtected(ctx context.Context, calendarID pgtype.UUID) (bool, error) {
	password, retrievalError := s.queries.GetCalendarPasswordAndSaltByID(ctx, calendarID)
	if retrievalError != nil {
		return false, fmt.Errorf("failed to check if calendar is password protected: %w", retrievalError)
	}

	return password.Password != nil && *password.Password != "", nil
}

func (s *CalendarService) VerifyCalendarPassword(ctx context.Context, calendarID pgtype.UUID, password string) (bool, error) {
	storedCredentials, retrievalError := s.queries.GetCalendarPasswordAndSaltByID(ctx, calendarID)
	if retrievalError != nil {
		return false, fmt.Errorf("failed to retrieve calendar password: %w", retrievalError)
	}

	if storedCredentials.Password == nil || *storedCredentials.Password == "" {
		return false, fmt.Errorf("calendar is not password protected")
	}

	hashedPassword, hashErr := s.passwordManager.HashPassword(password, *storedCredentials.Salt)
	if hashErr != nil {
		return false, fmt.Errorf("failed to hash provided password: %w", hashErr)
	}

	isMatch := storedCredentials.Password != nil && subtle.ConstantTimeCompare([]byte(*storedCredentials.Password), []byte(hashedPassword)) == 1
	return isMatch, nil
}

func (s *CalendarService) GetCalendarTimeSlotByID(ctx context.Context, timeSlotID pgtype.UUID) (*sqlc.CalendarTimeSlot, error) {
	timeSlot, retrievalError := s.queries.GetCalendarTimeSlotByID(ctx, timeSlotID)
	if retrievalError != nil {
		return nil, fmt.Errorf("failed to retrieve calendar time slot: %w", retrievalError)
	}

	return &timeSlot, nil
}

func (s *CalendarService) GetCalendarVotes(ctx context.Context, calendarID pgtype.UUID) ([]sqlc.GetVotesByCalendarIDRow, error) {
	votes, retrievalError := s.queries.GetVotesByCalendarID(ctx, calendarID)
	if retrievalError != nil {
		return nil, fmt.Errorf("failed to retrieve calendar votes: %w", retrievalError)
	}

	return votes, nil
}

func (s *CalendarService) Vote(ctx context.Context, calendarID pgtype.UUID, timeSlotID pgtype.UUID, username string) error {
	_, creationError := s.queries.CreateVote(ctx, sqlc.CreateVoteParams{
		CalendarID:         calendarID,
		Username:           username,
		CalendarTimeSlotID: timeSlotID,
	})

	if creationError != nil {
		return fmt.Errorf("failed to create calendar vote: %w", creationError)
	}

	return nil
}
