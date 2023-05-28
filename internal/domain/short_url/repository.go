package shorturl

import (
	"context"
	"math"
	"time"

	"github.com/WeiAnAn/url-shortener/internal/utils"
)

type ShortURLRepository interface {
	Save(context.Context, *ShortURLWithExpireTime) error
	FindByShortURL(context.Context, string) (*ShortURL, error)
}

type shortURLRepository struct {
	persistentStore PersistentStore
	cacheStore      CacheStore
	time            utils.TimeUtil
}

type ShortURL struct {
	ShortURL    string
	OriginalURL string
}

type ShortURLWithExpireTime struct {
	ShortUrl *ShortURL
	ExpireAt time.Time
}

func NewRepository(ps PersistentStore, cs CacheStore, t utils.TimeUtil) *shortURLRepository {
	repo := &shortURLRepository{ps, cs, t}

	return repo
}

func (repo *shortURLRepository) Save(c context.Context, shortURL *ShortURLWithExpireTime) error {
	err := repo.persistentStore.Save(c, shortURL)
	if err != nil {
		return err
	}
	return nil
}

func (repo *shortURLRepository) FindByShortURL(c context.Context, shortURL string) (*ShortURL, error) {
	originalURL, err := repo.cacheStore.Get(c, shortURL)
	if err != nil {
		return nil, err
	}
	if originalURL != "" {
		return &ShortURL{ShortURL: shortURL, OriginalURL: originalURL}, nil
	}

	url, err := repo.persistentStore.FindUnexpiredByShortURL(c, shortURL)
	if err != nil {
		return nil, err
	}

	if url == nil {
		return nil, nil
	}

	timeToExpired := repo.time.Until(url.ExpireAt).Seconds()
	cacheSecond := math.Min(timeToExpired, 300)
	err = repo.cacheStore.Set(c, url.ShortUrl.ShortURL, url.ShortUrl.OriginalURL, uint(cacheSecond))
	if err != nil {
		return nil, err
	}

	return url.ShortUrl, nil
}
