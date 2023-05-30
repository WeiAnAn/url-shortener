package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/WeiAnAn/url-shortener/internal/config"
	shorturl "github.com/WeiAnAn/url-shortener/internal/domain/short_url"
	"github.com/WeiAnAn/url-shortener/internal/middlewares"
	"github.com/WeiAnAn/url-shortener/internal/utils"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/redis/rueidis"
	"github.com/spf13/viper"
)

func main() {
	db := setupDB()
	defer db.Close()

	redisClient := setupRedis()
	defer redisClient.Close()

	r := setupRouter(db, redisClient)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func setupDB() *sql.DB {
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		viper.GetString("DATABASE_USER"),
		viper.GetString("DATABASE_PASSWORD"),
		viper.GetString("DATABASE_HOST"),
		viper.GetString("DATABASE_PORT"),
		viper.GetString("DATABASE_NAME"),
		viper.GetString("DATABASE_SSL_MODE"),
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func setupRedis() rueidis.Client {
	redisClient, err := rueidis.NewClient(rueidis.ClientOption{InitAddress: []string{viper.GetString("REDIS_HOST")}})
	if err != nil {
		log.Fatal(err)
	}
	return redisClient
}

func setupRouter(db *sql.DB, redisClient rueidis.Client) *gin.Engine {
	ps := shorturl.NewPostgresPersistentStore(db)
	cs := shorturl.NewRedisCacheStore(redisClient)
	sr := shorturl.NewRepository(ps, cs, &utils.RealTime{})
	sg := &utils.RandomBase62StringGenerator{}
	ss := shorturl.NewService(sr, sg)
	sc := shorturl.NewController(ss)

	r := gin.Default()
	r.Use(middlewares.ErrorHandler())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.POST("/api/v1/urls", sc.CreateShortURL)
	r.GET("/:url", sc.Redirect)

	return r
}
