package shorturl

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const COLLECTION_NAME = "short_urls"

type MongoPersistentStore struct {
	client   *mongo.Client
	database string
}

type ShortURLDocument struct {
	ShortURL    string    `bson:"short_url"`
	OriginalURL string    `bson:"original_url"`
	ExpireAt    time.Time `bson:"expire_at"`
}

func NewMongoPersistentStore(c *mongo.Client, d string) *MongoPersistentStore {
	unique := true
	index := mongo.IndexModel{
		Keys: bson.D{
			bson.E{Key: "short_url", Value: 1},
		},
		Options: &options.IndexOptions{Unique: &unique},
	}
	c.Database(d).Collection(COLLECTION_NAME).Indexes().CreateOne(context.Background(), index)
	return &MongoPersistentStore{c, d}
}

func (m *MongoPersistentStore) Save(c context.Context, shortUrl *ShortURLWithExpireTime) error {
	doc := ShortURLDocument{shortUrl.ShortUrl.ShortURL, shortUrl.ShortUrl.OriginalURL, shortUrl.ExpireAt}

	_, err := m.client.Database(m.database).Collection(COLLECTION_NAME).InsertOne(c, doc)

	return err
}

func (m *MongoPersistentStore) FindUnexpiredByShortURL(c context.Context, shortURL string) (*ShortURLWithExpireTime, error) {
	var doc ShortURLDocument
	err := m.client.Database(m.database).Collection(COLLECTION_NAME).FindOne(c, bson.M{
		"short_url": shortURL,
		"expire_at": bson.M{
			"$gt": time.Now(),
		},
	}).Decode(&doc)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &ShortURLWithExpireTime{
		ShortUrl: &ShortURL{ShortURL: doc.ShortURL, OriginalURL: doc.OriginalURL},
		ExpireAt: doc.ExpireAt,
	}, nil
}
