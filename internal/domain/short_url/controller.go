package shorturl

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	myerror "github.com/WeiAnAn/url-shortener/internal/my_error"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	service Service
	baseURL string
}

func NewController(service Service, baseURL string) *Controller {
	return &Controller{service, baseURL}
}

type CreateShortURLPayload struct {
	URL      string    `json:"url" binding:"required,url"`
	ExpireAt time.Time `json:"expireAt" binding:"required,gt"`
}

func (c *Controller) CreateShortURL(ctx *gin.Context) {
	var body CreateShortURLPayload
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		timeErr, isTimeParseErr := err.(*time.ParseError)
		if isTimeParseErr {
			err = myerror.NewValidationError("expireAt", timeErr.Value, "Invalid time format")
		}
		ctx.Error(err)
		return
	}
	if !strings.HasPrefix(body.URL, "http://") && !strings.HasPrefix(body.URL, "https://") {
		err = myerror.NewValidationError("url", body.URL, "url must be http or https URL")
		ctx.Error(err)
		return
	}
	if time.Now().AddDate(1, 0, 0).Before(body.ExpireAt) {
		err = myerror.NewValidationError("expireAt", body.ExpireAt.Format(time.RFC3339), "expireAt must be within one year")
		ctx.Error(err)
		return
	}

	shortUrl, err := c.service.CreateShortURL(ctx, body.URL, body.ExpireAt)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"url":      fmt.Sprintf("%s/%s", c.baseURL, shortUrl.ShortUrl.ShortURL),
		"expireAt": shortUrl.ExpireAt.Format(time.RFC3339),
	})
}

type RedirectParams struct {
	URL string `uri:"url" binding:"required"`
}

func (c *Controller) Redirect(ctx *gin.Context) {
	var params RedirectParams
	err := ctx.ShouldBindUri(&params)
	if err != nil {
		ctx.Error(err)
		return
	}

	if len(params.URL) != 7 {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	shortURL, err := c.service.GetOriginalURL(ctx, params.URL)
	if err != nil {
		ctx.Error(err)
		return
	}

	if shortURL == nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.Redirect(http.StatusFound, shortURL.OriginalURL)
}
