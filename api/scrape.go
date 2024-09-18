package handler

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/savioruz/simeru-scraper/config"
	"github.com/savioruz/simeru-scraper/internal/adapters/cache"
	"github.com/savioruz/simeru-scraper/internal/adapters/repositories"
	"github.com/savioruz/simeru-scraper/internal/cores/services"
	"log"
	"net/http"
	"time"
)

func ScrapeHandler(w http.ResponseWriter, r *http.Request) {
	conf, err := config.LoadConfig()
	if err != nil {
		http.Error(w, "Error loading config", http.StatusInternalServerError)
		log.Printf("Error loading config: %v", err)
		return
	}

	redis, err := cache.NewRedisCache(conf.Redis.Addr, conf.Redis.Password, conf.Redis.DB)
	if err != nil {
		http.Error(w, "Failed to connect to Redis", http.StatusInternalServerError)
		log.Printf("Failed to connect to Redis: %v", err)
		return
	}

	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent("Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0"),
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
	)

	repos := repositories.NewDB(redis)
	scrape := services.NewScrapeService(repos)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if err := scrape.ScrapeStudyPrograms(ctx, opts...); err != nil {
		http.Error(w, "Failed to scrape study programs", http.StatusInternalServerError)
		log.Printf("Failed to scrape study programs: %v", err)
		return
	}

	if err := scrape.ScrapeSchedule(ctx, opts...); err != nil {
		http.Error(w, "Failed to scrape schedule", http.StatusInternalServerError)
		log.Printf("Failed to scrape schedule: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("Scraping completed successfully at %s", time.Now().Format(time.RFC3339))
}
