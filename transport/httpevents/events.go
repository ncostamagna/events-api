package httpevents

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/ncostamagna/events-api/internal/events"
)

type (
	Endpoints struct {
		Get    fiber.Handler
		GetAll fiber.Handler
		Store  fiber.Handler
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

func MakeEventsEndpoints(s events.Service) Endpoints {
	return Endpoints{
		Get:    makeGet(s),
		GetAll: makeGetAll(s),
		Store:  makeStore(s),
	}
}

func makeGet(service events.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := GetReq{
			ID: c.Params("id"),
		}

		if req.ID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "ID is required",
			})
		}

		if _, err := uuid.Parse(req.ID); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid UUID format",
			})
		}

		event, err := service.GetEventByID(c.Context(), req.ID)
		if err != nil {
			if errors.Is(err, events.ErrEventNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"message": "Event not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Success",
			"data":    event,
		})
	}
}

func makeGetAll(service events.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		events, err := service.GetAllEvents(c.Context())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Success",
			"data":    events,
		})
	}
}

func makeStore(service events.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := StoreReq{}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid request body",
			})
		}

		if req.Title == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Title is required",
			})
		}

		if req.Description == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Description is required",
			})
		}

		if len(req.Title) > 100 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Title must be less than 100 characters",
			})
		}

		if req.StartTime.IsZero() {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Start time is required",
			})
		}

		if req.EndTime.IsZero() {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "End time is required",
			})
		}

		if req.StartTime.After(req.EndTime) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Start time must be before end time",
			})
		}

		event, err := service.CreateEvent(c.Context(), req.Title, req.Description, req.StartTime, req.EndTime)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Success",
			"data":    event,
		})
	}
}
