package httpservice

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vivekweb2013/gitnoter/internal/note"
	"gorm.io/gorm"
)

const (
	noteId = 123
)

func TestGetNote(t *testing.T) {
	t.Run("should return a note when the get request is valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := note.NewMockService(ctrl)

		router := gin.Default()
		n := note.Note{
			Model: gorm.Model{
				ID:        uint(noteId),
				CreatedAt: time.Date(2022, 01, 15, 11, 41, 29, 0, time.UTC),
				UpdatedAt: time.Date(2022, 01, 15, 11, 41, 29, 0, time.UTC),
			},
			Email:   "test@example.com",
			Title:   "Sample Note",
			Content: "This is a sample note!",
		}
		mockService.EXPECT().Get(noteId).Return(n, nil)
		handler := NewNoteHandler(mockService)

		router.GET("/api/v1/note/:id", handler.GetNote)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/note/%s", strconv.Itoa(noteId)), nil)

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.JSONEq(t, `{"id":123, "created_at":"2022-01-15T11:41:29Z", "updated_at":"2022-01-15T11:41:29Z", "email":"test@example.com", "title":"Sample Note", "content":"This is a sample note!"}`, response.Body.String())
	})

	t.Run("should return bad request error if request validation fails", func(t *testing.T) {
		tests := []struct {
			name             string
			noteIdParam      string
			expectedResponse string
		}{
			{
				name:             "with non integer note id",
				noteIdParam:      "abc",
				expectedResponse: `{"code": "validation_failed", "message": "must be an integer number"}`,
			},
		}
		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				mockService := note.NewMockService(ctrl)
				handler := NewNoteHandler(mockService)

				router := gin.Default()
				router.GET("/api/v1/note/:id", handler.GetNote)
				response := httptest.NewRecorder()
				req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/note/%s", test.noteIdParam), nil)

				router.ServeHTTP(response, req)
				assert.Equal(t, http.StatusBadRequest, response.Code)
				assert.JSONEq(t, test.expectedResponse, response.Body.String())
			})
		}

	})

}
