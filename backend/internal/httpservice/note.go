package httpservice

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/sirupsen/logrus"
	"github.com/vivekweb2013/gitnoter/internal/note"
)

type NoteHandler struct {
	noteService note.Service
}

type NoteResponse struct {
	ID        uint       `json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	Email     string     `json:"email"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
}

func NewNoteHandler(noteService note.Service) *NoteHandler {
	return &NoteHandler{noteService: noteService}
}

func (n *NoteHandler) GetNote(c *gin.Context) {
	id := c.Param("id")
	if err := validation.Validate(id, validation.Required, is.Int); err != nil {
		abortRequestWithError(c, NewAppError(ErrorCodeValidationFailed, err.Error()))
		return
	}
	logrus.Infof("request for note id: %s", id)

	noteId, _ := strconv.Atoi(id)
	note, err := n.noteService.Get(noteId)
	if err != nil {
		abortRequestWithError(c, err)
		return
	}
	c.JSON(http.StatusOK, NoteResponse{
		ID:        note.ID,
		CreatedAt: &note.CreatedAt,
		UpdatedAt: &note.UpdatedAt,
		Email:     note.Email,
		Title:     note.Title,
		Content:   note.Content,
	})
}
