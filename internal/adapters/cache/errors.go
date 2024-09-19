package cache

import "errors"

var (
	ErrCacheMiss   = errors.New("cache miss")
	ErrCacheFailed = errors.New("failed to get data from redis")
	ErrMarshal     = errors.New("failed to marshal data")
	ErrUnmarshal   = errors.New("failed to unmarshal data")
)
