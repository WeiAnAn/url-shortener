package shorturl

import "time"

type ShortURLGenerator interface {
	Generate(int) (string, error)
}

type Service interface {
	CreateShortURL(string, time.Time) (*ShortURLWithExpireTime, error)
	GetOriginalURL(string) (*ShortURL, error)
}

type service struct {
	shortURLRepository ShortURLRepository
	shortURLGenerator  ShortURLGenerator
}

func NewService(sr ShortURLRepository, sg ShortURLGenerator) *service {
	return &service{sr, sg}
}

func (s *service) CreateShortURL(originalURL string, expireAt time.Time) (*ShortURLWithExpireTime, error) {
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

	err = s.shortURLRepository.Save(shortURL)
	if err != nil {
		return nil, err
	}

	return shortURL, nil
}

func (s *service) GetOriginalURL(short string) (*ShortURL, error) {
	shortURL, err := s.shortURLRepository.FindByShortURL(short)
	if err != nil {
		return nil, err
	}
	return shortURL, nil
}
