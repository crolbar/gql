package util

import (
	"database/sql"
	"log"
)

func CheckMysql(uri string) error {
    db, err := sql.Open("mysql", uri)
	if err != nil {
		log.Fatal(err)
	}

    return db.Ping()
}
