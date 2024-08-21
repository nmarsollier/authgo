package log

import (
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

func new() *logrus.Entry {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	result := logger.WithField("Service", "authgo").WithField("Thread", uuid.NewV4().String())
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
