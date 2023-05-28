package shorturl_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	shorturl "github.com/WeiAnAn/url-shortener/internal/domain/short_url"
	mock_shorturl "github.com/WeiAnAn/url-shortener/internal/domain/short_url/mocks"
	myerror "github.com/WeiAnAn/url-shortener/internal/my_error"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

	url := "https://pkg.go.dev"
	expireAt, err := time.Parse(time.RFC3339, time.Now().AddDate(0, 0, 1).Format(time.RFC3339))
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
		CreateShortURL(ctx, url, expireAt).
		Return(shortURL, nil)

	controller.CreateShortURL(ctx)

	var resBody CreateShortURLResponse
	json.Unmarshal(w.Body.Bytes(), &resBody)

	if resBody.ExpireAt != body.ExpireAt || resBody.Id != "aaaaaaa" {
		t.Fail()
	}
}

func TestCreateShortURLResponseBadRequestIfBodyIsEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, controller := createController(ctrl)
	w := httptest.NewRecorder()
	ctx := createGinContext(w)
	setPostRequest(ctx, struct{}{})

	controller.CreateShortURL(ctx)

	if w.Code != http.StatusBadRequest {
		t.Error("Status is not BadRequest")
	}

	if len(ctx.Errors) != 1 {
		t.Error("context errors size is not equal to one")
	}

	err, ok := ctx.Errors[0].Err.(validator.ValidationErrors)
	if !ok {
		t.Error("context error is not ValidationError")
	}

	if len(err) != 2 {
		t.Error("ValidationError field is not equal to 2")
	}

	if err[0].Field() != "URL" && err[0].Error() != "'CreateShortURLPayload.URL' Error:Field validation for 'URL' failed on the 'required' tag" {
		t.Error("unexpected err[0]")
	}

	if err[1].Field() != "ExpireAt" && err[1].Error() != "'CreateShortURLPayload.ExpireAt' Error:Field validation for 'ExpireAt' failed on the 'required' tag" {
		t.Error("unexpected err[1]")
	}
}

func TestCreateShortURLResponseBadRequestIfExpireAtIsNotDateString(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, controller := createController(ctrl)
	w := httptest.NewRecorder()
	ctx := createGinContext(w)

	body := struct {
		URL      string `json:"url"`
		ExpireAt string `json:"expireAt"`
	}{"https://pkg.go.dev", "2023"}
	setPostRequest(ctx, body)

	controller.CreateShortURL(ctx)

	if w.Code != http.StatusBadRequest {
		t.Error("Status is not BadRequest")
	}

	if len(ctx.Errors) != 1 {
		t.Error("context errors size is not equal to one")
	}

	if errors.Is(ctx.Errors[0].Err, &time.ParseError{}) {
		t.Error("context error is not time.ParseError")
	}
}

func TestCreateShortURLResponseBadRequestIfTimeIsBeforeNow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, controller := createController(ctrl)
	w := httptest.NewRecorder()
	ctx := createGinContext(w)

	url := "https://pkg.go.dev"
	expireAt, err := time.Parse(time.RFC3339, time.Now().AddDate(0, 0, -1).Format(time.RFC3339))
	if err != nil {
		panic(err)
	}
	body := struct {
		URL      string `json:"url"`
		ExpireAt string `json:"expireAt"`
	}{url, expireAt.Format(time.RFC3339)}

	setPostRequest(ctx, body)

	controller.CreateShortURL(ctx)

	if w.Code != http.StatusBadRequest {
		t.Error("Status is not BadRequest")
	}

	if len(ctx.Errors) != 1 {
		t.Error("context errors size is not equal to one")
	}

	validationErr, ok := ctx.Errors[0].Err.(validator.ValidationErrors)
	if !ok {
		t.Error("context error is not ValidationError")
	}

	if len(validationErr) != 1 {
		t.Error("ValidationError field is not equal to 2")
	}

	if validationErr.Error() != "Key: 'CreateShortURLPayload.ExpireAt' Error:Field validation for 'ExpireAt' failed on the 'gt' tag" {
		t.Error("unexpected Error")
	}
}

func TestCreateShortURLResponseBadRequestIfURLFormatIsInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, controller := createController(ctrl)
	w := httptest.NewRecorder()
	ctx := createGinContext(w)

	url := "htps://pkg.go.dev"
	expireAt, err := time.Parse(time.RFC3339, time.Now().AddDate(0, 0, 1).Format(time.RFC3339))
	if err != nil {
		panic(err)
	}
	body := struct {
		URL      string `json:"url"`
		ExpireAt string `json:"expireAt"`
	}{url, expireAt.Format(time.RFC3339)}

	setPostRequest(ctx, body)

	controller.CreateShortURL(ctx)

	if w.Code != http.StatusBadRequest {
		t.Error("Status is not BadRequest")
	}

	if len(ctx.Errors) != 1 {
		t.Error("context errors size is not equal to one")
	}

	validationErr, ok := ctx.Errors[0].Err.(*myerror.ValidationError)
	if !ok {
		t.Error("context error is not ValidationError")
	}

	if validationErr.Error() != "Validation failed on url with value htps://pkg.go.dev. url must be http or https URL" {
		t.Error("unexpected Error")
	}
}

func TestCreateShortURLResponseBadRequestIfExpireAtIsAfterOneYear(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, controller := createController(ctrl)
	w := httptest.NewRecorder()
	ctx := createGinContext(w)

	url := "https://pkg.go.dev"
	expireAt, err := time.Parse(time.RFC3339, time.Now().AddDate(1, 0, 1).Format(time.RFC3339))
	if err != nil {
		panic(err)
	}
	body := struct {
		URL      string `json:"url"`
		ExpireAt string `json:"expireAt"`
	}{url, expireAt.Format(time.RFC3339)}

	setPostRequest(ctx, body)

	controller.CreateShortURL(ctx)

	if w.Code != http.StatusBadRequest {
		t.Error("Status is not BadRequest")
	}

	if len(ctx.Errors) != 1 {
		t.Error("context errors size is not equal to one")
	}

	validationErr, ok := ctx.Errors[0].Err.(*myerror.ValidationError)
	if !ok {
		t.Error("context error is not ValidationError")
	}

	expectErrorMessage := fmt.Sprintf(
		"Validation failed on expireAt with value %s. expireAt must be within one year",
		expireAt.Format(time.RFC3339),
	)
	if validationErr.Error() != expectErrorMessage {
		t.Error("unexpected Error")
	}
}

func TestShouldSetContextErrorIfServiceReturnAnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService, controller := createController(ctrl)
	w := httptest.NewRecorder()
	ctx := createGinContext(w)

	url := "https://pkg.go.dev"
	expireAt, err := time.Parse(time.RFC3339, time.Now().AddDate(0, 0, 1).Format(time.RFC3339))
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
	mockErr := errors.New("error")
	mockService.EXPECT().
		CreateShortURL(ctx, url, expireAt).
		Return(shortURL, mockErr)

	controller.CreateShortURL(ctx)

	if len(ctx.Errors) != 1 {
		t.Error("Context errors length is not 1")
	}

	if ctx.Errors[0].Err != mockErr {
		t.Error("Unexpected Error")
	}
}

func TestRedirectRedirectExpectedURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService, controller := createController(ctrl)
	w := httptest.NewRecorder()
	ctx := createGinContext(w)
	url := "aaaaaaa"
	setRedirectRequest(ctx, url)
	originalURL := "https://pkg.go.dev"
	shortURL := &shorturl.ShortURL{
		ShortURL:    url,
		OriginalURL: originalURL,
	}
	mockService.EXPECT().GetOriginalURL(ctx, url).Return(shortURL, nil)

	controller.Redirect(ctx)

	if w.Code != http.StatusFound {
		t.Error("Unexpected status code")
	}

	if w.Header().Get("location") != originalURL {
		t.Fail()
	}
}

func TestRedirectResponseNotFoundIfURLLengthIsNotEqual7(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	_, controller := createController(ctrl)
	w := httptest.NewRecorder()
	ctx := createGinContext(w)
	url := "aaaaaa"
	setRedirectRequest(ctx, url)

	controller.Redirect(ctx)

	if w.Code != http.StatusNotFound {
		t.Fail()
	}
}

func TestRedirectResponseNotFoundIfServiceReturnNilShortURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService, controller := createController(ctrl)
	w := httptest.NewRecorder()
	ctx := createGinContext(w)
	url := "aaaaaaa"
	setRedirectRequest(ctx, url)
	mockService.EXPECT().GetOriginalURL(ctx, url).Return(nil, nil)

	controller.Redirect(ctx)

	if w.Code != http.StatusNotFound {
		t.Error("Unexpected status code")
	}
}

func TestRedirectSetContextErrorIfServiceReturnError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService, controller := createController(ctrl)
	w := httptest.NewRecorder()
	ctx := createGinContext(w)
	url := "aaaaaaa"
	setRedirectRequest(ctx, url)
	mockErr := errors.New("error")
	mockService.EXPECT().GetOriginalURL(ctx, url).Return(nil, mockErr)

	controller.Redirect(ctx)

	if len(ctx.Errors) != 1 {
		t.Error("context errors length is not 1")
	}

	if ctx.Errors[0].Err != mockErr {
		t.Error("unexpected error")
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

func setRedirectRequest(ctx *gin.Context, url string) {
	ctx.Request.Method = http.MethodGet
	if url != "" {
		ctx.Params = []gin.Param{
			{
				Key:   "url",
				Value: url,
			},
		}
	}
}
