package logger

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

var Logger = NewLogger()

func NewLogger() *logrus.Logger {
	return &logrus.Logger{
		Out:   io.MultiWriter(os.Stderr),
		Level: logrus.DebugLevel,
		Formatter: &logrus.TextFormatter{
			FullTimestamp:          true,
			TimestampFormat:        "2006-01-02 15:04:05",
			ForceColors:            true,
			DisableLevelTruncation: true,
		},
	}
}