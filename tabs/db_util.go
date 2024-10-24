package tabs

import (
	"database/sql"
	"gql/mysql"
)

func (t *Tabs) UpdateDBTable(db *sql.DB) {

    whereClause := t.whereClauses["db"]
    cols, rows, err := mysql.GetDatabases(db, whereClause)

    t.Main.SetError(err)

    if err != nil {
        t.Main.Panes.Db.Table.SetRows(nil)

        t.Main.Panes.DbTables.Table.SetColumns(nil)
        t.Main.Panes.DbTables.Table.SetRows(nil)

        t.Main.Panes.Main.Table.SetColumns(nil)
        t.Main.Panes.Main.Table.SetRows(nil)
        return
    }

    t.Main.Panes.Db.Table.SetMaxWidth(cols[0].Width + 2)

    t.Main.Panes.Db.Table.SetColumns(cols)
    t.Main.Panes.Db.Table.SetRows(rows)

    t.UpdateDBTablesTable(db)
}

func (t *Tabs) UpdateDBTablesTable(db *sql.DB) {
    t.currDB = t.Main.Panes.Db.Table.GetSelectedRow()[0]

    whereClause := t.whereClauses[t.currDB]
    cols, rows, err := mysql.GetTables(db, t.currDB, whereClause)

    t.Main.SetError(err)

    if err != nil {
        t.Main.Panes.DbTables.Table.SetColumns(nil)
        t.Main.Panes.DbTables.Table.SetRows(nil)

        t.Main.Panes.Main.Table.SetColumns(nil)
        t.Main.Panes.Main.Table.SetRows(nil)
        return
    }

    t.Main.Panes.DbTables.Table.SetMaxWidth(cols[0].Width + 2)
    t.Main.Panes.DbTables.Table.SetColumns(cols)
    t.Main.Panes.DbTables.Table.SetRows(rows)

    t.UpdateMainTable(db)
}

func (t *Tabs) UpdateMainTable(db *sql.DB) {
    t.currDBTable = t.Main.Panes.DbTables.Table.GetSelectedRow()[0]

    whereClause := t.whereClauses[t.currDB + "/" + t.currDBTable]
    cols, rows, err := mysql.GetTable(db, t.currDB, t.currDBTable, whereClause)

    t.Main.SetError(err)

    if err != nil {
        t.Main.Panes.Main.Table.SetColumns(nil)
        t.Main.Panes.Main.Table.SetRows(nil)
        return
    }

    t.Main.Panes.Main.Table.SetColumns(cols)
    t.Main.Panes.Main.Table.SetRows(rows)
}

func (t *Tabs) UpdateDescribeTable(db *sql.DB) {

    cols, rows, err := mysql.GetDescribe(db, t.currDB, t.currDBTable)

    t.Main.SetError(err)

    if err != nil {
        t.Describe.Table.SetColumns(nil)
        t.Describe.Table.SetRows(nil)
        return
    }

    t.Describe.Table.SetColumns(cols)
    t.Describe.Table.SetRows(rows)
}

func (t *Tabs) DeleteSelectedDb(db *sql.DB) error {
    return mysql.DeleteDB(db, t.currDB)
}

func (t *Tabs) DeleteSelectedRow(db *sql.DB) error {
    return mysql.DeleteRow(
        db,
        t.currDB,
        t.currDBTable,
        t.Main.Panes.Main.Table.GetSelectedRow(),
        t.Main.Panes.Main.Table.GetCols(),
    )
}

func (t *Tabs) UpdateSelectedCell(db *sql.DB, value string) error {
    return mysql.UpdateCell(
        db,
        t.currDB,
        t.currDBTable,
        t.Main.Panes.Main.Table.GetSelectedRow(),
        t.Main.Panes.Main.Table.GetCols(),
        t.Main.Panes.Main.Table.GetCursor().X,
        value,
    )
}

func (t *Tabs) SendQuery(db *sql.DB, query string) {
    t.Main.SetError(mysql.SendQuery(db, query))
}
