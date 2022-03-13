package applicationconfig

import (
	"github.com/vivekweb2013/gitnoter/internal/config"
	"github.com/vivekweb2013/gitnoter/internal/note"
	"gorm.io/gorm"
)

type ApplicationConfig struct {
	Config      config.Config
	DB          *gorm.DB
	NoteService note.Service
}

func NewApplicationConfig(config config.Config, db *gorm.DB) *ApplicationConfig {
	noteRepo := note.NewRepository(db)
	noteService := note.NewService(noteRepo)
	return &ApplicationConfig{
		Config:      config,
		DB:          db,
		NoteService: noteService,
	}
}
