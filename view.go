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
    if (m.selectedPane != DB) {
        return m.dbTablesPane.table.View()
    } else {
        return lipgloss.JoinHorizontal(lipgloss.Top,
            m.dbPane.table.View(),
            m.dbTablesPane.table.View(),
        )
    }
}

func (m model) renderRight() string {
    var width int

    if (m.selectedPane != DB) {
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
        Render(m.getSelectedPane().table.GetSelectedCell())

    return header + "\n" + view
}

func (m model) mainView() string {
    leftTable     := m.renderLeftTable()
    mainTable     := m.mainPane.table.View()

    log := fmt.Sprintf(
        "Height: %d, Width: %d, yOff: %d, xOff: %d, cursor: %d, dbg: %s",
        m.mainPane.table.Height,
        m.mainPane.table.Width,
        m.mainPane.table.YOffset,
        m.mainPane.table.XOffset,
        m.mainPane.table.Cursor,
        m.mainPane.table.Dbg,
    )

    s := lipgloss.NewStyle()

    full := lipgloss.JoinHorizontal(lipgloss.Top, 
        s.Render(leftTable),
        s.Width(m.mainPane.table.Width).Render(mainTable),
        m.renderRight(),
    )


    return log + "\n" + full + "\n"// + m.table.HelpView() + "\n" 
}
