package user

import (
	"time"

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

	DisabledAt *time.Time
}
