package config

type Database struct {
	Host     string
	Port     string
	Name     string
	Username string
	Password string
	DebugLog bool
}

type HTTPServer struct {
	Host  string
	Port  string
	Debug bool
}

type Config struct {
	Database   Database
	HTTPServer HTTPServer
}
