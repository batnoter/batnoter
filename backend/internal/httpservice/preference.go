package httpservice

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/vivekweb2013/gitnoter/internal/preference"
)

type PreferenceHandler struct {
	preferenceService preference.Service
}

func NewPreferenceHandler(preferenceService preference.Service) *PreferenceHandler {
	return &PreferenceHandler{preferenceService: preferenceService}
}

func (p *PreferenceHandler) SaveDefaultRepo(c *gin.Context) {
	var repoPayload RepoPayload
	c.BindJSON(&repoPayload)
	if err := repoPayload.Validate(); err != nil {
		abortRequestWithError(c, NewAppError(ErrorCodeValidationFailed, err.Error()))
		return
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		logrus.Errorf("fetching user-id from context failed")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	logrus.WithField("user-id", userID).Infof("request to link default repo: %s", repoPayload.Name)
	dbDefaultRepo, err := p.preferenceService.GetByUserID(userID)
	if err != nil {
		logrus.Errorf("retrieving user's default repo failed")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	dbDefaultRepo.UserID = userID
	dbDefaultRepo.Name = repoPayload.Name
	dbDefaultRepo.Visibility = repoPayload.Visibility
	dbDefaultRepo.DefaultBranch = repoPayload.DefaultBranch

	if err := p.preferenceService.Save(dbDefaultRepo); err != nil {
		abortRequestWithError(c, err)
		return
	}
	c.Status(http.StatusOK)
	logrus.WithField("user-id", userID).Info("request to link default repo successful")
}
