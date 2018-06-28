package test

import (
	"testing"

	"github.com/nmarsollier/authgo/tests/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCursor(t *testing.T) {
	type testStruct struct {
		Name string
	}
	names := make([]interface{}, 2)
	names[0] = testStruct{Name: "uno"}
	names[1] = testStruct{Name: "dos"}

	cur := mocks.Cursor(names)

	result := testStruct{}
	assert.Equal(t, cur.Next(nil), true)
	cur.Decode(&result)
	assert.Equal(t, "uno", result.Name)

	assert.Equal(t, cur.Next(nil), true)
	cur.Decode(&result)
	assert.Equal(t, "dos", result.Name)

	assert.Equal(t, false, cur.Next(nil))
}
