package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func MonitorMiddleware(a *fiber.App) {
	a.Get("/metrics", monitor.New(
		monitor.Config{
			Title: "Short URL API Monitor",
		},
	))
	a.Use(healthcheck.New(
		healthcheck.Config{
			LivenessEndpoint:  "/api/v1/livez",
			ReadinessEndpoint: "/api/v1/readyz",
		},
	))
}
