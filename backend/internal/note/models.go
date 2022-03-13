package note

import (
	"gorm.io/gorm"
)

type Note struct {
	gorm.Model
	Email   string
	Title   string
	Content string
}
