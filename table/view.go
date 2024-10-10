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
    end := t.renderedColumns + t.XOffset

    for i := start; i < end; i++ {
        col := t.cols[i]
		if col.Width <= 0 {
			continue
		}

        style := t.generateStyleHeader(i, end)

        trunc := runewidth.Truncate(col.Title, col.Width, "…")
		renderedCell := style.Width(col.Width).Render(trunc)
		s = append(s, renderedCell)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, s...)
}

func (t *Table) renderRows() string {
	rows := make([]string, 0, len(t.rows))

    start := t.YOffset
    end := clamp(t.YOffset + t.Height / 2, 0, len(t.rows))

	for i := start; i < end; i++ {
		rows = append(rows, t.renderRow(i))
    }

    return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (t *Table) renderRow(r int) string {
	s := make([]string, 0, len(t.cols))

    start := t.XOffset
    end := t.renderedColumns + t.XOffset

    for i := start; i < end; i++ {
        value := t.rows[r][i]

		if t.cols[i].Width <= 0 {
			continue
		}

        style := t.generateStyleRow(i, end, r)

        trunc := runewidth.Truncate(value, t.cols[i].Width, "…")
		renderedCell := style.Render(trunc)
		s = append(s, renderedCell)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, s...)
}

func (t Table) generateStyleHeader(i, end int) lipgloss.Style {
    disable_right := i == end - 1

    style := lipgloss.NewStyle().
    Border(lipgloss.Border{
        Top:    "─",
        Left:   "│",
        Right:  "|",
        Bottom: "",
        TopLeft: "┼",
    }). 
    BorderBottom(false).
    BorderRight(disable_right).
    BorderForeground(lipgloss.Color("240"))
    return style
}

func (t Table) generateStyleRow(i, end, r int) lipgloss.Style {
    disable_right := i == end - 1

    style := lipgloss.NewStyle().
    Border(lipgloss.Border{
        Top:    "─",
        Left:   "│",
        Right:  "|",
        Bottom: "",
        TopLeft: "┼",
    }). 
    BorderBottom(false).
    BorderRight(disable_right).
    BorderForeground(lipgloss.Color("240")).
    Width(t.cols[i].Width)

    if t.Cursor.x == i && t.Cursor.y != r {
        style = style.Background(lipgloss.Color("57"))
    } else if t.Cursor.x == i && t.Cursor.y == r {
        style = style.Background(lipgloss.Color("57"))
    } else if t.Cursor.x != i && t.Cursor.y == r {
        style = style.Background(lipgloss.Color("57"))
    }

    return style
}
