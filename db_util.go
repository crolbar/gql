package main

import (
	"database/sql"
	"gql/dbms"
	"gql/mysql"
	"gql/util"
	"strings"
)

func (m *model) requiresAuth() bool {
    return m.uri == ""
}

func (m *model) changeCreds() {
    m.uri = ""
    m.auth.ResetAll()
}

func (m *model) onDBConnect(db *sql.DB) {
    m.dbms.SetDb(db)

    if (m.dbms.HasDb()) {
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
    if (uri == "") {
        return nil
    }

    if (strings.HasPrefix(uri, "postgresql")) {
        return nil // TODO
    }

    return &mysql.Model{}
}

func getDBUriFromCache() string {
    if util.CacheFileExists() {
        if uri := util.ReadFromCacheFile(); util.CheckMysql(uri) == nil {
            return uri
        }
    }

    return ""
}
