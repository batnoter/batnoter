package note

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

const (
	userID = uint(1234)
	noteID = 123
)

func TestServiceImpl_Get(t *testing.T) {
	t.Run("should retrieve a note", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := NewMockRepo(ctrl)
		n := Note{}

		service := NewService(mockRepo)
		mockRepo.EXPECT().Get(noteID).Return(n, nil)

		_, err := service.Get(noteID)
		assert.NoError(t, err)
	})

	t.Run("should return error when retrieving note fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := NewMockRepo(ctrl)
		n := Note{}

		service := NewService(mockRepo)
		mockRepo.EXPECT().Get(gomock.Any()).Return(n, errors.New("some error"))

		_, err := service.Get(noteID)
		assert.Error(t, err)
	})
}

func TestServiceImpl_Save(t *testing.T) {
	t.Run("should save a valid note", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := NewMockRepo(ctrl)
		n := Note{
			Model: gorm.Model{
				ID:        1,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
				DeletedAt: gorm.DeletedAt{
					Time:  time.Time{},
					Valid: false,
				},
			},
			UserID:  userID,
			Title:   "Sample Note",
			Content: "This is a sample note",
		}

		service := NewService(mockRepo)
		mockRepo.EXPECT().Save(n).Return(nil)

		err := service.Save(n)
		assert.NoError(t, err)
	})

	t.Run("should return error when saving note fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := NewMockRepo(ctrl)
		n := Note{}

		service := NewService(mockRepo)
		mockRepo.EXPECT().Save(gomock.Any()).Return(errors.New("some error"))

		err := service.Save(n)
		assert.Error(t, err)
	})
}

func TestServiceImpl_Delete(t *testing.T) {
	t.Run("should delete a note", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := NewMockRepo(ctrl)

		service := NewService(mockRepo)
		mockRepo.EXPECT().Delete(noteID).Return(nil)

		err := service.Delete(noteID)
		assert.NoError(t, err)
	})

	t.Run("should return error when deleting note fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := NewMockRepo(ctrl)

		service := NewService(mockRepo)
		mockRepo.EXPECT().Delete(gomock.Any()).Return(errors.New("some error"))

		err := service.Delete(noteID)
		assert.Error(t, err)
	})
}
