package shorturl

import "context"

type PersistentStore interface {
	Save(c context.Context, shortUrl *ShortURLWithExpireTime) error
	FindUnexpiredByShortURL(c context.Context, shortURL string) (*ShortURLWithExpireTime, error)
}
