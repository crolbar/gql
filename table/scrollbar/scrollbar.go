package scrollbar

type Scrollbar struct {
    scrollbarHeight int
    scrollbarStartPos int
    currScrollbarHeight int
    shouldShowScrollbar bool

    YOffset int
}


func New(
    height int,
    rowsLen int,
    end int,
    YOffset int,
) Scrollbar {
    rowsInView      := height / 2
	viewRatio       := float64(rowsInView) / float64(rowsLen)
	scrollbarHeight := int(viewRatio * float64(rowsInView))

	if scrollbarHeight < 1 {
		scrollbarHeight = 1
	}

    scrollbarStartPos := int((float64(YOffset) / float64(rowsLen - height / 2)) * float64((height / 2) - scrollbarHeight))
    return Scrollbar {
        scrollbarHeight:     scrollbarHeight,
        shouldShowScrollbar: end != rowsLen || YOffset > 0,
        scrollbarStartPos:   scrollbarStartPos,
        currScrollbarHeight: 0,
        YOffset:             YOffset,
    }
}

func (s *Scrollbar) IsScrollbarRow(i int) bool {
    if (!s.shouldShowScrollbar) {
        return false
    }
    isScrollbarRow := (i - s.YOffset == s.scrollbarStartPos) ||
        s.currScrollbarHeight > 0 && s.currScrollbarHeight < s.scrollbarHeight

    if (isScrollbarRow) {
        s.currScrollbarHeight++;
    }

    return isScrollbarRow
}
