package httpservice

import (
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vivekweb2013/gitnoter/internal/applicationconfig"
)

// Run starts the http server.
func Run(applicationconfig *applicationconfig.ApplicationConfig) error {
	gin.SetMode(gin.ReleaseMode)
	if applicationconfig.Config.HTTPServer.Debug {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()
	router.UseRawPath = true

	noteHandler := NewNoteHandler(applicationconfig.GithubService, applicationconfig.UserService)
	loginHandler := NewLoginHandler(applicationconfig.AuthService, applicationconfig.GithubService, applicationconfig.UserService)
	userHandler := NewUserHandler(applicationconfig.UserService)
	preferenceHandler := NewPreferenceHandler(applicationconfig.PreferenceService, applicationconfig.GithubService, applicationconfig.UserService)
	authMiddleware := NewMiddleware(applicationconfig.AuthService)

	v1 := router.Group("api/v1")
	v1.GET("/user/me", authMiddleware.AuthorizeToken(), userHandler.Profile)
	v1.GET("/user/preference/repo", authMiddleware.AuthorizeToken(), preferenceHandler.GetRepos)
	v1.POST("/user/preference/repo", authMiddleware.AuthorizeToken(), preferenceHandler.SaveDefaultRepo)
	v1.POST("/user/preference/auto/repo", authMiddleware.AuthorizeToken(), preferenceHandler.AutoSetupRepo)

	v1.GET("/search/notes", authMiddleware.AuthorizeToken(), noteHandler.SearchNotes)  // search notes (provide filters using query-params)
	v1.GET("/tree/notes", authMiddleware.AuthorizeToken(), noteHandler.GetNotesTree)   // get complete notes repo tree
	v1.GET("/notes", authMiddleware.AuthorizeToken(), noteHandler.GetAllNotes)         // get all notes from path (provide filters using query-params)
	v1.GET("/notes/:path", authMiddleware.AuthorizeToken(), noteHandler.GetNote)       // get single note
	v1.POST("/notes/:path", authMiddleware.AuthorizeToken(), noteHandler.SaveNote)     // create/update single note
	v1.DELETE("/notes/:path", authMiddleware.AuthorizeToken(), noteHandler.DeleteNote) // delete single note

	v1.GET("/oauth2/login/github", loginHandler.GithubLogin)
	v1.GET("/oauth2/github/callback", loginHandler.GithubOAuth2Callback)

	address := net.JoinHostPort(applicationconfig.Config.HTTPServer.Host, applicationconfig.Config.HTTPServer.Port)
	server := http.Server{
		Addr:           address,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   2 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}
	return server.ListenAndServe()
}
