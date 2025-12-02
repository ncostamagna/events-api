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
	event := &domain.Event{
		Title:       title,
		Description: description,
		StartTime:   startTime,
		EndTime:     endTime,
		CreatedAt:   time.Now(),
	}
	err := s.db.CreateEvent(ctx, event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (s *service) GetEventByID(ctx context.Context, id string) (*domain.Event, error) {
	event, err := s.db.GetEventByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, ErrEventNotFound
	}
	return event, nil
}

func (s *service) GetAllEvents(ctx context.Context) ([]domain.Event, error) {
	return s.db.GetAllEvents(ctx, s.logger)
}
