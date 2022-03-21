package httpservice

import (
	"errors"
	"net/http"
	"strconv"

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
	userID, err := getUserIDFromContext(c)
	if err != nil {
		logrus.Errorf("fetching user-id from context failed")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	logrus.WithField("user-id", userID).Info("request for profile started")
	dbUser, err := u.userService.Get(userID)
	if err != nil {
		abortRequestWithError(c, err)
		return
	}
	userResp := UserResponsePayload{
		Email:      dbUser.Email,
		Name:       dbUser.Name,
		Location:   dbUser.Location,
		AvatarURL:  dbUser.AvatarURL,
		DisabledAt: dbUser.DisabledAt,
	}
	if dbUser.DefaultRepo != nil {
		userResp.DefaultRepo = &RepoPayload{
			Name:          dbUser.DefaultRepo.Name,
			Visibility:    dbUser.DefaultRepo.Visibility,
			DefaultBranch: dbUser.DefaultRepo.DefaultBranch,
		}
	}
	c.JSON(http.StatusOK, userResp)
	logrus.WithField("user-id", userID).Info("request for profile successful")
}

func getUserIDFromContext(c *gin.Context) (uint, error) {
	claims, _ := c.Get("claims")
	if claims == nil {
		return 0, errors.New("fatching claims from context failed")
	}
	userIDString := claims.(jwt.MapClaims)["sub"].(string)
	userID, err := strconv.ParseUint(userIDString, 10, 64)
	if err != nil {
		return 0, errors.New("parsing user-id from token failed")
	}
	return uint(userID), nil
}
