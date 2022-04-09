package httpservice

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vivekweb2013/gitnoter/internal/github"
	"github.com/vivekweb2013/gitnoter/internal/preference"
	"github.com/vivekweb2013/gitnoter/internal/user"
)

func TestGetRepos(t *testing.T) {
	t.Run("should return repos when the request is valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPreferenceService := preference.NewMockService(ctrl)
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
		mockGithubService.EXPECT().GetRepos(gomock.Any(), getOAuth2Token(u.GithubToken)).Return(repos, nil)
		mockUserService.EXPECT().Get(userID).Return(u, nil)
		handler := NewPreferenceHandler(mockPreferenceService, mockGithubService, mockUserService)

		router.GET("/api/v1/user/preference/repo", getClaimsHandler(), handler.GetRepos)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/user/preference/repo", nil)

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.JSONEq(t, `[{"default_branch":"main", "name":"testrepo", "visibility":"private"},{"default_branch":"main2", "name":"testrepo2", "visibility":"private"}]`, response.Body.String())
	})

	t.Run("should return internal server error when fatching repos fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPreferenceService := preference.NewMockService(ctrl)
		mockGithubService := github.NewMockService(ctrl)
		mockUserService := user.NewMockService(ctrl)

		router := getRouter()
		u := validUser()
		mockUserService.EXPECT().Get(userID).Return(u, nil)
		mockGithubService.EXPECT().GetRepos(gomock.Any(), getOAuth2Token(u.GithubToken)).Return([]github.GitRepo{}, errors.New("some error"))
		handler := NewPreferenceHandler(mockPreferenceService, mockGithubService, mockUserService)

		router.GET("/api/v1/user/preference/repo", getClaimsHandler(), handler.GetRepos)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/user/preference/repo", nil)

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.JSONEq(t, internalServerErrJSON, response.Body.String())
	})

	t.Run("should return unauthorized error when retrieving user fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPreferenceService := preference.NewMockService(ctrl)
		mockGithubService := github.NewMockService(ctrl)
		mockUserService := user.NewMockService(ctrl)

		router := getRouter()
		mockUserService.EXPECT().Get(userID).Return(user.User{}, errors.New("some error"))
		handler := NewPreferenceHandler(mockPreferenceService, mockGithubService, mockUserService)

		router.GET("/api/v1/user/preference/repo", getClaimsHandler(), handler.GetRepos)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/user/preference/repo", nil)

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusUnauthorized, response.Code)
		assert.Equal(t, "", response.Body.String())
	})
}

func TestSaveDefaultRepo(t *testing.T) {
	t.Run("should save default repo against user when the request is valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPreferenceService := preference.NewMockService(ctrl)
		mockGithubService := github.NewMockService(ctrl)
		mockUserService := user.NewMockService(ctrl)

		router := getRouter()
		repoPayload := fmt.Sprintf(`{
			"name":"%s",
			"visibility":"%s",
			"default_branch":"%s"
		}`, repository, visibility, branch)
		dbDefaultRepo := preference.DefaultRepo{
			UserID:        userID,
			Name:          repository,
			Visibility:    visibility,
			DefaultBranch: branch,
		}
		mockPreferenceService.EXPECT().GetByUserID(userID).Return(dbDefaultRepo, nil)
		mockPreferenceService.EXPECT().Save(dbDefaultRepo).Return(nil)
		handler := NewPreferenceHandler(mockPreferenceService, mockGithubService, mockUserService)

		router.POST("/api/v1/user/preference/repo", getClaimsHandler(), handler.SaveDefaultRepo)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/user/preference/repo", strings.NewReader(repoPayload))

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("should return internal server error when saving default repo fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPreferenceService := preference.NewMockService(ctrl)
		mockGithubService := github.NewMockService(ctrl)
		mockUserService := user.NewMockService(ctrl)

		router := getRouter()
		repoPayload := fmt.Sprintf(`{
			"name":"%s",
			"visibility":"%s",
			"default_branch":"%s"
		}`, repository, visibility, branch)
		dbDefaultRepo := preference.DefaultRepo{
			UserID:        userID,
			Name:          repository,
			Visibility:    visibility,
			DefaultBranch: branch,
		}
		mockPreferenceService.EXPECT().GetByUserID(userID).Return(dbDefaultRepo, nil)
		mockPreferenceService.EXPECT().Save(dbDefaultRepo).Return(errors.New("some error"))
		handler := NewPreferenceHandler(mockPreferenceService, mockGithubService, mockUserService)

		router.POST("/api/v1/user/preference/repo", getClaimsHandler(), handler.SaveDefaultRepo)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/user/preference/repo", strings.NewReader(repoPayload))

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.JSONEq(t, internalServerErrJSON, response.Body.String())
	})

	t.Run("should return internal server error when retring default repo from preference fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPreferenceService := preference.NewMockService(ctrl)
		mockGithubService := github.NewMockService(ctrl)
		mockUserService := user.NewMockService(ctrl)

		router := getRouter()
		repoPayload := fmt.Sprintf(`{
			"name":"%s",
			"visibility":"%s",
			"default_branch":"%s"
		}`, repository, visibility, branch)
		mockPreferenceService.EXPECT().GetByUserID(userID).Return(preference.DefaultRepo{}, errors.New("some error"))
		handler := NewPreferenceHandler(mockPreferenceService, mockGithubService, mockUserService)

		router.POST("/api/v1/user/preference/repo", getClaimsHandler(), handler.SaveDefaultRepo)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/user/preference/repo", strings.NewReader(repoPayload))

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.JSONEq(t, internalServerErrJSON, response.Body.String())
	})

	t.Run("should return bad request error when repo request payload is invalid (missing repo-name)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPreferenceService := preference.NewMockService(ctrl)
		mockGithubService := github.NewMockService(ctrl)
		mockUserService := user.NewMockService(ctrl)

		router := getRouter()
		repoPayload := fmt.Sprintf(`{
			"visibility":"%s",
			"default_branch":"%s"
		}`, visibility, branch)
		handler := NewPreferenceHandler(mockPreferenceService, mockGithubService, mockUserService)

		router.POST("/api/v1/user/preference/repo", getClaimsHandler(), handler.SaveDefaultRepo)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/user/preference/repo", strings.NewReader(repoPayload))

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.JSONEq(t, `{"code":"validation_failed", "message":"name: cannot be blank."}`, response.Body.String())
	})

	t.Run("should return unauthorized error when claims missing in context", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPreferenceService := preference.NewMockService(ctrl)
		mockGithubService := github.NewMockService(ctrl)
		mockUserService := user.NewMockService(ctrl)

		router := getRouter()
		repoPayload := fmt.Sprintf(`{
			"name":"%s",
			"visibility":"%s",
			"default_branch":"%s"
		}`, repository, visibility, branch)
		handler := NewPreferenceHandler(mockPreferenceService, mockGithubService, mockUserService)

		router.POST("/api/v1/user/preference/repo", handler.SaveDefaultRepo)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/user/preference/repo", strings.NewReader(repoPayload))

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusUnauthorized, response.Code)
		assert.Equal(t, "", response.Body.String())
	})
}

func TestAutoSetupRepo(t *testing.T) {
	t.Run("should auto setup notes repo for user when the request is valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPreferenceService := preference.NewMockService(ctrl)
		mockGithubService := github.NewMockService(ctrl)
		mockUserService := user.NewMockService(ctrl)

		router := getRouter()
		u := validUser()
		dbDefaultRepo := preference.DefaultRepo{
			UserID:        userID,
			Name:          repository,
			Visibility:    visibility,
			DefaultBranch: branch,
		}
		mockUserService.EXPECT().Get(userID).Return(u, nil)
		mockGithubService.EXPECT().CreateRepo(gomock.Any(), gomock.Any(), repository).Return(nil)
		mockPreferenceService.EXPECT().Save(dbDefaultRepo).Return(nil)
		handler := NewPreferenceHandler(mockPreferenceService, mockGithubService, mockUserService)

		router.POST("/api/v1/user/preference/repo/auto", getClaimsHandler(), handler.AutoSetupRepo)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/user/preference/repo/auto?repoName=%s", repository), nil)

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("should return bad request error response when repo name query param is not provided", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPreferenceService := preference.NewMockService(ctrl)
		mockGithubService := github.NewMockService(ctrl)
		mockUserService := user.NewMockService(ctrl)

		router := getRouter()
		handler := NewPreferenceHandler(mockPreferenceService, mockGithubService, mockUserService)

		router.POST("/api/v1/user/preference/repo/auto", getClaimsHandler(), handler.AutoSetupRepo)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/user/preference/repo/auto", nil)

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.JSONEq(t, `{"code":"validation_failed", "message":"repoName: cannot be blank"}`, response.Body.String())
	})

	t.Run("should return error response when creating new repo fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPreferenceService := preference.NewMockService(ctrl)
		mockGithubService := github.NewMockService(ctrl)
		mockUserService := user.NewMockService(ctrl)

		router := getRouter()
		u := validUser()
		mockUserService.EXPECT().Get(userID).Return(u, nil)
		mockGithubService.EXPECT().CreateRepo(gomock.Any(), gomock.Any(), repository).Return(errors.New("some error"))
		handler := NewPreferenceHandler(mockPreferenceService, mockGithubService, mockUserService)

		router.POST("/api/v1/user/preference/repo/auto", getClaimsHandler(), handler.AutoSetupRepo)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/user/preference/repo/auto?repoName=%s", repository), nil)

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.JSONEq(t, internalServerErrJSON, response.Body.String())
	})

	t.Run("should return error response when storing default repo fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPreferenceService := preference.NewMockService(ctrl)
		mockGithubService := github.NewMockService(ctrl)
		mockUserService := user.NewMockService(ctrl)

		router := getRouter()
		u := validUser()
		dbDefaultRepo := preference.DefaultRepo{
			UserID:        userID,
			Name:          repository,
			Visibility:    visibility,
			DefaultBranch: branch,
		}
		mockUserService.EXPECT().Get(userID).Return(u, nil)
		mockGithubService.EXPECT().CreateRepo(gomock.Any(), gomock.Any(), repository).Return(nil)
		mockPreferenceService.EXPECT().Save(dbDefaultRepo).Return(errors.New(("some error")))
		handler := NewPreferenceHandler(mockPreferenceService, mockGithubService, mockUserService)

		router.POST("/api/v1/user/preference/repo/auto", getClaimsHandler(), handler.AutoSetupRepo)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/user/preference/repo/auto?repoName=%s", repository), nil)

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.JSONEq(t, internalServerErrJSON, response.Body.String())
	})

	t.Run("should return unauthorized error when retrieving user fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPreferenceService := preference.NewMockService(ctrl)
		mockGithubService := github.NewMockService(ctrl)
		mockUserService := user.NewMockService(ctrl)

		router := getRouter()
		mockUserService.EXPECT().Get(userID).Return(user.User{}, errors.New("some error"))
		handler := NewPreferenceHandler(mockPreferenceService, mockGithubService, mockUserService)

		router.POST("/api/v1/user/preference/repo/auto", getClaimsHandler(), handler.AutoSetupRepo)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/user/preference/repo/auto?repoName=%s", repository), nil)

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusUnauthorized, response.Code)
		assert.Equal(t, "", response.Body.String())
	})
}
