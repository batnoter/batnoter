package httpservice

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	noteId = "123"
)

func TestGetNote(t *testing.T) {
	t.Run("should return a note when the request is valid", func(t *testing.T) {
		router := gin.Default()
		handler := NewNoteHandler()

		router.GET("/api/v1/note/:id", handler.GetNote)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/note/%s", noteId), nil)

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusOK, response.Code)
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
				router := gin.Default()
				handler := NewNoteHandler()

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
