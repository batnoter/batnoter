package httpservice

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/vivekweb2013/gitnoter/internal/config"
)

const (
	email     = "john.doe@example.com"
	name      = "John Doe"
	location  = "New York"
	avatarURL = "http://example.com/avatar"
)

func TestGithubLogin(t *testing.T) {
	t.Run("should redirect to provider when the github login request is valid", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.Default()
		handler := NewGithubHandler(nil, nil, config.OAuth2{})

		router.GET("/api/v1/oauth2/login/github", handler.GithubLogin)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/oauth2/login/github", nil)

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusTemporaryRedirect, response.Code)
	})
}
