package shorturl

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service Service
}

func NewController(service Service) *Controller {
	return &Controller{service}
}

type CreateShortURLPayload struct {
	Url      string    `json:"url" binding:"required"`
	ExpireAt time.Time `json:"expireAt" binding:"required,gt"`
}

func (c *Controller) CreateShortURL(ctx *gin.Context) {
	var body CreateShortURLPayload
	err := ctx.BindJSON(&body)
	if err != nil {
		ctx.Error(err)
		return
	}
	shortUrl, err := c.service.CreateShortURL(body.Url, body.ExpireAt)
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
		ctx.Error(err)
		return
	}

	if len(params.URL) != 7 {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	shortURL, err := c.service.GetOriginalURL(params.URL)
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
