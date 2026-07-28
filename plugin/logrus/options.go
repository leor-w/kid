package logrus

import (
	"github.com/lsjhtang/kid/logger"
	"github.com/sirupsen/logrus"
)

type Options struct {
	logger.Options
	formatter  logrus.Formatter
	reportCall bool
}

func WithLogLevel(logLevel uint32) logger.Option {
	return func(options *logger.Options) {
		options.Level = logger.Level(logLevel)
	}
}

type formatterKey struct{}

func WithFormatter(formatter logrus.Formatter) logger.Option {
	return logger.SetOption(formatterKey{}, formatter)
}

type reportCallKey struct{}

func WithReportCall(reportCall bool) logger.Option {
	return logger.SetOption(reportCallKey{}, reportCall)
}
