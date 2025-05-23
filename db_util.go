package main

import (
	"database/sql"
	"gql/dbms"
	"gql/mysql"
	"gql/postgres"
	"gql/sqlite"
	"gql/util"
	"strings"
)

func (m *model) requiresAuth() bool {
	return m.dbms == nil || !m.dbms.HasUri()
}

func (m *model) changeCreds() {
	uri := m.dbms.GetUri()
	m.dbms.SetUri("")
	m.dbms.CloseDbConnection()
	m.auth.Reset(uri)
}

func (m *model) onDBConnect(db *sql.DB) {
	m.dbms.SetDb(db)

	if m.dbms.HasDb() {
		m.tabs.UpdateDBTable(m.dbms)
		m.tabs.Main.Panes.Db.Table.Focus()
		user, err := m.dbms.GetUser()
		if err != nil {
			m.tabs.Main.SetError(err)
		} else {
			m.user = user
		}
	}
}

func InitDBMS(uri string) dbms.DBMS {
	if uri == "" {
		return nil
	}

	if strings.HasPrefix(uri, "postgresql") {
		return &postgres.Model{Uri: uri}
	}

	if strings.HasPrefix(uri, "file:") {
		return &sqlite.Model{Uri: uri}
	}

	return &mysql.Model{Uri: uri}
}

func getDBUriFromCache() string {
	if util.CacheFileExists() {
		if uri := util.ReadFromCacheFile(); util.CheckDBMS(uri) == nil {
			return uri
		}
	}

	return ""
}
