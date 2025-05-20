package gometrix

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

const LOG_FORMAT = "2006-01-02T15:04:05.999Z07:00"

func setupMetricsLogger(config LoggingMetricsData) *logrus.Logger {
	logger := logrus.New()

	if config.LogFilePath != "" {
		if err := os.MkdirAll(config.LogFilePath, 0755); err != nil {
			logrus.Errorf("Could not create metrics log folder: %v", err)
			return logger
		}

		fileLogger := &lumberjack.Logger{
			Filename:   strings.Join([]string{config.LogFilePath, "metrics.log"}, "/"),
			MaxSize:    int(config.MaxFileSize),
			MaxBackups: config.MaxFiles,
			Compress:   false,
		}

		logger.SetOutput(fileLogger)
	} else {
		logger.SetOutput(os.Stdout)
	}

	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:   true,
		FullTimestamp:   true,
		DisableQuote:    true,
		DisableSorting:  true,
		TimestampFormat: LOG_FORMAT,
	})
	logger.SetLevel(logrus.InfoLevel)

	return logger
}
