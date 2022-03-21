package preference

import "gorm.io/gorm"

type DefaultRepo struct {
	gorm.Model
	UserID uint

	Name          string
	Visibility    string
	DefaultBranch string
}
