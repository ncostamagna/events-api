package events

import (
	"context"
	"log/slog"
	"time"

	"github.com/ncostamagna/events-api/adapters/db"
	"github.com/ncostamagna/events-api/domain"
)

type (
	Service interface {
		CreateEvent(ctx context.Context, title string, description string, startTime time.Time, endTime time.Time) (*domain.Event, error)
		GetEventByID(ctx context.Context, id string) (*domain.Event, error)
		GetAllEvents(ctx context.Context) ([]domain.Event, error)
	}
	service struct {
		db     *db.DB
		logger *slog.Logger
	}
)

func NewService(db *db.DB, logger *slog.Logger) Service {
	return &service{db: db, logger: logger}
}

func (s *service) CreateEvent(ctx context.Context, title string, description string, startTime time.Time, endTime time.Time) (*domain.Event, error) {
	s.logger.Debug("Creating event", "title", title, "description", description, "startTime", startTime, "endTime", endTime)
	event := &domain.Event{
		Title:       title,
		Description: description,
		StartTime:   startTime,
		EndTime:     endTime,
		CreatedAt:   time.Now(),
	}
	if err := s.db.CreateEvent(ctx, event); err != nil {
		s.logger.Error("Error creating event", "error", err)
		return nil, err
	}
	return event, nil
}

func (s *service) GetEventByID(ctx context.Context, id string) (*domain.Event, error) {
	s.logger.Debug("Getting event by ID", "id", id)
	event, err := s.db.GetEventByID(ctx, id)
	if err != nil {
		s.logger.Error("Error event id", "id", id, "error", err)
		return nil, err
	}
	if event == nil {
		s.logger.Debug("Event not found", "id", id)
		return nil, ErrEventNotFound
	}
	return event, nil
}

func (s *service) GetAllEvents(ctx context.Context) ([]domain.Event, error) {
	s.logger.Debug("Getting all events")
	events, err := s.db.GetAllEvents(ctx, s.logger)
	if err != nil {
		s.logger.Error("Error getting all events", "error", err)
		return nil, err
	}
	return events, nil
}
