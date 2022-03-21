package httpservice

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	gh "github.com/google/go-github/v43/github"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/vivekweb2013/gitnoter/internal/auth"
	"github.com/vivekweb2013/gitnoter/internal/config"
	"github.com/vivekweb2013/gitnoter/internal/user"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
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

type UserResponsePayload struct {
	Email       string       `json:"email"`
	Name        string       `json:"name,omitempty"`
	Location    string       `json:"location,omitempty"`
	AvatarURL   string       `json:"avatar_url,omitempty"`
	DisabledAt  *time.Time   `json:"disabled_at,omitempty"`
	DefaultRepo *RepoPayload `json:"default_repo,omitempty"`
}

type GithubHandler struct {
	authService  auth.Service
	userService  user.Service
	oauth2Config oauth2.Config
}

func NewGithubHandler(authService auth.Service, userService user.Service, config config.OAuth2) *GithubHandler {
	oauth2Config := oauth2.Config{
		RedirectURL:  config.Github.RedirectURL,
		ClientID:     config.Github.ClientID,
		ClientSecret: config.Github.ClientSecret,
		Scopes:       []string{"read:user", "user:email", "repo"},
		Endpoint:     github.Endpoint,
	}
	return &GithubHandler{
		authService:  authService,
		userService:  userService,
		oauth2Config: oauth2Config,
	}
}

func (g *GithubHandler) GithubLogin(c *gin.Context) {
	state := uuid.NewString()
	c.SetCookie("state", state, 600, "/", "", true, true)

	// AuthCodeURL receive state that is a token to protect the user from CSRF attacks.
	// Generate a random `state` string and validate that it matches the `state` query parameter
	// on redirect callback
	url := g.oauth2Config.AuthCodeURL(state)

	// trigger authorization code grant flow
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (g *GithubHandler) GithubOAuth2Callback(c *gin.Context) {
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

	githubToken, err := g.oauth2Config.Exchange(c, code)
	if err != nil {
		logrus.Errorf("auth code exchange for token failed: %s", err.Error())
		c.Redirect(http.StatusTemporaryRedirect, failRedirectPath)
		return
	}

	client := gh.NewClient(g.oauth2Config.Client(c, githubToken))

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
	dbUser, err := g.userService.GetByEmail(*githubUser.Email)
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
	mapUserAttributes(&dbUser, *githubUser.Email, string(githubTokenJSON), githubUser)

	// create/update the user record
	userID, err := g.userService.Save(dbUser)
	if err != nil {
		logrus.Errorf("saving user to db failed: %s", err.Error())
		c.Redirect(http.StatusTemporaryRedirect, failRedirectPath)
		return
	}

	appToken, err := g.authService.GenerateToken(userID)
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

func (g *GithubHandler) GithubUserRepos(c *gin.Context) {
	dbUser, err := g.getUser(c)
	if err != nil {
		logrus.Errorf("retrieving user failed: %s", err.Error())
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	token := dbUser.GithubToken
	oauth2Token := oauth2.Token{}
	if err := json.Unmarshal([]byte(token), &oauth2Token); err != nil {
		logrus.Errorf("parsing user's github token failed: %s", err.Error())
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	client := gh.NewClient(g.oauth2Config.Client(c, &oauth2Token))
	repos, _, err := client.Repositories.List(c, "", nil)
	if err != nil {
		logrus.Errorf("retrieving user repos from github failed: %s", err.Error())
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if repos == nil {
		c.JSON(http.StatusOK, []gh.Repository{})
	}

	repoResp := make([]RepoPayload, 0, len(repos))
	for _, ghRepo := range repos {
		repoPayload := mapRepoAttributes(ghRepo)
		repoResp = append(repoResp, repoPayload)
	}

	c.JSON(http.StatusOK, repoResp)
}

func mapRepoAttributes(ghRepo *gh.Repository) RepoPayload {
	repoPayload := RepoPayload{}
	if ghRepo.Name != nil {
		repoPayload.Name = *ghRepo.Name
	}
	if ghRepo.Visibility != nil {
		repoPayload.Visibility = *ghRepo.Visibility
	}
	if ghRepo.DefaultBranch != nil {
		repoPayload.DefaultBranch = *ghRepo.DefaultBranch
	}
	return repoPayload
}

func (g *GithubHandler) getUser(c *gin.Context) (user.User, error) {
	userID, err := getUserIDFromContext(c)
	logrus.Infof("context has user-id: %d", userID)
	if err != nil {
		logrus.Errorf("fetching user-id from context failed")
		c.AbortWithStatus(http.StatusUnauthorized)
		return user.User{}, err
	}
	dbUser, err := g.userService.Get(userID)
	return dbUser, err
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
