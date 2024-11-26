package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	SERVER_PORT          string `mapstructure:"SERVER_PORT"`
	SERVER_NAME          string `mapstructure:"SERVER_NAME"`
	SERVER_ENV           string `mapstructure:"SERVER_ENV"`
	DB_HOST              string `mapstructure:"DB_HOST"`
	DB_PORT              string `mapstructure:"DB_PORT"`
	DB_USER              string `mapstructure:"DB_USER"`
	DB_PASSWORD          string `mapstructure:"DB_PASSWORD"`
	DB_NAME              string `mapstructure:"DB_NAME"`
	REDIS_HOST           string `mapstructure:"REDIS_HOST"`
	REDIS_PORT           string `mapstructure:"REDIS_PORT"`
	KAFKA_BROKERS        string `mapstructure:"KAFKA_BROKERS"`
	KAFKA_TOPIC          string `mapstructure:"KAFKA_TOPIC"`
	USER_SERVER_NAME     string `mapstructure:"USER_SERVER_NAME"`
	USER_SERVER_PORT     string `mapstructure:"USER_SERVER_PORT"`
	TWEET_SERVER_NAME     string `mapstructure:"TWEET_SERVER_NAME"`
	TWEET_SERVER_PORT     string `mapstructure:"TWEET_SERVER_PORT"`
}

func LoadConfig() (Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	var config Config

	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
