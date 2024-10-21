package panes

import (
	"gql/table"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Pane struct {
    Table     table.Table
    KeyMap    help.KeyMap
    Help      help.Model
    updatef   func(m Panes, msg tea.Msg) (Panes, tea.Cmd)
    helpViewf func(m Panes) string
}

func NewPane(
    table     table.Table,
    keyMap    help.KeyMap,
    updatef   func(m Panes, msg tea.Msg) (Panes, tea.Cmd),
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
}

type Opts func(*Panes)

func New(opts ...Opts) Panes {
    p := Panes {
        selected: DB,
    }

    for _, opt := range opts {
		opt(&p)
	}

    return p
}

type RequireMainTableUpdateMsg struct{}
func RequireMainTableUpdate() tea.Msg {
    return RequireMainTableUpdateMsg{}
}

type RequireDBTablesUpdateMsg struct{}
func RequireDBTablesUpdate() tea.Msg {
    return RequireDBTablesUpdateMsg{}
}

func (p Panes) Update(msg tea.Msg) (Panes, tea.Cmd) {
    switch (p.selected) {
    case DB:
        return p.Db.updatef(p, msg)
    case DBTables:
        return p.DbTables.updatef(p, msg)
    case Main:
        return p.Main.updatef(p, msg)
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
