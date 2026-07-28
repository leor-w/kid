package simple

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/lsjhtang/kid/plugin/logrus/formatters"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"path"
	"time"
)

func NewSimpleRotate(opts ...Option) (logrus.Hook, error) {
	options := Options{
		path:    "./logs",
		logName: "kid",
		link:    "",
		rotate:  24,
		maxAge:  15,
	}
	for _, o := range opts {
		o(&options)
	}
	writer, err := rotatelogs.New(
		path.Join(options.path, options.logName+"_%Y%m%d.log"),
		rotatelogs.WithLinkName(options.link),
		rotatelogs.WithRotationTime(options.rotate*time.Hour),
		rotatelogs.WithMaxAge(options.maxAge*time.Hour*24),
	)
	if err != nil {
		return nil, err
	}
	return lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, formatters.NewSimpleFormatter()), nil
}
