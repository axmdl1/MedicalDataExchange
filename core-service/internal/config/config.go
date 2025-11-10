package config

import "github.com/spf13/viper"

type Config struct {
	REST struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"rest"`
	GRPC struct {
		DataTransfer string `mapstructure:"data_transfer"`
	} `mapstructure:"grpc"`
}

func Load() Config {
	viper.SetConfigName("local")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	viper.SetDefault("rest.port", 8080)
	viper.SetDefault("grpc.data_transfer", "localhost:50052")

	_ = viper.ReadInConfig()

	var c Config
	_ = viper.Unmarshal(&c)
	return c
}
