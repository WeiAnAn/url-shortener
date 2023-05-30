package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	_ "github.com/WeiAnAn/url-shortener/internal/config"
	shorturl "github.com/WeiAnAn/url-shortener/internal/domain/short_url"
	"github.com/WeiAnAn/url-shortener/internal/middlewares"
	"github.com/WeiAnAn/url-shortener/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/rueidis"
	"github.com/spf13/viper"
)

func main() {
	pool := setupDB()
	defer pool.Close()

	redisClient := setupRedis()
	defer redisClient.Close()

	r := setupRouter(pool, redisClient)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func setupDB() *pgxpool.Pool {
	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		viper.GetString("DATABASE_USER"),
		viper.GetString("DATABASE_PASSWORD"),
		viper.GetString("DATABASE_HOST"),
		viper.GetString("DATABASE_PORT"),
		viper.GetString("DATABASE_NAME"),
		viper.GetString("DATABASE_SSL_MODE"),
	)
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}
	return pool
}

func setupRedis() rueidis.Client {
	redisClient, err := rueidis.NewClient(rueidis.ClientOption{InitAddress: []string{viper.GetString("REDIS_HOST")}})
	if err != nil {
		log.Fatal(err)
	}
	return redisClient
}

func setupRouter(pool *pgxpool.Pool, redisClient rueidis.Client) *gin.Engine {
	ps := shorturl.NewPgxPersistentStore(pool)
	cs := shorturl.NewRedisCacheStore(redisClient)
	sr := shorturl.NewRepository(ps, cs, &utils.RealTime{})
	sg := &utils.RandomBase62StringGenerator{}
	ss := shorturl.NewService(sr, sg)
	sc := shorturl.NewController(ss, viper.GetString("BASE_URL"))

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
