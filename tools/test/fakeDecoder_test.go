package test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecoder(t *testing.T) {
	err := errors.New("Test")

	decoder := FakeDecoder(func(v interface{}) error {
		return err
	})

	assert.Equal(t, err, decoder.Decode(""))

	decoder = FakeDecoder(func(v interface{}) error {
		return nil
	})

	assert.Nil(t, decoder.Decode(""))
}
