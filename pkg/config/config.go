package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Log      LogConfig
}

type ServerConfig struct {
	Port int
	Host string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret string
	Expire int
}

type LogConfig struct {
	Level string
	File  string
}

var cfg *Config

func Init() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../configs")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	cfg = &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return err
	}

	return nil
}

func GetConfig() *Config {
	return cfg
}
