package util

import (
	"database/sql"
	"strings"

	_ "github.com/lib/pq"
)

func CheckDBMS(uri string) error {
	var db *sql.DB
	var err error

	if strings.HasPrefix(uri, "postgresql") {
		db, err = sql.Open("postgres", uri)
	} else {
		db, err = sql.Open("mysql", uri)
	}

	if err != nil {
		return err
	}

	return db.Ping()
}
