package httpservice

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/sirupsen/logrus"
	"github.com/vivekweb2013/gitnoter/internal/note"
	"gorm.io/gorm"
)

type NoteHandler struct {
	noteService note.Service
}

type NoteRequestPayload struct {
	Email   string `json:"email"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (n NoteRequestPayload) Validate() error {
	return validation.ValidateStruct(&n,
		validation.Field(&n.Email, validation.Required, is.Email),
		validation.Field(&n.Title, validation.Required, validation.Length(1, 255)),
		validation.Field(&n.Content, validation.Required, validation.Length(1, 5000)),
	)
}

type NoteResponsePayload struct {
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
		abortRequestWithError(c, NewAppError(ErrorCodeValidationFailed, fmt.Sprintf("id: %s", err.Error())))
		return
	}
	logrus.WithField("note_id", id).Info("request to retrieve note")

	noteId, _ := strconv.Atoi(id)
	note, err := n.noteService.Get(noteId)
	if err != nil {
		abortRequestWithError(c, err)
		return
	}
	c.JSON(http.StatusOK, NoteResponsePayload{
		ID:        note.ID,
		CreatedAt: &note.CreatedAt,
		UpdatedAt: &note.UpdatedAt,
		Email:     note.Email,
		Title:     note.Title,
		Content:   note.Content,
	})
	logrus.WithField("note_id", id).Info("request to retrieve note successful")
}

func (n *NoteHandler) CreateNote(c *gin.Context) {
	var noteReqPayload NoteRequestPayload
	c.BindJSON(&noteReqPayload)
	if err := noteReqPayload.Validate(); err != nil {
		abortRequestWithError(c, NewAppError(ErrorCodeValidationFailed, err.Error()))
		return
	}
	logrus.Info("request to create note")
	note := note.Note{
		Email:   noteReqPayload.Email,
		Title:   noteReqPayload.Title,
		Content: noteReqPayload.Content,
	}
	if err := n.noteService.Save(note); err != nil {
		abortRequestWithError(c, err)
		return
	}
	c.Status(http.StatusOK)
	logrus.Info("request to create note successful")
}

func (n *NoteHandler) UpdateNote(c *gin.Context) {
	id := c.Param("id")
	if err := validation.Validate(id, validation.Required, is.Int); err != nil {
		abortRequestWithError(c, NewAppError(ErrorCodeValidationFailed, fmt.Sprintf("id: %s", err.Error())))
		return
	}
	noteId, _ := strconv.Atoi(id)
	var noteReqPayload NoteRequestPayload
	c.BindJSON(&noteReqPayload)
	if err := noteReqPayload.Validate(); err != nil {
		abortRequestWithError(c, NewAppError(ErrorCodeValidationFailed, err.Error()))
		return
	}
	logrus.WithField("note_id", noteId).Info("request to update note")
	note := note.Note{
		Model: gorm.Model{
			ID: uint(noteId),
		},
		Email:   noteReqPayload.Email,
		Title:   noteReqPayload.Title,
		Content: noteReqPayload.Content,
	}
	if err := n.noteService.Save(note); err != nil {
		abortRequestWithError(c, err)
		return
	}
	c.Status(http.StatusOK)
	logrus.WithField("note_id", noteId).Info("request to update note successful")
}

func (n *NoteHandler) DeleteNote(c *gin.Context) {
	id := c.Param("id")
	if err := validation.Validate(id, validation.Required, is.Int); err != nil {
		abortRequestWithError(c, NewAppError(ErrorCodeValidationFailed, fmt.Sprintf("id: %s", err.Error())))
		return
	}
	logrus.WithField("note_id", id).Info("request to delete note")

	noteId, _ := strconv.Atoi(id)
	if err := n.noteService.Delete(noteId); err != nil {
		abortRequestWithError(c, err)
		return
	}
	c.Status(http.StatusOK)
	logrus.WithField("note_id", id).Info("request to delete note successful")
}
