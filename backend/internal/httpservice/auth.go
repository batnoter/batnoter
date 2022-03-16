package httpservice

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vivekweb2013/gitnoter/internal/auth"
)

type LoginResponsePayload struct {
	Token string `json:"token"`
}

type AuthHandler struct {
	authService auth.Service
}

func NewAuthHandler(authService auth.Service) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (a *AuthHandler) Login(c *gin.Context) {
	token, err := a.authService.Login("john.doe@example.com") // FIXME
	if err != nil {
		abortRequestWithError(c, NewAppError(ErrorCodeValidationFailed, err.Error()))
		return
	}
	c.JSON(http.StatusOK, LoginResponsePayload{
		Token: token,
	})
}
