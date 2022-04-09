package httpservice

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/sirupsen/logrus"
	"github.com/vivekweb2013/gitnoter/internal/github"
	"github.com/vivekweb2013/gitnoter/internal/preference"
	"github.com/vivekweb2013/gitnoter/internal/user"
)

// RepoPayload represents the http request/response payload of repository entity.
type RepoPayload struct {
	Name          string `json:"name"`
	Visibility    string `json:"visibility"`
	DefaultBranch string `json:"default_branch"`
}

// Validate validates the repo http request payload.
func (r RepoPayload) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(1, 50)),
		validation.Field(&r.Visibility, validation.Required, validation.Length(1, 20)),
		validation.Field(&r.Visibility, validation.Length(0, 50)),
	)
}

// PreferenceHandler represents http handler for managing user preferences.
type PreferenceHandler struct {
	preferenceService preference.Service
	githubService     github.Service
	userService       user.Service
}

// NewPreferenceHandler creates and returns a new preference handler.
func NewPreferenceHandler(preferenceService preference.Service, githubService github.Service, userService user.Service) *PreferenceHandler {
	return &PreferenceHandler{
		preferenceService: preferenceService,
		githubService:     githubService,
		userService:       userService,
	}
}

// GetRepos returns user repositories of logged in user.
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

// SaveDefaultRepo stores the requested repo as user's default repo.
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
		logrus.Errorf("saving user's default repo failed")
		abortRequestWithError(c, err)
		return
	}
	c.Status(http.StatusOK)
	logrus.WithField("user-id", userID).Info("request to link default repo successful")
}

// AutoSetupRepo creates a new notes repo in user's github account and stores it as user's default repo preference.
func (p *PreferenceHandler) AutoSetupRepo(c *gin.Context) {
	repoName := c.Query("repoName")
	if err := validation.Validate(repoName, validation.Required); err != nil {
		abortRequestWithError(c, NewAppError(ErrorCodeValidationFailed, fmt.Sprintf("repoName: %s", err.Error())))
		return
	}
	user, err := p.getUser(c)
	if err != nil {
		logrus.Errorf("fetching user from context failed")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	logrus.WithField("user-id", user.ID).Infof("request to auto setup default repo")
	err = p.githubService.CreateRepo(c, parseOAuth2Token(user.GithubToken), repoName)
	if err != nil {
		logrus.Errorf("creating a new notes repo in github failed")
		abortRequestWithError(c, err)
		return
	}

	defaultRepo := preference.DefaultRepo{
		UserID:        user.ID,
		Name:          repoName,
		Visibility:    "private",
		DefaultBranch: "main",
	}

	if err := p.preferenceService.Save(defaultRepo); err != nil {
		logrus.Errorf("saving user's default repo failed")
		abortRequestWithError(c, err)
		return
	}
	c.Status(http.StatusOK)
	logrus.WithField("user-id", user.ID).Info("request to auto setup default repo successful")
}

func (p *PreferenceHandler) getUser(c *gin.Context) (user.User, error) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		logrus.Errorf("fetching user-id from context failed")
		return user.User{}, err
	}
	return p.userService.Get(userID)
}
