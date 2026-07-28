package middleware

import (
	"net/http"
	"runtime"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/lsjhtang/kid/logger"
)

func Recovery(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 2048)
			_ = runtime.Stack(buf, true)
			logger.Errorf("panic: %v\n堆栈: %s", err, buf[:])
			debug.PrintStack()
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    50000,
				"message": "内部错误",
			})
		}
	}()
	ctx.Next()
}
