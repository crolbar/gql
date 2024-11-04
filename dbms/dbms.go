package dbms

import (
	"database/sql"
	"gql/table"

	tea "github.com/charmbracelet/bubbletea"
)

type DbConnectMsg struct {Db *sql.DB}

type DBMS interface {
    // Uri has to be set before calling Open
    Open() tea.Cmd
    HasDb() bool
    SetDb(*sql.DB)
    CloseDbConnection()

    HasUri() bool
    SetUri(string)
    GetUri() string

    GetDatabases(
        whereClause string,
    ) ([]table.Column, []table.Row, error) 

    GetDBTables(
        dbName,
        whereClause string,
    ) ([]table.Column, []table.Row, error)

    GetTable(
        currDB,
        selTable,
        whereClause string,
    ) ([]table.Column, []table.Row, error)

    GetDescribe(
        currDB,
        selTable string,
    ) ([]table.Column, []table.Row, error)

    GetUser() (string, error)

    DeleteDB(dbName string) error
    DeleteDBTable(dbName, selTable string) error
    DeleteRow(
        dbName,
        tableName string,
        row table.Row,
        cols []table.Column,
    ) error

    UpdateCell(
        dbName,
        tableName string,
        row table.Row,
        cols []table.Column,
        selectedCol int,
        value string,
    ) error

    ChangeDbTableName(
        dbName,
        tableName string,
        value string,
    ) error

    SendQuery(query string) error
}
