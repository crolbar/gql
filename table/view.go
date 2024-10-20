package table

import (
    "gql/table/scrollbar"

    "github.com/charmbracelet/bubbles/key"
    "github.com/charmbracelet/lipgloss"
)


func (km KeyMap) ShortHelp() []key.Binding {
    return []key.Binding { km.LineUp, km.LineDown, km.LineLeft, km.LineRight }
}

func (km KeyMap) FullHelp() [][]key.Binding {
    return [][]key.Binding {
        {km.LineUp, km.LineDown},
        {km.LineLeft, km.LineRight},
        {km.HalfPageUp, km.HalfPageDown},
        {km.SelectColumn, km.SelectRow},
    }
}

func (t Table) HelpView() string {
    return t.help.View(t.keyMap);
}

func (t Table) View() string {
    return t.headersView() + "\n" + t.renderRows()
}

func (t Table) headersView() string {
    s := make([]string, 0, len(t.cols))

    start := clamp(t.xOffset, 0, len(t.cols))
    end   := clamp(t.renderedColumns + t.xOffset, 0, len(t.cols))

    for i := start; i < end; i++ {
        col := t.cols[i]
        if col.Width <= 0 {
            continue
        }

        value        := truncate(col.Title, t.cols[i].Width)
        style        := t.generateStyleHeader(i, end)
        renderedCell := style.Width(col.Width).Render(value)

        s = append(s, renderedCell)
    }

    return lipgloss.JoinHorizontal(lipgloss.Top, s...)
}

func (t *Table) renderRows() string {
    rows := make([]string, 0, len(t.rows))

    // 1 for the bottom border 2 for the header
    // and / 2 because one item takes up 2 rows because of the top border
    itemsInView := (t.height - 1 - 2) / 2

    start := clamp(t.yOffset, 0, len(t.rows))
    end   := clamp(t.yOffset + itemsInView, 0, len(t.rows))

    vScrollbar := scrollbar.New(itemsInView, len(t.rows), t.yOffset)

    for i := start; i < end; i++ {
        rows = append(rows, t.renderRow(i, end, &vScrollbar))
    }

    return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (t *Table) renderRow(r, rEnd int, vScrollbar *scrollbar.Scrollbar) string {
    s := make([]string, 0, len(t.cols))

    start := clamp(t.xOffset, 0, len(t.cols))
    end   := clamp(t.renderedColumns + t.xOffset, 0, len(t.cols))

    isScrollbarRow := vScrollbar.IsScrollbarItem(r)

    hScrollbar := scrollbar.New(t.renderedColumns, len(t.cols), t.xOffset)

    for i := start; i < end; i++ {
        if t.cols[i].Width <= 0 {
            continue
        }

        isScrollbarCol := hScrollbar.IsScrollbarItem(i)

        value        := truncate(t.rows[r][i], t.cols[i].Width)
        style        := t.generateStyleRow(i, end, r, rEnd, isScrollbarRow, isScrollbarCol)
        renderedCell := style.Width(t.cols[i].Width).Render(value)

        s = append(s, renderedCell)
    }

    return lipgloss.JoinHorizontal(lipgloss.Top, s...)
}

func (t Table) generateStyleHeader(colI, end int) lipgloss.Style {
    topLeftBorder       := iff(colI == t.xOffset, "┌", "┬")
    topRightBorder      := iff(colI == end - 1, "┐", "")
    enableRightBorder   := colI == end - 1


    return lipgloss.NewStyle().
    Border(lipgloss.Border{
        Top:      "─",
        Left:     "│",
        Right:    "│",
        TopLeft:  topLeftBorder,
        TopRight: topRightBorder,
    }). 
    BorderBottom(false).
    BorderRight(enableRightBorder).
    BorderForeground(lipgloss.Color("240")).
    Bold(true)
}

func (t Table) generateStyleRow(
    colI,
    cEnd,
    rowI,
    rEnd int,
    isScrollbarRow,
    isScrollbarCol bool,
) lipgloss.Style {
    enableRightBorder   := colI == cEnd - 1
    enableBottomBorder  := rowI == rEnd - 1

    topLeftBorder    := iff(colI == t.xOffset, "├", "┼")
    topRightBorder   := iff(colI == cEnd - 1, "┤", "")
    BottomLeftBorder := iff(colI == t.xOffset, "└", "┴")
    RightBorder      := "│"
    LeftBorder       := "│"
    BottomBorder     := "─"
    TopBorder        := "─"

    if (isScrollbarRow) {
        topRightBorder = "█"
        RightBorder    = "█"
    }

    if (isScrollbarCol) {
        BottomBorder     = "▄"
        BottomLeftBorder = "▄"
    }

    style := lipgloss.NewStyle().
    Border(lipgloss.Border{
        Top:         TopBorder,
        Left:        LeftBorder,
        Right:       RightBorder,
        Bottom:      BottomBorder,
        BottomRight: "┘",
        BottomLeft : BottomLeftBorder,
        TopLeft:     topLeftBorder,
        TopRight:    topRightBorder,
    }). 
    BorderBottom(enableBottomBorder).
    BorderRight(enableRightBorder).
    BorderForeground(lipgloss.Color("240"))

    if !t.wouldBeFocused {
        return style
    }

    cursorX := t.cursor.X
    cursorY := t.cursor.Y

    if (t.selectionStart != -1 ) {
        if (t.rowSelect && isBetween(rowI, cursorY, t.selectionStart)) {
            style = style.Background(lipgloss.Color("58"))
        }
        if (t.columnSelect && isBetween(colI, cursorX, t.selectionStart)) {
            style = style.Background(lipgloss.Color("58"))
        }
    }

    if cursorX == colI && cursorY == rowI {
        style = style.Background(
            lipgloss.Color(iff(t.IsFocused(), "57", "240")),
        )
    }

    return style
}
