package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	SERVER_PORT          string `mapstructure:"SERVER_PORT"`
	SERVER_NAME          string `mapstructure:"SERVER_NAME"`
	SERVER_ENV           string `mapstructure:"SERVER_ENV"`
	KAFKA_BROKERS        string `mapstructure:"KAFKA_BROKERS"`
	KAFKA_TOPIC_FOLLOWS string `mapstructure:"KAFKA_TOPIC_FOLLOWS"`
	KAFKA_TOPIC_LIKES    string `mapstructure:"KAFKA_TOPIC_LIKES"`
	KAFKA_TOPIC_COMMENTS  string `mapstructure:"KAFKA_TOPIC_COMMENTS"`
	KAFKA_TOPIC_DIRECTS   string `mapstructure:"KAFKA_TOPIC_DIRECTS"`
	KAFKA_CONSUMER_GROUP  string `mapstructure:"KAFKA_CONSUMER_GROUP"`
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