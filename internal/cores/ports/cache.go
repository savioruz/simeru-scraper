package ports

import "time"

type CacheRepository interface {
	Get(key string, value interface{}) error
	Set(key string, value interface{}, expiration time.Duration) error
}
