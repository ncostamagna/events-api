package httpevents

import (
	"context"

	"github.com/ncostamagna/events-api/internal/events"
	"github.com/ncostamagna/go-http-utils/response"

	"time"

	"github.com/ncostamagna/events-api/pkg/httputil"
)

type (
	Endpoints struct {
		Get    httputil.Endpoint
		GetAll httputil.Endpoint
		Store  httputil.Endpoint
	}

	StoreReq struct {
		Title       string    `json:"title"`
		Description string    `json:"description"`
		StartTime   time.Time `json:"start_time"`
		EndTime     time.Time `json:"end_time"`
	}

	GetReq struct {
		ID string `json:"id"`
	}
)

func MakeEndpoints(s events.Service) Endpoints {
	return Endpoints{
		Get:    makeGet(s),
		GetAll: makeGetAll(s),
		Store:  makeStore(s),
	}
}

func makeGet(service events.Service) httputil.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetReq)

		event, err := service.GetEventByID(ctx, req.ID)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("Success", event, nil), nil
	}
}

func makeGetAll(service events.Service) httputil.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		events, err := service.GetAllEvents(ctx)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("Success", events, nil), nil
	}
}

func makeStore(service events.Service) httputil.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(StoreReq)

		if req.Title == "" {
			return nil, response.BadRequest("Title is required")
		}

		if req.Description == "" {
			return nil, response.BadRequest("Description is required")
		}

		event, err := service.CreateEvent(ctx, req.Title, req.Description, req.StartTime, req.EndTime)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}
		return response.Created("Success", event, nil), nil
	}
}
