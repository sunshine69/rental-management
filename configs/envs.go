package configs

import (
	_ "github.com/joho/godotenv/autoload"
	u "github.com/sunshine69/golang-tools/utils"
)

type Config struct {
	DBPath   string
	Port     string
	PathBase string
}

func InitConfig() Config {
	return Config{
		DBPath:   u.Getenv("DB_PATH", "test.sqlite3"),
		Port:     u.Getenv("PORT", "8080"),
		PathBase: u.Getenv("PATH_BASE", ""),
	}
}
