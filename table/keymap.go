package table

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type KeyMap struct {
	LineUp       key.Binding
	LineDown     key.Binding
	LineLeft     key.Binding
	LineRight    key.Binding

	PageUp       key.Binding
	PageDown     key.Binding
	HalfPageUp   key.Binding
	HalfPageDown key.Binding
	GotoTop      key.Binding
	GotoBottom   key.Binding

	ScrollDown   key.Binding
    ScrollUp     key.Binding

    SelectRow    key.Binding
    SelectColumn key.Binding

    EndSelection key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		LineUp: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		LineDown: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		LineLeft: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "left"),
		),
		LineRight: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "right"),
		),

		HalfPageUp: key.NewBinding(
			key.WithKeys("pgup", "ctrl+u"),
			key.WithHelp("ctrl+u/pgup", "½ page up"),
		),
		HalfPageDown: key.NewBinding(
			key.WithKeys("pgdown", "ctrl+d"),
			key.WithHelp("ctrl+d/pgdown", "½ page down"),
		),

		GotoTop: key.NewBinding(
			key.WithKeys("home", "g"),
			key.WithHelp("g/home", "go to start"),
		),
		GotoBottom: key.NewBinding(
			key.WithKeys("end", "G"),
			key.WithHelp("G/end", "go to end"),
		),

        ScrollUp: key.NewBinding(key.WithKeys("ctrl+y")),
        ScrollDown: key.NewBinding(key.WithKeys("ctrl+e")),

        SelectRow: key.NewBinding(key.WithKeys("V")),
        SelectColumn: key.NewBinding(key.WithKeys("ctrl+v")),

        EndSelection: key.NewBinding(key.WithKeys("esc")),
	}
}

func (t Table) Update(msg tea.Msg) (Table, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, t.keyMap.LineUp):
			t.MoveUp(1)
		case key.Matches(msg, t.keyMap.LineDown):
			t.MoveDown(1)
		case key.Matches(msg, t.keyMap.LineLeft):
			t.MoveLeft(1)
		case key.Matches(msg, t.keyMap.LineRight):
			t.MoveRight(1)
		case key.Matches(msg, t.keyMap.HalfPageUp):
			t.MoveUp(t.height / 4)
		case key.Matches(msg, t.keyMap.HalfPageDown):
			t.MoveDown(t.height / 4)
		case key.Matches(msg, t.keyMap.LineDown):
			t.MoveDown(1)
		case key.Matches(msg, t.keyMap.GotoTop):
			t.GotoTop()
		case key.Matches(msg, t.keyMap.GotoBottom):
			t.GotoBottom()

        case key.Matches(msg, t.keyMap.ScrollUp):
            t.ScrollUp()
		case key.Matches(msg, t.keyMap.ScrollDown):
			t.ScrollDown()

        case key.Matches(msg, t.keyMap.SelectRow):
            t.SelectRow()
		case key.Matches(msg, t.keyMap.SelectColumn):
            t.SelectColumn()

		case key.Matches(msg, t.keyMap.EndSelection):
            t.columnSelect, t.rowSelect = false, false
		}
	}

	return t, nil
}

func (t *Table) ScrollUp() {
    t.yOffset = max(t.yOffset - 1, 0) 
}
func (t *Table) ScrollDown() {
    t.yOffset = min(t.yOffset + 1, max(len(t.rows) - (t.height - 1 - 2) / 2, 0))
}

func (t *Table) GotoTop() {
    t.MoveUp(t.cursor.Y)
}

func (t *Table) GotoBottom() {
    t.MoveDown(len(t.rows))
}

func (t *Table) MoveUp(i int) {
    if (len(t.rows) == 0) {
        return;
    }

    t.cursor.Y = clamp(t.cursor.Y-i, 0, len(t.rows)-1)
    t.UpdateOffset()
}

func (t *Table) MoveDown(i int) {
    if (len(t.rows) == 0) {
        return;
    }

    t.cursor.Y = clamp(t.cursor.Y+i, 0, len(t.rows)-1)
    t.UpdateOffset()
}

func (t *Table) MoveLeft(i int) {
    if (len(t.rows) == 0) {
        return;
    }

    t.cursor.X = clamp(t.cursor.X-i, 0, len(t.cols)-1)
    t.UpdateOffset()
}

func (t *Table) MoveRight(i int) {
    if (len(t.rows) == 0) {
        return;
    }

    t.cursor.X = clamp(t.cursor.X+i, 0, len(t.cols)-1)
    t.UpdateOffset()
}

func (t *Table) SelectRow() {
    t.rowSelect    = !t.rowSelect
    t.columnSelect = false

    if (t.rowSelect) {
        t.selectionStart = t.cursor.Y
    } else if (t.selectionStart >= 0) {
        t.selectionStart = -1
    }
}

func (t *Table) SelectColumn() {
    t.columnSelect = !t.columnSelect
    t.rowSelect    = false

    if (t.columnSelect) {
        t.selectionStart = t.cursor.X
    } else if (t.selectionStart >= 0) {
        t.selectionStart = -1
    }
}
