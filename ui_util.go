package main

type Pane int
const (
    DB Pane = iota
    DBTables
    Main
)

func (m *model) getSelectedPane() pane {
    switch (m.selectedPane) {
    case DB:
        return m.dbPane
    case DBTables:
        return m.dbTablesPane
    case Main:
        return m.mainPane
    }

    panic("No pane for the table ?")
}

func (m *model) selectDBPane() {
    m.selectedPane = DB

    m.dbPane.table.Focus()
    m.dbTablesPane.table.DeFocus()
}

func (m *model) selectDBTablesPane() {
    m.selectedPane = DBTables

    m.dbPane.table.DeFocus()
    m.mainPane.table.DeFocus()
    m.dbTablesPane.table.Focus()
}

func (m *model) selectMainPane() {
    m.selectedPane = Main

    m.dbTablesPane.table.DeFocus()
    m.mainPane.table.Focus()
}
