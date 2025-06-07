package util

import (
	log "github.com/sirupsen/logrus"
)

func Logger() *log.Logger {
	logger := log.New()
	logger.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	return logger
}
