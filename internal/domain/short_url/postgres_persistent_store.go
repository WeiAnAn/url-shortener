package shorturl

import (
	"database/sql"
	"time"
)

type PostgresPersistentStore struct {
	db *sql.DB
}

func NewPostgresPersistentStore(db *sql.DB) *PostgresPersistentStore {
	return &PostgresPersistentStore{db}
}

func (p *PostgresPersistentStore) Save(shortUrl *ShortURLWithExpireTime) error {
	_, err := p.db.Exec("INSERT INTO short_urls(short_url, original_url, expireAt) VALUES($1, $2, $3)",
		shortUrl.ShortUrl.ShortURL,
		shortUrl.ShortUrl.OriginalURL,
		shortUrl.ExpireAt,
	)

	return err
}

func (p *PostgresPersistentStore) FindUnexpiredByShortURL(shortURL string) (*ShortURLWithExpireTime, error) {
	var (
		short       string
		originalURL string
		expireAt    time.Time
	)
	err := p.db.QueryRow("SELECT short_url, original_url, expireAt FROM short_urls WHERE short_url = $1 AND expireAt > NOW()", shortURL).
		Scan(&short, &originalURL, &expireAt)

	if err != nil {
		return nil, err
	}

	return &ShortURLWithExpireTime{
		ShortUrl: &ShortURL{ShortURL: short, OriginalURL: originalURL},
		ExpireAt: expireAt,
	}, nil
}
