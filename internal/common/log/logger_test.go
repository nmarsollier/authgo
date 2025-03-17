package log

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLogRusEntry_WithField(t *testing.T) {
	entry := logrus.NewEntry(logrus.New())
	logEntry := &logRusEntry{entry: entry}

	key := "testKey"
	value := "testValue"

	newLogEntry := logEntry.WithField(key, value)

	assert.NotNil(t, newLogEntry)
	assert.Equal(t, value, newLogEntry.Data()[key])
	assert.Equal(t, value, logEntry.entry.Data[key])
}
