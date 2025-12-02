package bootstrap

import (
	"database/sql"
	"log/slog"
	"os"

	"github.com/ncostamagna/events-api/adapters/db"
	"github.com/ncostamagna/events-api/internal/events"
)

func NewEventsService(db *db.DB, logger *slog.Logger) events.Service {
	return events.NewService(db, logger)
}

func NewDatabase(logger *slog.Logger) *db.DB {
	database, err := sql.Open("postgres", os.Getenv("DB_DNS"))
	if err != nil {
		slog.Error("Error connecting to database", "error", err)
		os.Exit(1)
	}
	return db.New(database, logger)
}
