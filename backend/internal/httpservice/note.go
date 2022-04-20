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

// NoteRequestPayload represents the http request payload of note entity.
type NoteRequestPayload struct {
	SHA     string `json:"sha"`
	Content string `json:"content"`
}

// NoteResponsePayload represents the http response payload of note entity.
type NoteResponsePayload struct {
	SHA     string `json:"sha"`
	Path    string `json:"path"`
	Content string `json:"content"`
	Size    int    `json:"size"`
	IsDir   bool   `json:"is_dir"`
}

// NoteSearchResponsePayload represents the http response payload for note search operation.
// Total is the count of total results found.
// Notes are the subset of search result as requested with pagination attributes.
type NoteSearchResponsePayload struct {
	Total int                   `json:"total"`
	Notes []NoteResponsePayload `json:"notes"`
}

// NoteHandler represents http handler for managing note entities.
type NoteHandler struct {
	githubService github.Service
	userService   user.Service
}

// NewNoteHandler creates and returns a new note handler.
func NewNoteHandler(githubService github.Service, userService user.Service) *NoteHandler {
	return &NoteHandler{githubService: githubService, userService: userService}
}

// SearchNotes performs a note search operation with specified filter criteria.
// It returns the result of search operation as a http response.
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

// GetNotesTree returns a complete tree of note repository as a http response.
func (n *NoteHandler) GetNotesTree(c *gin.Context) {
	user, err := n.getUser(c)
	if err != nil {
		logrus.Errorf("fetching user from context failed")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	logrus.WithField("user-id", user.ID).Info("request to retrieve tree")
	fileProps := makeFileProps(user, NoteRequestPayload{}, "")
	gitFiles, err := n.githubService.GetTree(c, parseOAuth2Token(user.GithubToken), fileProps)
	if err != nil {
		logrus.Errorf("retrieving tree on github failed")
		abortRequestWithError(c, err)
		return
	}
	notes := make([]NoteResponsePayload, 0, len(gitFiles))
	for _, gitFile := range gitFiles {
		note := makeNoteResponsePayload(gitFile)
		notes = append(notes, note)
	}
	c.JSON(http.StatusOK, notes)
	logrus.WithField("user-id", user.ID).Info("request to retrieve tree successful")
}

// GetAllNotes returns all the notes from a path as a http response.
func (n *NoteHandler) GetAllNotes(c *gin.Context) {
	path := c.Query("path")
	user, err := n.getUser(c)
	if err != nil {
		logrus.Errorf("fetching user from context failed")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	logrus.WithField("user-id", user.ID).WithField("note_path", path).Info("request to retrieve notes started")
	fileProps := makeFileProps(user, NoteRequestPayload{}, path)
	gitFiles, err := n.githubService.GetAllFiles(c, parseOAuth2Token(user.GithubToken), fileProps)
	if err != nil {
		logrus.Errorf("retrieving notes from github failed")
		abortRequestWithError(c, err)
		return
	}
	notes := make([]NoteResponsePayload, 0, len(gitFiles))
	for _, gitFile := range gitFiles {
		note := makeNoteResponsePayload(gitFile)
		notes = append(notes, note)
	}
	c.JSON(http.StatusOK, notes)
	logrus.WithField("user-id", user.ID).WithField("note_path", path).Info("request to retrieve notes successful")
}

// GetNote returns a note with requested path as a http response.
func (n *NoteHandler) GetNote(c *gin.Context) {
	path := c.Param("path")
	if err := validation.Validate(path, validation.Required, validation.Match(regexp.MustCompile(github.ValidFilePathRegex))); err != nil {
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

// SaveNote stores the note and returns the metadata as a http response.
func (n *NoteHandler) SaveNote(c *gin.Context) {
	path := c.Param("path")
	if err := validation.Validate(path, validation.Required, validation.Match(regexp.MustCompile(github.ValidFilePathRegex))); err != nil {
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

// DeleteNote deletes a note with requested path.
func (n *NoteHandler) DeleteNote(c *gin.Context) {
	path := c.Param("path")
	if err := validation.Validate(path, validation.Required, validation.Match(regexp.MustCompile(github.ValidFilePathRegex))); err != nil {
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
	authorName := user.Name
	if authorName == "" {
		authorName = user.GithubUsername
	}
	return github.GitFileProps{
		SHA:         noteReqPayload.SHA,
		Content:     noteReqPayload.Content,
		Path:        path,
		AuthorName:  authorName,
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
