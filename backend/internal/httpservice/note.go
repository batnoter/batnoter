package httpservice

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/sirupsen/logrus"
)

type NoteHandler struct {
	// TODO: inject service
}

type NoteResponse struct {
	ID        uint       `json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	Email     string     `json:"email"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
}

func NewNoteHandler() *NoteHandler {
	return &NoteHandler{}
}

func (n *NoteHandler) GetNote(c *gin.Context) {
	id := c.Param("id")
	if err := validation.Validate(id, validation.Required, is.Int); err != nil {
		abortRequestWithError(c, NewAppError(ErrorCodeValidationFailed, err.Error()))
		return
	}
	logrus.Infof("request for note id: %s", id)

	// FIXME: Retrieve the note from db
	c.JSON(http.StatusOK, NoteResponse{
		ID:        0,
		CreatedAt: &time.Time{},
		UpdatedAt: &time.Time{},
		DeletedAt: &time.Time{},
		Email:     "",
		Title:     "",
		Content:   "",
	})
}
