package user

import (
	"time"

	"github.com/vivekweb2013/gitnoter/internal/preference"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Email          string
	Name           string
	Location       string
	AvatarURL      string
	GithubID       int64
	GithubUsername string
	GithubToken    string
	DisabledAt     *time.Time

	DefaultRepo *preference.DefaultRepo `gorm:"foreignkey:UserID"`
}
