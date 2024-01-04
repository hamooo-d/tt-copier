package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func Init() {
	log.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	log.Out = os.Stdout

	file, err := os.OpenFile("./logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err == nil {
		log.Out = io.MultiWriter(file, os.Stdout)
	} else {
		Warn("Failed to log to file, using default stderr", "INIT", "FAILED")
	}

}

func Info(message string, action string, status string) {
	log.WithFields(logrus.Fields{
		"message": message,
		"status":  status,
		"action":  action,
	}).Info()
}

func Warn(message string, action string, status string) {
	log.WithFields(logrus.Fields{
		"message": message,
		"status":  status,
		"action":  action,
	}).Warn()
}

func Error(message string, err error) {
	log.WithFields(logrus.Fields{
		"message": message,
		"error":   err,
	}).Error()
}
