package httpservice

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/vivekweb2013/gitnoter/internal/user"
)

type UserHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (u *UserHandler) Profile(c *gin.Context) {
	claims, _ := c.Get("claims")
	if claims == nil {
		logrus.Errorf("failed to get claims from context")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	email := claims.(jwt.MapClaims)["sub"].(string)
	dbUser, err := u.userService.GetByEmail(email)
	if err != nil {
		abortRequestWithError(c, err)
		return
	}
	c.JSON(http.StatusOK, UserResponsePayload{
		Email:      dbUser.Email,
		Name:       dbUser.Name,
		Location:   dbUser.Location,
		AvatarURL:  dbUser.AvatarURL,
		DisabledAt: dbUser.DisabledAt,
	})
}
