package httpservice

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	gh "github.com/google/go-github/v43/github"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/vivekweb2013/gitnoter/internal/auth"
	"github.com/vivekweb2013/gitnoter/internal/github"
	"github.com/vivekweb2013/gitnoter/internal/user"
)

type UserResponsePayload struct {
	Email       string       `json:"email"`
	Name        string       `json:"name,omitempty"`
	Location    string       `json:"location,omitempty"`
	AvatarURL   string       `json:"avatar_url,omitempty"`
	DisabledAt  *time.Time   `json:"disabled_at,omitempty"`
	DefaultRepo *RepoPayload `json:"default_repo,omitempty"`
}

type LoginHandler struct {
	authService  auth.Service
	githubServie github.Service
	userService  user.Service
}

func NewLoginHandler(authService auth.Service, githubServie github.Service, userService user.Service) *LoginHandler {
	return &LoginHandler{
		authService:  authService,
		githubServie: githubServie,
		userService:  userService,
	}
}

func (l *LoginHandler) GithubLogin(c *gin.Context) {
	state := uuid.NewString()
	c.SetCookie("state", state, 600, "/", "", true, true)

	url := l.githubServie.GetAuthCodeURL(state)

	// trigger authorization code grant flow
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (l *LoginHandler) GithubOAuth2Callback(c *gin.Context) {
	logrus.Info("github oauth2 callback started")
	state, _ := c.Cookie("state")
	stateFromCallback := c.Query("state")
	code := c.Query("code")
	failRedirectPath := "/?login_error=true"

	if stateFromCallback != state {
		logrus.Error("invalid oauth state")
		c.Redirect(http.StatusTemporaryRedirect, failRedirectPath)
		return
	}

	githubToken, err := l.githubServie.GetToken(c, code)
	if err != nil {
		logrus.Errorf("auth code exchange for token failed: %s", err.Error())
		c.Redirect(http.StatusTemporaryRedirect, failRedirectPath)
		return
	}

	githubUser, err := l.githubServie.GetUser(c, githubToken)
	if err != nil {
		logrus.Errorf("retrieving user from github failed: %s", err.Error())
		c.Redirect(http.StatusTemporaryRedirect, failRedirectPath)
		return
	}

	// get user from db if exists
	dbUser, err := l.userService.GetByEmail(*githubUser.Email)
	if err != nil {
		logrus.Errorf("retrieving user from db using email failed: %s", err.Error())
		c.Redirect(http.StatusTemporaryRedirect, failRedirectPath)
		return
	}
	githubTokenJSON, err := json.Marshal(githubToken)
	if err != nil {
		logrus.Errorf("converting github token to json failed: %s", err.Error())
		c.Redirect(http.StatusTemporaryRedirect, failRedirectPath)
		return
	}
	mapUserAttributes(&dbUser, string(githubTokenJSON), githubUser)

	// create/update the user record
	userID, err := l.userService.Save(dbUser)
	if err != nil {
		logrus.Errorf("saving user to db failed: %s", err.Error())
		c.Redirect(http.StatusTemporaryRedirect, failRedirectPath)
		return
	}

	appToken, err := l.authService.GenerateToken(userID)
	if err != nil {
		logrus.Errorf("token generation failed: %s", err.Error())
		c.Redirect(http.StatusTemporaryRedirect, failRedirectPath)
		return
	}

	// for security reasons, avoid using cookies to send the token to client
	// instead use html with a script that stores the token to localstorage and redirects to homepage
	// this is the only workaround to send token to client without using cookies
	// since the client(frontend) can only read headers/response with ajax request, and this call is not ajax
	c.Header("Content-Type", "text/html")
	c.String(200, `<!DOCTYPE html><html><body><script>(function(){localStorage.setItem("token","%s");location.replace("/");}());</script></body></html>`, appToken)
	logrus.Info("github oauth2 callback finished")
}

func mapUserAttributes(dbUser *user.User, ghToken string, githubUser gh.User) {
	dbUser.GithubToken = ghToken
	dbUser.Email = githubUser.GetEmail()
	dbUser.Name = githubUser.GetName()
	dbUser.Location = githubUser.GetLocation()
	dbUser.AvatarURL = githubUser.GetAvatarURL()
	dbUser.GithubID = githubUser.GetID()
	dbUser.GithubUsername = githubUser.GetLogin()
}
