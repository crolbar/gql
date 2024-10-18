package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)


var style lipgloss.Style = lipgloss.NewStyle()

func (m model) View() string {
    if (m.requiresAuth()) {
        return m.auth.View()
    }

    return m.mainView()
}

func (m model) renderLeftTable() string {
    dbTable       := m.panes.Db.Table
    dbTablesTable := m.panes.DbTables.Table

    dbTables := style.Width(dbTablesTable.GetWidth()).
        Height(dbTablesTable.GetHeight()).
        MaxHeight(dbTablesTable.GetHeight()).
        Render(dbTablesTable.View())

    if (!m.panes.IsDbSelected()) {
        return dbTables
    }

    db := style.Width(dbTable.GetWidth()).
        Height(dbTable.GetHeight()).
        MaxHeight(dbTable.GetHeight()).
        Render(dbTable.View())

    return lipgloss.JoinHorizontal(lipgloss.Top, db, dbTables)
}

func (m model) renderRight() string {
    selectedCell := m.panes.GetSelected().Table.GetSelectedCell()

    dbWidth := 0
    if (m.panes.IsDbSelected()) {
        dbWidth = m.panes.Db.Table.GetWidth()
    }

    dbTablesWidth := m.panes.DbTables.Table.GetWidth()
    mainWidth     := m.panes.Main.Table.GetWidth()

    tablesWidth := dbWidth + dbTablesWidth + mainWidth

    width := (m.width - tablesWidth) - (1 + 1)

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

    height := perc(80, m.height) - 2

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

func (m model) renderTables() string {
    mainTable := m.panes.Main.Table.View()

    mid := style.
        Height(m.panes.Main.Table.GetHeight()).
        MaxHeight(m.panes.Main.Table.GetHeight()).
        Width(m.panes.Main.Table.GetWidth()).
        Render(mainTable)

    return lipgloss.JoinHorizontal(lipgloss.Top, 
        m.renderLeftTable(),
        mid,
        m.renderRight(),
    )
}

func (m model) renederTop() string {
    dbg := fmt.Sprintf(
        "Height: %d, Width: %d, yOff: %d, xOff: %d, cursor: %d, fullHeight: %d, fullWidth: %d, dbg: %d",
        m.panes.GetSelected().Table.GetHeight(),
        m.panes.GetSelected().Table.GetWidth(),
        m.panes.GetSelected().Table.GetYOffset(),
        m.panes.GetSelected().Table.GetXOffset(),
        m.panes.GetSelected().Table.GetCursor(),
        m.height,
        m.width,
        perc(80, m.height),
    )

    height := perc(20, m.height) - 2

    top := style.
        Width(m.width - 2).
        Height(height).
        Border(lipgloss.NormalBorder()).
        BorderForeground(lipgloss.Color("240")).
        Render(dbg)

    return top
}

func (m model) mainView() string {
    full := lipgloss.JoinVertical(lipgloss.Top,
        m.renederTop(),
        m.renderTables(),
    )

    return full
}
