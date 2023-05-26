package shorturl_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	shorturl "github.com/WeiAnAn/url-shortener/internal/domain/short_url"
	mock_shorturl "github.com/WeiAnAn/url-shortener/internal/domain/short_url/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

type CreateShortURLResponse struct {
	Id       string `json:"id"`
	ExpireAt string `json:"expireAt"`
}

func TestCreateShortURLResponseCreatedData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService, controller := createController(ctrl)
	w := httptest.NewRecorder()
	ctx := createGinContext(w)

	d, err := time.ParseDuration("24h")
	if err != nil {
		panic(err)
	}

	url := "https://pkg.go.dev"
	expireAt, err := time.Parse(time.RFC3339, time.Now().Add(d).Format(time.RFC3339))
	if err != nil {
		panic(err)
	}
	body := struct {
		URL      string `json:"url"`
		ExpireAt string `json:"expireAt"`
	}{url, expireAt.Format(time.RFC3339)}
	setPostRequest(ctx, body)

	shortURL := &shorturl.ShortURLWithExpireTime{
		ShortUrl: &shorturl.ShortURL{
			OriginalURL: url,
			ShortURL:    "aaaaaaa",
		},
		ExpireAt: expireAt,
	}
	mockService.EXPECT().
		CreateShortURL(url, expireAt).
		Return(shortURL, nil)

	controller.CreateShortURL(ctx)

	var resBody CreateShortURLResponse
	json.Unmarshal(w.Body.Bytes(), &resBody)

	if resBody.ExpireAt != body.ExpireAt || resBody.Id != "aaaaaaa" {
		t.Fail()
	}
}

func createController(ctrl *gomock.Controller) (*mock_shorturl.MockService, shorturl.Controller) {
	mockService := mock_shorturl.NewMockService(ctrl)
	controller := shorturl.NewController(mockService)

	return mockService, *controller
}

func createGinContext(w *httptest.ResponseRecorder) *gin.Context {
	gin.SetMode(gin.TestMode)

	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
	}

	return ctx
}

func setPostRequest(ctx *gin.Context, content interface{}) {
	ctx.Request.Method = http.MethodPost
	ctx.Request.Header.Set("Content-Type", "application/json")
	jsonBytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
}
