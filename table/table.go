package table

import (
	"math"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
)

type Cursor struct {
	X int
	Y int
}

type Table struct {
	cols []Column
	rows []Row

	cursor Cursor

	height    int
	width     int
	maxHeight int
	maxWidth  int

	keyMap KeyMap
	help   help.Model

	xOffset int
	yOffset int

	renderedColumns int

	columnSelect   bool
	rowSelect      bool
	selectionStart int

	focused        bool
	wouldBeFocused bool
}

type Row []string

type Column struct {
	Title string
	Width int
}

func New(cols []Column, rows []Row, height int, width int) Table {
	help := help.New()

	help.ShowAll = true
	help.Styles.FullKey = lipgloss.NewStyle().Bold(true)
	help.Styles.FullDesc = lipgloss.NewStyle().Italic(true)
	help.Styles.FullSeparator = lipgloss.NewStyle()
	help.FullSeparator = ""

	return Table{
		cols: cols,
		rows: rows,

		cursor: Cursor{0, 0},

		height:    height - 2,
		width:     width,
		maxHeight: height - 2,
		maxWidth:  width,

		xOffset: 0,
		yOffset: 0,

		renderedColumns: 1,

		keyMap: DefaultKeyMap(),
		help:   help,

		columnSelect:   false,
		rowSelect:      false,
		selectionStart: -1,

		focused:        false,
		wouldBeFocused: true,
	}
}

func (t Table) IsFocused() bool {
	return t.focused
}

func (t *Table) Focus() {
	t.focused = true
}

func (t *Table) DeFocus() {
	t.focused = false
}

// uses the table only to display something
// currently only difference is the unfocused cursor pos
func (t *Table) SetDisplayOnly() {
	t.wouldBeFocused = false
}

func (t *Table) SetColumns(cols []Column) {
	if cols == nil {
		cols = []Column{{Title: "Empty", Width: 20}}
	}
	t.cols = cols
	if t.cursor.X >= len(t.cols) {
		t.cursor.X = 0
	}

	t.UpdateOffset()
}

func (t *Table) SetRows(rows []Row) {
	t.rows = rows
	if t.cursor.Y >= len(t.rows) {
		t.cursor.Y = 0
	}

	t.UpdateOffset()
}

func (t *Table) GetColumns() []Column {
	return t.cols
}

func (t *Table) GetRows() []Row {
	return t.rows
}

func (t *Table) GetCursor() Cursor {
	return t.cursor
}

func (t *Table) GetWidth() int {
	return t.width
}

func (t *Table) GetHeight() int {
	return t.height
}

func (t *Table) SetMaxSize(width, height int) {
	t.maxHeight = height
	t.maxWidth = width
}

func (t *Table) SetMaxHeight(height int) {
	t.maxHeight = height
}

func (t *Table) SetMaxWidth(width int) {
	t.maxWidth = width
}

func (t *Table) GetXOffset() int {
	return t.xOffset
}

func (t *Table) GetYOffset() int {
	return t.yOffset
}

func (t *Table) IsSelectingRows() bool {
	return t.rowSelect
}

func (t *Table) IsSelectingCols() bool {
	return t.columnSelect
}

func (t *Table) GetSelectionStart() int {
	return t.selectionStart
}

func (t Table) GetCols() []Column {
	return t.cols
}

func (t Table) GetSelColumnName() string {
	if len(t.cols) == 0 {
		return ""
	}

	return t.cols[t.cursor.X].Title
}

func (t Table) GetSelectedCell() string {
	if len(t.rows) == 0 {
		return ""
	}

	selCell := t.rows[t.cursor.Y][t.cursor.X]

	return strings.ReplaceAll(selCell, "\\n", "\n")
}

func (t *Table) GetSelectedRow() Row {
	if len(t.rows) == 0 {
		return nil
	}

	return t.rows[t.cursor.Y]
}

// returns nil if we aren't in selection
func (t *Table) GetSelectedRows() []Row {
	if t.rowSelect && t.selectionStart >= 0 {
		start := min(t.selectionStart, t.cursor.Y)
		end := max(t.selectionStart, t.cursor.Y)
		return t.rows[start : end+1]
	}
	return []Row{t.GetSelectedRow()}
}

// returns nil if we aren't in selection
func (t *Table) GetSelectedColumns() []Column {
	if t.columnSelect && t.selectionStart >= 0 {
		start := min(t.selectionStart, t.cursor.X)
		end := max(t.selectionStart, t.cursor.X)
		return t.cols[start : end+1]
	}
	return nil
}

// Must be called when the width of the terminal changes
// or there is an update to cursor.x
func (t *Table) UpdateRenderedColums() {
	width := 0
	for i := t.xOffset; i < len(t.cols); i++ {
		currColWidth := t.cols[i].Width + 1
		width += currColWidth

		if width >= t.maxWidth {
			t.renderedColumns = i - t.xOffset
			t.width = min((width-currColWidth)+1, t.maxWidth)
			return
		}
	}

	t.renderedColumns = len(t.cols) - t.xOffset

	t.width = min(width+1, t.maxWidth)
}

// Must be called when the width of the terminal changes
// or there is an update to the cursor
func (t *Table) UpdateOffset() {
	// we need 2 rows for each item (so * 2)
	// + 1 for bottom border & + 2 for the header
	t.height = min((len(t.rows)*2)+1+2, t.maxHeight)

	itemsInView := (t.height - 1 - 2) / 2

	lines_till_sow := t.cursor.Y - t.yOffset
	lines_till_eow := (itemsInView - 1) - lines_till_sow

	if lines_till_eow < 0 {
		t.yOffset += int(math.Abs(float64(lines_till_eow)))
	}

	if lines_till_sow < 0 {
		t.yOffset -= int(math.Abs(float64(lines_till_sow)))
	}

	t.UpdateRenderedColums()

	cols_till_sow := t.cursor.X - t.xOffset
	cols_till_eow := (t.renderedColumns - 1) - cols_till_sow

	if cols_till_eow < 0 {
		t.xOffset += int(math.Abs(float64(cols_till_eow)))
	}

	if cols_till_sow < 0 {
		t.xOffset -= int(math.Abs(float64(cols_till_sow)))
	}

	t.UpdateRenderedColums()
}
