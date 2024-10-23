package panes

import "gql/table"

func (p *Panes) getSelectedPane(paneType paneType) *Pane {
    switch (paneType) {
    case DB:
        return &p.Db
    case DBTables:
        return &p.DbTables
    case Main:
        return &p.Main
    }


    panic("No pane for the table ?")
}

func (p *Panes) GetSelectedTable() *table.Table {
    switch (p.selected) {
    case Dialog:
        return &p.getSelectedPane(p.prev).Table
    default:
        return &p.getSelectedPane(p.selected).Table
    }
}

func (p *Panes) SelectDB() {
    p.selected = DB

    p.Db.Table.Focus()
    p.DbTables.Table.DeFocus()
}

func (p *Panes) SelectDBTables() {
    p.selected = DBTables

    p.Db.Table.DeFocus()
    p.Main.Table.DeFocus()
    p.DbTables.Table.Focus()
}

func (p *Panes) SelectMain() {
    p.selected = Main

    p.DbTables.Table.DeFocus()
    p.Main.Table.Focus()
}

func (p *Panes) SelectDialog() {
    p.prev     = p.selected
    p.selected = Dialog

    p.DbTables.Table.DeFocus()
    p.Main.Table.DeFocus()
    p.Db.Table.DeFocus()
}

func (p *Panes) ShouldShowDB() bool {
    return p.selected == DB || p.prev == DB
}

func (p *Panes) IsDialogSelected() bool {
    return p.selected == Dialog
}

func (p *Panes) DeSelectDialog() {
    switch p.prev {
    case DB:
        p.SelectDB()
    case DBTables:
        p.SelectDBTables()
    case Main:
        p.SelectMain()
    }

    p.prev = Dialog
}
