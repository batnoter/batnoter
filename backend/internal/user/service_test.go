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
	userID    = uint(1234)
	email     = "john.doe@example.com"
	name      = "John Doe"
	location  = "New York"
	avatarURL = "http://example.com/avatar"
)

func TestServiceImpl_Get(t *testing.T) {
	t.Run("should retrieve a user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := NewMockRepo(ctrl)
		n := User{}

		service := NewService(mockRepo)
		mockRepo.EXPECT().Get(userID).Return(n, nil)

		_, err := service.Get(userID)
		assert.NoError(t, err)
	})

	t.Run("should return error when retrieving user fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := NewMockRepo(ctrl)
		n := User{}

		service := NewService(mockRepo)
		mockRepo.EXPECT().Get(gomock.Any()).Return(n, errors.New("some error"))

		_, err := service.Get(userID)
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
				ID:        userID,
				CreatedAt: date,
				UpdatedAt: date,
			},
			Email:          email,
			Name:           name,
			Location:       location,
			AvatarURL:      avatarURL,
			GithubID:       12345,
			GithubUsername: "johndoe",
			GithubToken:    "cOLcG0c0gJn2iv25hrTPY3A%3D3D",
			DisabledAt:     nil,
		}

		service := NewService(mockRepo)
		mockRepo.EXPECT().Save(n).Return(userID, nil)

		id, err := service.Save(n)
		assert.Equal(t, userID, id)
		assert.NoError(t, err)
	})

	t.Run("should return error when saving user fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := NewMockRepo(ctrl)
		n := User{}

		service := NewService(mockRepo)
		mockRepo.EXPECT().Save(gomock.Any()).Return(uint(0), errors.New("some error"))

		_, err := service.Save(n)
		assert.Error(t, err)
	})
}

func TestServiceImpl_Delete(t *testing.T) {
	t.Run("should delete a user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := NewMockRepo(ctrl)

		service := NewService(mockRepo)
		mockRepo.EXPECT().Delete(userID).Return(nil)

		err := service.Delete(userID)
		assert.NoError(t, err)
	})

	t.Run("should return error when deleting user fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := NewMockRepo(ctrl)

		service := NewService(mockRepo)
		mockRepo.EXPECT().Delete(gomock.Any()).Return(errors.New("some error"))

		err := service.Delete(userID)
		assert.Error(t, err)
	})
}
