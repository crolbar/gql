package main

import (
	"database/sql"
	"gql/util"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

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
        m.panes.UpdateDBTable(db)
        m.panes.Db.Table.Focus()
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
