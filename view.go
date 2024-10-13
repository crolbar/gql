package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
    dbTableView      := m.DBTable.View()
    dbTableTableView := m.DBTablesTable.View()
    mainTableView    := m.mainTable.View()

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
        baseStyle.Render(dbTableView),
        baseStyle.Render(dbTableTableView),
        baseStyle.Render(mainTableView),
    )


    return log + "\n" + full + "\n"// + m.table.HelpView() + "\n" 
}

