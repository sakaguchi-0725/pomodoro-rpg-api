package db

import (
	"fmt"
	"log"

	"pomodoro-rpg-api/pkg/config"

	"github.com/cockroachdb/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(cfg *config.DBConfig) (*gorm.DB, error) {
	dsn := genDSN(cfg)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	log.Println("Successfully connected to the database!")
	return db, nil
}

func genDSN(cfg *config.DBConfig) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)
}
