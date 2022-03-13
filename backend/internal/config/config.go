package config

type Database struct {
	Host       string
	Port       string
	DBName     string
	Username   string
	Password   string
	DriverName string
	Debug      bool
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
