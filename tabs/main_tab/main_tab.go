package main_tab

import (
	"gql/tabs/main_tab/panes"
	tea "github.com/charmbracelet/bubbletea"
)

type MainTab struct {
    Panes panes.Panes

    width  int
    height int
}

func New(
    Db       panes.Pane,
    DbTables panes.Pane,
    Main     panes.Pane,
) MainTab {
    return MainTab{ 
        Panes: panes.New(
            panes.WithDBPane(Db),
            panes.WithDBTablesPane(DbTables),
            panes.WithMainPane(Main),
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
