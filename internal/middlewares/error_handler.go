package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	myerror "github.com/WeiAnAn/url-shortener/internal/my_error"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, err := range c.Errors {
			var msg string
			status := http.StatusInternalServerError

			switch err.Err.(type) {
			case validator.ValidationErrors:
				myErr := err.Err.(validator.ValidationErrors)
				status = http.StatusBadRequest
				errorMessage := make([]string, len(myErr))
				for i, fieldErr := range myErr {
					errorMessage[i] = fmt.Sprintf("on %s with %v", fieldErr.Field(), fieldErr.Value())
				}
				msg = fmt.Sprintf("Validation errors: %s", strings.Join(errorMessage, ", "))
			case *myerror.ValidationError:
				myErr := err.Err.(*myerror.ValidationError)
				status = http.StatusBadRequest
				msg = myErr.Error()
			default:
				msg = "Internal server error"
			}

			if !c.Writer.Written() {
				c.JSON(status, gin.H{
					"message": msg,
				})
			}
		}
	}
}
