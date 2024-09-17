package repositories

import (
	"github.com/savioruz/simeru-scraper/internal/adapters/cache"
	"github.com/savioruz/simeru-scraper/internal/cores/ports"
)

type DB struct {
	cache ports.CacheRepository
}

func NewDB(c *cache.RedisCache) *DB {
	return &DB{
		cache: c,
	}
}
