package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
	Database struct {
		URL string `mapstructure:"url"`
	} `mapstructure:"database"`
	Auth struct {
		JWTSecret string `mapstructure:"jwt_secret"`
	} `mapstructure:"auth"`
}

func Load() Config {
	viper.SetConfigName("local")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	// дефолты
	viper.SetDefault("server.port", 50052)
	viper.SetDefault("database.url", "host=localhost user=postgres password=postgres dbname=medical_exchange port=5432 sslmode=disable")
	viper.SetDefault("auth.jwt_secret", "your-secret-key-change-in-production")

	_ = viper.ReadInConfig() // игнорим ошибку — упадёт только если структура не совпадает

	var c Config
	_ = viper.Unmarshal(&c)
	return c
}
