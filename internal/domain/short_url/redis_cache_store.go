package shorturl

import (
	"context"

	"github.com/redis/rueidis"
)

type RedisCacheStore struct {
	client rueidis.Client
}

func NewRedisCacheStore(client rueidis.Client) *RedisCacheStore {
	return &RedisCacheStore{client}
}

func (r *RedisCacheStore) Get(key string) (string, error) {
	ctx := context.Background()

	v, err := r.client.Do(ctx, r.client.B().Get().Key(key).Build()).ToString()

	if err != nil {
		if rueidis.IsRedisNil(err) {
			return "", nil
		}
		return "", err
	}

	return v, nil
}

func (r *RedisCacheStore) Set(key, value string, expireSecond uint) error {
	ctx := context.Background()

	cmd := r.client.B().Set().Key(key).Value(value).ExSeconds(int64(expireSecond)).Build()
	err := r.client.Do(ctx, cmd).Error()

	if err != nil {
		return err
	}
	return nil
}
