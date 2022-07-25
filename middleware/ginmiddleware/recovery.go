package ginmiddleware

import (
	"github.com/gin-gonic/gin"
	"github.com/leor-w/kid/logger"
	"net/http"
	"runtime"
	"runtime/debug"
)

func Recovery(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 2048)
			_ = runtime.Stack(buf, true)
			logger.Errorf("panic: %v \n panic 堆栈信息: %s", err, buf[:])
			debug.PrintStack()
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    50000,
				"message": "内部错误",
			})
		}
	}()
	c.Next()
}
