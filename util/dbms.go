package util

import (
	"database/sql"
	"strings"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func CheckDBMS(uri string) error {
	var db *sql.DB
	var err error

	if strings.HasPrefix(uri, "postgresql") {
		db, err = sql.Open("postgres", uri)
	} else if strings.HasPrefix(uri, "file:") {
		db, err = sql.Open("sqlite3", strings.ReplaceAll(uri, "file:", ""))
	} else {
		db, err = sql.Open("mysql", uri)
	}

	if err != nil {
		return err
	}

	return db.Ping()
}
