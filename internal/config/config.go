package config

import (
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("DATABASE_NAME", "short_urls")
	viper.SetDefault("DATABASE_HOST", "localhost")
	viper.SetDefault("DATABASE_USER", "short_url")
	viper.SetDefault("DATABASE_PASSWORD", "")
	viper.SetDefault("DATABASE_PORT", "5432")
	viper.SetDefault("DATABASE_SSL_MODE", "disable")
	viper.SetDefault("REDIS_HOST", "localhost:6379")
	viper.AllowEmptyEnv(true)
	viper.AutomaticEnv()
}
