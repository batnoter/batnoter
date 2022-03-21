package note

import (
	"gorm.io/gorm"
)

type Note struct {
	gorm.Model
	UserID  uint
	Title   string
	Content string
}
