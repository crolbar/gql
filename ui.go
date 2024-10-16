package main

import "gql/table"

type Pane int
const (
    DB Pane = iota
    DBTables
    Main
)

func (m *model) selectedTable() table.Table {
    switch (m.selectedPane) {
    case DB:
        return m.DBTable
    case DBTables:
        return m.DBTablesTable
    case Main:
        return m.mainTable
    }

    panic("No pane for the table ?")
}

func (m *model) selectDBpane() {
    m.selectedPane = DB

    m.DBTable.Focus()
    m.DBTablesTable.DeFocus()
}
func (m *model) selectDBTablespane() {
    m.selectedPane = DBTables

    m.DBTable.DeFocus()
    m.mainTable.DeFocus()
    m.DBTablesTable.Focus()
}
func (m *model) selectMainpane() {
    m.selectedPane = Main

    m.DBTablesTable.DeFocus()
    m.mainTable.Focus()
}
