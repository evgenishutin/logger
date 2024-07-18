package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type NotificatorInterface interface {
	SendMessage(info map[string]interface{})
}

type LoggerInterface interface {
	MapFields(fields map[string]interface{}) logrus.Fields
	WrapperNotifyError(msg map[string]interface{}) logrus.Fields
	WithFields(fields logrus.Fields) *logrus.Entry
}

type Log struct {
	*logrus.Logger
	notify NotificatorInterface
}

func New(level string, notify NotificatorInterface) LoggerInterface {
	var log = logrus.New()
	var lev logrus.Level

	switch strings.ToLower(level) {
	case "error":
		lev = 2
	case "warn":
		lev = 3
	case "info":
		lev = 4
	case "debug":
		lev = 5
	default:
		lev = 4
	}

	jsonFormatter := &logrus.JSONFormatter{
		PrettyPrint:     true,
		TimestampFormat: "2006-01-02 15:04:05",
	}

	log.SetLevel(lev)
	log.SetFormatter(jsonFormatter)
	log.SetOutput(os.Stdout)

	return &Log{
		log,
		notify,
	}
}

func (l *Log) MapFields(fields map[string]interface{}) logrus.Fields {
	log := logrus.Fields{}

	for k, v := range fields {
		log[k] = v
	}

	return log
}

func (l *Log) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.Logger.WithFields(fields)
}

func (l *Log) WrapperNotifyError(msg map[string]interface{}) logrus.Fields {
	l.notify.SendMessage(msg)

	return l.MapFields(msg)
}
