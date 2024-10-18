package panes

import (
	"database/sql"
	"fmt"
	"gql/mysql"
	"gql/table"
)

func (p *Panes) UpdateDBTable(db *sql.DB) {
    dbs := mysql.GetDatabases(db)

    rows := make([]table.Row, 0, len(dbs))
    cols := []table.Column { {Title: "Databases", Width: 20}, }

    for i := 0; i < len(dbs); i++ {
        rows = append(rows, []string{dbs[i]})
    }

    p.Db.Table.SetMaxWidth(cols[0].Width + 2)

    p.Db.Table.SetColumns(cols)
    p.Db.Table.SetRows(rows)

    p.UpdateDBTablesTable(db)
}

func (p *Panes) UpdateDBTablesTable(db *sql.DB) {
    p.currDB = p.Db.Table.GetSelectedRow()[0]

    tables := mysql.GetTables(db, p.currDB)
    rows := make([]table.Row, 0, len(tables))
    cols := []table.Column { {Title: fmt.Sprintf("tables in %s", p.currDB), Width: 20}, }

    for i := 0; i < len(tables); i++ {
        rows = append(rows, []string{tables[i]})
    }

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

