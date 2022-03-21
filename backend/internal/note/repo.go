package note

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

//go:generate mockgen -source=repo.go -package=note -destination=mock_repo.go
type Repo interface {
	GetAll(userID uint) ([]Note, error)
	Get(noteId int) (Note, error)
	Save(note Note) error
	Delete(noteId int) error
}

type repoImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repo {
	return &repoImpl{
		db: db,
	}
}

func (r *repoImpl) GetAll(userID uint) ([]Note, error) {
	var notes []Note
	if err := r.db.Where("user_id = ?", userID).Find(&notes).Error; err != nil {
		return notes, errors.Wrap(err, "retrieving notes from database failed")
	}
	return notes, nil
}

func (r *repoImpl) Get(noteId int) (Note, error) {
	var note Note
	if err := r.db.Where("id = ?", noteId).First(&note).Error; err != nil {
		return note, errors.Wrap(err, "retrieving note from database failed")
	}
	return note, nil
}

func (r *repoImpl) Save(note Note) error {
	if err := r.db.Save(&note).Error; err != nil {
		return errors.Wrap(err, "storing note to database failed")
	}
	return nil
}
func (r *repoImpl) Delete(noteId int) error {
	var note Note
	if err := r.db.Where("id = ?", noteId).Delete(&note).Error; err != nil {
		return errors.Wrap(err, "deleting note from database failed")
	}
	return nil
}
