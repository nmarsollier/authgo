package log

import "github.com/sirupsen/logrus"

type LogRusEntry interface {
	Data() logrus.Fields
	Info(args ...interface{})
	Error(args ...interface{})
	Warn(args ...interface{})
	Fatal(args ...interface{})
	WithField(key string, value interface{}) LogRusEntry
}

type logRusEntry struct {
	entry *logrus.Entry
}

func (l logRusEntry) Error(args ...interface{}) {
	l.entry.Error(args...)
}

func (l logRusEntry) WithField(key string, value interface{}) LogRusEntry {
	l.entry.WithField(key, value)
	return l
}

func (l logRusEntry) Info(args ...interface{}) {
	l.entry.Info(args...)
}

func (l logRusEntry) Warn(args ...interface{}) {
	l.entry.Warn(args...)
}

func (l logRusEntry) Fatal(args ...interface{}) {
	l.entry.Fatal(args...)
}

func (l logRusEntry) Data() logrus.Fields {
	return l.entry.Data
}
