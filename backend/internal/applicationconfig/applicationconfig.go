package applicationconfig

import (
	"github.com/vivekweb2013/gitnoter/internal/config"
	"gorm.io/gorm"
)

type ApplicationConfig struct {
	Config config.Config
	DB     *gorm.DB
}

func NewApplicationConfig(config config.Config, db *gorm.DB) *ApplicationConfig {
	return &ApplicationConfig{
		Config: config,
		DB:     db,
	}
}
