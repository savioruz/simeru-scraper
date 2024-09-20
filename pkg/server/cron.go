package server

import (
	"context"
	"fmt"
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
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
	)
	repos := repositories.NewDB(redis)
	scrape := services.NewScrapeService(repos)
	_, err = c.cron.AddFunc("0 */12 * * *", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		log.Printf("scraping started at %s", time.Now().Format(time.RFC3339))

		if err := retries(ctx, func(ctx context.Context) error {
			return scrape.ScrapeStudyPrograms(ctx, opts...)
		}, "study programs"); err != nil {
			log.Printf("failed to scrape study programs after retries: %v", err)
		} else {
			log.Printf("study programs scraped successfully at %s", time.Now().Format(time.RFC3339))
		}

		if err := retries(ctx, func(ctx context.Context) error {
			return scrape.ScrapeSchedule(ctx, opts...)
		}, "schedule"); err != nil {
			log.Printf("failed to scrape schedule after retries: %v", err)
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

func retries(ctx context.Context, scrapeFn func(context.Context) error, taskName string) error {
	maxRetries := 3
	var err error

	for i := 0; i < maxRetries; i++ {
		err = scrapeFn(ctx)
		if err == nil {
			return nil
		}

		log.Printf("failed to scrape %s (attempt %d/%d): %v", taskName, i+1, maxRetries, err)

		// If this is not the last attempt, wait before retrying
		if i < maxRetries-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Duration(i+1) * time.Second):
				// Exponential backoff
			}
		}
	}

	return fmt.Errorf("failed to scrape %s after %d attempts: %v", taskName, maxRetries, err)
}
