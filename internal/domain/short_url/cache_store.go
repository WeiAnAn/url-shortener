package shorturl

import "context"

type CacheStore interface {
	Get(c context.Context, key string) (string, error)
	Set(c context.Context, key, value string, expireSecond uint) error
}
