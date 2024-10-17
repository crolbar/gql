package table

import (
	"gql/table/scrollbar"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

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

        style := t.generateStyleHeader(i, end)

        trunc        := runewidth.Truncate(col.Title, col.Width, "…")
		renderedCell := style.Width(col.Width).Render(trunc)
		s = append(s, renderedCell)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, s...)
}

func (t *Table) renderRows() string {
	rows := make([]string, 0, len(t.rows))

    start := clamp(t.yOffset, 0, len(t.rows))
    end   := clamp(t.yOffset + t.height / 2, 0, len(t.rows))

    vScrollbar := scrollbar.New(t.height / 2, len(t.rows), t.yOffset)

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
        value := t.rows[r][i]

		if t.cols[i].Width <= 0 {
			continue
		}

        isScrollbarCol := hScrollbar.IsScrollbarItem(i)

        style        := t.generateStyleRow(i, end, r, rEnd, isScrollbarRow, isScrollbarCol)
        trunc        := runewidth.Truncate(value, t.cols[i].Width, "…")
		renderedCell := style.Width(t.cols[i].Width).Render(trunc)

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
