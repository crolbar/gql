package main

import (
	"fmt"
	"gql/table"
	"gql/util"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

var style lipgloss.Style = lipgloss.NewStyle()

func (m model) View() string {
	if m.requiresAuth() {
		return m.auth.View()
	}

	if m.dbms.HasDb() {
		return m.mainView()
	}

	return style.
		Height(m.height).
		Width(m.width).
		AlignVertical(lipgloss.Center).
		AlignHorizontal(lipgloss.Center).
		Render("loading..")
}

func (m model) renderDbg() string {
	width := m.width

	dbg := fmt.Sprintf(
		"Height: %d, Width: %d, yOff: %d, xOff: %d, fullHeight: %d, fullWidth: %d, dbg: %d",
		m.tabs.Main.Panes.GetSelectedTable().GetHeight(),
		m.tabs.Main.Panes.GetSelectedTable().GetWidth(),
		m.tabs.Main.Panes.GetSelectedTable().GetYOffset(),
		m.tabs.Main.Panes.GetSelectedTable().GetXOffset(),
		m.height,
		m.width,
		perc(80, m.height),
	)

	if len(dbg) > width {
		//dbg = ""
	} else {
		//dbg = lipgloss.JoinHorizontal( lipgloss.Top,
		//    strings.Repeat(" ", width - lipgloss.Width(dbg)),
		//    dbg,
		//)
	}

	return dbg
}

func (km KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{}
}

func (km KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{km.ChangeCreds, km.Quit},
	}
}

func (m model) renderHelp(infoLen int) string {
	helpMsg := lipgloss.JoinHorizontal(lipgloss.Right,
		m.tabs.HelpView(),
		m.help.View(m.keyMap),
	)

	helpMsgSplit := strings.Split(helpMsg, "\n")

	width := util.MaxLine(helpMsg)

	// 2 for border
	if width+2 > m.width-infoLen {
		return ""
	}

	help := table.New(
		[]table.Column{
			{Title: "Help",
				Width: width},
		},
		[]table.Row{
			{helpMsgSplit[0]},
			{helpMsgSplit[1]},
		},
		100, 100,
	)
	help.SetDisplayOnly()

	return help.View()
}

func (m model) renderTopInfo() string {
	selectedTable := m.tabs.Main.Panes.GetSelectedTable()

	dbName := m.tabs.GetCurrDB()
	userName := m.user

	columnsNum := fmt.Sprintf("%d", len(selectedTable.GetColumns()))
	rowsNum := fmt.Sprintf("%d", len(selectedTable.GetRows()))

	selCol := fmt.Sprintf("%d", selectedTable.GetCursor().X)

	if selectedTable.IsSelectingCols() {
		selectionStart := selectedTable.GetSelectionStart()

		if selectedTable.GetCursor().X > selectionStart {
			selCol = fmt.Sprintf("%d-%s", selectionStart, selCol)
		} else {
			selCol = fmt.Sprintf("%s-%d", selCol, selectionStart)
		}
	}

	selRow := fmt.Sprintf("%d", selectedTable.GetCursor().Y)

	if selectedTable.IsSelectingRows() {
		selectionStart := selectedTable.GetSelectionStart()

		if selectedTable.GetCursor().Y > selectionStart {
			selRow = fmt.Sprintf("%d-%s", selectionStart, selRow)
		} else {
			selRow = fmt.Sprintf("%s-%d", selRow, selectionStart)
		}
	}

	info := table.New(
		[]table.Column{
			{Title: "Columns",
				Width: max(len("Columns"), len(columnsNum)+len(selCol)+3)},
			{Title: "Rows",
				Width: max(len("Rows"), len(rowsNum)+len(selRow)+3)},
			{Title: "Current Database",
				Width: max(len("Current Database"), len(dbName))},
			{Title: "User",
				Width: max(len("User"), len(userName))},
		},
		[]table.Row{
			{
				selCol + " : " + columnsNum,
				selRow + " : " + rowsNum,
				dbName,
				userName,
			},
		},
		100, 100,
	)

	info.SetDisplayOnly()
	info.UpdateOffset()

	return info.View()
}

func (m model) renederTop() string {
	topInfo := m.renderTopInfo()
	topHelp := m.renderHelp(lipgloss.Width(topInfo))

	InfoTabs := lipgloss.JoinVertical(lipgloss.Left,
		topInfo,
		m.tabs.View(),
	)

	whereClause := "> where " + m.tabs.GetWhereClausePrefix() + m.tabs.GetWhereClause()

	if m.tabs.Main.Panes.IsFilterSelected() {
		whereClause = m.tabs.Main.Panes.Filter.View()
	}

	var right string

	if !m.tabs.Main.Panes.IsCmdSelected() {
		right = lipgloss.JoinVertical(lipgloss.Left,
			topHelp,
			whereClause,
		)
	} else {
		right = lipgloss.JoinVertical(lipgloss.Left,
			topHelp,
			m.tabs.Main.Panes.Cmd.View(),
		)
	}

	full := lipgloss.JoinHorizontal(lipgloss.Left,
		InfoTabs,
		right,
	)

	width := m.width
	height := perc(20, m.height)

	return style.
		Width(width).
		Height(height).
		Render(full)
}

func (m model) mainView() string {
	full := lipgloss.JoinVertical(lipgloss.Top,
		m.renederTop(),
		m.tabs.SelectedTabView(),
	)

	return full
}
