package httpservice

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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
	id, _ := strconv.Atoi(c.Param("id"))
	logrus.Infof("Request for note id: %d", id)

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
