package db

import (
	"context"
	"log/slog"

	"database/sql"
	_ "embed"
	"errors"

	"github.com/ncostamagna/events-api/domain"
)

var (
	//go:embed queries/create_event.sql
	createEventQuery string

	//go:embed queries/get_event_by_id.sql
	getEventByIDQuery string

	//go:embed queries/get_all_events.sql
	getAllEventsQuery string
)

func (db *DB) CreateEvent(ctx context.Context, event *domain.Event) error {
	err := db.db.QueryRowContext(
		ctx,
		createEventQuery,
		event.Title,
		event.Description,
		event.StartTime,
		event.EndTime,
	).Scan(&event.ID)

	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetEventByID(ctx context.Context, id string) (*domain.Event, error) {
	var event domain.Event
	err := db.db.QueryRowContext(ctx, getEventByIDQuery, id).Scan(
		&event.ID,
		&event.Title,
		&event.Description,
		&event.StartTime,
		&event.EndTime,
		&event.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &event, nil
}

func (db *DB) GetAllEvents(ctx context.Context, logger *slog.Logger) ([]domain.Event, error) {
	rows, err := db.db.QueryContext(ctx, getAllEventsQuery)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error("failed to close rows", "error", err)
		}
	}()

	var events []domain.Event
	for rows.Next() {
		var event domain.Event
		err := rows.Scan(&event.ID, &event.Title, &event.Description, &event.StartTime, &event.EndTime, &event.CreatedAt)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}
