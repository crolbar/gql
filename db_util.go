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

func (m *model) updateDBTable() {
    dbs := mysql.GetDatabases(m.db)

    rows := make([]table.Row, 0, len(dbs))
    cols := []table.Column { {Title: "Databases", Width: 20}, }

    for i := 0; i < len(dbs); i++ {
        rows = append(rows, []string{dbs[i]})
    }

    m.dbPane.table.SetColumns(cols)
    m.dbPane.table.SetRows(rows)

    m.updateDBTablesTable()
}

func (m *model) updateDBTablesTable() {
    m.currDB = m.dbPane.table.GetSelectedRow()[0]

    tables := mysql.GetTables(m.db, m.currDB)
    rows := make([]table.Row, 0, len(tables))
    cols := []table.Column { {Title: fmt.Sprintf("tables in %s", m.currDB), Width: 20}, }

    for i := 0; i < len(tables); i++ {
        rows = append(rows, []string{tables[i]})
    }

    m.dbTablesPane.table.SetColumns(cols)
    m.dbTablesPane.table.SetRows(rows)

    m.updateMainTable()
}

func (m *model) updateMainTable() {
    m.currDBTable = m.dbTablesPane.table.GetSelectedRow()[0]

    cols, rows := mysql.GetTable(m.db, m.currDB, m.currDBTable)

    m.mainPane.table.SetColumns(cols)
    m.mainPane.table.SetRows(rows)
}

func (m model) openMysql() tea.Msg {
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
        m.updateDBTable()
        m.dbPane.table.Focus()
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
