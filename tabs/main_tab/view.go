package main_tab

import (
	"github.com/charmbracelet/lipgloss"
)

var style lipgloss.Style = lipgloss.NewStyle()

func (t MainTab) renderLeftTable() string {
    dbTable       := t.Panes.Db.Table
    dbTablesTable := t.Panes.DbTables.Table

    dbTables := style.Width(dbTablesTable.GetWidth()).
        Height(dbTablesTable.GetHeight()).
        MaxHeight(dbTablesTable.GetHeight()).
        Render(dbTablesTable.View())

    if (!t.Panes.ShouldShowDB()) {
        return dbTables
    }

    db := style.Width(dbTable.GetWidth()).
        Height(dbTable.GetHeight()).
        MaxHeight(dbTable.GetHeight()).
        Render(dbTable.View())

    return lipgloss.JoinHorizontal(lipgloss.Top, db, dbTables)
}

func (t MainTab) renderDialog() string {
    dbWidth := 0
    if (t.Panes.ShouldShowDB()) {
        dbWidth = t.Panes.Db.Table.GetWidth()
    }

    dbTablesWidth := t.Panes.DbTables.Table.GetWidth()
    mainWidth     := t.Panes.Main.Table.GetWidth()

    tablesWidth := dbWidth + dbTablesWidth + mainWidth

    width := (t.width - tablesWidth) - (1 + 1)

    style := style.
        Align(lipgloss.Left).
        Width(width)

    border := lipgloss.NormalBorder()

    header := style.
        Align(lipgloss.Center).
        Bold(true).
        Border(border).
        BorderBottom(false).
        Render("Dialog")

    height := t.height - 2

    // down by one on even to match the max main table height
    if height & 1 == 0 {
        height--;
    }

    border.TopLeft  = "├"
    border.TopRight = "┤"

    view := style.
        Height(height - lipgloss.Height(header)).
        MaxHeight((height + 2) - lipgloss.Height(header)).
        Border(border).
        Render(t.Panes.Dialog.TextInputView())

    return header + "\n" + view
}

func (t MainTab) renderRight() string {
    if t.Panes.IsDialogSelected() {
        return t.renderDialog()
    }

    selectedCell := t.Panes.GetSelectedTable().GetSelectedCell()

    dbWidth := 0
    if (t.Panes.ShouldShowDB()) {
        dbWidth = t.Panes.Db.Table.GetWidth()
    }

    dbTablesWidth := t.Panes.DbTables.Table.GetWidth()
    mainWidth     := t.Panes.Main.Table.GetWidth()

    tablesWidth := dbWidth + dbTablesWidth + mainWidth

    width := (t.width - tablesWidth) - (1 + 1)

    style := style.
        Align(lipgloss.Left).
        Width(width)

    border := lipgloss.NormalBorder()

    header := style.
        Align(lipgloss.Center).
        Bold(true).
        Border(border).
        BorderBottom(false).
        BorderForeground(lipgloss.Color("240")).
        Render("Selected Cell")

    height := t.height - 2

    // down by one on even to match the max main table height
    if height & 1 == 0 {
        height--;
    }

    border.TopLeft  = "├"
    border.TopRight = "┤"

    view := style.
        Height(height - lipgloss.Height(header)).
        MaxHeight((height + 2) - lipgloss.Height(header)).
        Border(border).
        BorderForeground(lipgloss.Color("240")).
        Render(selectedCell)

    return header + "\n" + view
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
