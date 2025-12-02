package main

import (
	"context"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/ncostamagna/events-api/pkg/bootstrap"
	"github.com/ncostamagna/events-api/pkg/httputil"
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

	ctx := context.Background()

	eventsSrv := bootstrap.NewEventsService(bootstrap.NewDatabase(logger), logger)

	h := httpevents.NewHTTPServer(ctx, httpevents.MakeEndpoints(eventsSrv))

	url := os.Getenv("APP_URL")
	if url == "" {
		url = defaultURL
	}

	srv := &http.Server{
		Handler: httputil.AccessControl(h),
		Addr:    url,
	}

	errs := make(chan error)

	go func() {
		logger.Info("Listening", "url", url)
		errs <- srv.ListenAndServe()
	}()

	err := <-errs
	if err != nil {
		logger.Error("Error server", "err", err)
	}
}
