package tabs

import (
	"fmt"
	"gql/dbms"
	"gql/tabs/describe_tab"
	"gql/tabs/main_tab"
	"gql/tabs/main_tab/panes"
	"gql/tabs/main_tab/panes/cmd_pane"
	"gql/tabs/main_tab/panes/dialog_pane"
	"gql/tabs/main_tab/panes/filter_pane"

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

    whereClauses map[string]string
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

        whereClauses: make(map[string]string),
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

type DeleteSelectedDBTableMsg struct{}
func DeleteSelectedDBTable() tea.Msg {
    return DeleteSelectedDBTableMsg{}
}

type DeleteSelectedRowMsg struct{}
func DeleteSelectedRow() tea.Msg {
    return DeleteSelectedRowMsg{}
}

type UpdateSelectedCellMsg struct {}
func UpdateSelectedCell() tea.Msg {
    return UpdateSelectedCellMsg{}
}

type ChangeDbTableNameMsg struct {}
func ChangeDbTableName() tea.Msg {
    return ChangeDbTableNameMsg{}
}

type FocusFilterMsg struct {}
func FocusFilter() tea.Msg {
    return FocusFilterMsg{}
}

type FocusCmdMsg struct {}
func FocusCmd() tea.Msg {
    return FocusCmdMsg{}
}

func (t Tabs) Update(db dbms.DBMS, msg tea.Msg) (Tabs, tea.Cmd) {
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


    case FocusCmdMsg:
        t.Main.Panes.SelectCmd()

    case cmd_pane.CancelMsg:
        t.Main.Panes.DeSelectDialogFilterCmd()

    case cmd_pane.AcceptMsg:
        t.SendQuery(db, msg.Query)
        t.Main.Panes.DeSelectDialogFilterCmd()


    case FocusFilterMsg:
        t.Main.Panes.Filter.UpdateValue(t.GetWhereClause())
        t.Main.Panes.Filter.UpdatePrefix(t.GetWhereClausePrefix())
        t.Main.Panes.SelectFilter()
        t.Main.Panes.Filter.SetWidth(t.Main.GetWidth())

    case filter_pane.AcceptMsg:
        t.Main.Panes.DeSelectDialogFilterCmd()
        t.UpdateCurrentWhereClause(db, msg.Txt) // after DeSelect !!!

    case filter_pane.CancelMsg:
        t.Main.Panes.DeSelectDialogFilterCmd()
        t.UpdateCurrentWhereClause(db, "")      // after DeSelect !!!


    case dialog_pane.CancelMsg:
        t.Main.Panes.DeSelectDialogFilterCmd()
        t.UpdateMainTable(db)

    case dialog_pane.RequestConfirmationMsg:
        t.Main.Panes.SelectDialog()

        t.Main.Panes.Dialog.SetupConfirmation(
            msg.Cmd,
            t.generateDialogHelpMsg(msg.Cmd()),
        )

        t.Main.Panes.Dialog.SetWidth(
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

        t.Main.Panes.Dialog.SetWidth(
            t.Main.GetWidth(),
            t.Main.Panes.Db.Table.GetWidth(),
            t.Main.Panes.DbTables.Table.GetWidth(),
            t.Main.Panes.Main.Table.GetWidth(),
        )


    case DeleteSelectedDBMsg:
        if t.handleError(
            t.DeleteSelectedDb(db),
        ) {
            t.UpdateDBTable(db)
        }

    case DeleteSelectedDBTableMsg:
        if t.handleError(
            t.DeleteSelectedDbTable(db),
        ) {
            t.UpdateDBTablesTable(db)
        }

    case DeleteSelectedRowMsg:
        if t.handleError(
            t.DeleteSelectedRows(db),
        ) {
            t.UpdateMainTable(db)
        }

    case dialog_pane.AcceptValueUpdateMsg:
        switch msg.Cmd().(type) {
        case UpdateSelectedCellMsg:
            if t.handleError(
                t.UpdateSelectedCell(db, msg.Value),
            ) {
                t.UpdateMainTable(db)
            }
        case ChangeDbTableNameMsg:
            if t.handleError(
                t.ChangeDbTableName(db, msg.Value),
            ) {
                t.UpdateDBTablesTable(db)
            }
        }



    case tea.KeyMsg:
        if !t.IsTyting() {
            switch {
            case key.Matches(msg, t.keyMap.SelectMain):
                t.selected = Main

            case key.Matches(msg, t.keyMap.SelectDescribe):
                t.selected = Describe
                t.UpdateDescribeTable(db)
            }
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
    case ChangeDbTableNameMsg:
        return t.Main.Panes.DbTables.Table.GetSelectedCell()
    }
    return ""
}

func (t *Tabs) handleError(err error) bool {
    if err != nil {
        if t.Main.Panes.IsDialogSelected() {
            t.Main.Panes.Dialog.SetError(err.Error())

            return false
        }

        return false
    }

    if t.Main.Panes.IsDialogSelected() {
        t.Main.Panes.DeSelectDialogFilterCmd()
        t.Main.Panes.Dialog.Reset()
    }

    return true
}

func (t *Tabs) UpdateCurrentWhereClause(db dbms.DBMS, value string) {
    switch t.Main.Panes.GetSelected() {
    case panes.Main:
        t.whereClauses[t.currDB + "/" + t.currDBTable] = value
        t.UpdateMainTable(db)

    case panes.DBTables:
        t.whereClauses[t.currDB] = value
        t.UpdateDBTablesTable(db)

    case panes.DB:
        t.whereClauses["db"] = value
        t.UpdateDBTable(db)
    }
}

func (t *Tabs) GetWhereClause() string {
    switch t.Main.Panes.GetSelectedOnlyTables() {
    case panes.Main:
        return t.whereClauses[t.currDB + "/" + t.currDBTable]
    case panes.DBTables:
        return t.whereClauses[t.currDB]
    case panes.DB:
        return t.whereClauses["db"]
    }

    return ""
}

func (t *Tabs) GetWhereClausePrefix() string {
    switch t.Main.Panes.GetSelectedOnlyTables() {
    case panes.DBTables:
        return fmt.Sprintf("Tables_in_%s = ", t.currDB)
    case panes.DB:
        return "Database = "
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

func (t *Tabs) IsTyting() bool {
    return t.Main.Panes.IsDialogSelected() ||
    t.Main.Panes.IsFilterSelected() ||
    t.Main.Panes.IsCmdSelected()
}
