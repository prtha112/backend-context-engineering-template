package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func New(level string) *logrus.Logger {
	logger := logrus.New()

	// Set log level
	logLevel, err := logrus.ParseLevel(strings.ToLower(level))
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	logger.SetLevel(logLevel)

	// Set formatter
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	})

	// Set output
	logger.SetOutput(os.Stdout)

	return logger
}
