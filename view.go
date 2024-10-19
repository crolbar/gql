package main

import (
	"fmt"
	"gql/table"
	"strings"

	"github.com/charmbracelet/lipgloss"
)


var style lipgloss.Style = lipgloss.NewStyle()

func (m model) View() string {
    if (m.requiresAuth()) {
        return m.auth.View()
    }

    if (m.db != nil) {
        return m.mainView()
    }

    return ""
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

func (m model) renderDbg() string {
    width  := m.width
    height := perc(20, m.height)

    dbg := fmt.Sprintf(
        "Height: %d, Width: %d, yOff: %d, xOff: %d, fullHeight: %d, fullWidth: %d, dbg: %d",
        m.panes.GetSelected().Table.GetHeight(),
        m.panes.GetSelected().Table.GetWidth(),
        m.panes.GetSelected().Table.GetYOffset(),
        m.panes.GetSelected().Table.GetXOffset(),
        m.height,
        m.width,
        perc(80, m.height),
    )


    if (len(dbg) > width) {
        //dbg = ""
    } else {
        dbg = lipgloss.JoinHorizontal( lipgloss.Top,
            strings.Repeat(" ", width - lipgloss.Width(dbg)),
            dbg,
        )

        dbg = strings.Repeat("\n", max(height - (1 + 5), 0)) + dbg
    }

    return dbg
}

func (m model) renderTopInfo() string {
    dbName     := m.panes.GetCurrDB()
    userName   := m.user

    columnsNum := fmt.Sprintf("%d", len(m.panes.GetSelected().Table.GetColumns()))
    rowsNum    := fmt.Sprintf("%d", len(m.panes.GetSelected().Table.GetRows()))

    selCol     := fmt.Sprintf("%d", m.panes.GetSelected().Table.GetCursor().X)
    selRow     := fmt.Sprintf("%d", m.panes.GetSelected().Table.GetCursor().Y)

    info := table.New(
        []table.Column {
            { Title: "Columns",
                Width: max(len("Columns"), len(columnsNum) + len(selCol) + 3), },
            { Title: "Rows",
                Width: max(len("Rows"), len(rowsNum) + len(selRow) + 3), },
            { Title: "Current Database",
                Width: max(len("Current Database"), len(dbName)), },
            { Title: "User",
                Width: max(len("User"), len(userName)), },
        },
        []table.Row {
            {
                selCol + " : " + columnsNum,
                selRow + " : " + rowsNum,
                dbName,
                userName,
            },
        },
        100, 100,
    )

    info.SetDisplayOnly()
    info.UpdateOffset()

    return info.View()
}

func (m model) renederTop() string {
    horizontal := lipgloss.JoinHorizontal(lipgloss.Left,
        m.renderTopInfo(),
    )

    full := lipgloss.JoinVertical(lipgloss.Left,
        horizontal,
        m.renderDbg(),
    )

    width  := m.width
    height := perc(20, m.height)

    return style.
        Width(width).
        Height(height).
        Render(full)
}

func (m model) mainView() string {
    full := lipgloss.JoinVertical(lipgloss.Top,
        m.renederTop(),
        m.renderTables(),
    )

    return full
}
