package tabs

import (
	"gql/dbms"
	"errors"
)

func (t *Tabs) UpdateDBTable(db dbms.DBMS) {
    whereClause := t.whereClauses["db"]
    cols, rows, err := db.GetDatabases(whereClause)

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

func (t *Tabs) UpdateDBTablesTable(db dbms.DBMS) {
    selRow := t.Main.Panes.Db.Table.GetSelectedRow()
    if selRow == nil {
        return
    }

    t.currDB = selRow[0]

    whereClause := t.whereClauses[t.currDB]
    cols, rows, err := db.GetDBTables(t.currDB, whereClause)

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

func (t *Tabs) UpdateMainTable(db dbms.DBMS) {
    selRow := t.Main.Panes.DbTables.Table.GetSelectedRow()
    if selRow == nil {
        t.Main.Panes.Main.Table.SetColumns(nil)
        t.Main.Panes.Main.Table.SetRows(nil)

        t.Main.SetError(errors.New("No tables in database"))
        return
    }

    t.currDBTable = selRow[0]

    whereClause := t.whereClauses[t.currDB + "/" + t.currDBTable]
    cols, rows, err := db.GetTable(t.currDB, t.currDBTable, whereClause)

    t.Main.SetError(err)

    if err != nil {
        t.Main.Panes.Main.Table.SetColumns(nil)
        t.Main.Panes.Main.Table.SetRows(nil)
        return
    }

    t.Main.Panes.Main.Table.SetColumns(cols)
    t.Main.Panes.Main.Table.SetRows(rows)
}

func (t *Tabs) UpdateDescribeTable(db dbms.DBMS) {

    cols, rows, err := db.GetDescribe(t.currDB, t.currDBTable)

    t.Main.SetError(err)

    if err != nil {
        t.Describe.Table.SetColumns(nil)
        t.Describe.Table.SetRows(nil)
        return
    }

    t.Describe.Table.SetColumns(cols)
    t.Describe.Table.SetRows(rows)
}

func (t *Tabs) DeleteSelectedDb(db dbms.DBMS) error {
    return db.DeleteDB(t.currDB)
}

func (t *Tabs) DeleteSelectedDbTable(db dbms.DBMS) error {
    return db.DeleteDBTable(t.currDB, t.currDBTable)
}

func (t *Tabs) DeleteSelectedRows(db dbms.DBMS) error {
    rows := t.Main.Panes.Main.Table.GetSelectedRows()

    for i := 0; i < len(rows); i++ {
        err := db.DeleteRow(
            t.currDB,
            t.currDBTable,
            rows[i],
            t.Main.Panes.Main.Table.GetCols(),
        )

        if err != nil {
            return err
        }
    }

    return nil
}

func (t *Tabs) UpdateSelectedCell(db dbms.DBMS, value string) error {
    selRow := t.Main.Panes.Main.Table.GetSelectedRow()
    if selRow == nil {
        return errors.New("Empty table")
    }

    return db.UpdateCell(
        t.currDB,
        t.currDBTable,
        selRow,
        t.Main.Panes.Main.Table.GetCols(),
        t.Main.Panes.Main.Table.GetCursor().X,
        value,
    )
}

func (t *Tabs) ChangeDbTableName(db dbms.DBMS, value string) error {
    return db.ChangeDbTableName(t.currDB, t.currDBTable, value)
}

func (t *Tabs) SendQuery(db dbms.DBMS, query string) {
    t.Main.SetError(db.SendQuery(query))
}
