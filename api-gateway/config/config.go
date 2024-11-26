package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	SERVER_PORT          string `mapstructure:"SERVER_PORT"`
	SERVER_NAME          string `mapstructure:"SERVER_NAME"`
	SERVER_ENV           string `mapstructure:"SERVER_ENV"`
	JWT_SECRET           string `mapstructure:"JWT_SECRET"`
	TWEET_SERVER_NAME   string `mapstructure:"TWEET_SERVER_NAME"`
	TWEET_SERVER_PORT   string `mapstructure:"TWEET_SERVER_PORT"`
	LIKE_SERVER_NAME     string `mapstructure:"LIKE_SERVER_NAME"`
	LIKE_SERVER_PORT     string `mapstructure:"LIKE_SERVER_PORT"`
	COMMENT_SERVER_NAME  string `mapstructure:"COMMENT_SERVER_NAME"`
	COMMENT_SERVER_PORT  string `mapstructure:"COMMENT_SERVER_PORT"`
	DIRECT_SERVER_NAME   string `mapstructure:"DIRECT_SERVER_NAME"`
	DIRECT_SERVER_PORT   string `mapstructure:"DIRECT_SERVER_PORT"`
	USER_SERVER_NAME     string `mapstructure:"USER_SERVER_NAME"`
	USER_SERVER_PORT     string `mapstructure:"USER_SERVER_PORT"`
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
