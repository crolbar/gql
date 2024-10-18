package panes

import (
	"database/sql"
	"gql/mysql"
)

func (p *Panes) UpdateDBTable(db *sql.DB) {
    cols, rows := mysql.GetDatabases(db)

    p.Db.Table.SetMaxWidth(cols[0].Width + 2)

    p.Db.Table.SetColumns(cols)
    p.Db.Table.SetRows(rows)

    p.UpdateDBTablesTable(db)
}

func (p *Panes) UpdateDBTablesTable(db *sql.DB) {
    p.currDB = p.Db.Table.GetSelectedRow()[0]

    cols, rows := mysql.GetTables(db, p.currDB)

    p.DbTables.Table.SetMaxWidth(cols[0].Width + 2)

    p.DbTables.Table.SetColumns(cols)
    p.DbTables.Table.SetRows(rows)

    p.UpdateMainTable(db)
}

func (p *Panes) UpdateMainTable(db *sql.DB) {
    p.currDBTable = p.DbTables.Table.GetSelectedRow()[0]

    cols, rows := mysql.GetTable(db, p.currDB, p.currDBTable)

    p.Main.Table.SetColumns(cols)
    p.Main.Table.SetRows(rows)
}

