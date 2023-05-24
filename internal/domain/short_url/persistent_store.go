package shorturl

type PersistentStore interface {
	Save(shortUrl *ShortURLWithExpireTime) error
	FindUnexpiredByShortURL(shortURL string) (*ShortURLWithExpireTime, error)
}
