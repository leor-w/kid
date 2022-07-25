package ginmiddleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func Logger() func(ctx *gin.Context) {
	return gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		log := fmt.Sprintf("[%s] %s %s %s %d %s \"%s\" %s",
			params.Method,
			params.ClientIP,
			params.Path,
			params.Request.Proto,
			params.StatusCode,
			params.Latency,
			params.Request.UserAgent(),
			params.ErrorMessage,
		)
		return log
	})
}
