package httpservice

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vivekweb2013/gitnoter/internal/auth"
)

func TestAuthorizeToken(t *testing.T) {
	t.Run("should validate token & store claims in context when the request has valid token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := auth.NewMockService(ctrl)

		tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE4OTk4NzcxMzgsImlhdCI6MTY0NzQwOTI0NiwiaXNzIjoidGVzdCJ9.fF2-2_Y9rby5FgjnDp14pLI-RlFisT4fnDqJM4ek35M"
		token, _ := getToken(tokenString)
		var claims map[string]interface{}
		expectedClaims := map[string]interface{}(map[string]interface{}{"exp": 1.899877138e+09, "iat": 1.647409246e+09, "iss": "test"})

		gin.SetMode(gin.TestMode)
		router := gin.Default()
		middleware := NewMiddleware(mockService)
		mockService.EXPECT().ValidateToken(tokenString).Return(token, nil)

		// Creating a new handler to test the modified context.
		// Since this is the only way to verify context
		// https://github.com/gin-gonic/gin/blob/master/auth_test.go#L91-L126
		router.GET("/", middleware.AuthorizeToken(), func(c *gin.Context) {
			claimsVal, _ := c.Get("claims")
			claims = claimsVal.(jwt.MapClaims)
		})
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, expectedClaims, claims)
	})

	t.Run("should abort with unauthorized status when the request does not have auth header", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := auth.NewMockService(ctrl)
		var hasClaims bool

		gin.SetMode(gin.TestMode)
		router := gin.Default()
		middleware := NewMiddleware(mockService)

		// Creating a new handler to test the modified context.
		// Since this is the only way to verify context
		// https://github.com/gin-gonic/gin/blob/master/auth_test.go#L91-L126
		router.GET("/", middleware.AuthorizeToken(), func(c *gin.Context) {
			claimsVal, _ := c.Get("claims")
			_, hasClaims = claimsVal.(jwt.MapClaims)
		})
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/", nil)

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusUnauthorized, response.Code)
		assert.Equal(t, false, hasClaims)
	})

	t.Run("should abort with unauthorized status when request has expired token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockService := auth.NewMockService(ctrl)

		expiredToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDc0MDkyOTksImlhdCI6MTY0NzQwOTI0NiwiaXNzIjoidGVzdCJ9.TFUqmSuWnXY1f9RR8QuJ2CYK95c_JeCSQiYZoTdgBvI"
		token, _ := getToken(expiredToken)
		var hasClaims bool

		gin.SetMode(gin.TestMode)
		router := gin.Default()
		middleware := NewMiddleware(mockService)
		mockService.EXPECT().ValidateToken(expiredToken).Return(token, nil)

		// Creating a new handler to test the modified context.
		// Since this is the only way to verify context
		// https://github.com/gin-gonic/gin/blob/master/auth_test.go#L91-L126
		router.GET("/", middleware.AuthorizeToken(), func(c *gin.Context) {
			claimsVal, _ := c.Get("claims")
			_, hasClaims = claimsVal.(jwt.MapClaims)
		})
		response := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", expiredToken))

		router.ServeHTTP(response, req)
		assert.Equal(t, http.StatusUnauthorized, response.Code)
		assert.Equal(t, false, hasClaims)
	})
}

func getToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("test"), nil
	})
}
