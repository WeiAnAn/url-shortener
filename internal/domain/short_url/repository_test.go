package shorturl_test

import (
	"context"
	"errors"
	"testing"
	"time"

	shorturl "github.com/WeiAnAn/url-shortener/internal/domain/short_url"
	mock_shorturl "github.com/WeiAnAn/url-shortener/internal/domain/short_url/mocks"
	mock_utils "github.com/WeiAnAn/url-shortener/internal/utils/mocks"
	"github.com/golang/mock/gomock"
)

func TestSaveCallPersistentStoreSave(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ps, cs, tu := createMock(mockCtrl)
	repo := shorturl.NewRepository(ps, cs, tu)

	now := time.Now()
	url := &shorturl.ShortURLWithExpireTime{
		ShortUrl: &shorturl.ShortURL{
			ShortURL:    "short",
			OriginalURL: "https://example.com/long",
		},
		ExpireAt: now,
	}
	c := context.Background()
	ps.EXPECT().Save(c, gomock.Eq(url))

	repo.Save(c, url)
}

func TestSaveReturnError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ps, cs, tu := createMock(mockCtrl)
	repo := shorturl.NewRepository(ps, cs, tu)

	now := time.Now()
	url := &shorturl.ShortURLWithExpireTime{
		ShortUrl: &shorturl.ShortURL{
			ShortURL:    "short",
			OriginalURL: "https://example.com/long",
		},
		ExpireAt: now,
	}

	mockErr := errors.New("error")
	c := context.Background()
	ps.EXPECT().Save(c, gomock.Eq(url)).Return(mockErr)

	err := repo.Save(c, url)

	if err != mockErr {
		t.Error("save not return same error")
	}
}

func TestFindByShortURLGetFromCacheFirst(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ps, cs, tu := createMock(mockCtrl)
	repo := shorturl.NewRepository(ps, cs, tu)

	now := time.Now()
	url := &shorturl.ShortURLWithExpireTime{
		ShortUrl: &shorturl.ShortURL{
			ShortURL:    "short",
			OriginalURL: "https://example.com/long",
		},
		ExpireAt: now,
	}
	c := context.Background()
	cs.EXPECT().Get(c, gomock.Eq(url.ShortUrl.ShortURL)).Return(&url.ShortUrl.OriginalURL, nil)

	result, err := repo.FindByShortURL(c, url.ShortUrl.ShortURL)
	if err != nil {
		t.Fail()
	}
	if result.OriginalURL != url.ShortUrl.OriginalURL {
		t.Fail()
	}
}
func TestFindByShortURLReturnNilIfCacheReturnEmptyString(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ps, cs, tu := createMock(mockCtrl)
	repo := shorturl.NewRepository(ps, cs, tu)

	s := ""
	shortURL := "short"
	c := context.Background()
	cs.EXPECT().Get(c, gomock.Eq(shortURL)).Return(&s, nil)

	result, err := repo.FindByShortURL(c, shortURL)
	if err != nil {
		t.Fail()
	}
	if result != nil {
		t.Fail()
	}
}

func TestFindByShortURLGetFromPersistent(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ps, cs, tu := createMock(mockCtrl)
	repo := shorturl.NewRepository(ps, cs, tu)

	d, _ := time.ParseDuration("24h")
	expireAt := time.Now().Add(d)
	url := &shorturl.ShortURLWithExpireTime{
		ShortUrl: &shorturl.ShortURL{
			ShortURL:    "short",
			OriginalURL: "https://example.com/long",
		},
		ExpireAt: expireAt,
	}
	c := context.Background()
	cs.EXPECT().Get(c, gomock.Eq(url.ShortUrl.ShortURL)).Return(nil, nil)
	ps.EXPECT().FindUnexpiredByShortURL(c, gomock.Eq(url.ShortUrl.ShortURL)).Return(url, nil)
	cs.EXPECT().Set(c, url.ShortUrl.ShortURL, url.ShortUrl.OriginalURL, uint(300)).Return(nil)
	tu.EXPECT().Until(expireAt).Return(d)

	result, err := repo.FindByShortURL(c, url.ShortUrl.ShortURL)
	if err != nil {
		t.Fail()
	}
	if result.OriginalURL != url.ShortUrl.OriginalURL {
		t.Fail()
	}
}

func TestFindByShortURLSetCacheIfTimeToExpireIsLessThan300Second(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ps, cs, tu := createMock(mockCtrl)
	repo := shorturl.NewRepository(ps, cs, tu)

	d, _ := time.ParseDuration("200s")
	expireAt := time.Now().Add(d)
	url := &shorturl.ShortURLWithExpireTime{
		ShortUrl: &shorturl.ShortURL{
			ShortURL:    "short",
			OriginalURL: "https://example.com/long",
		},
		ExpireAt: expireAt,
	}
	c := context.Background()
	cs.EXPECT().Get(c, gomock.Eq(url.ShortUrl.ShortURL)).Return(nil, nil)
	ps.EXPECT().FindUnexpiredByShortURL(c, gomock.Eq(url.ShortUrl.ShortURL)).Return(url, nil)
	tu.EXPECT().Until(expireAt).Return(d)
	cs.EXPECT().Set(c, url.ShortUrl.ShortURL, url.ShortUrl.OriginalURL, uint(d.Seconds())).Return(nil)

	result, err := repo.FindByShortURL(c, url.ShortUrl.ShortURL)
	if err != nil {
		t.Fail()
	}
	if result.OriginalURL != url.ShortUrl.OriginalURL {
		t.Fail()
	}
}

func TestFindByShortURLSetCacheEmptyStringIfPersistentStoreReturnNil(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ps, cs, tu := createMock(mockCtrl)
	repo := shorturl.NewRepository(ps, cs, tu)

	expireAt := time.Now().AddDate(0, 0, 1)
	url := &shorturl.ShortURLWithExpireTime{
		ShortUrl: &shorturl.ShortURL{
			ShortURL:    "short",
			OriginalURL: "",
		},
		ExpireAt: expireAt,
	}
	c := context.Background()
	cs.EXPECT().Get(c, gomock.Eq(url.ShortUrl.ShortURL)).Return(nil, nil)
	ps.EXPECT().FindUnexpiredByShortURL(c, gomock.Eq(url.ShortUrl.ShortURL)).Return(nil, nil)
	cs.EXPECT().Set(c, url.ShortUrl.ShortURL, "", uint(300)).Return(nil)

	result, err := repo.FindByShortURL(c, url.ShortUrl.ShortURL)
	if err != nil {
		t.Fail()
	}
	if result != nil {
		t.Fail()
	}
}

func TestFindByShortURLReturnErrorIfCacheGetReturnError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ps, cs, tu := createMock(mockCtrl)
	repo := shorturl.NewRepository(ps, cs, tu)

	expireAt := time.Now()
	url := &shorturl.ShortURLWithExpireTime{
		ShortUrl: &shorturl.ShortURL{
			ShortURL:    "short",
			OriginalURL: "",
		},
		ExpireAt: expireAt,
	}
	mockErr := errors.New("Error")
	c := context.Background()
	cs.EXPECT().Get(c, gomock.Eq(url.ShortUrl.ShortURL)).Return(nil, mockErr)

	_, err := repo.FindByShortURL(c, url.ShortUrl.ShortURL)
	if err != mockErr {
		t.Fail()
	}
}

func TestFindByShortURLReturnErrorIfPersistentGetReturnError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ps, cs, tu := createMock(mockCtrl)
	repo := shorturl.NewRepository(ps, cs, tu)

	expireAt := time.Now()
	url := &shorturl.ShortURLWithExpireTime{
		ShortUrl: &shorturl.ShortURL{
			ShortURL:    "short",
			OriginalURL: "",
		},
		ExpireAt: expireAt,
	}
	mockErr := errors.New("Error")
	c := context.Background()
	cs.EXPECT().Get(c, gomock.Eq(url.ShortUrl.ShortURL)).Return(nil, nil)
	ps.EXPECT().FindUnexpiredByShortURL(c, gomock.Eq(url.ShortUrl.ShortURL)).Return(nil, mockErr)

	_, err := repo.FindByShortURL(c, url.ShortUrl.ShortURL)
	if err != mockErr {
		t.Fail()
	}
}

func TestFindByShortURLReturnErrorIfCacheSetReturnError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ps, cs, tu := createMock(mockCtrl)
	repo := shorturl.NewRepository(ps, cs, tu)

	d, _ := time.ParseDuration("200s")
	expireAt := time.Now().Add(d)
	url := &shorturl.ShortURLWithExpireTime{
		ShortUrl: &shorturl.ShortURL{
			ShortURL:    "short",
			OriginalURL: "https://example.com/long",
		},
		ExpireAt: expireAt,
	}
	mockErr := errors.New("Error")
	c := context.Background()
	cs.EXPECT().Get(c, gomock.Eq(url.ShortUrl.ShortURL)).Return(nil, nil)
	ps.EXPECT().FindUnexpiredByShortURL(c, gomock.Eq(url.ShortUrl.ShortURL)).Return(url, nil)
	tu.EXPECT().Until(expireAt).Return(d)
	cs.EXPECT().Set(c, url.ShortUrl.ShortURL, url.ShortUrl.OriginalURL, uint(d.Seconds())).Return(mockErr)

	_, err := repo.FindByShortURL(c, url.ShortUrl.ShortURL)
	if err != mockErr {
		t.Fail()
	}
}

func createMock(ctrl *gomock.Controller) (*mock_shorturl.MockPersistentStore, *mock_shorturl.MockCacheStore, *mock_utils.MockTimeUtil) {
	mockCacheStore := mock_shorturl.NewMockCacheStore(ctrl)
	mockPersistentStore := mock_shorturl.NewMockPersistentStore(ctrl)
	mockTime := mock_utils.NewMockTimeUtil(ctrl)

	return mockPersistentStore, mockCacheStore, mockTime
}
