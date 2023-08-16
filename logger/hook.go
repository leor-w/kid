package logger

import "github.com/sirupsen/logrus"

type Hook interface {
	logrus.Hook
}
