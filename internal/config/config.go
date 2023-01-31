package config

type Config struct {
	LogLevel   string
	DB         Database
	ListenPort int
	ApiURL     string
}

type Database struct {
	DSN string
}
