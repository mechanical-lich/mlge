package utility

import (
	"math/rand"
)

func GetRandom(low int, high int) int {
	if high < low {
		return low
	}
	if low == high {
		return low
	}
	return (rand.Intn((high - low))) + low
}

func Distance(x1 int, y1 int, x2 int, y2 int) int {
	var dy int
	if y1 > y2 {
		dy = y1 - y2
	} else {
		dy = (y2 - y1)
	}

	var dx int
	if x1 > x2 {
		dx = x1 - x2
	} else {
		dx = x2 - x1
	}

	var d int
	if dy > dx {
		d = dy + (dx >> 1)
	} else {
		d = dx + (dy >> 1)
	}

	return d
}

func Sgn(a int) int {
	switch {
	case a < 0:
		return -1
	case a > 0:
		return +1
	}
	return 0
}

func Clamp(value, min, max int) int {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}

func Wrap(value, max int) int {
	for value > max {
		value = value - max
	}

	return value
}

func Abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func RectsOverlap(x1, y1, w1, h1, x2, y2, w2, h2 int) bool {
	if x1 < x2+w2 && x1+w1 > x2 && y1 < y2+h2 && y1+h1 > y2 {
		return true
	}
	return false
}

func RectContains(x, y, w, h, px, py int) bool {
	if px >= x && px <= x+w && py >= y && py <= y+h {
		return true
	}
	return false
}
