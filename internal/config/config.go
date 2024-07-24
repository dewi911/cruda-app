package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	DB     Postgres
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
}

type Postgres struct {
	Host     string
	Port     int
	Username string
	Name     string
	SSLMode  string
	Password string
}

func New(folder, filename string) (*Config, error) {
	cfg := new(Config)

	viper.AddConfigPath(folder)
	viper.SetConfigName(filename)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	if err := envconfig.Process("DB", &cfg.DB); err != nil {
		log.Printf("Error processing DB_HOST: %v", err)
		return nil, err
	}

	return cfg, nil
}
