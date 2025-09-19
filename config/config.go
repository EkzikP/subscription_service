package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func InitLogrus(logLevel string) (file *os.File, err error) {
	// установим уровень логирования
	switch logLevel {
	case "trace":
		log.SetLevel(logrus.TraceLevel)
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	case "fatal":
		log.SetLevel(logrus.FatalLevel)
	case "panic":
		log.SetLevel(logrus.PanicLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}

	// установим форматирование логов в джейсоне
	log.SetFormatter(&logrus.JSONFormatter{})

	// установим вывод логов в файл
	file, err = os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
		return file, nil
	} else {
		log.Info("Не удалось открыть файл логов, используется стандартный stderr")
		return nil, err
	}
}

func InitConfig() {}
