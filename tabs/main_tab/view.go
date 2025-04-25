package main_tab

import (
	"github.com/charmbracelet/lipgloss"
)

var style lipgloss.Style = lipgloss.NewStyle()

func (t MainTab) renderLeftTable() string {
	dbTable := t.Panes.Db.Table
	dbTablesTable := t.Panes.DbTables.Table

	dbTables := style.Width(dbTablesTable.GetWidth()).
		Height(dbTablesTable.GetHeight()).
		MaxHeight(dbTablesTable.GetHeight()).
		Render(dbTablesTable.View())

	if !t.Panes.ShouldShowDB() {
		return dbTables
	}

	db := style.Width(dbTable.GetWidth()).
		Height(dbTable.GetHeight()).
		MaxHeight(dbTable.GetHeight()).
		Render(dbTable.View())

	return lipgloss.JoinHorizontal(lipgloss.Top, db, dbTables)
}

func (t MainTab) generateRightBoard(
	borderColor,
	headerTxt string,
	viewfn interface{},
	width,
	height int,
) string {

	style := style.
		Align(lipgloss.Left).
		Width(width)

	border := lipgloss.NormalBorder()

	header := style.
		Align(lipgloss.Center).
		Bold(true).
		Border(border).
		BorderBottom(false).
		BorderForeground(lipgloss.Color(borderColor)).
		Render(headerTxt)

	border.TopLeft = "├"
	border.TopRight = "┤"

	viewTxt := ""

	switch viewfn.(type) {
	case func(lipgloss.Style) string:
		viewTxt = viewfn.(func(lipgloss.Style) string)(style)
	case func() string:
		viewTxt = viewfn.(func() string)()
	}

	view := style.
		Height(height - lipgloss.Height(header)).
		MaxHeight((height + 2) - lipgloss.Height(header)).
		Border(border).
		BorderForeground(lipgloss.Color(borderColor)).
		Render(viewTxt)

	return header + "\n" + view

}

func (t MainTab) getRightSize() (int, int) {
	dbWidth := 0
	if t.Panes.ShouldShowDB() {
		dbWidth = t.Panes.Db.Table.GetWidth()
	}

	dbTablesWidth := t.Panes.DbTables.Table.GetWidth()
	mainWidth := t.Panes.Main.Table.GetWidth()
	tablesWidth := dbWidth + dbTablesWidth + mainWidth
	width := t.width - tablesWidth

	height := t.height

	// down by one on even to match the max main table height
	if height&1 == 0 {
		height--
	}

	return width, height
}

func (t MainTab) renderDialog() string {
	width, height := t.getRightSize()

	return t.generateRightBoard("255",
		t.Panes.Dialog.GetHelpMsg(),
		t.Panes.Dialog.TextInputView,
		width-2, height-2,
	)
}

func (t MainTab) renderRight() string {
	if t.Panes.IsDialogSelected() {
		return t.renderDialog()
	}

	width, height := t.getRightSize()
	width -= 2

	err := style.
		Width(width).
		Render(t.GetErrorStr())

	renderOnlyError := false

	if t.HasError() {
		if lipgloss.Height(err)+2+2 > height/2 {
			renderOnlyError = true
		} else {
			height -= lipgloss.Height(err) + 2 + 2 // make space for the error + header + borders
		}
	}

	selCellHeight := height / 2
	selColHeight := height / 2

	if height&1 > 0 {
		selCellHeight++
	}

	full := ""

	if !renderOnlyError {
		full = lipgloss.JoinVertical(lipgloss.Left,
			t.generateRightBoard("240",
				"Selected Cell",
				t.Panes.GetSelectedTable().GetSelectedCell,
				width, selCellHeight-2,
			),
			t.generateRightBoard("240",
				"Selected Column",
				func() string {
					return t.Panes.GetSelectedTable().GetSelColumnName()
				},
				width, selColHeight-2,
			),
		)
	}

	if t.HasError() {
		full = lipgloss.JoinVertical(lipgloss.Left,
			full,
			t.generateRightBoard("1",
				"Error",
				func() string { return err },
				width, min(lipgloss.Height(err)+2, height), // + 2 for the header
			),
		)
	}

	return full
}

func (t MainTab) RenderTables() string {
	mainTable := t.Panes.Main.Table

	mid := style.
		Height(mainTable.GetHeight()).
		MaxHeight(mainTable.GetHeight()).
		Width(mainTable.GetWidth()).
		Render(mainTable.View())

	return lipgloss.JoinHorizontal(lipgloss.Top,
		t.renderLeftTable(),
		mid,
		t.renderRight(),
	)
}
