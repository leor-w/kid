package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Cros 跨域
func Cros(conf cors.Config) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, AccessToken, X-CSRF-Tokens, Authorization, Tokens")
		ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		ctx.Header("Access-Control-Allow-Credentials", "true")
		// 允许放行OPTIONS请求
		if method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
		}
		ctx.Next()
	}
}

func CrosWithConfig(conf cors.Config) gin.HandlerFunc {
	return cors.New(conf)
}

func DefualtCorsConfig() cors.Config {
	return cors.Config{
		AllowAllOrigins: true,
		AllowOrigins:    []string{"*"},
		AllowMethods:    []string{http.MethodPost, http.MethodGet, http.MethodDelete, http.MethodPut, http.MethodHead, http.MethodOptions, http.MethodPatch},
		AllowHeaders:    []string{"Authorization", "Content-Length", "X-CSRF-Tokens", "Tokens", "sign", "X-Custom-Header"},
		MaxAge:          12 * time.Hour,
	}
}
