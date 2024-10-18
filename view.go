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
    mainWidth     := m.panes.Main.Table.GetWidth()

    tablesWidth := dbWidth + dbTablesWidth + mainWidth

    width := (m.width - tablesWidth) - (1 + 1)

    style := lipgloss.NewStyle().
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

    height := perc(80, m.height)

    border.TopLeft  = "├"
    border.TopRight = "┤"

    if height & 1 == 0 {
        height++
    }
    view := style.
        Height(height - lipgloss.Height(header)).
        MaxHeight((height + 2) - lipgloss.Height(header)).
        Border(border).
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
