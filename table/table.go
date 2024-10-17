package table

import (
	"fmt"
	"log"
	"math"
	"strings"
)

type Cursor struct {
    X int
    Y int
}

type Table struct {
    cols []Column
    rows []Row

	cursor Cursor

    height int
    width  int

    keyMap KeyMap

    xOffset int
    yOffset int

    renderedColumns int

    columnSelect   bool
    rowSelect      bool
    selectionStart int

    focused bool
}


type Row []string

type Column struct {
	Title string
	Width int
}

func New(cols []Column, rows []Row, height int, width int) Table {
    return Table {
        cols:   cols,
        rows:   rows,
        height: height - 2,
        width:  width,

        xOffset: 0,
        yOffset: 0,

        renderedColumns: 1,

        keyMap: DefaultKeyMap(),
        cursor: Cursor{0, 0},

        columnSelect:   false,
        rowSelect:      false,
        selectionStart: -1,


        focused: false,
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

func (t *Table) SetColumns(cols []Column) {
    t.cols     = cols
    t.cursor.X = 0

    t.UpdateOffset()
}

func (t *Table) SetRows(rows []Row) {
    t.rows     = rows
    t.cursor.Y = 0

    t.UpdateOffset()
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

func (t *Table) GetXOffset() int {
    return t.xOffset
}

func (t *Table) GetYOffset() int {
    return t.yOffset
}

func (t Table) GetSelectedCell() string {
    if (len(t.rows) == 0) {
        return ""
    }

    selCell := t.rows[t.cursor.Y][t.cursor.X]

    return strings.ReplaceAll(selCell, "\\n", "\n")
}

func (t *Table) GetSelectedRow() Row {
    return t.rows[t.cursor.Y]
}

// returns nil if we aren't in selection
func (t *Table) GetSelectedRows() []Row {
    if (t.rowSelect && t.selectionStart >= 0) {
        start := min(t.selectionStart, t.cursor.Y)
        end := max(t.selectionStart, t.cursor.Y)
        return t.rows[start:end+1]
    }
    return nil
}

// returns nil if we aren't in selection
func (t *Table) GetSelectedColumns() []Column {
    if (t.columnSelect && t.selectionStart >= 0) {
        start := min(t.selectionStart, t.cursor.X)
        end := max(t.selectionStart, t.cursor.X)
        return t.cols[start:end+1]
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

        if currColWidth > t.width {
            log.Fatal(fmt.Sprintf("Column %d's width is bigger than the table's width %d", i, t.width))
        }

        if width > t.width {
            t.renderedColumns = i - t.xOffset
            return;
        }
    }

    t.renderedColumns = len(t.cols) - t.xOffset
}

// Must be called when the width of the terminal changes 
// or there is an update to the cursor
func (t *Table) UpdateOffset()  {
    lines_till_sow := t.cursor.Y - t.yOffset
    lines_till_eow := ((t.height / 2) - 1) - lines_till_sow

    if (lines_till_eow < 0) {
        t.yOffset += int(math.Abs(float64(lines_till_eow)))
    } 

    if (lines_till_sow < 0) {
        t.yOffset -= int(math.Abs(float64(lines_till_sow)))
    }

    t.UpdateRenderedColums()

    cols_till_sow := t.cursor.X - t.xOffset
    cols_till_eow := (t.renderedColumns - 1) - cols_till_sow

    //t.Dbg = fmt.Sprintf("cols_t_sow: %d, cols_t_eow: %d, rc: %d", cols_till_sow, cols_till_eow, t.renderedColumns)

    for cols_till_eow < 0 || cols_till_sow < 0 {
        if (cols_till_eow < 0) {
            t.xOffset += int(math.Abs(float64(cols_till_eow)))
        }

        if (cols_till_sow < 0) {
            t.xOffset -= int(math.Abs(float64(cols_till_sow)))
        }

        t.UpdateRenderedColums()

        cols_till_sow = t.cursor.X - t.xOffset
        cols_till_eow = (t.renderedColumns - 1) - cols_till_sow
    }
}
