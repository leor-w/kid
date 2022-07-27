package ginmiddleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/leor-w/kid/logger"
	"net/http"
)

func Logger() func(ctx *gin.Context) {
	return gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		log := fmt.Sprintf("[%s] [ROUTE] [%s] %s %s %s %d %s \"%s\" %s",
			params.TimeStamp.Format("2006-01-02 15:04:05:000"),
			params.Method,
			params.ClientIP,
			params.Path,
			params.Request.Proto,
			params.StatusCode,
			params.Latency,
			params.Request.UserAgent(),
			params.ErrorMessage,
		)
		if params.StatusCode == http.StatusOK {
			logger.Info(log)
		} else {
			logger.Error(log)
		}
		return log + "\n"
	})
}
