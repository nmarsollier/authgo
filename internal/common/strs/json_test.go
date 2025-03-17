package strs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToJson(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected string
	}{
		{struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}{"John", 30}, `{"name":"John","age":30}`},
		{map[string]int{"one": 1, "two": 2}, `{"one":1,"two":2}`},
		{[]string{"apple", "banana"}, `["apple","banana"]`},
		{"hello", `"hello"`},
	}

	for _, test := range tests {
		result := ToJson(test.input)
		assert.JSONEq(t, test.expected, result)
	}
}
