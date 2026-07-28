package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lsjhtang/kid/logger"
)

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()
		latency := time.Now().Sub(start)
		log := fmt.Sprintf("[ROUTE] %s %s %s %s %d %s \"%s\" %s",
			ctx.Request.Method,
			ctx.ClientIP(),
			ctx.Request.URL.Path,
			ctx.Request.Proto,
			ctx.Writer.Status(),
			latency,
			ctx.Request.UserAgent(),
			ctx.Errors.ByType(gin.ErrorTypePrivate).String(),
		)
		if ctx.Writer.Status() == 200 {
			logger.Info(log)
		} else {
			logger.Error(log)
		}
	}
}
