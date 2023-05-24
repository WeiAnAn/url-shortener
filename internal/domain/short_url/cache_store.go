package shorturl

type CacheStore interface {
	Get(key string) (string, error)
	Set(key, value string, expireSecond uint) error
}
