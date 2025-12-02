package db

import (
	"database/sql"
	"log/slog"
)

type (
	DB struct {
		db     *sql.DB
		logger *slog.Logger
	}
)

func New(db *sql.DB, logger *slog.Logger) *DB {
	return &DB{db: db, logger: logger}
}
