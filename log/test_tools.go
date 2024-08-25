package log

import "github.com/sirupsen/logrus"

func NewTestLogger() *logrus.Entry {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	return logrus.NewEntry(logger)
}
