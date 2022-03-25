package httpservice

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vivekweb2013/gitnoter/internal/github"
	"github.com/vivekweb2013/gitnoter/internal/preference"
	"github.com/vivekweb2013/gitnoter/internal/user"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

const (
	sha                   = "5ab2f8a4323abafb10abb68657d9d39f1a775057"
	content               = "Hello"
	size                  = 5
	notePath              = "foo/bar.md"
	repository            = "testrepo"
	visibility            = "private"
	owner                 = "johndoe"
	branch                = "main"
	authorName            = "john doe"
	authorEmail           = "john.doe@example.com"
	searchQuery           = "birthday"
	pageNumber            = 2
	token                 = "token"
	internalServerErrJson = `{"code":"internal_server_error", "message":"something went wrong. contact support"}`
)

func TestGetRepos(t *testing.T) {
	t.Run("should return repos when the request is valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockGithubService := github.NewMockService(ctrl)
		mockUserService := user.NewMockService(ctrl)

		router := getRouter()
		u := validUser()
		repos := []github.GitRepo{{
			Name:          repository,
			Visibility:    visibility,
			DefaultBranch: branch,
		}, {
			Name:          repository + "2",
			Visibility:    visibility,
			DefaultBranch: branch + "2",
		}}
		mockUserService.EXPECT().Get(userID).Return(u, nil)
		mockGithubService.EXPECT().GetRepos(gomock.Any(), getOAuth2Token(u.GithubToken)).Return(repos, nil)
		handler := NewNoteHandler(mockGithubService, mockUserService)

		router.GET("/api/v1/user/github/repo", getClaimsHandler(), handler.GetRepos)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/user/github/repo", nil)

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.JSONEq(t, `[{"default_branch":"main", "name":"testrepo", "visibility":"private"},{"default_branch":"main2", "name":"testrepo2", "visibility":"private"}]`, response.Body.String())
	})

	t.Run("should return internal server error fatching repos fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockGithubService := github.NewMockService(ctrl)
		mockUserService := user.NewMockService(ctrl)

		router := getRouter()
		u := validUser()
		mockUserService.EXPECT().Get(userID).Return(u, nil)
		mockGithubService.EXPECT().GetRepos(gomock.Any(), getOAuth2Token(u.GithubToken)).Return([]github.GitRepo{}, errors.New("some error"))
		handler := NewNoteHandler(mockGithubService, mockUserService)

		router.GET("/api/v1/user/github/repo", getClaimsHandler(), handler.GetRepos)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/user/github/repo", nil)

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.JSONEq(t, internalServerErrJson, response.Body.String())
	})
}

func TestSearchNotes(t *testing.T) {
	t.Run("should return notes when the search request is valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockGithubService := github.NewMockService(ctrl)
		mockUserService := user.NewMockService(ctrl)

		router := getRouter()
		u := validUser()
		fp := github.GitFileProps{SHA: "", Path: notePath, Content: "", AuthorName: authorName, AuthorEmail: authorEmail, RepoDetails: github.GitRepoProps{Repository: repository, DefaultBranch: branch, Owner: owner}}
		gitFiles := []github.GitFile{{
			SHA:     sha,
			Path:    notePath,
			Content: content,
			Size:    size,
			IsDir:   false,
		}}
		mockUserService.EXPECT().Get(userID).Return(u, nil)
		mockGithubService.EXPECT().SearchFiles(gomock.Any(), getOAuth2Token(u.GithubToken), fp, searchQuery, pageNumber).Return(gitFiles, 1, nil)
		handler := NewNoteHandler(mockGithubService, mockUserService)

		router.GET("/api/v1/note", getClaimsHandler(), handler.SearchNotes)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/note?path=%s&query=%s&page=%d", notePath, searchQuery, pageNumber), nil)

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.JSONEq(t, `{"notes":[{"content":"Hello", "is_dir":false, "path":"foo/bar.md", "sha":"5ab2f8a4323abafb10abb68657d9d39f1a775057", "size":5}], "total":1}`, response.Body.String())
	})

	t.Run("should return internal server error when the search fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockGithubService := github.NewMockService(ctrl)
		mockUserService := user.NewMockService(ctrl)

		router := getRouter()
		u := validUser()
		mockUserService.EXPECT().Get(userID).Return(u, nil)
		mockGithubService.EXPECT().SearchFiles(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]github.GitFile{}, 0, errors.New("some error"))
		handler := NewNoteHandler(mockGithubService, mockUserService)

		router.GET("/api/v1/note", getClaimsHandler(), handler.SearchNotes)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/note?path=%s&query=%s&page=%d", notePath, searchQuery, pageNumber), nil)

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.JSONEq(t, internalServerErrJson, response.Body.String())
	})
}

func TestGetNote(t *testing.T) {
	t.Run("should return a note when the get request is valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockGithubService := github.NewMockService(ctrl)
		mockUserService := user.NewMockService(ctrl)

		router := getRouter()
		u := validUser()
		fp := github.GitFileProps{SHA: "", Path: notePath, Content: "", AuthorName: authorName, AuthorEmail: authorEmail, RepoDetails: github.GitRepoProps{Repository: repository, DefaultBranch: branch, Owner: owner}}
		f := validGitFile()
		mockUserService.EXPECT().Get(userID).Return(u, nil)
		mockGithubService.EXPECT().GetFile(gomock.Any(), getOAuth2Token(u.GithubToken), fp).Return(f, nil)
		handler := NewNoteHandler(mockGithubService, mockUserService)

		router.GET("/api/v1/note/:path", getClaimsHandler(), handler.GetNote)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/note/%s", url.QueryEscape(notePath)), nil)

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.JSONEq(t, fmt.Sprintf(`{"content":"%s", "is_dir":%t, "path":"%s", "sha":"%s", "size":%d}`, f.Content, f.IsDir, f.Path, f.SHA, f.Size), response.Body.String())
	})

	t.Run("should return error when retrieving a note fails due to missing user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockGithubService := github.NewMockService(ctrl)
		mockUserService := user.NewMockService(ctrl)

		router := getRouter()
		mockUserService.EXPECT().Get(gomock.Any()).Return(user.User{}, errors.New("some error"))
		handler := NewNoteHandler(mockGithubService, mockUserService)

		router.GET("/api/v1/note/:path", getClaimsHandler(), handler.GetNote)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/note/%s", url.QueryEscape(notePath)), nil)

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusUnauthorized, response.Code)
		assert.Equal(t, "", response.Body.String())
	})

	t.Run("should return error when retrieving a note fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockGithubService := github.NewMockService(ctrl)
		mockUserService := user.NewMockService(ctrl)

		router := getRouter()
		u := validUser()
		mockUserService.EXPECT().Get(gomock.Any()).Return(u, nil)
		mockGithubService.EXPECT().GetFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(github.GitFile{}, errors.New("some error"))
		handler := NewNoteHandler(mockGithubService, mockUserService)

		router.GET("/api/v1/note/:path", getClaimsHandler(), handler.GetNote)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/note/%s", url.QueryEscape(notePath)), nil)

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.JSONEq(t, internalServerErrJson, response.Body.String())
	})

	t.Run("should return bad request error when get request has invalid path param", func(t *testing.T) {
		for _, invalidPath := range getInvalidNotePaths() {
			t.Run("with invalid path: "+invalidPath, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				handler := NewNoteHandler(nil, nil)

				router := getRouter()
				router.GET("/api/v1/note/:path", handler.GetNote)
				response := httptest.NewRecorder()
				req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/note/%s", url.QueryEscape(invalidPath)), nil)

				router.ServeHTTP(response, req)
				assert.Equal(t, http.StatusBadRequest, response.Code)
				assert.JSONEq(t, `{"code":"validation_failed", "message":"path: must be in a valid format"}`, response.Body.String())
			})
		}
	})
}

func TestSaveNote(t *testing.T) {
	t.Run("should save(create) a new note when the save request payload does not have the sha value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockGithubService := github.NewMockService(ctrl)
		mockUserService := user.NewMockService(ctrl)

		router := getRouter()
		u := validUser()
		fp := github.GitFileProps{SHA: "", Path: notePath, Content: content, AuthorName: authorName, AuthorEmail: authorEmail, RepoDetails: github.GitRepoProps{Repository: repository, DefaultBranch: branch, Owner: owner}}
		f := validGitFile()
		n := NoteRequestPayload{
			Content: content,
		}
		noteJson, _ := json.Marshal(n)
		mockUserService.EXPECT().Get(userID).Return(u, nil)
		mockGithubService.EXPECT().SaveFile(gomock.Any(), getOAuth2Token(u.GithubToken), fp).Return(f, nil)
		handler := NewNoteHandler(mockGithubService, mockUserService)

		router.POST("/api/v1/note/:path", getClaimsHandler(), handler.SaveNote)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/note/%s", url.QueryEscape(notePath)), strings.NewReader(string(noteJson)))

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.JSONEq(t, fmt.Sprintf(`{"content":"%s", "is_dir":%t, "path":"%s", "sha":"%s", "size":%d}`, f.Content, f.IsDir, f.Path, f.SHA, f.Size), response.Body.String())
	})

	t.Run("should save(update) a new note when the save request payload has the sha value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockGithubService := github.NewMockService(ctrl)
		mockUserService := user.NewMockService(ctrl)

		router := getRouter()
		u := validUser()
		fp := github.GitFileProps{SHA: sha, Path: notePath, Content: content, AuthorName: authorName, AuthorEmail: authorEmail, RepoDetails: github.GitRepoProps{Repository: repository, DefaultBranch: branch, Owner: owner}}
		f := validGitFile()
		n := NoteRequestPayload{
			SHA:     sha,
			Content: content,
		}
		noteJson, _ := json.Marshal(n)
		mockUserService.EXPECT().Get(userID).Return(u, nil)
		mockGithubService.EXPECT().SaveFile(gomock.Any(), getOAuth2Token(u.GithubToken), fp).Return(f, nil)
		handler := NewNoteHandler(mockGithubService, mockUserService)

		router.POST("/api/v1/note/:path", getClaimsHandler(), handler.SaveNote)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/note/%s", url.QueryEscape(notePath)), strings.NewReader(string(noteJson)))

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.JSONEq(t, fmt.Sprintf(`{"content":"%s", "is_dir":%t, "path":"%s", "sha":"%s", "size":%d}`, f.Content, f.IsDir, f.Path, f.SHA, f.Size), response.Body.String())
	})

	t.Run("should return internal server error when saving note fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockGithubService := github.NewMockService(ctrl)
		mockUserService := user.NewMockService(ctrl)

		router := getRouter()
		u := validUser()
		n := NoteRequestPayload{
			Content: content,
		}
		noteJson, _ := json.Marshal(n)
		mockUserService.EXPECT().Get(userID).Return(u, nil)
		mockGithubService.EXPECT().SaveFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(github.GitFile{}, errors.New("some error"))
		handler := NewNoteHandler(mockGithubService, mockUserService)

		router.POST("/api/v1/note/:path", getClaimsHandler(), handler.SaveNote)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/note/%s", url.QueryEscape(notePath)), strings.NewReader(string(noteJson)))

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.JSONEq(t, internalServerErrJson, response.Body.String())
	})

	t.Run("should return bad request error when save request payload validation fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		handler := NewNoteHandler(nil, nil)

		router := getRouter()
		router.POST("/api/v1/note/:path", handler.SaveNote)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/note/%s", url.QueryEscape(notePath)), strings.NewReader("{}"))

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.JSONEq(t, `{"code":"validation_failed", "message":"content: cannot be blank"}`, response.Body.String())
	})

	t.Run("should return bad request error when save request has invalid path param", func(t *testing.T) {
		for _, invalidPath := range getInvalidNotePaths() {
			t.Run("with invalid path: "+invalidPath, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				handler := NewNoteHandler(nil, nil)

				router := getRouter()
				router.POST("/api/v1/note/:path", handler.SaveNote)
				response := httptest.NewRecorder()
				req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/note/%s", url.QueryEscape(invalidPath)), nil)

				router.ServeHTTP(response, req)
				assert.Equal(t, http.StatusBadRequest, response.Code)
				assert.JSONEq(t, `{"code":"validation_failed", "message":"path: must be in a valid format"}`, response.Body.String())
			})
		}
	})
}

func TestDeleteNote(t *testing.T) {
	t.Run("should delete a note when the delete request is valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockGithubService := github.NewMockService(ctrl)
		mockUserService := user.NewMockService(ctrl)

		router := getRouter()
		u := validUser()
		fp := github.GitFileProps{SHA: sha, Path: notePath, Content: "", AuthorName: authorName, AuthorEmail: authorEmail, RepoDetails: github.GitRepoProps{Repository: repository, DefaultBranch: branch, Owner: owner}}
		n := NoteRequestPayload{
			SHA: sha,
		}
		noteJson, _ := json.Marshal(n)
		mockUserService.EXPECT().Get(userID).Return(u, nil)
		mockGithubService.EXPECT().DeleteFile(gomock.Any(), getOAuth2Token(u.GithubToken), fp).Return(nil)
		handler := NewNoteHandler(mockGithubService, mockUserService)

		router.DELETE("/api/v1/note/:path", getClaimsHandler(), handler.DeleteNote)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/note/%s", url.QueryEscape(notePath)), strings.NewReader(string(noteJson)))

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "", response.Body.String())
	})

	t.Run("should return internal server error when deleting a note fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockGithubService := github.NewMockService(ctrl)
		mockUserService := user.NewMockService(ctrl)

		router := getRouter()
		u := validUser()
		n := NoteRequestPayload{
			SHA: sha,
		}
		noteJson, _ := json.Marshal(n)
		mockUserService.EXPECT().Get(userID).Return(u, nil)
		mockGithubService.EXPECT().DeleteFile(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("some error"))
		handler := NewNoteHandler(mockGithubService, mockUserService)

		router.DELETE("/api/v1/note/:path", getClaimsHandler(), handler.DeleteNote)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/note/%s", url.QueryEscape(notePath)), strings.NewReader(string(noteJson)))

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.JSONEq(t, internalServerErrJson, response.Body.String())
	})

	t.Run("should return bad request error when delete request payload validation fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		handler := NewNoteHandler(nil, nil)

		router := getRouter()
		router.DELETE("/api/v1/note/:path", handler.DeleteNote)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/note/%s", url.QueryEscape(notePath)), strings.NewReader("{}"))

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.JSONEq(t, `{"code":"validation_failed", "message":"sha: cannot be blank"}`, response.Body.String())
	})

	t.Run("should return bad request error when delete request has invalid path param", func(t *testing.T) {
		for _, invalidPath := range getInvalidNotePaths() {
			t.Run("with invalid path: "+invalidPath, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				handler := NewNoteHandler(nil, nil)

				router := getRouter()
				router.DELETE("/api/v1/note/:path", handler.DeleteNote)
				response := httptest.NewRecorder()
				req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/note/%s", url.QueryEscape(invalidPath)), nil)

				router.ServeHTTP(response, req)
				assert.Equal(t, http.StatusBadRequest, response.Code)
				assert.JSONEq(t, `{"code":"validation_failed", "message":"path: must be in a valid format"}`, response.Body.String())
			})
		}
	})
}

func validUser() user.User {
	return user.User{
		Model: gorm.Model{
			ID: userID,
		},
		Email:          authorEmail,
		Name:           authorName,
		Location:       location,
		GithubUsername: owner,
		GithubToken:    oauth2TokenJSON,
		DefaultRepo: &preference.DefaultRepo{
			Name:          repository,
			Visibility:    visibility,
			DefaultBranch: branch,
		},
	}
}

func validGitFile() github.GitFile {
	return github.GitFile{
		SHA:     sha,
		Path:    notePath,
		Content: content,
		Size:    size,
		IsDir:   false,
	}
}

func getInvalidNotePaths() []string {
	return []string{".md", "/", "/foo", "/.md", "/bar.md", "foo", "foo/bar", "foo/.md", "foo/bar.md/foo", "foo/bar.md/foo.md"}
}

func getRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.UseRawPath = true
	return router
}

func getClaimsHandler() func(c *gin.Context) {
	// this test handler simulates auth middleware
	return func(c *gin.Context) {
		claims := jwt.MapClaims{"sub": strconv.FormatUint(uint64(userID), 10)}
		c.Set("claims", claims)
	}
}

func getOAuth2Token(s string) oauth2.Token {
	var oauth2Token oauth2.Token
	json.Unmarshal([]byte(s), &oauth2Token)
	return oauth2Token
}
