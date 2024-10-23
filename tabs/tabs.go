package tabs

import (
	"database/sql"
	"gql/tabs/describe_tab"
	"gql/tabs/main_tab"
	"gql/tabs/main_tab/panes"
	"gql/tabs/main_tab/panes/dialog_pane"

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

func New(
    DbPane       panes.Pane,
    DbTablesPane panes.Pane,
    MainPane     panes.Pane,
) Tabs {
    return Tabs {
        selected: Main,

        Main: main_tab.New(
            DbPane,
            DbTablesPane,
            MainPane,
        ),
        Describe: describe_tab.New(),

        keyMap: defaultKeyMap(),

        currDB:      "",
        currDBTable: "",
    }
}

type RequireMainTableUpdateMsg struct{}
func RequireMainTableUpdate() tea.Msg {
    return RequireMainTableUpdateMsg{}
}

type RequireDBTablesUpdateMsg struct{}
func RequireDBTablesUpdate() tea.Msg {
    return RequireDBTablesUpdateMsg{}
}

type DeleteSelectedDBMsg struct{}
func DeleteSelectedDB() tea.Msg {
    return DeleteSelectedDBMsg{}
}

type DeleteSelectedRowMsg struct{}
func DeleteSelectedRow() tea.Msg {
    return DeleteSelectedRowMsg{}
}

type UpdateSelectedCellMsg struct {}
func UpdateSelectedCell() tea.Msg {
    return UpdateSelectedCellMsg{}
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
    case RequireDBTablesUpdateMsg:
        t.UpdateDBTablesTable(db)
    case RequireMainTableUpdateMsg:
        t.UpdateMainTable(db)

    case dialog_pane.CancelMsg:
        t.Main.Panes.DeSelectDialog()

    case dialog_pane.RequestConfirmationMsg:
        t.Main.Panes.SelectDialog()

        t.Main.Panes.Dialog.SetupConfirmation(
            msg.Cmd,
            t.generateDialogHelpMsg(msg.Cmd()),
        )

        t.Main.Panes.Dialog.OnWindowResize(
            t.Main.GetHight(),
            t.Main.GetWidth(),
            t.Main.Panes.Db.Table.GetWidth(),
            t.Main.Panes.DbTables.Table.GetWidth(),
            t.Main.Panes.Main.Table.GetWidth(),
        )
    case dialog_pane.RequestValueUpdateMsg:
        t.Main.Panes.SelectDialog()

        t.Main.Panes.Dialog.SetupValueUpdate(
            msg.Cmd,
            t.generateDialogHelpMsg(msg.Cmd()),
            t.getSelectedValue(msg.Cmd()),
        )

        t.Main.Panes.Dialog.OnWindowResize(
            t.Main.GetHight(),
            t.Main.GetWidth(),
            t.Main.Panes.Db.Table.GetWidth(),
            t.Main.Panes.DbTables.Table.GetWidth(),
            t.Main.Panes.Main.Table.GetWidth(),
        )


    case DeleteSelectedDBMsg:
        t.DeleteSelectedDb(db)
        t.Main.Panes.DeSelectDialog()
        t.UpdateDBTable(db)

    case DeleteSelectedRowMsg:
        t.DeleteSelectedRow(db)
        t.Main.Panes.DeSelectDialog()
        t.UpdateMainTable(db)

    case dialog_pane.AcceptValueUpdateMsg:
        switch msg.Cmd().(type) {
        case UpdateSelectedCellMsg:
            t.Main.Panes.DeSelectDialog()
            t.UpdateSelectedCell(db, msg.Value)
            t.UpdateMainTable(db)
        }


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

func (t Tabs) generateDialogHelpMsg(msg tea.Msg) string {
    switch msg.(type) {
    case DeleteSelectedDBMsg:
        return "Are you sure you want to delete database " + t.currDB
    case DeleteSelectedRowMsg:
        return "Are you sure you want to delete this row"
    case UpdateSelectedCellMsg:
        return "Set new value for the selected cell"
    }
    return ""
}

func (t Tabs) getSelectedValue(msg tea.Msg) string {
    switch msg.(type) {
    case UpdateSelectedCellMsg:
        return t.Main.Panes.Main.Table.GetSelectedCell()
    }
    return ""
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
