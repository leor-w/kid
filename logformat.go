package kid

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

func LoggerWithFormatter(f LogFormatter) Middleware {
	return func(ctx *Context) {
		gin.LoggerWithConfig(gin.LoggerConfig{
			Formatter: convertLogFormatter(f),
		})
	}
}

type LogFormatter func(params LogFormatterParams) string

type LogFormatterParams struct {
	Request      *http.Request
	TimeStamp    time.Time
	StatusCode   int
	Latency      time.Duration
	ClientIP     string
	Method       string
	Path         string
	ErrorMessage string
	isTerm       bool
	BodySize     int
	Keys         map[string]interface{}
}

type LoggerConfig struct {
	Formatter LogFormatter
	Output    io.Writer
	SkipPaths []string
}

func LoggerWithConfig(conf LoggerConfig) Middleware {
	return func(ctx *Context) {
		gin.LoggerWithConfig(gin.LoggerConfig{
			Formatter: convertLogFormatter(conf.Formatter),
			Output:    conf.Output,
			SkipPaths: conf.SkipPaths,
		})
	}
}

func convertLogFormatter(formatter LogFormatter) gin.LogFormatter {
	return func(params gin.LogFormatterParams) string {
		return formatter(LogFormatterParams{
			Request:      params.Request,
			TimeStamp:    params.TimeStamp,
			StatusCode:   params.StatusCode,
			Latency:      params.Latency,
			ClientIP:     params.ClientIP,
			Method:       params.Method,
			Path:         params.Method,
			ErrorMessage: params.ErrorMessage,
			BodySize:     params.BodySize,
			Keys:         params.Keys,
		})
	}
}
