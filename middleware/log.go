package middleware

import (
	"fmt"
	"github.com/leor-w/kid"
	"github.com/leor-w/kid/logger"
)

// RequestLog 输出请求日志到文件
func RequestLog() kid.Middleware {
	return kid.LoggerWithFormatter(func(params kid.LogFormatterParams) string {
		log := fmt.Sprintf("%s - [%s] %s %s %s %d %s \"%s\" %s\n",
			params.TimeStamp.Format("2006-01-02 15:04:05.999"),
			params.Method,
			params.ClientIP,
			params.Path,
			params.Request.Proto,
			params.StatusCode,
			params.Latency,
			params.Request.UserAgent(),
			params.ErrorMessage,
		)
		logger.Info(log)
		return log
	})
}
