// Package logger используется для инициализации логгера
package logger

import "github.com/sirupsen/logrus"

// New - Возвращает инстанс нового логгера
func New() *logrus.Logger {
	return logrus.New()
}
