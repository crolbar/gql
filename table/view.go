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

func (t Table) generateStyleHeader(i, end int) lipgloss.Style {
    topLeftBorder       :=  iff(i == t.XOffset, "┌", "┬")
    topRightBorder      :=  iff(i == end - 1, "┐", "")
    disableRightBorder  :=  i == end - 1


    return lipgloss.NewStyle().
    Border(lipgloss.Border{
        Top:      "─",
        Left:     "│",
        Right:    "│",
        TopLeft:  topLeftBorder,
        TopRight: topRightBorder,
    }). 
    BorderBottom(false).
    BorderRight(disableRightBorder).
    BorderForeground(lipgloss.Color("240")).
    Bold(true)
}

func (t Table) generateStyleRow(i, end, r, rEnd int) lipgloss.Style {
    disableRightBorder  := i == end - 1
    disableBottomBorder := r == rEnd - 1

    topLeftBorder    := iff(i == t.XOffset, "├", "┼")
    topRightBorder   := iff(i == end - 1, "┤", "")
    BottomLeftBorder := iff(i == t.XOffset, "└", "┴")
    RightBorder      := iff(end == len(t.cols), "│", ">")
    LeftBorder       := iff(i == t.XOffset && i != 0, "<", "│")


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
    BorderBottom(disableBottomBorder).
    BorderRight(disableRightBorder).
    BorderForeground(lipgloss.Color("240"))
    

    cursorX := t.Cursor.x
    cursorY := t.Cursor.y

    switch {
    case cursorX == i && cursorY != r && t.columnSelect:// column
        style = style.Background(lipgloss.Color("57"))

    case cursorX == i && cursorY == r:                  // single
        style = style.Background(lipgloss.Color("57")) 

    case cursorX != i && cursorY == r && t.rowSelect:   // row
        style = style.Background(lipgloss.Color("57"))
    }

    return style
}

func iff(cond bool, t, f string) string {
    if (cond) {
        return t
    }
    return f
}
