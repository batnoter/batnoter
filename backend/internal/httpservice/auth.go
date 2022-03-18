package httpservice

import (
	"net/http"

	"github.com/gin-gonic/gin"
	gh "github.com/google/go-github/v43/github"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/vivekweb2013/gitnoter/internal/auth"
	"github.com/vivekweb2013/gitnoter/internal/config"
	"github.com/vivekweb2013/gitnoter/internal/user"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type AuthHandler struct {
	authService  auth.Service
	userService  user.Service
	oauth2Config oauth2.Config
}

func NewAuthHandler(authService auth.Service, userService user.Service, config config.OAuth2) *AuthHandler {
	oauth2Config := oauth2.Config{
		RedirectURL:  config.Github.RedirectURL,
		ClientID:     config.Github.ClientID,
		ClientSecret: config.Github.ClientSecret,
		Scopes:       []string{"read:user", "user:email", "repo"},
		Endpoint:     github.Endpoint,
	}
	return &AuthHandler{
		authService:  authService,
		userService:  userService,
		oauth2Config: oauth2Config,
	}
}

func (a *AuthHandler) GithubLogin(c *gin.Context) {
	state := uuid.NewString()
	c.SetCookie("state", state, 600, "/", "", true, true)

	// AuthCodeURL receive state that is a token to protect the user from CSRF attacks.
	// Generate a random `state` string and validate that it matches the `state` query parameter
	// on redirect callback
	url := a.oauth2Config.AuthCodeURL(state)

	// trigger authorization code grant flow
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (a *AuthHandler) GithubOAuth2Callback(c *gin.Context) {
	state, _ := c.Cookie("state")
	stateFromCallback := c.Query("state")
	code := c.Query("code")
	failRedirectPath := "/?error=true"

	if stateFromCallback != state {
		logrus.Error("invalid oauth state")
		c.Redirect(http.StatusTemporaryRedirect, failRedirectPath)
		return
	}

	githubToken, err := a.oauth2Config.Exchange(c, code)
	if err != nil {
		logrus.Errorf("auth code exchange for token failed: %s", err.Error())
		c.Redirect(http.StatusTemporaryRedirect, failRedirectPath)
		return
	}

	client := gh.NewClient(a.oauth2Config.Client(c, githubToken))

	githubUser, _, err := client.Users.Get(c, "")
	if err != nil {
		logrus.Errorf("auth code exchange for token failed: %s", err.Error())
		c.Redirect(http.StatusTemporaryRedirect, failRedirectPath)
		return
	}
	if githubUser == nil || githubUser.Email == nil {
		logrus.Errorf("failed to process github user object: %s", err.Error())
		c.Redirect(http.StatusTemporaryRedirect, failRedirectPath)
		return
	}

	// check if user record already exist
	dbUser, err := a.userService.GetByEmail(*githubUser.Email)
	if err != nil {
		logrus.Errorf("retrieving user with email failed: %s", err.Error())
		c.Redirect(http.StatusTemporaryRedirect, failRedirectPath)
		return
	}
	mapUserAttributes(&dbUser, *githubUser.Email, githubToken.AccessToken, githubUser)

	// create/update the user record
	a.userService.Save(dbUser)

	appToken, err := a.authService.GenerateToken(dbUser.Email)
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
}

func mapUserAttributes(dbUser *user.User, email string, githubToken string, githubUser *gh.User) {
	dbUser.Email = email
	dbUser.GithubToken = githubToken

	if githubUser.Name != nil {
		dbUser.Name = *githubUser.Name
	}
	if githubUser.Location != nil {
		dbUser.Location = *githubUser.Location
	}
	if githubUser.AvatarURL != nil {
		dbUser.AvatarURL = *githubUser.AvatarURL
	}
	if githubUser.ID != nil {
		dbUser.GithubID = *githubUser.ID
	}
	if githubUser.Login != nil {
		dbUser.GithubUsername = *githubUser.Login
	}
}
