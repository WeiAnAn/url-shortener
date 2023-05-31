package main

import (
	"context"
	"log"
	"time"

	_ "github.com/WeiAnAn/url-shortener/internal/config"
	shorturl "github.com/WeiAnAn/url-shortener/internal/domain/short_url"
	"github.com/WeiAnAn/url-shortener/internal/middlewares"
	"github.com/WeiAnAn/url-shortener/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/rueidis"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	c := setupMongo()
	defer c.Disconnect(context.Background())

	redisClient := setupRedis()
	defer redisClient.Close()

	r := setupRouter(c, redisClient)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func setupRedis() rueidis.Client {
	redisClient, err := rueidis.NewClient(rueidis.ClientOption{InitAddress: []string{viper.GetString("REDIS_HOST")}})
	if err != nil {
		log.Fatal(err)
	}
	return redisClient
}

func setupRouter(c *mongo.Client, redisClient rueidis.Client) *gin.Engine {
	ps := shorturl.NewMongoPersistentStore(c, "short_urls")
	cs := shorturl.NewRedisCacheStore(redisClient)
	sr := shorturl.NewRepository(ps, cs, &utils.RealTime{})
	sg := &utils.RandomBase62StringGenerator{}
	ss := shorturl.NewService(sr, sg)
	sc := shorturl.NewController(ss, viper.GetString("BASE_URL"))

	r := gin.Default()
	r.Use(middlewares.ErrorHandler())

	r.POST("/api/v1/urls", sc.CreateShortURL)
	r.GET("/:url", sc.Redirect)

	return r
}

func setupMongo() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(viper.GetString("MONGODB_URI")))
	if err != nil {
		log.Fatal(err)
	}

	return client
}
