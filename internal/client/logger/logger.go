// Package logger используется для инициализации логгера
package logger

import (
	"path/filepath"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// New - Возвращает инстанс нового логгера
func New(root string) *logrus.Logger {
	log := logrus.New()
	infoLevel := filepath.Join(root, "info.log")
	errorLevel := filepath.Join(root, "error.log")
	pathMap := lfshook.PathMap{
		logrus.InfoLevel:  infoLevel,
		logrus.ErrorLevel: errorLevel,
		logrus.FatalLevel: errorLevel,
	}
	log.Hooks.Add(lfshook.NewHook(
		pathMap,
		&logrus.JSONFormatter{},
	))

	return log
}
