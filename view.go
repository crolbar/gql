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
    if (!m.panes.IsDbSelected()) {
        return m.panes.DbTables.Table.View()
    }

    return lipgloss.JoinHorizontal(lipgloss.Top,
        m.panes.Db.Table.View(),
        m.panes.DbTables.Table.View(),
    )
}

func (m model) renderRight() string {
    var width int

    if (!m.panes.IsDbSelected()) {
        width = 40
    } else {
        width = 20
    }

    style := lipgloss.NewStyle().
        Align(lipgloss.Left).
        MarginLeft(3).
        Width(width)

    header := style.Align(lipgloss.Center).
        Render("Selected Cell")

    view := style.
        Border(lipgloss.NormalBorder()).
        BorderForeground(lipgloss.Color("240")).
        Render(m.panes.GetSelected().Table.GetSelectedCell())

    return header + "\n" + view
}

func (m model) mainView() string {
    leftTable     := m.renderLeftTable()
    mainTable     := m.panes.Main.Table.View()

    log := fmt.Sprintf(
        "Height: %d, Width: %d, yOff: %d, xOff: %d, cursor: %d",
        m.panes.Main.Table.GetHeight(),
        m.panes.Main.Table.GetWidth(),
        m.panes.Main.Table.GetYOffset(),
        m.panes.Main.Table.GetXOffset(),
        m.panes.Main.Table.GetCursor(),
    )

    s := lipgloss.NewStyle()

    full := lipgloss.JoinHorizontal(lipgloss.Top, 
        s.Render(leftTable),
        s.Width(m.panes.Main.Table.GetWidth()).Render(mainTable),
        m.renderRight(),
    )


    return log + "\n" + full + "\n"
}
