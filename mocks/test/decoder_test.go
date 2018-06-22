package test

import (
	"errors"
	"testing"

	"github.com/nmarsollier/authgo/mocks"
	"github.com/stretchr/testify/assert"
)

func TestDecoder(t *testing.T) {
	err := errors.New("Test")

	decoder := mocks.Decoder(func(v interface{}) error {
		return err
	})

	assert.Equal(t, err, decoder.Decode(""))

	decoder = mocks.Decoder(func(v interface{}) error {
		return nil
	})

	assert.Nil(t, decoder.Decode(""))
}
