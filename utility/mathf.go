package utility

import "math"

func ClampF(value, min, max float64) float64 {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}

func WrapF(value, max float64) float64 {
	for value > max {
		value = value - max
	}

	return value
}

func AbsF(a float64) float64 {
	if a < 0 {
		return -a
	}
	return a
}

func MaxF(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func MinF(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func Lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

func LerpAngle(a, b, t float64) float64 {
	diff := b - a
	for diff > math.Pi {
		diff -= 2 * math.Pi
	}
	for diff < -math.Pi {
		diff += 2 * math.Pi
	}
	return a + diff*t
}
