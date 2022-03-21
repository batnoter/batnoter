package user

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

//go:generate mockgen -source=repo.go -package=user -destination=mock_repo.go
type Repo interface {
	Get(userId uint) (User, error)
	GetByEmail(email string) (User, error)
	Save(user User) error
	Delete(userId uint) error
}

type repoImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repo {
	return &repoImpl{
		db: db,
	}
}

func (r *repoImpl) Get(userId uint) (User, error) {
	var user User
	if err := r.db.Where("id = ?", userId).First(&user).Error; err != nil {
		return user, errors.Wrap(err, "failed to retrieve user from database")
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
		return user, errors.Wrap(err, "failed to retrieve user from database")
	}
	return user, nil
}

func (r *repoImpl) Save(user User) error {
	if err := r.db.Save(&user).Error; err != nil {
		return errors.Wrap(err, "failed to store user to database")
	}
	return nil
}

func (r *repoImpl) Delete(userId uint) error {
	var user User
	if err := r.db.Where("id = ?", userId).Delete(&user).Error; err != nil {
		return errors.Wrap(err, "failed to delete user from database")
	}
	return nil
}
