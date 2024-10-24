package util

import (
	"database/sql"
)

func CheckMysql(uri string) error {
    db, err := sql.Open("mysql", uri)
	if err != nil {
        return err
	}

    return db.Ping()
}
