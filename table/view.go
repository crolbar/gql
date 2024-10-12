package table

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

func (t Table) View() string {
	return t.headersView() + "\n" + t.renderRows()
}

func (t Table) headersView() string {
	s := make([]string, 0, len(t.cols))

    start := t.XOffset
    end   := t.renderedColumns + t.XOffset

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

    start := t.YOffset
    end   := clamp(t.YOffset + t.Height / 2, 0, len(t.rows))

	for i := start; i < end; i++ {
		rows = append(rows, t.renderRow(i, end))
    }

    return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (t *Table) renderRow(r, rEnd int) string {
	s := make([]string, 0, len(t.cols))

    start := t.XOffset
    end   := t.renderedColumns + t.XOffset

    for i := start; i < end; i++ {
        value := t.rows[r][i]

		if t.cols[i].Width <= 0 {
			continue
		}

        style := t.generateStyleRow(i, end, r, rEnd)

        trunc        := runewidth.Truncate(value, t.cols[i].Width, "…")
		renderedCell := style.Width(t.cols[i].Width).Render(trunc)

		s = append(s, renderedCell)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, s...)
}

func (t Table) generateStyleHeader(colI, end int) lipgloss.Style {
    topLeftBorder       :=  iff(colI == t.XOffset, "┌", "┬")
    topRightBorder      :=  iff(colI == end - 1, "┐", "")
    enableRightBorder   :=  colI == end - 1


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

func (t Table) generateStyleRow(colI, cEnd, rowI, rEnd int) lipgloss.Style {
    enableRightBorder   := colI == cEnd - 1
    enableBottomBorder  := rowI == rEnd - 1

    topLeftBorder    := iff(colI == t.XOffset, "├", "┼")
    topRightBorder   := iff(colI == cEnd - 1, "┤", "")
    BottomLeftBorder := iff(colI == t.XOffset, "└", "┴")
    RightBorder      := iff(cEnd == len(t.cols), "│", ">")
    LeftBorder       := iff(colI == t.XOffset && colI != 0, "<", "│")


    style := lipgloss.NewStyle().
    Border(lipgloss.Border{
        Top:         "─",
        Left:        LeftBorder,
        Right:       RightBorder,
        Bottom:      "─",
        BottomRight: "┘",
        BottomLeft : BottomLeftBorder,
        TopLeft:     topLeftBorder,
        TopRight:    topRightBorder,
    }). 
    BorderBottom(enableBottomBorder).
    BorderRight(enableRightBorder).
    BorderForeground(lipgloss.Color("240"))
    

    cursorX := t.Cursor.X
    cursorY := t.Cursor.Y

    if (t.selectionStart != -1 ) {
        if (t.rowSelect && isBetween(rowI, cursorY, t.selectionStart)) {
            style = style.Background(lipgloss.Color("58"))
        }
        if (t.columnSelect && isBetween(colI, cursorX, t.selectionStart)) {
            style = style.Background(lipgloss.Color("58"))
        }
    }

    if cursorX == colI && cursorY == rowI {
        style = style.Background(lipgloss.Color("57")) 
    }

    return style
}
