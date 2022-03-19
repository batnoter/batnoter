package httpservice

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
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

type UserResponsePayload struct {
	Email      string     `json:"email"`
	Name       string     `json:"name,omitempty"`
	Location   string     `json:"location,omitempty"`
	AvatarURL  string     `json:"avatar_url,omitempty"`
	DisabledAt *time.Time `json:"disabled_at,omitempty"`
}

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

	githubToken, err := a.oauth2Config.Exchange(c, code)
	if err != nil {
		logrus.Errorf("auth code exchange for token failed: %s", err.Error())
		c.Redirect(http.StatusTemporaryRedirect, failRedirectPath)
		return
	}

	client := gh.NewClient(a.oauth2Config.Client(c, githubToken))

	githubUser, _, err := client.Users.Get(c, "")
	if err != nil {
		logrus.Errorf("get user from github failed: %s", err.Error())
		c.Redirect(http.StatusTemporaryRedirect, failRedirectPath)
		return
	}
	if githubUser == nil || githubUser.Email == nil {
		logrus.Errorf("processing github user object failed: %s", err.Error())
		c.Redirect(http.StatusTemporaryRedirect, failRedirectPath)
		return
	}

	// check if user record already exist
	dbUser, err := a.userService.GetByEmail(*githubUser.Email)
	if err != nil {
		logrus.Errorf("retrieving user from db using email failed: %s", err.Error())
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
	logrus.Info("github oauth2 callback finished")
}

func (a *AuthHandler) Profile(c *gin.Context) {
	claims, _ := c.Get("claims")
	if claims == nil {
		logrus.Errorf("failed to get claims from context")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	email := claims.(jwt.MapClaims)["sub"].(string)
	u, err := a.userService.GetByEmail(email)
	if err != nil {
		abortRequestWithError(c, err)
		return
	}
	c.JSON(http.StatusOK, UserResponsePayload{
		Email:      u.Email,
		Name:       u.Name,
		Location:   u.Location,
		AvatarURL:  u.AvatarURL,
		DisabledAt: u.DisabledAt,
	})
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
