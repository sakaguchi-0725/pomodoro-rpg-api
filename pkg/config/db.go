package config

import "os"

type DBConfig struct {
	Port     string
	Host     string
	Name     string
	User     string
	Password string
	SSLMode  string
}

func newDBConfig() *DBConfig {
	return &DBConfig{
		Port:     os.Getenv("DB_PORT"),
		Host:     os.Getenv("DB_HOST"),
		Name:     os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		SSLMode:  os.Getenv("SSL_MODE"),
	}
}
