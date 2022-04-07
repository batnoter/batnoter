package config

// App represents configuration properties specific to the application.
type App struct {
	SecretKey string
}

// Database represents configuration properties required to connect to a database.
type Database struct {
	Host       string
	Port       string
	DBName     string
	Username   string
	Password   string
	DriverName string
	SSLMode    string
	Debug      bool
}

// HTTPServer represents configuration properties required for starting http server.
type HTTPServer struct {
	Host  string
	Port  string
	Debug bool
}

// OAuth2 represents configuration grouped by the oauth2 provider.
type OAuth2 struct {
	Github Github
}

// Github represents configuration properties required consume github oauth2 api.
type Github struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// Config represents all the application configurations grouped as per their category.
type Config struct {
	App        App
	Database   Database
	HTTPServer HTTPServer
	OAuth2     OAuth2
}
