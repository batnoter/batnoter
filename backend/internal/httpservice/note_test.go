package httpservice

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vivekweb2013/gitnoter/internal/note"
	"gorm.io/gorm"
)

const (
	noteId                = 123
	title                 = "Sample Note"
	content               = "This is a sample note!"
	internalServerErrJson = `{"code":"internal_server_error", "message":"something went wrong. contact support"}`
)

func TestGetNote(t *testing.T) {
	t.Run("should return a note when the get request is valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := note.NewMockService(ctrl)

		router := gin.Default()
		gin.SetMode(gin.TestMode)
		n := note.Note{
			Model: gorm.Model{
				ID:        uint(noteId),
				CreatedAt: time.Date(2022, 01, 15, 11, 41, 29, 0, time.UTC),
				UpdatedAt: time.Date(2022, 01, 15, 11, 41, 29, 0, time.UTC),
			},
			UserID:  userID,
			Title:   title,
			Content: content,
		}
		mockService.EXPECT().Get(noteId).Return(n, nil)
		handler := NewNoteHandler(mockService)

		router.GET("/api/v1/note/:id", handler.GetNote)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/note/%s", strconv.Itoa(noteId)), nil)

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.JSONEq(t, fmt.Sprintf(`{"id":123, "created_at":"2022-01-15T11:41:29Z", "updated_at":"2022-01-15T11:41:29Z", "title":"%s", "content":"%s"}`, title, content), response.Body.String())
	})

	t.Run("should return internal server error when retrieving a note fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := note.NewMockService(ctrl)

		router := gin.Default()
		gin.SetMode(gin.TestMode)
		n := note.Note{}
		mockService.EXPECT().Get(noteId).Return(n, errors.New("some error"))
		handler := NewNoteHandler(mockService)

		router.GET("/api/v1/note/:id", handler.GetNote)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/note/%s", strconv.Itoa(noteId)), nil)

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.JSONEq(t, internalServerErrJson, response.Body.String())
	})

	t.Run("should return bad request error when note id param is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := note.NewMockService(ctrl)
		handler := NewNoteHandler(mockService)

		router := gin.Default()
		gin.SetMode(gin.TestMode)
		router.GET("/api/v1/note/:id", handler.GetNote)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/note/%s", "abc"), nil)

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.JSONEq(t, `{"code": "validation_failed", "message": "id: must be an integer number"}`, response.Body.String())
	})
}

func TestCreateNote(t *testing.T) {
	t.Run("should create a note when the create request is valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := note.NewMockService(ctrl)

		router := gin.Default()
		gin.SetMode(gin.TestMode)
		n := note.Note{
			UserID:  userID,
			Title:   title,
			Content: content,
		}
		mockService.EXPECT().Save(n).Return(nil)
		handler := NewNoteHandler(mockService)
		claims := jwt.MapClaims{"sub": strconv.FormatUint(uint64(userID), 10)}

		router.POST("/api/v1/note", func(c *gin.Context) { c.Set("claims", claims) }, handler.CreateNote)
		response := httptest.NewRecorder()
		validNoteJson := getNotePayloadJsonString(title, content)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/note", strings.NewReader(validNoteJson))

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "", response.Body.String())
	})

	t.Run("should return internal server error when creating a note fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := note.NewMockService(ctrl)

		router := gin.Default()
		gin.SetMode(gin.TestMode)
		mockService.EXPECT().Save(gomock.Any()).Return(errors.New("some error"))
		handler := NewNoteHandler(mockService)
		claims := jwt.MapClaims{"sub": strconv.FormatUint(uint64(userID), 10)}

		router.POST("/api/v1/note", func(c *gin.Context) { c.Set("claims", claims) }, handler.CreateNote)
		response := httptest.NewRecorder()
		validNoteJson := getNotePayloadJsonString(title, content)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/note", strings.NewReader(validNoteJson))

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.JSONEq(t, internalServerErrJson, response.Body.String())
	})

	t.Run("should return bad request error when note payload validation fails", func(t *testing.T) {
		tests := getNotePayloadValidations()

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				mockService := note.NewMockService(ctrl)
				handler := NewNoteHandler(mockService)
				claims := jwt.MapClaims{"sub": strconv.FormatUint(uint64(userID), 10)}

				router := gin.Default()
				gin.SetMode(gin.TestMode)
				router.POST("/api/v1/note", func(c *gin.Context) { c.Set("claims", claims) }, handler.CreateNote)
				response := httptest.NewRecorder()
				req, _ := http.NewRequest(http.MethodPost, "/api/v1/note", strings.NewReader(test.notePayloadJson))

				router.ServeHTTP(response, req)
				assert.Equal(t, http.StatusBadRequest, response.Code)
				assert.JSONEq(t, test.expectedResponse, response.Body.String())
			})
		}
	})
}

func TestUpdateNote(t *testing.T) {
	t.Run("should update a note when the update request is valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := note.NewMockService(ctrl)

		router := gin.Default()
		gin.SetMode(gin.TestMode)
		n := note.Note{
			Model: gorm.Model{
				ID: uint(noteId),
			},
			UserID:  userID,
			Title:   title,
			Content: content,
		}
		mockService.EXPECT().Save(n).Return(nil)
		handler := NewNoteHandler(mockService)
		claims := jwt.MapClaims{"sub": strconv.FormatUint(uint64(userID), 10)}

		router.PUT("/api/v1/note/:id", func(c *gin.Context) { c.Set("claims", claims) }, handler.UpdateNote)
		response := httptest.NewRecorder()
		validNoteJson := getNotePayloadJsonString(title, content)
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/note/%s", strconv.Itoa(noteId)), strings.NewReader(validNoteJson))

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "", response.Body.String())
	})

	t.Run("should return internal server error when updating a note fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := note.NewMockService(ctrl)

		router := gin.Default()
		gin.SetMode(gin.TestMode)
		mockService.EXPECT().Save(gomock.Any()).Return(errors.New("some error"))
		handler := NewNoteHandler(mockService)
		claims := jwt.MapClaims{"sub": strconv.FormatUint(uint64(userID), 10)}

		router.PUT("/api/v1/note/:id", func(c *gin.Context) { c.Set("claims", claims) }, handler.UpdateNote)
		response := httptest.NewRecorder()
		validNoteJson := getNotePayloadJsonString(title, content)
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/note/%s", strconv.Itoa(noteId)), strings.NewReader(validNoteJson))

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.JSONEq(t, internalServerErrJson, response.Body.String())
	})

	t.Run("should return bad request error when note id param is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := note.NewMockService(ctrl)

		router := gin.Default()
		gin.SetMode(gin.TestMode)
		handler := NewNoteHandler(mockService)
		claims := jwt.MapClaims{"sub": strconv.FormatUint(uint64(userID), 10)}

		router.PUT("/api/v1/note/:id", func(c *gin.Context) { c.Set("claims", claims) }, handler.UpdateNote)
		response := httptest.NewRecorder()
		validNoteJson := getNotePayloadJsonString(title, content)
		req, _ := http.NewRequest(http.MethodPut, "/api/v1/note/abc", strings.NewReader(validNoteJson))

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.JSONEq(t, `{"code": "validation_failed", "message": "id: must be an integer number"}`, response.Body.String())
	})

	t.Run("should return bad request error when note payload validation fails", func(t *testing.T) {
		tests := getNotePayloadValidations()

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				mockService := note.NewMockService(ctrl)
				handler := NewNoteHandler(mockService)
				claims := jwt.MapClaims{"sub": strconv.FormatUint(uint64(userID), 10)}

				router := gin.Default()
				gin.SetMode(gin.TestMode)
				router.PUT("/api/v1/note/:id", func(c *gin.Context) { c.Set("claims", claims) }, handler.UpdateNote)
				response := httptest.NewRecorder()
				req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/note/%s", strconv.Itoa(noteId)), strings.NewReader(test.notePayloadJson))

				router.ServeHTTP(response, req)
				assert.Equal(t, http.StatusBadRequest, response.Code)
				assert.JSONEq(t, test.expectedResponse, response.Body.String())
			})
		}
	})
}

func TestDeleteNote(t *testing.T) {
	t.Run("should delete a note when the delete request is valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := note.NewMockService(ctrl)

		router := gin.Default()
		gin.SetMode(gin.TestMode)
		mockService.EXPECT().Delete(noteId).Return(nil)
		handler := NewNoteHandler(mockService)

		router.DELETE("/api/v1/note/:id", handler.DeleteNote)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/note/%s", strconv.Itoa(noteId)), nil)

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "", response.Body.String())
	})

	t.Run("should return internal server error when deleting a note fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := note.NewMockService(ctrl)

		router := gin.Default()
		gin.SetMode(gin.TestMode)
		mockService.EXPECT().Delete(noteId).Return(errors.New("some error"))
		handler := NewNoteHandler(mockService)

		router.DELETE("/api/v1/note/:id", handler.DeleteNote)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/note/%s", strconv.Itoa(noteId)), nil)

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.JSONEq(t, internalServerErrJson, response.Body.String())
	})

	t.Run("should return bad request error when note id param is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := note.NewMockService(ctrl)
		handler := NewNoteHandler(mockService)

		router := gin.Default()
		gin.SetMode(gin.TestMode)
		router.DELETE("/api/v1/note/:id", handler.DeleteNote)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/note/%s", "abc"), nil)

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.JSONEq(t, `{"code": "validation_failed", "message": "id: must be an integer number"}`, response.Body.String())
	})
}

func getNotePayloadValidations() []struct {
	name             string
	notePayloadJson  string
	expectedResponse string
} {
	return []struct {
		name             string
		notePayloadJson  string
		expectedResponse string
	}{
		{
			name:             "with blank title",
			notePayloadJson:  getNotePayloadJsonString("", content),
			expectedResponse: `{"code":"validation_failed", "message":"title: cannot be blank."}`,
		},
		{
			name:             "with title length more than 255 chars",
			notePayloadJson:  getNotePayloadJsonString(randString(256), content),
			expectedResponse: `{"code":"validation_failed", "message":"title: the length must be between 1 and 255."}`,
		},
		{
			name:             "with blank content",
			notePayloadJson:  getNotePayloadJsonString(title, ""),
			expectedResponse: `{"code":"validation_failed", "message":"content: cannot be blank."}`,
		},
		{
			name:             "with content length more than 5000 chars",
			notePayloadJson:  getNotePayloadJsonString(title, randString(5001)),
			expectedResponse: `{"code":"validation_failed", "message":"content: the length must be between 1 and 5000."}`,
		},
	}
}

func getNotePayloadJsonString(title string, content string) string {
	return fmt.Sprintf(`{"title":"%s", "content":"%s"}`, title, content)
}

func randString(n int) string {
	rand.Seed(time.Now().UnixNano())

	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
