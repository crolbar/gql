package tabs

import (
	"database/sql"
	"gql/tabs/describe_tab"
	"gql/tabs/main_tab"
	"gql/tabs/main_tab/panes"

	"github.com/charmbracelet/lipgloss"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type KeyMap struct {
    SelectMain     key.Binding
    SelectDescribe key.Binding
}

func defaultKeyMap() KeyMap {
    return KeyMap { 
        SelectMain:     key.NewBinding(key.WithKeys("1")),
        SelectDescribe: key.NewBinding(key.WithKeys("2")),
    }
}

type tabType int
const (
    Main tabType = iota
    Describe
)

type Tabs struct {
    selected tabType

    Main     main_tab.MainTab
    Describe describe_tab.DescribeTab

    keyMap KeyMap

    currDB      string
    currDBTable string
}

func New() Tabs {
    return Tabs {
        selected: Main,

        Main:     main_tab.New(),
        Describe: describe_tab.New(),

        keyMap: defaultKeyMap(),

        currDB:      "",
        currDBTable: "",
    }
}

func (t Tabs) Update(db *sql.DB, msg tea.Msg) (Tabs, tea.Cmd) {
    var cmd tea.Cmd

    switch t.selected {
    case Main:
        t.Main.Panes, cmd = t.Main.Panes.Update(msg)
    case Describe:
        t.Describe, cmd = t.Describe.Update(msg)
    }

    switch msg := msg.(type) {
    case panes.RequireDBTablesUpdateMsg:
        t.UpdateDBTablesTable(db)
    case panes.RequireMainTableUpdateMsg:
        t.UpdateMainTable(db)

    case tea.KeyMsg:
        switch {
        case key.Matches(msg, t.keyMap.SelectMain):
            t.selected = Main

        case key.Matches(msg, t.keyMap.SelectDescribe):
            t.selected = Describe
            t.UpdateDescribeTable(db)
        }
    }


    return t, cmd
}

func (t Tabs) SelectedTabView() string {
    switch t.selected {
    case Main:
        return t.Main.RenderTables()
    case Describe:
        return t.Describe.View()
    }

    return ""
}

func (t Tabs) HelpView() string {
    switch t.selected {
    case Main:
        return lipgloss.JoinHorizontal(lipgloss.Right,
            t.Main.Panes.HelpView(),
        )
    }

    return ""
}

func (t *Tabs) OnWindowResize(msg tea.WindowSizeMsg, isConnected bool) {
    t.Main.OnWindowResize(msg, isConnected)
    t.Describe.OnWindowResize(msg, isConnected)
}

func (t *Tabs) GetCurrDB() string {
    return t.currDB
}
