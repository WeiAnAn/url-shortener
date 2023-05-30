package shorturl

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxPersistentStore struct {
	pool *pgxpool.Pool
}

func NewPgxPersistentStore(p *pgxpool.Pool) *PgxPersistentStore {
	return &PgxPersistentStore{p}
}

func (p *PgxPersistentStore) Save(c context.Context, shortUrl *ShortURLWithExpireTime) error {
	conn, err := p.pool.Acquire(c)
	defer conn.Release()
	if err != nil {
		return err
	}

	_, err = conn.Exec(c, "INSERT INTO short_urls(short_url, original_url, expireAt) VALUES($1, $2, $3)",
		shortUrl.ShortUrl.ShortURL,
		shortUrl.ShortUrl.OriginalURL,
		shortUrl.ExpireAt,
	)

	return err
}

func (p *PgxPersistentStore) FindUnexpiredByShortURL(c context.Context, shortURL string) (*ShortURLWithExpireTime, error) {
	conn, err := p.pool.Acquire(c)
	defer conn.Release()
	if err != nil {
		return nil, err
	}
	var (
		short       string
		originalURL string
		expireAt    time.Time
	)

	err = p.pool.QueryRow(c, "SELECT short_url, original_url, expireAt FROM short_urls WHERE short_url = $1 AND expireAt > NOW()", shortURL).
		Scan(&short, &originalURL, &expireAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &ShortURLWithExpireTime{
		ShortUrl: &ShortURL{ShortURL: short, OriginalURL: originalURL},
		ExpireAt: expireAt,
	}, nil
}
