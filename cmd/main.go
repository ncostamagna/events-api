package main

import (
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/ncostamagna/events-api/pkg/bootstrap"
	"github.com/ncostamagna/events-api/pkg/log"
	"github.com/ncostamagna/events-api/transport/httpevents"
)

const defaultURL = "0.0.0.0:80"

func main() {
	_ = godotenv.Load()

	logger := log.New(log.Config{
		AppName:   "events-service",
		Level:     os.Getenv("LOG_LEVEL"),
		AddSource: true,
	})

	defer func() {
		if r := recover(); r != nil {
			logger.Error("Application panicked", "err", r)
		}
	}()

	eventsSrv := bootstrap.NewEventsService(bootstrap.NewDatabase(logger), logger)

	app := httpevents.NewHTTPServer(httpevents.MakeEventsEndpoints(eventsSrv))

	url := os.Getenv("APP_URL")
	if url == "" {
		url = defaultURL
	}

	errs := make(chan error)

	go func() {
		logger.Info("Listening", "url", url)
		errs <- app.Listen(url)
	}()

	err := <-errs
	if err != nil {
		logger.Error("Error server", "err", err)
	}
}
