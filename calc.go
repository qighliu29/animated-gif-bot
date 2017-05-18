package main

func mustInRange(l, r, num int) int {
	switch {
	case num < l:
		return l
	case num > r:
		return r
	default:
		return num
	}
}

func intersec(l1, r1, l2, r2 int) (l, r int) {
	l = mustInRange(l1, r1, l2)
	r = mustInRange(l, r1, r2)
	return
}
