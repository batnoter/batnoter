package httpservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/sirupsen/logrus"
	"github.com/vivekweb2013/gitnoter/internal/github"
	"github.com/vivekweb2013/gitnoter/internal/user"
	"golang.org/x/oauth2"
)

type RepoPayload struct {
	Name          string `json:"name"`
	Visibility    string `json:"visibility"`
	DefaultBranch string `json:"default_branch"`
}

func (r RepoPayload) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(1, 50)),
		validation.Field(&r.Visibility, validation.Required, validation.Length(1, 20)),
		validation.Field(&r.Visibility, validation.Length(0, 50)),
	)
}

type NoteRequestPayload struct {
	SHA     string `json:"sha"`
	Content string `json:"content"`
}

type NoteResponsePayload struct {
	SHA     string `json:"sha"`
	Path    string `json:"path"`
	Content string `json:"content"`
	Size    int    `json:"size"`
	IsDir   bool   `json:"is_dir"`
}

type NoteSearchResponsePayload struct {
	Total int                   `json:"total"`
	Notes []NoteResponsePayload `json:"notes"`
}

type NoteHandler struct {
	githubService github.Service
	userService   user.Service
}

func NewNoteHandler(githubService github.Service, userService user.Service) *NoteHandler {
	return &NoteHandler{githubService: githubService, userService: userService}
}

const (
	notePathRegex = `(?m)^[^/][/a-zA-Z0-9-]+([^/]\.md)$`
)

func (n *NoteHandler) SearchNotes(c *gin.Context) {
	// get note-path, query, page from query-params as a filter criteria
	path := c.Query("path")
	query := c.Query("query")
	page, _ := strconv.Atoi(c.Query("page"))
	user, err := n.getUser(c)
	if err != nil {
		logrus.Errorf("fetching user from context failed")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	logrus.WithField("user-id", user.ID).WithField("note_path", path).
		WithField("query", query).WithField("page", page).Info("request to search & retrieve notes")
	fileProps := makeFileProps(user, NoteRequestPayload{}, path)
	gitFiles, total, err := n.githubService.SearchFiles(c, parseOAuth2Token(user.GithubToken), fileProps, query, page)
	if err != nil {
		logrus.Errorf("searching notes on github failed")
		abortRequestWithError(c, err)
		return
	}
	noteSearchPayload := makeNoteSearchResponsePayload(gitFiles, total)
	c.JSON(http.StatusOK, noteSearchPayload)
	logrus.WithField("user-id", user.ID).WithField("note_path", path).
		WithField("query", query).WithField("page", page).Info("request to search & retrieve notes successful")
}

func (n *NoteHandler) GetNote(c *gin.Context) {
	path := c.Param("path")
	if err := validation.Validate(path, validation.Required, validation.Match(regexp.MustCompile(notePathRegex))); err != nil {
		abortRequestWithError(c, NewAppError(ErrorCodeValidationFailed, fmt.Sprintf("path: %s", err.Error())))
		return
	}
	user, err := n.getUser(c)
	if err != nil {
		logrus.Errorf("fetching user from context failed")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	logrus.WithField("user-id", user.ID).WithField("note_path", path).Info("request to retrieve note started")
	fileProps := makeFileProps(user, NoteRequestPayload{}, path)
	gitFile, err := n.githubService.GetFile(c, parseOAuth2Token(user.GithubToken), fileProps)
	if err != nil {
		abortRequestWithError(c, err)
		return
	}
	note := makeNoteResponsePayload(gitFile)
	c.JSON(http.StatusOK, note)
	logrus.WithField("user-id", user.ID).WithField("note_path", path).Info("request to retrieve note successful")
}

func (n *NoteHandler) SaveNote(c *gin.Context) {
	path := c.Param("path")
	if err := validation.Validate(path, validation.Required, validation.Match(regexp.MustCompile(notePathRegex))); err != nil {
		abortRequestWithError(c, NewAppError(ErrorCodeValidationFailed, fmt.Sprintf("path: %s", err.Error())))
		return
	}
	var noteReqPayload NoteRequestPayload
	c.BindJSON(&noteReqPayload)
	if err := validation.Validate(noteReqPayload.Content, validation.Required); err != nil {
		abortRequestWithError(c, NewAppError(ErrorCodeValidationFailed, fmt.Sprintf("content: %s", err.Error())))
		return
	}
	user, err := n.getUser(c)
	if err != nil {
		logrus.Errorf("fetching user from context failed")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	logrus.WithField("user-id", user.ID).WithField("note_path", path).Info("request to save note started")
	fileProps := makeFileProps(user, noteReqPayload, path)
	gitFile, err := n.githubService.SaveFile(c, parseOAuth2Token(user.GithubToken), fileProps)
	note := makeNoteResponsePayload(gitFile)
	if err != nil {
		abortRequestWithError(c, err)
		return
	}
	c.JSON(http.StatusOK, note)
	logrus.WithField("user-id", user.ID).WithField("note_path", path).Info("request to save note successful")
}

func (n *NoteHandler) DeleteNote(c *gin.Context) {
	path := c.Param("path")
	if err := validation.Validate(path, validation.Required, validation.Match(regexp.MustCompile(notePathRegex))); err != nil {
		abortRequestWithError(c, NewAppError(ErrorCodeValidationFailed, fmt.Sprintf("path: %s", err.Error())))
		return
	}
	var noteReqPayload NoteRequestPayload
	c.BindJSON(&noteReqPayload)
	if err := validation.Validate(noteReqPayload.SHA, validation.Required); err != nil {
		abortRequestWithError(c, NewAppError(ErrorCodeValidationFailed, fmt.Sprintf("sha: %s", err.Error())))
		return
	}
	user, err := n.getUser(c)
	if err != nil {
		logrus.Errorf("fetching user from context failed")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	logrus.WithField("user-id", user.ID).WithField("note_path", path).Info("request to delete note started")
	fileProps := makeFileProps(user, noteReqPayload, path)
	err = n.githubService.DeleteFile(c, parseOAuth2Token(user.GithubToken), fileProps)
	if err != nil {
		abortRequestWithError(c, err)
		return
	}
	c.Status(http.StatusOK)
	logrus.WithField("user-id", user.ID).WithField("note_path", path).Info("request to delete note successful")
}

func (n *NoteHandler) getUser(c *gin.Context) (user.User, error) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		logrus.Errorf("fetching user-id from context failed")
		return user.User{}, err
	}
	return n.userService.Get(userID)
}

func makeFileProps(user user.User, noteReqPayload NoteRequestPayload, path string) github.GitFileProps {
	return github.GitFileProps{
		SHA:         noteReqPayload.SHA,
		Content:     noteReqPayload.Content,
		Path:        path,
		AuthorName:  user.Name,
		AuthorEmail: user.Email,
		RepoDetails: github.GitRepoProps{
			Repository:    user.DefaultRepo.Name,
			DefaultBranch: user.DefaultRepo.DefaultBranch,
			Owner:         user.GithubUsername,
		},
	}
}

func makeNoteSearchResponsePayload(gitFiles []github.GitFile, total int) NoteSearchResponsePayload {
	notes := make([]NoteResponsePayload, 0, len(gitFiles))
	for _, gitFile := range gitFiles {
		note := makeNoteResponsePayload(gitFile)
		notes = append(notes, note)
	}
	return NoteSearchResponsePayload{
		Total: total,
		Notes: notes,
	}
}

func makeNoteResponsePayload(gitFile github.GitFile) NoteResponsePayload {
	return NoteResponsePayload{
		SHA:     gitFile.SHA,
		Path:    gitFile.Path,
		Content: gitFile.Content,
		Size:    gitFile.Size,
		IsDir:   gitFile.IsDir,
	}
}

func parseOAuth2Token(ghToken string) oauth2.Token {
	oauth2Token := oauth2.Token{}
	if err := json.Unmarshal([]byte(ghToken), &oauth2Token); err != nil {
		logrus.Warn("failed to parse token json to oauth2 token")
	}
	return oauth2Token
}
