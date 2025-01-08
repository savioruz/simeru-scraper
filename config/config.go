package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

// Config holds the configuration settings for the application
type Config struct {
	Server ServerConfig
	Redis  RedisConfig
	Scrape ScrapeConfig
}

// ServerConfig holds the configuration settings for the server
type ServerConfig struct {
	Host string
	Port string
}

// RedisConfig holds the configuration settings for Redis
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// ScrapeConfig holds the configuration settings for Scrape
type ScrapeConfig struct {
	Interval int
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Get server configuration
	host := getEnv("APP_HOST", "localhost")
	port := getEnv("APP_PORT", "3000")

	// Get Redis configuration
	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")
	redisDB, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_DB value: %v", err)
	}

	scrapeInterval, err := strconv.Atoi(getEnv("SCRAPE_INTERVAL", "12"))
	if err != nil {
		return nil, fmt.Errorf("invalid SCRAPE_INTERVAL value: %v", err)
	}

	return &Config{
		Server: ServerConfig{
			Host: host,
			Port: port,
		},
		Redis: RedisConfig{
			Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
			Password: redisPassword,
			DB:       redisDB,
		},
		Scrape: ScrapeConfig{
			Interval: scrapeInterval,
		},
	}, nil
}

// getEnv retrieves the value of an environment variable or returns a default value if not set
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	log.Printf("Using default value for %s: %s", key, defaultValue)
	return defaultValue
}
