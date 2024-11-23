package db

import (
	"fmt"
	"log"

	"pomodoro-rpg-api/pkg/config"

	"github.com/cockroachdb/errors"
	migrate "github.com/rubenv/sql-migrate"
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

func Migration(db *gorm.DB, cfg *config.DBConfig) error {
	sqlDB, err := db.DB()
	if err != nil {
		return errors.WithStack(errors.Wrap(err, "get sql.DB failed"))
	}

	defer sqlDB.Close()

	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}

	n, err := migrate.Exec(sqlDB, "postgres", migrations, migrate.Up)
	if err != nil {
		return errors.WithStack(errors.Wrap(err, "failed to apply migrations"))
	}

	log.Printf("Applied %v migrations!\n", n)
	return nil
}

func genDSN(cfg *config.DBConfig) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)
}
