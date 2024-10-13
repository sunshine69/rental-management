package configs

import (
	"fmt"
	"strconv"

	"github.com/joho/godotenv"
	u "github.com/sunshine69/golang-tools/utils"
	"github.com/sunshine69/rental-management/utils"
)

type Config struct {
	PublicHost             string
	Port                   string
	PathBase               string
	DBUser                 string
	DBPassword             string
	DBAddress              string
	DBName                 string
	JWTSecret              string
	JWTExpirationInSeconds int64
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		PublicHost:             u.Getenv("PUBLIC_HOST", "http://localhost"),
		Port:                   u.Getenv("PORT", "8080"),
		DBUser:                 u.Getenv("DB_USER", "root"),
		DBPassword:             u.Getenv("DB_PASSWORD", "mypassword"),
		DBAddress:              fmt.Sprintf("%s:%s", u.Getenv("DB_HOST", "127.0.0.1"), u.Getenv("DB_PORT", "3306")),
		DBName:                 u.Getenv("DB_NAME", "ecom"),
		JWTSecret:              u.Getenv("JWT_SECRET", "not-so-secret-now-is-it?"),
		JWTExpirationInSeconds: utils.Must(strconv.ParseInt(u.Getenv("JWT_EXPIRATION_IN_SECONDS", "3600"), 10, 64)),
	}
}
