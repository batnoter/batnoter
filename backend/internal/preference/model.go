package preference

import "gorm.io/gorm"

// DefaultRepo represents an entity model used to store & retrieve user's default repo to/from database.
type DefaultRepo struct {
	gorm.Model
	UserID uint

	Name          string
	Visibility    string
	DefaultBranch string
}
