package preference

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

//go:generate mockgen -source=repo.go -package=preference -destination=mock_repo.go
type Repo interface {
	Save(defaultRepo DefaultRepo) error
	GetByUserID(userID uint) (DefaultRepo, error)
}

type repoImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repo {
	return &repoImpl{
		db: db,
	}
}

func (r *repoImpl) GetByUserID(userID uint) (DefaultRepo, error) {
	var defaultRepo DefaultRepo
	err := r.db.Where("user_id = ?", userID).First(&defaultRepo).Error
	if err == gorm.ErrRecordNotFound {
		return DefaultRepo{}, nil
	}
	if err != nil {
		return defaultRepo, errors.Wrap(err, "retrieving user's default repo from database failed")
	}
	return defaultRepo, nil
}

func (r *repoImpl) Save(defaultRepo DefaultRepo) error {
	if err := r.db.Save(&defaultRepo).Error; err != nil {
		return errors.Wrap(err, "storing user's default repo to database failed")
	}
	return nil
}
