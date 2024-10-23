package tabs

import (
	"database/sql"
	"gql/mysql"
	"log"
)

func (t *Tabs) UpdateDBTable(db *sql.DB) {
    cols, rows := mysql.GetDatabases(db)

    t.Main.Panes.Db.Table.SetMaxWidth(cols[0].Width + 2)

    t.Main.Panes.Db.Table.SetColumns(cols)
    t.Main.Panes.Db.Table.SetRows(rows)
    t.UpdateDBTablesTable(db)
}

func (t *Tabs) UpdateDBTablesTable(db *sql.DB) {
    t.currDB = t.Main.Panes.Db.Table.GetSelectedRow()[0]

    cols, rows := mysql.GetTables(db, t.currDB)

    t.Main.Panes.DbTables.Table.SetMaxWidth(cols[0].Width + 2)
    t.Main.Panes.DbTables.Table.SetColumns(cols)
    t.Main.Panes.DbTables.Table.SetRows(rows)

    t.UpdateMainTable(db)
}

func (t *Tabs) UpdateMainTable(db *sql.DB) {
    t.currDBTable = t.Main.Panes.DbTables.Table.GetSelectedRow()[0]

    cols, rows := mysql.GetTable(db, t.currDB, t.currDBTable)

    t.Main.Panes.Main.Table.SetColumns(cols)
    t.Main.Panes.Main.Table.SetRows(rows)
}

func (t *Tabs) UpdateDescribeTable(db *sql.DB) {

    cols, rows := mysql.GetDescribe(db, t.currDB, t.currDBTable)

    t.Describe.Table.SetColumns(cols)
    t.Describe.Table.SetRows(rows)
}

func (t *Tabs) DeleteSelectedDb(db *sql.DB) {
    err := mysql.DeleteDB(db, t.currDB)
    // TODO: handle error
    if err != nil {
        log.Fatal(err)
    }
}

func (t *Tabs) DeleteSelectedRow(db *sql.DB) {
    err := mysql.DeleteRow(
        db,
        t.currDB,
        t.currDBTable,
        t.Main.Panes.Main.Table.GetSelectedRow(),
        t.Main.Panes.Main.Table.GetCols(),
    )

    // TODO: handle error
    if err != nil {
        log.Fatal(err)
    }
}

func (t *Tabs) UpdateSelectedCell(db *sql.DB, value string) {
    err := mysql.UpdateCell(
        db,
        t.currDB,
        t.currDBTable,
        t.Main.Panes.Main.Table.GetSelectedRow(),
        t.Main.Panes.Main.Table.GetCols(),
        t.Main.Panes.Main.Table.GetCursor().X,
        value,
    )

    // TODO: handle error
    if err != nil {
        log.Fatal(err)
    }
}
