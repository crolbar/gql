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

    start := clamp(t.XOffset, 0, len(t.cols))
    end   := clamp(t.renderedColumns + t.XOffset, 0, len(t.cols))

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

    start := clamp(t.YOffset, 0, len(t.rows))
    end   := clamp(t.YOffset + t.Height / 2, 0, len(t.rows))

    scrollbar := scrollbar.New(
        t.Height,
        len(t.rows),
        end,
        t.YOffset,
    )

	for i := start; i < end; i++ {
		rows = append(rows, t.renderRow(i, end, &scrollbar))
    }

    return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (t *Table) renderRow(r, rEnd int, sb *scrollbar.Scrollbar) string {
	s := make([]string, 0, len(t.cols))

    start := clamp(t.XOffset, 0, len(t.cols))
    end   := clamp(t.renderedColumns + t.XOffset, 0, len(t.cols))

    for i := start; i < end; i++ {
        value := t.rows[r][i]

		if t.cols[i].Width <= 0 {
			continue
		}

        style        := t.generateStyleRow(i, end, r, rEnd, sb.IsScrollbarRow(r))
        trunc        := runewidth.Truncate(value, t.cols[i].Width, "…")
		renderedCell := style.Width(t.cols[i].Width).Render(trunc)

		s = append(s, renderedCell)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, s...)
}

func (t Table) generateStyleHeader(colI, end int) lipgloss.Style {
    topLeftBorder       := iff(colI == t.XOffset, "┌", "┬")
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

func (t Table) generateStyleRow(colI, cEnd, rowI, rEnd int, isScrollbarRow bool) lipgloss.Style {
    enableRightBorder   := colI == cEnd - 1
    enableBottomBorder  := rowI == rEnd - 1

    topLeftBorder    := iff(colI == t.XOffset, "├", "┼")
    topRightBorder   := iff(colI == cEnd - 1, "┤", "")
    BottomLeftBorder := iff(colI == t.XOffset, "└", "┴")
    RightBorder      := iff(isScrollbarRow, "█", "│")
    LeftBorder       := iff(colI == t.XOffset && colI != 0, "<", "│")
    BottomBorder     := iff(rowI == rEnd - 1 && rEnd != len(t.rows), "˯", "─") // kidna ugly. replace with scrollbar ?
    TopBorder        := iff(rowI == t.YOffset && t.YOffset != 0, "˄", "─")


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
