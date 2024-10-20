package panes

import (
	"database/sql"
	"gql/table"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Pane struct {
    Table     table.Table
    KeyMap    help.KeyMap
    Help      help.Model
    updatef   func(m Panes, db *sql.DB, msg tea.Msg) (Panes, tea.Cmd)
    helpViewf func(m Panes) string
}

func NewPane(
    table     table.Table,
    keyMap    help.KeyMap,
    updatef   func(m Panes, db *sql.DB, msg tea.Msg) (Panes, tea.Cmd),
    helpViewf func(m Panes) string,
) Pane {
    help                       := help.New()
    help.ShowAll               = true
    help.Styles.FullKey        = lipgloss.NewStyle().Bold(true)
    help.Styles.FullDesc       = lipgloss.NewStyle().Italic(true)
    help.Styles.FullSeparator  = lipgloss.NewStyle()
    help.FullSeparator         = ""

    return Pane {
        Table:     table,
        KeyMap:    keyMap,
        Help:      help,
        updatef:   updatef,
        helpViewf: helpViewf,
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

func (p Panes) HelpView() string {
    switch (p.selected) {
    case DB:
        return p.Db.helpViewf(p)
    case DBTables:
        return p.DbTables.helpViewf(p)
    case Main:
        return p.Main.helpViewf(p)
    }

    panic("No help view for the pane ?")
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
