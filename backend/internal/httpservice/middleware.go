package httpservice

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/vivekweb2013/gitnoter/internal/auth"
)

type Middleware struct {
	authService auth.Service
}

func NewMiddleware(authservice auth.Service) *Middleware {
	return &Middleware{authService: authservice}
}

func (m *Middleware) AuthorizeToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) < len("Bearer ") {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenString := authHeader[len("Bearer "):]

		token, err := m.authService.ValidateToken(tokenString)

		if err != nil || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("claims", claims)
		c.Next()
	}
}
