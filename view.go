package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)


func (m model) View() string {
    if (m.requiresAuth()) {
        return m.auth.View()
    }

    return m.mainView()
}

func (m model) renderLeftTable() string {
    s := lipgloss.NewStyle()
    dbTablesTable := m.panes.DbTables.Table
    dbTable := m.panes.Db.Table

    dbTables := s.Width(dbTablesTable.GetWidth()).
            Render(dbTablesTable.View())

    if (!m.panes.IsDbSelected()) {
        return dbTables
    }

    db := s.Width(dbTable.GetWidth()).
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
    mainWidth := m.panes.Main.Table.GetWidth()

    tablesWidth := dbWidth + dbTablesWidth + mainWidth

    width := (m.width - tablesWidth) - (3 + 1 + 1 + 3)

    style := lipgloss.NewStyle().
        Align(lipgloss.Left).
        MarginLeft(3).
        Width(width)

    header := style.Align(lipgloss.Center).
        Render("Selected Cell")

    view := style.
        Height(perc(80, m.height)).
        MaxHeight(perc(80, m.height) + 2).
        Border(lipgloss.NormalBorder()).
        BorderForeground(lipgloss.Color("240")).
        Render(selectedCell)

    return header + "\n" + view
}

func (m model) mainView() string {
    leftTable     := m.renderLeftTable()
    mainTable     := m.panes.Main.Table.View()

    dbg := fmt.Sprintf(
        "Height: %d, Width: %d, yOff: %d, xOff: %d, cursor: %d",
        m.panes.GetSelected().Table.GetHeight(),
        m.panes.GetSelected().Table.GetWidth(),
        m.panes.GetSelected().Table.GetYOffset(),
        m.panes.GetSelected().Table.GetXOffset(),
        m.panes.GetSelected().Table.GetCursor(),
    )

    s := lipgloss.NewStyle()

    full := lipgloss.JoinHorizontal(lipgloss.Top, 
        s.Render(leftTable),
        s.Width(m.panes.Main.Table.GetWidth()).Render(mainTable),
        m.renderRight(),
    )


    return dbg + "\n" + full + "\n"
}
