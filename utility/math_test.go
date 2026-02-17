package utility

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrap(t *testing.T) {
	tests := []struct {
		value, max, expected int
	}{
		{5, 10, 5},
		{15, 10, 5},
		{4, 3, 1},
		{34, 32, 2},
	}

	for _, test := range tests {
		result := Wrap(test.value, test.max)
		if result != test.expected {
			assert.Equal(t, test.expected, result, "Wrap(%d, %d) = %d; want %d", test.value, test.max, result, test.expected)
		}
	}
}
