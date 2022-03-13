package note

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

//go:generate mockgen -source=repo.go -package=note -destination=mock_repo.go
type Repo interface {
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

func (r *repoImpl) Get(noteId int) (Note, error) {
	var note Note
	if err := r.db.Where("id = ?", noteId).First(&note).Error; err != nil {
		return note, errors.Wrap(err, "failed to retrieve note from database")
	}
	return note, nil
}

func (r *repoImpl) Save(note Note) error {
	if err := r.db.Save(&note).Error; err != nil {
		return errors.Wrap(err, "failed to store note to database")
	}
	return nil
}
func (r *repoImpl) Delete(noteId int) error {
	var note Note
	if err := r.db.Where("id = ?", noteId).Delete(&note).Error; err != nil {
		return errors.Wrap(err, "failed to delete note from database")
	}
	return nil
}
