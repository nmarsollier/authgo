package strs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAtoiDefault(t *testing.T) {
	tests := []struct {
		input    string
		def      int
		expected int
	}{
		{"123", 0, 123},
		{"abc", 5, 5},
		{"", 10, 10},
		{"123abc", 7, 7},
	}

	for _, test := range tests {
		result := AtoiDefault(test.input, test.def)
		assert.Equal(t, test.expected, result)
	}
}
