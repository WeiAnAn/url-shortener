package shorturl_test

import (
	"errors"
	"testing"
	"time"

	shorturl "github.com/WeiAnAn/url-shortener/internal/domain/short_url"
	mock_shorturl "github.com/WeiAnAn/url-shortener/internal/domain/short_url/mocks"
	"github.com/golang/mock/gomock"
)

func TestCreateShortURLGenerateShortURLAndCallRepoSave(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo, mockShortURLGenerator, service := createService(ctrl)

	originalURL := "https://pkg.go.dev/"
	expireAt := time.Now()
	shortURL := "aaaaaaa"
	mockShortURLGenerator.EXPECT().Generate(7).Return(shortURL, nil)
	mockRepo.EXPECT().Save(&shorturl.ShortURLWithExpireTime{
		ShortUrl: &shorturl.ShortURL{
			ShortURL:    shortURL,
			OriginalURL: originalURL,
		},
		ExpireAt: expireAt,
	}).Return(nil)

	result, err := service.CreateShortURL(originalURL, expireAt)
	if err != nil {
		t.Fail()
	}

	if !result.ExpireAt.Equal(expireAt) || result.ShortUrl.OriginalURL != originalURL || result.ShortUrl.ShortURL != shortURL {
		t.Fail()
	}
}

func TestCreateShortURLReturnErrorIfGeneratorReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, mockShortURLGenerator, service := createService(ctrl)

	originalURL := "https://pkg.go.dev/"
	expireAt := time.Now()

	mockErr := errors.New("error")
	mockShortURLGenerator.EXPECT().Generate(7).Return("", mockErr)

	_, err := service.CreateShortURL(originalURL, expireAt)

	if err != mockErr {
		t.Fail()
	}
}

func TestCreateShortURLReturnErrorIfRepoSaveReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo, mockShortURLGenerator, service := createService(ctrl)

	originalURL := "https://pkg.go.dev/"
	expireAt := time.Now()
	shortURL := "aaaaaaa"
	mockShortURLGenerator.EXPECT().Generate(7).Return(shortURL, nil)
	mockErr := errors.New("error")
	mockRepo.EXPECT().Save(&shorturl.ShortURLWithExpireTime{
		ShortUrl: &shorturl.ShortURL{
			ShortURL:    shortURL,
			OriginalURL: originalURL,
		},
		ExpireAt: expireAt,
	}).Return(mockErr)

	_, err := service.CreateShortURL(originalURL, expireAt)
	if err != mockErr {
		t.Fail()
	}
}

func TestGetOriginalURLReturnExpectedURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo, _, service := createService(ctrl)

	shortURL := &shorturl.ShortURL{
		ShortURL:    "aaaaaaa",
		OriginalURL: "https://pkg.go.dev",
	}
	mockRepo.EXPECT().FindByShortURL(shortURL.ShortURL).Return(shortURL, nil)

	result, err := service.GetOriginalURL(shortURL.ShortURL)
	if err != nil {
		t.Fail()
	}
	if result != shortURL {
		t.Fail()
	}
}

func TestGetOriginalURLReturnErrorIfRepoFindReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo, _, service := createService(ctrl)

	shortURL := &shorturl.ShortURL{
		ShortURL:    "aaaaaaa",
		OriginalURL: "https://pkg.go.dev",
	}
	mockErr := errors.New("error")
	mockRepo.EXPECT().FindByShortURL(shortURL.ShortURL).Return(nil, mockErr)

	_, err := service.GetOriginalURL(shortURL.ShortURL)
	if err != mockErr {
		t.Fail()
	}
}

func createService(ctrl *gomock.Controller) (*mock_shorturl.MockShortURLRepository, *mock_shorturl.MockShortURLGenerator, shorturl.Service) {
	mockRepo := mock_shorturl.NewMockShortURLRepository(ctrl)
	mockShortURLGenerator := mock_shorturl.NewMockShortURLGenerator(ctrl)
	service := shorturl.NewService(mockRepo, mockShortURLGenerator)
	return mockRepo, mockShortURLGenerator, service
}
