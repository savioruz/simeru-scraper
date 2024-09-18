package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func FiberMiddleware(a *fiber.App) {
	a.Use(logger.New())
	a.Use(recover.New())
	a.Use(healthcheck.New())
	a.Get("/metrics", monitor.New(
		monitor.Config{
			Title: "short metrics",
		},
	))
	a.Use(cors.New(cors.Config{
		AllowOrigins: "https://*.savioruz.me, https://simeru.vercel.app",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET,POST",
	}))
	a.Use(requestid.New(requestid.Config{
		Header: "x-request-id",
	}))
}
