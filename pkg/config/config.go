package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

func InitLogrus(logLevel string, logger *logrus.Logger) (file *os.File, err error) {
	// установим уровень логирования
	switch logLevel {
	case "trace":
		logger.SetLevel(logrus.TraceLevel)
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logger.SetLevel(logrus.FatalLevel)
	case "panic":
		logger.SetLevel(logrus.PanicLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	// установим форматирование логов в джейсоне
	logger.SetFormatter(&logrus.JSONFormatter{})

	// установим вывод логов в файл
	file, err = os.OpenFile("logs/subscription.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.SetOutput(file)
		return file, nil
	} else {
		logger.Info("Не удалось открыть файл логов, используется стандартный stderr")
		return nil, err
	}
}
