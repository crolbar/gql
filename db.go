package main

import (
	"database/sql"
	"fmt"
	"gql/table"
	"gql/util"
    "gql/mysql"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) updateCurrDB() {
    m.currDB = m.DBTable.GetSelectedRow()[0]
    m.UpdateDBTablesTable()
    m.updateMainTable()
}

func (m *model) updateMainTable() {
    m.currTable = m.DBTablesTable.GetSelectedRow()[0]

    cols, rows := mysql.GetTable(m.db, m.currDB, m.currTable)

    m.mainTable.SetColumns(cols)
    m.mainTable.SetRows(rows)
}

func (m *model) UpdateDBTablesTable() {
    tables := mysql.GetTables(m.db, m.currDB)
    rows := make([]table.Row, 0, len(tables))
    cols := []table.Column { {Title: fmt.Sprintf("tables in %s", m.currDB), Width: 20}, }

    for i := 0; i < len(tables); i++ {
        rows = append(rows, []string{tables[i]})
    }

    m.DBTablesTable.SetColumns(cols)
    m.DBTablesTable.SetRows(rows)
}

func (m *model) UpdateDBTable() {
    dbs := mysql.GetDatabases(m.db)
    rows := make([]table.Row, 0, len(dbs))
    cols := []table.Column { {Title: "Databases", Width: 20}, }

    for i := 0; i < len(dbs); i++ {
        rows = append(rows, []string{dbs[i]})
    }

    m.DBTable.SetColumns(cols)
    m.DBTable.SetRows(rows)
}

func (m model) OpenMysql() tea.Msg {
    if m.requiresAuth() {
        return nil
    }

    db, err := sql.Open("mysql", m.uri)
	if err != nil {
		log.Fatal(err)
	}

    err = db.Ping()
	if err != nil {
		return dbConnectMsg {nil}
	}
    return dbConnectMsg {db}
}

func (m *model) requiresAuth() bool {
    return m.uri == ""
}

func (m *model) changeCreds() {
    m.uri = ""
    m.auth.ResetAll()
}

func (m *model) onDBConnect(db *sql.DB) {
    m.db = db
    if (m.db != nil) {
        m.UpdateDBTable()
        m.updateCurrDB()
        m.DBTable.Focus()
    }
}

func getDBUriFromCache() string {
    if util.CacheFileExists() {
        if uri := util.ReatFromCacheFile(); util.CheckMysql(uri) == nil {
            return uri
        }
    }

    return ""
}
