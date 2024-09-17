package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/savioruz/simeru-scraper/config"
	"github.com/savioruz/simeru-scraper/internal/adapters/cache"
	"github.com/savioruz/simeru-scraper/internal/adapters/handlers/rest"
	"github.com/savioruz/simeru-scraper/internal/adapters/repositories"
	"github.com/savioruz/simeru-scraper/internal/cores/services"
	"github.com/savioruz/simeru-scraper/pkg/middlewares"
	"github.com/savioruz/simeru-scraper/pkg/routes"
	"github.com/savioruz/simeru-scraper/pkg/utils"
	"net/http"
	"os"
	"os/signal"
)

type Fiber struct {
	app  *fiber.App
	conf *config.Config
	cron *CronAdapter
}

func NewFiberServer(conf *config.Config) Fiber {
	a := fiber.New()
	c := NewCronAdapter(conf)

	// Middleware
	middlewares.FiberMiddleware(a)
	middlewares.LimiterMiddleware(a)
	middlewares.MonitorMiddleware(a)

	return Fiber{
		app:  a,
		conf: conf,
		cron: c,
	}
}

func (s *Fiber) ServerStart() {
	s.cron.Start()
	s.initializeScheduleHandler()
	s.initializeRoutes()
	s.startServerWithGrafeculShutdown()
}

func (s *Fiber) Adaptor() http.HandlerFunc {
	s.cron.Start()
	s.initializeScheduleHandler()
	s.initializeRoutes()

	return adaptor.FiberApp(s.app)
}

func (s *Fiber) initializeScheduleHandler() {
	redis, err := cache.NewRedisCache(s.conf.Redis.Addr, s.conf.Redis.Password, s.conf.Redis.DB)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	validator := utils.NewValidator()
	repos := repositories.NewDB(redis)
	service := services.NewScheduleService(repos)
	handler := rest.NewScheduleHandler(service, validator)

	r := s.app.Group("/api/v1")
	r.Post("/schedule", handler.GetSchedule)
	r.Post("/study-programs", handler.GetStudyPrograms)
}

func (s *Fiber) initializeRoutes() {
	routes.SwaggerRoute(s.app)
	routes.NotFoundRoute(s.app)
}

func (s *Fiber) startServerWithGrafeculShutdown() {
	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := s.app.Shutdown(); err != nil {
			log.Errorf("Fiber shutdown error: %v", err)
		}

		close(idleConnsClosed)
	}()

	fiberConnectionURL := fmt.Sprintf("%s:%s", s.conf.Server.Host, s.conf.Server.Port)

	if err := s.app.Listen(fiberConnectionURL); err != nil {
		log.Errorf("Fiber listen error: %v", err)
	}

	<-idleConnsClosed
}
