package config

import (
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("MONGODB_URI", "mongodb://short_url@localhost:27017")
	viper.SetDefault("REDIS_HOST", "localhost:6379")
	viper.SetDefault("BASE_URL", "http://localhost:8080")
	viper.AllowEmptyEnv(true)
	viper.AutomaticEnv()
}
