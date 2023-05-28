package main

import (
	"database/sql"
	"log"
	"net/http"

	shorturl "github.com/WeiAnAn/url-shortener/internal/domain/short_url"
	"github.com/WeiAnAn/url-shortener/internal/utils"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/redis/rueidis"
)

func main() {

	db, err := sql.Open("postgres", "user=postgres dbname=short_url password=password sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	redisClient, err := rueidis.NewClient(rueidis.ClientOption{InitAddress: []string{"127.0.0.1:6379"}})
	if err != nil {
		log.Fatal(err)
	}
	defer redisClient.Close()

	ps := shorturl.NewPostgresPersistentStore(db)
	cs := shorturl.NewRedisCacheStore(redisClient)
	sr := shorturl.NewRepository(ps, cs, &utils.RealTime{})
	sg := &utils.RandomBase62StringGenerator{}
	ss := shorturl.NewService(sr, sg)
	sc := shorturl.NewController(ss)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.POST("/api/v1/urls", sc.CreateShortURL)
	r.GET("/:url", sc.Redirect)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
