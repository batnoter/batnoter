package user

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

//go:generate mockgen -source=repo.go -package=user -destination=mock_repo.go
type Repo interface {
	Get(userID uint) (User, error)
	GetByEmail(email string) (User, error)
	Save(user User) (uint, error)
	Delete(userID uint) error
}

type repoImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repo {
	return &repoImpl{
		db: db,
	}
}

func (r *repoImpl) Get(userID uint) (User, error) {
	var user User
	if err := r.db.Where("id = ?", userID).Preload("DefaultRepo").First(&user).Error; err != nil {
		return user, errors.Wrap(err, "retrieving user from database failed")
	}
	return user, nil
}

func (r *repoImpl) GetByEmail(email string) (User, error) {
	var user User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return User{}, nil
	}
	if err != nil {
		return user, errors.Wrap(err, "retrieving user from database failed")
	}
	return user, nil
}

func (r *repoImpl) Save(user User) (uint, error) {
	if err := r.db.Save(&user).Error; err != nil {
		return 0, errors.Wrap(err, "storing user to database failed")
	}
	return user.ID, nil
}

func (r *repoImpl) Delete(userID uint) error {
	var user User
	if err := r.db.Where("id = ?", userID).Delete(&user).Error; err != nil {
		return errors.Wrap(err, "deleting user from database failed")
	}
	return nil
}
