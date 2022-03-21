package httpservice

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vivekweb2013/gitnoter/internal/user"
)

func TestProfile(t *testing.T) {
	t.Run("should get user profile when the request is valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserService := user.NewMockService(ctrl)

		gin.SetMode(gin.TestMode)
		router := gin.Default()
		handler := NewUserHandler(mockUserService)
		mockUserService.EXPECT().GetByEmail(email).Return(user.User{
			Email:      email,
			Name:       name,
			Location:   location,
			AvatarURL:  avatarURL,
			DisabledAt: nil,
		}, nil)

		// simulate auth middleware with custom handler
		router.GET("/api/v1/user/me", func(c *gin.Context) { c.Set("claims", jwt.MapClaims{"sub": email}) }, handler.Profile)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/user/me", nil)

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.JSONEq(t, `{"avatar_url":"http://example.com/avatar", "email":"john.doe@example.com", "location":"New York", "name":"John Doe"}`, response.Body.String())
	})

	t.Run("should fail with unauthorized response when the claims are not available in context", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserService := user.NewMockService(ctrl)

		gin.SetMode(gin.TestMode)
		router := gin.Default()
		handler := NewUserHandler(mockUserService)

		router.GET("/api/v1/user/me", handler.Profile)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/user/me", nil)

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusUnauthorized, response.Code)
		assert.Empty(t, response.Body.String())
	})

	t.Run("should fail with internal error when user does not exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockUserService := user.NewMockService(ctrl)

		gin.SetMode(gin.TestMode)
		router := gin.Default()
		handler := NewUserHandler(mockUserService)
		mockUserService.EXPECT().GetByEmail(email).Return(user.User{}, errors.New("some error"))
		claims := jwt.MapClaims{"sub": email}

		router.GET("/api/v1/user/me", func(c *gin.Context) { c.Set("claims", claims) }, handler.Profile)
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/user/me", nil)

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusInternalServerError, response.Code)
		assert.JSONEq(t, internalServerErrJson, response.Body.String())
	})

}
