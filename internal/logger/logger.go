package logger

import (
	"os"
	"tt-copier/config"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var log = logrus.New()

func Init(cfg *config.Config) {
	log.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	if cfg.Env != "Prod" {
		log.SetOutput(os.Stdout)
	} else {
		fileOutput := &lumberjack.Logger{
			Filename:   cfg.LogPath,
			MaxSize:    10, // megabytes
			MaxBackups: 3,
			MaxAge:     28, // days
			Compress:   true,
		}

		log.SetOutput(fileOutput)
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
