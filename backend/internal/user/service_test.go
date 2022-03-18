package user

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const (
	userId = 123
	email  = "john.doe@example.com"
)

func TestServiceImpl_Get(t *testing.T) {
	t.Run("should retrieve a user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := NewMockRepo(ctrl)
		n := User{}

		service := NewService(mockRepo)
		mockRepo.EXPECT().Get(userId).Return(n, nil)

		_, err := service.Get(userId)
		assert.NoError(t, err)
	})

	t.Run("should return error when retrieving user fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := NewMockRepo(ctrl)
		n := User{}

		service := NewService(mockRepo)
		mockRepo.EXPECT().Get(gomock.Any()).Return(n, errors.New("some error"))

		_, err := service.Get(userId)
		assert.Error(t, err)
	})
}

func TestServiceImpl_GetByEmail(t *testing.T) {
	t.Run("should retrieve a user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := NewMockRepo(ctrl)
		n := User{}

		service := NewService(mockRepo)
		mockRepo.EXPECT().GetByEmail(email).Return(n, nil)

		_, err := service.GetByEmail(email)
		assert.NoError(t, err)
	})

	t.Run("should return error when retrieving user fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := NewMockRepo(ctrl)
		n := User{}

		service := NewService(mockRepo)
		mockRepo.EXPECT().GetByEmail(gomock.Any()).Return(n, errors.New("some error"))

		_, err := service.GetByEmail(email)
		assert.Error(t, err)
	})
}

func TestServiceImpl_Save(t *testing.T) {
	t.Run("should save a valid user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := NewMockRepo(ctrl)
		date := time.Date(2022, 02, 12, 4, 45, 55, 0, time.UTC)
		n := User{
			Model: gorm.Model{
				ID:        1,
				CreatedAt: date,
				UpdatedAt: date,
			},
			Email:          "john.doe@example.com",
			Name:           "John Doe",
			Location:       "New York",
			AvatarURL:      "http://example.com/avatar",
			GithubID:       12345,
			GithubUsername: "johndoe",
			GithubToken:    "cOLcG0c0gJn2iv25hrTPY3A%3D3D",
			DisabledAt:     nil,
		}

		service := NewService(mockRepo)
		mockRepo.EXPECT().Save(n).Return(nil)

		err := service.Save(n)
		assert.NoError(t, err)
	})

	t.Run("should return error when saving user fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := NewMockRepo(ctrl)
		n := User{}

		service := NewService(mockRepo)
		mockRepo.EXPECT().Save(gomock.Any()).Return(errors.New("some error"))

		err := service.Save(n)
		assert.Error(t, err)
	})
}

func TestServiceImpl_Delete(t *testing.T) {
	t.Run("should delete a user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := NewMockRepo(ctrl)

		service := NewService(mockRepo)
		mockRepo.EXPECT().Delete(userId).Return(nil)

		err := service.Delete(userId)
		assert.NoError(t, err)
	})

	t.Run("should return error when deleting user fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := NewMockRepo(ctrl)

		service := NewService(mockRepo)
		mockRepo.EXPECT().Delete(gomock.Any()).Return(errors.New("some error"))

		err := service.Delete(userId)
		assert.Error(t, err)
	})
}
