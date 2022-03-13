package db

import (
	"fmt"

	"github.com/vivekweb2013/gitnoter/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(config config.Database) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		config.Host, config.Username, config.Password, config.DBName, config.Port)

	gormConfig := &gorm.Config{}
	if config.Debug {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DriverName: config.DriverName,
		DSN:        dsn,
	}), gormConfig)
	if err != nil {
		return nil, err
	}

	return db, nil
}
