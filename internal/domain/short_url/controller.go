package shorturl

import (
	"net/http"
	"strings"
	"time"

	myerror "github.com/WeiAnAn/url-shortener/internal/my_error"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	service Service
}

func NewController(service Service) *Controller {
	return &Controller{service}
}

type CreateShortURLPayload struct {
	URL      string    `json:"url" binding:"required"`
	ExpireAt time.Time `json:"expireAt" binding:"required,gt"`
}

func (c *Controller) CreateShortURL(ctx *gin.Context) {
	var body CreateShortURLPayload
	err := ctx.BindJSON(&body)
	if err != nil {
		return
	}
	if !strings.HasPrefix(body.URL, "http://") && !strings.HasPrefix(body.URL, "https://") {
		err = myerror.NewValidationError("url", body.URL, "url must be http or https URL")
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if time.Now().AddDate(1, 0, 0).Before(body.ExpireAt) {
		err = myerror.NewValidationError("expireAt", body.ExpireAt.Format(time.RFC3339), "expireAt must be within one year")
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	shortUrl, err := c.service.CreateShortURL(ctx, body.URL, body.ExpireAt)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":       shortUrl.ShortUrl.ShortURL,
		"expireAt": shortUrl.ExpireAt.Format(time.RFC3339),
	})
}

type RedirectParams struct {
	URL string `uri:"url" binding:"required"`
}

func (c *Controller) Redirect(ctx *gin.Context) {
	var params RedirectParams
	err := ctx.BindUri(&params)
	if err != nil {
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
