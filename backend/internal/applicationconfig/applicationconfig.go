package applicationconfig

import (
	"github.com/vivekweb2013/gitnoter/internal/auth"
	"github.com/vivekweb2013/gitnoter/internal/config"
	"github.com/vivekweb2013/gitnoter/internal/github"
	"github.com/vivekweb2013/gitnoter/internal/preference"
	"github.com/vivekweb2013/gitnoter/internal/user"
	"golang.org/x/oauth2"
	gh "golang.org/x/oauth2/github"
	"gorm.io/gorm"
)

type ApplicationConfig struct {
	Config            config.Config
	DB                *gorm.DB
	OAuth2Config      oauth2.Config
	AuthService       auth.Service
	UserService       user.Service
	PreferenceService preference.Service
	GithubService     github.Service
}

func NewApplicationConfig(config config.Config, db *gorm.DB) *ApplicationConfig {
	// create github oauth2 config
	oauth2Config := oauth2.Config{
		RedirectURL:  config.OAuth2.Github.RedirectURL,
		ClientID:     config.OAuth2.Github.ClientID,
		ClientSecret: config.OAuth2.Github.ClientSecret,
		Scopes:       []string{"read:user", "user:email", "repo"},
		Endpoint:     gh.Endpoint,
	}

	authService := auth.NewService(auth.TokenConfig{
		SecretKey: config.App.SecretKey,
		Issuer:    "https://gitnoter.com",
	})
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	preferenceRepo := preference.NewRepository(db)
	preferenceService := preference.NewService(preferenceRepo)

	githubClientBuilder := github.NewClientBuilder(&oauth2Config)
	githubService := github.NewService(githubClientBuilder)

	return &ApplicationConfig{
		Config:            config,
		DB:                db,
		OAuth2Config:      oauth2Config,
		AuthService:       authService,
		UserService:       userService,
		PreferenceService: preferenceService,
		GithubService:     githubService,
	}
}
