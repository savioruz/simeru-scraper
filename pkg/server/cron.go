package server

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/robfig/cron/v3"
	"github.com/savioruz/simeru-scraper/config"
	"github.com/savioruz/simeru-scraper/internal/adapters/cache"
	"github.com/savioruz/simeru-scraper/internal/adapters/repositories"
	"github.com/savioruz/simeru-scraper/internal/cores/services"
	"log"
	"time"
)

type CronAdapter struct {
	conf *config.Config
	cron *cron.Cron
}

func NewCronAdapter(conf *config.Config) *CronAdapter {
	return &CronAdapter{
		conf: conf,
		cron: cron.New(
			cron.WithLogger(cron.DefaultLogger),
			cron.WithChain(cron.Recover(cron.DefaultLogger)),
		),
	}
}

func (c *CronAdapter) Start() {
	redis, err := cache.NewRedisCache(c.conf.Redis.Addr, c.conf.Redis.Password, c.conf.Redis.DB)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent("Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0"),
		chromedp.Headless,
		chromedp.DisableGPU,
		chromedp.IgnoreCertErrors,
	)
	repos := repositories.NewDB(redis)
	scrape := services.NewScrapeService(repos)
	_, err = c.cron.AddFunc("@every 1m", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		log.Printf("scraping started at %s", time.Now().Format(time.RFC3339))
		if err := scrape.ScrapeStudyPrograms(ctx, opts...); err != nil {
			log.Printf("failed to scrape study programs: %v", err)
		} else {
			log.Printf("study programs scraped successfully at %s", time.Now().Format(time.RFC3339))
		}

		if err := scrape.ScrapeSchedule(ctx, opts...); err != nil {
			log.Printf("failed to scrape schedule: %v", err)
		} else {
			log.Printf("schedule scraped successfully at %s", time.Now().Format(time.RFC3339))
		}
	})
	if err != nil {
		log.Fatalf("failed to add job: %v", err)
	}

	log.Printf("starting scheduler at %s", time.Now().Format(time.RFC3339))
	c.cron.Start()
}

func (c *CronAdapter) Stop() {
	c.cron.Stop()
}
