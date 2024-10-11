package table

func isBetween(v, i, j int) bool {
    min := min(i, j)
    max := max(i, j)
    return v >= min && v <= max
}

func iff(cond bool, t, f string) string {
    if (cond) {
        return t
    }
    return f
}

func clamp(v, low, high int) int {
	return min(max(v, low), high)
}
