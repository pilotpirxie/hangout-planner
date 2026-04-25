package handlers

import (
	"meeting-planner/backend/internal/db"
	"meeting-planner/backend/internal/services"
)

type Handler struct {
	DB              *db.DB
	CalendarService *services.CalendarService
}

func New(database *db.DB, passwordManager *services.PasswordManager) *Handler {
	return &Handler{
		DB:              database,
		CalendarService: services.NewCalendarService(database, passwordManager),
	}
}
