package log

import (
	"github.com/sirupsen/logrus"
)

func new() *logrus.Entry {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	result := logger.WithField("Service", "authgo")
	return result
}

func Get(ctx ...interface{}) *logrus.Entry {
	for _, o := range ctx {
		if tc, ok := o.(*logrus.Entry); ok {
			return tc
		}
	}
	return new()
}
