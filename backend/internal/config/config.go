package config

type App struct {
	SecretKey string
}

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
	App        App
	Database   Database
	HTTPServer HTTPServer
}
