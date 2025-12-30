package httpevents

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func NewHTTPServer(endpoints Endpoints) *fiber.App {

	app := fiber.New()
	app.Use(recover.New())
	app.Get("/events", endpoints.GetAll)
	app.Post("/events", endpoints.Store)
	app.Get("/events/:id", endpoints.Get)
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	return app
}
