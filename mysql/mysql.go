package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"gql/table"

	_ "github.com/go-sql-driver/mysql"
)

func GetDatabases(db *sql.DB) []string {
    rows, err := db.Query("show databases;")
	if err != nil {
		log.Fatal(err)
	}
    var databases []string
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			log.Fatal(err)
		}
		databases = append(databases, dbName)
	}

    return databases
}

func GetTables(db *sql.DB, dbName string) []string {
    rows, err := db.Query(fmt.Sprintf("show tables from %s;", dbName))
	if err != nil {
		log.Fatal(err)
	}

    var tables []string
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			log.Fatal(err)
		}
		tables = append(tables, dbName)
	}

    return tables
}

func GetTable(db *sql.DB, currDB, selTable string) ([]table.Column, []table.Row) {
    rowsRes, err := db.Query(fmt.Sprintf("select * from %s.%s", currDB, selTable))
	if err != nil {
		log.Fatal(err)
	}

	columnsRes, err := rowsRes.Columns()
	if err != nil {
		log.Fatal(err)
	}

	values        := make([]interface{}, len(columnsRes))
	valuePointers := make([]interface{}, len(columnsRes))

    var rows []table.Row

    for rowsRes.Next() {
        for i := range columnsRes {
			valuePointers[i] = &values[i]
		}

        var currRow table.Row

        if err := rowsRes.Scan(valuePointers...); err != nil {
			log.Fatal(err)
		}

        for i := range columnsRes {
			switch val := values[i].(type) {
			case nil:
                currRow = append(currRow, "NULL")
			case []byte:
                text := string(val)
                text = strings.ReplaceAll(text, "\n", "\\n")
                currRow = append(currRow, text)
            case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
                currRow = append(currRow, fmt.Sprintf("%d", val))
            case float32, float64:
                currRow = append(currRow, fmt.Sprintf("%f", val))
			default:
                log.Fatalf("Found type that's not supported, val: %v, type: %T", val, val)
			}
        }

        rows = append(rows, currRow)
    }

    columns := make([]table.Column, 0, len(columnsRes))
    for _, col := range columnsRes {
        columns = append(columns, table.Column {
            Title: col,
            Width: 10,
        })
    }


    return columns, rows
}
