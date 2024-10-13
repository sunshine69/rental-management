package model

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mutecomm/go-sqlcipher/v4"
	u "github.com/sunshine69/golang-tools/utils"
)

// DbConn - Global DB connection
var DB *sqlx.DB

func init() {
	DB = sqlx.MustConnect("sqlite3", u.Getenv("DB_PATH", "test.sqlite3"))
}

type Searchable struct {
	Where string
}
