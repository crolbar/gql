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
        baseStyle.Render(mainTableView),
    )


    return log + "\n" + full + "\n"// + m.table.HelpView() + "\n" 
}
