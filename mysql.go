package main

import (
	"database/sql"
	"fmt"
	"log"

	"gql/table"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/go-sql-driver/mysql"
)

func (m *model) UpdateMainTable() {
    cols, rows := GetTable(m.db, m.currDB, m.currTable)

    m.mainTable.SetColumns(cols)
    m.mainTable.SetRows(rows)
}

func (m *model) UpdateDBTablesTable() {
    tables := GetTables(m.db, m.currDB)
    rows := make([]table.Row, 0, len(tables))
    cols := []table.Column { {Title: fmt.Sprintf("tables in %s", m.currDB), Width: 20}, }

    for i := 0; i < len(tables); i++ {
        rows = append(rows, []string{tables[i]})
    }

    m.DBTablesTable.SetColumns(cols)
    m.DBTablesTable.SetRows(rows)
}

func (m *model) UpdateDBTable() {
    dbs := GetDatabases(m.db)
    rows := make([]table.Row, 0, len(dbs))
    cols := []table.Column { {Title: "Databases", Width: 20}, }

    for i := 0; i < len(dbs); i++ {
        rows = append(rows, []string{dbs[i]})
    }

    m.DBTable.SetColumns(cols)
    m.DBTable.SetRows(rows)
}

const uri = "crolbar:@tcp(127.0.0.1:3306)/"

// TODO: make error msgs (wrong creds for eg)
func OpenMysql() tea.Msg {
    db, err := sql.Open("mysql", uri)
	if err != nil {
		log.Fatal(err)
	}

    err = db.Ping()
	if err != nil {
		return dbConnectMsg {nil}
	}
    return dbConnectMsg {db}
}

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
                currRow = append(currRow, string(val))
			default:
                currRow = append(currRow, "")
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
