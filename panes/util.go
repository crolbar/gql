package panes

func (p *Panes) GetSelected() Pane {
    switch (p.selected) {
    case DB:
        return p.Db
    case DBTables:
        return p.DbTables
    case Main:
        return p.Main
    }

    panic("No pane for the table ?")
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

func (p *Panes) IsDbSelected() bool {
    return p.selected == DB
}
