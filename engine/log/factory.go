package log

import (
	"fmt"
	"net"

	"github.com/nmarsollier/authgo/engine/env"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

const LOG_FIELD_CORRELATION_ID = "correlation_id"
const LOG_FIELD_CONTROLLER = "controller"
const LOG_FIELD_RABBIT_ACTION = "rabbit_action"
const LOG_FIELD_RABBIT_EXCHANGE = "exchange"
const LOG_FIELD_RABBIT_QUEUE = "queue"
const LOG_FIELD_HTTP_METHOD = "http_method"
const LOG_FIELD_HTTP_PATH = "http_path"
const LOG_FIELD_HTTP_STATUS = "http_status"
const LOG_FIELD_SERVER = "server"
const LOG_FIELD_USER_ID = "user_id"
const LOG_FIELD_THREAD = "thread"

func Get() LogRusEntry {
	logger := logrus.New()
	configureFluent(logger)

	logger.SetLevel(logrus.DebugLevel)
	logEntry := logger.WithField(LOG_FIELD_SERVER, "authgo").WithField(LOG_FIELD_THREAD, uuid.NewV4().String())

	return logRusEntry{entry: logEntry}
}

func configureFluent(logger *logrus.Logger) {
	conn, err := net.Dial("tcp", env.Get().FluentUrl)
	if err == nil {
		logger.SetFormatter(&logrus.JSONFormatter{})
		logger.SetOutput(conn)
		logger.AddHook(&logrusConnectionHook{
			conn: conn,
			fmt: &logrus.TextFormatter{
				FullTimestamp:  true,
				ForceColors:    true,
				DisableSorting: false,
			},
		})
	}
}

type logrusConnectionHook struct {
	conn net.Conn
	fmt  logrus.Formatter
}

func (hook *logrusConnectionHook) Fire(entry *logrus.Entry) error {
	msg, err := hook.fmt.Format(entry)
	if err == nil {
		fmt.Println(string(msg))
	}
	return nil
}

func (hook *logrusConnectionHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *logrusConnectionHook) Close() error {
	hook.conn.Close()
	return nil
}
