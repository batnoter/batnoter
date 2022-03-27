package httpservice

import (
	"errors"
	"net/http"
	"net/http/httptest"
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

	t.Run("should return internal server error fatching repos fails", func(t *testing.T) {
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
		assert.JSONEq(t, internalServerErrJson, response.Body.String())
	})
}
