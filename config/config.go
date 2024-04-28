package config

import (
	"os"
)

type Config struct {
	DB *DBConfig
}

type DBConfig struct {
	Dialect  string
	Host     string
	Port     int
	Username string
	Password string
	Name     string
	Charset  string
}

func GetConfig() *Config {
	return &Config{
		DB: &DBConfig{
			Dialect:  "mysql",
			Host:     "viaduct.proxy.rlwy.net",
			Port:     29329, // Assign port obtained from environment variable
			Username: "root",
			Password: "OxDiuBlXdQCkjPzWJABTljdaHYUwOIdG",
			Name:     "railway",
			Charset:  "utf8",
		},
	}
}
