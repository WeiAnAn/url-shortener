package shorturl_test

import (
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

	ps.EXPECT().Save(gomock.Eq(url))

	repo.Save(url)
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
	ps.EXPECT().Save(gomock.Eq(url)).Return(mockErr)

	err := repo.Save(url)

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
	cs.EXPECT().Get(gomock.Eq(url.ShortUrl.ShortURL)).Return(url.ShortUrl.OriginalURL, nil)

	result, err := repo.FindByShortURL(url.ShortUrl.ShortURL)
	if err != nil {
		t.Fail()
	}
	if result.OriginalURL != url.ShortUrl.OriginalURL {
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
	cs.EXPECT().Get(gomock.Eq(url.ShortUrl.ShortURL)).Return("", nil)
	ps.EXPECT().FindUnexpiredByShortURL(gomock.Eq(url.ShortUrl.ShortURL)).Return(url, nil)
	cs.EXPECT().Set(url.ShortUrl.ShortURL, url.ShortUrl.OriginalURL, uint(300)).Return(nil)

	result, err := repo.FindByShortURL(url.ShortUrl.ShortURL)
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
	cs.EXPECT().Get(gomock.Eq(url.ShortUrl.ShortURL)).Return("", nil)
	ps.EXPECT().FindUnexpiredByShortURL(gomock.Eq(url.ShortUrl.ShortURL)).Return(url, nil)
	tu.EXPECT().Until(expireAt).Return(d)
	cs.EXPECT().Set(url.ShortUrl.ShortURL, url.ShortUrl.OriginalURL, uint(d.Seconds())).Return(nil)

	result, err := repo.FindByShortURL(url.ShortUrl.ShortURL)
	if err != nil {
		t.Fail()
	}
	if result.OriginalURL != url.ShortUrl.OriginalURL {
		t.Fail()
	}
}

func TestFindByShortURLDoNotSetCacheIfOriginalURLIsEmptyString(t *testing.T) {
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
	cs.EXPECT().Get(gomock.Eq(url.ShortUrl.ShortURL)).Return("", nil)
	ps.EXPECT().FindUnexpiredByShortURL(gomock.Eq(url.ShortUrl.ShortURL)).Return(nil, nil)
	cs.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	result, err := repo.FindByShortURL(url.ShortUrl.ShortURL)
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
	cs.EXPECT().Get(gomock.Eq(url.ShortUrl.ShortURL)).Return("", mockErr)

	_, err := repo.FindByShortURL(url.ShortUrl.ShortURL)
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
	cs.EXPECT().Get(gomock.Eq(url.ShortUrl.ShortURL)).Return("", nil)
	ps.EXPECT().FindUnexpiredByShortURL(gomock.Eq(url.ShortUrl.ShortURL)).Return(nil, mockErr)

	_, err := repo.FindByShortURL(url.ShortUrl.ShortURL)
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
	cs.EXPECT().Get(gomock.Eq(url.ShortUrl.ShortURL)).Return("", nil)
	ps.EXPECT().FindUnexpiredByShortURL(gomock.Eq(url.ShortUrl.ShortURL)).Return(url, nil)
	tu.EXPECT().Until(expireAt).Return(d)
	cs.EXPECT().Set(url.ShortUrl.ShortURL, url.ShortUrl.OriginalURL, uint(d.Seconds())).Return(mockErr)

	_, err := repo.FindByShortURL(url.ShortUrl.ShortURL)
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
