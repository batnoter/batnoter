package applicationconfig

import (
	"github.com/vivekweb2013/gitnoter/internal/auth"
	"github.com/vivekweb2013/gitnoter/internal/config"
	"github.com/vivekweb2013/gitnoter/internal/note"
	"github.com/vivekweb2013/gitnoter/internal/preference"
	"github.com/vivekweb2013/gitnoter/internal/user"
	"gorm.io/gorm"
)

type ApplicationConfig struct {
	Config            config.Config
	DB                *gorm.DB
	AuthService       auth.Service
	UserService       user.Service
	PreferenceService preference.Service
	NoteService       note.Service
}

func NewApplicationConfig(config config.Config, db *gorm.DB) *ApplicationConfig {
	authService := auth.NewService(auth.TokenConfig{
		SecretKey: config.App.SecretKey,
		Issuer:    "https://gitnoter.com",
	})
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	preferenceRepo := preference.NewRepository(db)
	preferenceService := preference.NewService(preferenceRepo)

	noteRepo := note.NewRepository(db)
	noteService := note.NewService(noteRepo)

	return &ApplicationConfig{
		Config:            config,
		DB:                db,
		AuthService:       authService,
		UserService:       userService,
		PreferenceService: preferenceService,
		NoteService:       noteService,
	}
}
