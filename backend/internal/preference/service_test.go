package preference

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const (
	userID = uint(1001)
)

func TestGetByUserID(t *testing.T) {
	t.Run("should retrieve the default repo preference associated with user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := NewMockRepo(ctrl)
		d := DefaultRepo{}

		service := NewService(mockRepo)
		mockRepo.EXPECT().GetByUserID(userID).Return(d, nil)

		_, err := service.GetByUserID(userID)
		assert.NoError(t, err)
	})

	t.Run("should return error when retrieving repo associated with user fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := NewMockRepo(ctrl)

		service := NewService(mockRepo)
		mockRepo.EXPECT().GetByUserID(userID).Return(DefaultRepo{}, errors.New("some error"))

		_, err := service.GetByUserID(userID)
		assert.Error(t, err)
	})
}
func TestSave(t *testing.T) {
	t.Run("should save a valid repo as user's default repo preference", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := NewMockRepo(ctrl)
		d := DefaultRepo{
			UserID:        userID,
			Name:          "test-repo",
			Visibility:    "private",
			DefaultBranch: "main",
		}

		service := NewService(mockRepo)
		mockRepo.EXPECT().Save(d).Return(nil)

		err := service.Save(d)
		assert.NoError(t, err)
	})

	t.Run("should return error when saving repo preference fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := NewMockRepo(ctrl)
		d := DefaultRepo{}

		service := NewService(mockRepo)
		mockRepo.EXPECT().Save(d).Return(errors.New("some error"))

		err := service.Save(d)
		assert.Error(t, err)
	})
}
