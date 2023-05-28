package shorturl

import (
	"context"
	"time"
)

type ShortURLGenerator interface {
	Generate(int) (string, error)
}

type Service interface {
	CreateShortURL(context.Context, string, time.Time) (*ShortURLWithExpireTime, error)
	GetOriginalURL(context.Context, string) (*ShortURL, error)
}

type service struct {
	shortURLRepository ShortURLRepository
	shortURLGenerator  ShortURLGenerator
}

func NewService(sr ShortURLRepository, sg ShortURLGenerator) *service {
	return &service{sr, sg}
}

func (s *service) CreateShortURL(c context.Context, originalURL string, expireAt time.Time) (*ShortURLWithExpireTime, error) {
	short, err := s.shortURLGenerator.Generate(7)
	if err != nil {
		return nil, err
	}

	shortURL := &ShortURLWithExpireTime{
		ShortUrl: &ShortURL{
			OriginalURL: originalURL,
			ShortURL:    short,
		},
		ExpireAt: expireAt,
	}

	err = s.shortURLRepository.Save(c, shortURL)
	if err != nil {
		return nil, err
	}

	return shortURL, nil
}

func (s *service) GetOriginalURL(c context.Context, short string) (*ShortURL, error) {
	shortURL, err := s.shortURLRepository.FindByShortURL(c, short)
	if err != nil {
		return nil, err
	}
	return shortURL, nil
}
