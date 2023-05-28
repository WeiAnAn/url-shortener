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

func (r *RedisCacheStore) Get(c context.Context, key string) (*string, error) {
	v, err := r.client.Do(c, r.client.B().Get().Key(key).Build()).ToString()

	if err != nil {
		if rueidis.IsRedisNil(err) {
			return nil, nil
		}
		return nil, err
	}

	return &v, nil
}

func (r *RedisCacheStore) Set(c context.Context, key, value string, expireSecond uint) error {
	cmd := r.client.B().Set().Key(key).Value(value).ExSeconds(int64(expireSecond)).Build()
	err := r.client.Do(c, cmd).Error()

	if err != nil {
		return err
	}
	return nil
}
