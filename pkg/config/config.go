package config

import (
	"log"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}
}

type Config struct {
	DB  *DBConfig
	AWS *AWS
}

func NewConfig() *Config {
	return &Config{
		DB:  newDBConfig(),
		AWS: newAWSConfig(),
	}
}
