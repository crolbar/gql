package main_tab

import (
	"gql/tabs/main_tab/panes/db_pane"
	"gql/tabs/main_tab/panes/db_tables_pane"
	"gql/tabs/main_tab/panes/main_pane"
	"gql/tabs/main_tab/panes"
	tea "github.com/charmbracelet/bubbletea"
)

type MainTab struct {
    Panes panes.Panes

    width  int
    height int
}

func New() MainTab {
    return MainTab{ 
        Panes: panes.New(
            panes.WithDBPane(db_pane.New()),
            panes.WithDBTablesPane(db_tables_pane.New()),
            panes.WithMainPane(main_pane.New()),
        ),
    }
}

func (t *MainTab) OnWindowResize(msg tea.WindowSizeMsg, isConnected bool) {
    width  := msg.Width
    height := msg.Height

    height = perc(80, height)

    t.SetWidth(msg.Width)
    t.SetHeight(height)

    t.Panes.Main.Table.SetMaxSize(perc(60, width), height)

    // only height because we are using only one column
    // and are setting the max width in from UpdateDBTablesTable(): db_util
    t.Panes.Db.Table.SetMaxHeight(height)
    t.Panes.DbTables.Table.SetMaxHeight(height)

    if isConnected {
        t.Panes.Db.Table.UpdateOffset()
        t.Panes.DbTables.Table.UpdateOffset()
        t.Panes.Main.Table.UpdateOffset()
    }
}

func (t *MainTab) SetHeight(height int) {
    t.height = height
}

func (t *MainTab) SetWidth(width int) {
    t.width = width
}

func perc(per, num int) int {
    return int(float32(num) * (float32(per) / float32(100)))
}
