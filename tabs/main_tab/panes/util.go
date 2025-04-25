package panes

import "gql/table"

func (p *Panes) getSelectedTable(paneType PaneType) *table.Table {
	switch paneType {
	case DB:
		return &p.Db.Table
	case DBTables:
		return &p.DbTables.Table
	case Main:
		return &p.Main.Table
	case Dialog, Filter, Cmd:
		return p.getSelectedTable(p.prev)
	}

	panic("No pane for the table ?")
}

func (p *Panes) GetSelectedTable() *table.Table {
	return p.getSelectedTable(p.selected)
}

func (p *Panes) GetSelected() PaneType {
	return p.selected
}

func (p *Panes) GetSelectedOnlyTables() PaneType {
	if p.IsDialogSelected() || p.IsFilterSelected() {
		return p.prev
	}

	return p.selected
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
	p.prev = p.selected
	p.selected = Dialog

	p.DbTables.Table.DeFocus()
	p.Main.Table.DeFocus()
	p.Db.Table.DeFocus()
}

func (p *Panes) SelectFilter() {
	p.prev = p.selected
	p.selected = Filter
	p.Filter.Focus()

	p.DbTables.Table.DeFocus()
	p.Main.Table.DeFocus()
	p.Db.Table.DeFocus()
}

func (p *Panes) SelectCmd() {
	p.prev = p.selected
	p.selected = Cmd
	p.Cmd.Focus()

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

func (p *Panes) IsFilterSelected() bool {
	return p.selected == Filter
}

func (p *Panes) IsCmdSelected() bool {
	return p.selected == Cmd
}

func (p *Panes) DeSelectDialogFilterCmd() {
	switch p.selected {
	case Filter:
		p.Filter.DeFocus()
	case Cmd:
		p.Cmd.DeFocus()
	}

	switch p.prev {
	case DB:
		p.SelectDB()
	case DBTables:
		p.SelectDBTables()
	case Main:
		p.SelectMain()
	}
}
