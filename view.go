package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)


func (m model) renderLeftTable() string {
    if (m.selectedPane != DB) {
        return m.DBTablesTable.View()
    } else {
        return lipgloss.JoinHorizontal(lipgloss.Top,
            m.DBTable.View(),
            m.DBTablesTable.View(),
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
        Render(m.selectedTable().GetSelectedCell())

    return header + "\n" + view
}

func (m model) View() string {
    leftTable     := m.renderLeftTable()
    mainTableView := m.mainTable.View()

    log := fmt.Sprintf(
        "Height: %d, Width: %d, yOff: %d, xOff: %d, cursor: %d, dbg: %s",
        m.mainTable.Height,
        m.mainTable.Width,
        m.mainTable.YOffset,
        m.mainTable.XOffset,
        m.mainTable.Cursor,
        m.mainTable.Dbg,
    )

    full := lipgloss.JoinHorizontal(lipgloss.Top, 
        baseStyle.Render(leftTable),
        baseStyle.Width(m.mainTable.Width).Render(mainTableView),
        m.renderRight(),
    )


    return log + "\n" + full + "\n"// + m.table.HelpView() + "\n" 
}
