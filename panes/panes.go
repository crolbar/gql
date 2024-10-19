package panes

import (
	"database/sql"
	"gql/table"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

type Pane struct {
    Table  table.Table
    KeyMap interface{}
    Help   help.Model
    updatef func(m Panes, db *sql.DB, msg tea.Msg) (Panes, tea.Cmd)
}

func NewPane(
    table table.Table,
    keyMap interface{},
    help help.Model,
    updatef func(m Panes, db *sql.DB, msg tea.Msg) (Panes, tea.Cmd),
) Pane {
    return Pane {
        Table: table,
        KeyMap: keyMap,
        Help: help,
        updatef: updatef,
    }
}

type paneType int
const (
    DB paneType = iota
    DBTables
    Main
)

type Panes struct {
    selected paneType

    Db       Pane
    DbTables Pane
    Main     Pane

    currDB      string
    currDBTable string
}

type Opts func(*Panes)

func New(opts ...Opts) Panes {
    p := Panes {
        selected: DB,
        currDB: "",
        currDBTable: "",
    }

    for _, opt := range opts {
		opt(&p)
	}

    return p
}

func (p Panes) Update(db *sql.DB, msg tea.Msg) (Panes, tea.Cmd) {
    switch (p.selected) {
    case DB:
        return p.Db.updatef(p, db, msg)
    case DBTables:
        return p.DbTables.updatef(p, db, msg)
    case Main:
        return p.Main.updatef(p, db, msg)
    }

    panic("No update for the pane ?")
}

func WithDBPane(dbPane Pane) Opts {
    return func(p *Panes) {
        p.Db = dbPane
    }
}

func WithDBTablesPane(dbTablesPane Pane) Opts {
    return func(p *Panes) {
        p.DbTables = dbTablesPane
    }
}

func WithMainPane(main Pane) Opts {
    return func(p *Panes) {
        p.Main = main
    }
}

func (p Panes) GetCurrDB() string {
    return p.currDB
}
