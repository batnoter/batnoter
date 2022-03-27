package httpservice

import (
	"net/http"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/sirupsen/logrus"
	"github.com/vivekweb2013/gitnoter/internal/github"
	"github.com/vivekweb2013/gitnoter/internal/preference"
	"github.com/vivekweb2013/gitnoter/internal/user"
)

type RepoPayload struct {
	Name          string `json:"name"`
	Visibility    string `json:"visibility"`
	DefaultBranch string `json:"default_branch"`
}

func (r RepoPayload) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(1, 50)),
		validation.Field(&r.Visibility, validation.Required, validation.Length(1, 20)),
		validation.Field(&r.Visibility, validation.Length(0, 50)),
	)
}

type PreferenceHandler struct {
	preferenceService preference.Service
	githubService     github.Service
	userService       user.Service
}

func NewPreferenceHandler(preferenceService preference.Service, githubService github.Service, userService user.Service) *PreferenceHandler {
	return &PreferenceHandler{
		preferenceService: preferenceService,
		githubService:     githubService,
		userService:       userService,
	}
}

func (p *PreferenceHandler) GetRepos(c *gin.Context) {
	user, err := p.getUser(c)
	if err != nil {
		logrus.Errorf("fetching user from context failed")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	logrus.WithField("user-id", user.ID).Info("request to retrieve repos started")
	gitRepos, err := p.githubService.GetRepos(c, parseOAuth2Token(user.GithubToken))
	if err != nil {
		logrus.Errorf("retrieving repos from github failed")
		abortRequestWithError(c, err)
		return
	}

	repos := make([]RepoPayload, 0, len(gitRepos))
	for _, gitRepo := range gitRepos {
		repo := RepoPayload{
			Name:          gitRepo.Name,
			Visibility:    gitRepo.Visibility,
			DefaultBranch: gitRepo.DefaultBranch,
		}
		repos = append(repos, repo)
	}

	c.JSON(http.StatusOK, repos)
	logrus.WithField("user-id", user.ID).Info("request to retrieve repos successful")
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
		abortRequestWithError(c, err)
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

func (p *PreferenceHandler) getUser(c *gin.Context) (user.User, error) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		logrus.Errorf("fetching user-id from context failed")
		return user.User{}, err
	}
	return p.userService.Get(userID)
}
