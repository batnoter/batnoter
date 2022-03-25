package httpservice

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	gh "github.com/google/go-github/v43/github"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vivekweb2013/gitnoter/internal/auth"
	"github.com/vivekweb2013/gitnoter/internal/github"
	"github.com/vivekweb2013/gitnoter/internal/preference"
	"github.com/vivekweb2013/gitnoter/internal/user"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

const (
	userID          = uint(1012)
	email           = "john.doe@example.com"
	name            = "John Doe"
	location        = "New York"
	avatarURL       = "http://example.com/avatar"
	oauth2TokenJSON = `{"access_token":"gho_token","token_type":"bearer","expiry":"0001-01-01T00:00:00Z"}`
)

const ()

func TestGithubLogin(t *testing.T) {
	t.Run("should redirect to provider when the github login request is valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		githubService := github.NewMockService(ctrl)

		gin.SetMode(gin.TestMode)
		router := gin.Default()
		handler := NewLoginHandler(nil, githubService, nil)
		githubService.EXPECT().GetAuthCodeURL(gomock.Any()).Return("/")

		router.GET("/api/v1/oauth2/login/github", handler.GithubLogin)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/oauth2/login/github", nil)

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusTemporaryRedirect, response.Code)
	})
}

func TestGithubOAuth2Callback(t *testing.T) {
	t.Run("should save the user(with token) & return token response when callback invoked", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		authService := auth.NewMockService(ctrl)
		githubService := github.NewMockService(ctrl)
		userService := user.NewMockService(ctrl)
		state := uuid.NewString()
		authCode := "abcd"
		appToken := "app_token"
		var oauthToken oauth2.Token
		json.Unmarshal([]byte(oauth2TokenJSON), &oauthToken)
		githubUser := validGithubUser()
		dbUser := makeDBUser(githubUser, oauth2TokenJSON)

		gin.SetMode(gin.TestMode)
		router := gin.Default()
		handler := NewLoginHandler(authService, githubService, userService)
		githubService.EXPECT().GetToken(gomock.Any(), authCode).Return(oauthToken, nil)
		githubService.EXPECT().GetUser(gomock.Any(), oauthToken).Return(githubUser, nil)
		userService.EXPECT().GetByEmail(email).Return(dbUser, nil)
		userService.EXPECT().Save(dbUser).Return(uint(1), nil)
		authService.EXPECT().GenerateToken(uint(1)).Return(appToken, nil)

		router.GET("/oauth2/github/callback", handler.GithubOAuth2Callback)
		response := httptest.NewRecorder()
		cookie := http.Cookie{
			Name:     "state",
			Value:    state,
			Path:     "/",
			Expires:  time.Now().Add(10 * time.Minute),
			HttpOnly: true,
		}
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/oauth2/github/callback?code=%s&state=%s", authCode, state), nil)
		req.AddCookie(&cookie)

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Contains(t, response.Body.String(), appToken)
	})
}

func makeDBUser(githubUser gh.User, token string) user.User {
	return user.User{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Email:          *githubUser.Email,
		Name:           *githubUser.Name,
		Location:       *githubUser.Location,
		AvatarURL:      *githubUser.AvatarURL,
		GithubID:       *githubUser.ID,
		GithubUsername: *githubUser.Login,
		GithubToken:    token,
		DefaultRepo: &preference.DefaultRepo{
			Model: gorm.Model{
				ID:        1,
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
			UserID:        1,
			Name:          repository,
			Visibility:    visibility,
			DefaultBranch: branch,
		},
	}
}

func validGithubUser() gh.User {
	testEmail := email
	testName := name
	testLocation := location
	testAvatarURL := avatarURL
	githubID := int64(12345)
	githubUsername := "johndoe"
	return gh.User{
		ID:        &githubID,
		Login:     &githubUsername,
		Email:     &testEmail,
		Name:      &testName,
		Location:  &testLocation,
		AvatarURL: &testAvatarURL,
	}
}
