package scrollbar

type Scrollbar struct {
	scrollbarSize       int
	scrollbarStartPos   int
	currScrollbarSize   int
	shouldShowScrollbar bool

	offset int
}

func New(
	itemsInView int,
	itemsLen int,
	offset int,
) Scrollbar {
	viewRatio := float64(itemsInView) / float64(itemsLen)
	scrollbarHeight := int(viewRatio * float64(itemsInView))

	if scrollbarHeight < 1 {
		scrollbarHeight = 1
	}

	scrollbarStartPos := int(
		(float64(offset) / float64(itemsLen-itemsInView)) *
			float64(itemsInView-scrollbarHeight),
	)

	lastItemInViewIdx := clamp(itemsInView+offset, 0, itemsLen)
	return Scrollbar{
		shouldShowScrollbar: offset > 0 || lastItemInViewIdx != itemsLen,
		scrollbarSize:       scrollbarHeight,
		currScrollbarSize:   0,
		scrollbarStartPos:   scrollbarStartPos,
		offset:              offset,
	}
}

func (s *Scrollbar) IsScrollbarItem(i int) bool {
	if !s.shouldShowScrollbar {
		return false
	}

	isScrollbarRow := (i-s.offset == s.scrollbarStartPos) ||
		s.currScrollbarSize > 0 && s.currScrollbarSize < s.scrollbarSize

	if isScrollbarRow {
		s.currScrollbarSize++
	}

	return isScrollbarRow
}

func clamp(v, low, high int) int {
	return min(max(v, low), high)
}
