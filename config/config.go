package config

import (
	"github.com/joho/godotenv"
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
	godotenv.Load(".env")

	return &Config{
		DB: &DBConfig{
			Dialect:  "mysql",
			Host:     os.Getenv("HOST"),
			Port:     29329,
			Username: "root",
			Password: "OxDiuBlXdQCkjPzWJABTljdaHYUwOIdG",
			Name:     "railway",
			Charset:  "utf8",
		},
	}
}
