package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SwaggerRoute(a *fiber.App) {
	r := a.Group("/swagger")
	r.Get("*", swagger.HandlerDefault)

	a.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/swagger/index.html")
	})
}
