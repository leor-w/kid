package kid

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type RouterGroup struct {
	group *gin.RouterGroup
}

type HandleFunc func(*Context) interface{}

type Middleware func(*Context)

func (group *RouterGroup) POST(path string, handler HandleFunc, middlewares ...Middleware) {

	group.group.POST(path, convert(handler, middlewares...)...)
}

func (group *RouterGroup) GET(path string, handler HandleFunc, middlewares ...Middleware) {
	group.group.GET(path, convert(handler, middlewares...)...)
}

func (group *RouterGroup) PUT(path string, handler HandleFunc, middlewares ...Middleware) {
	group.group.PUT(path, convert(handler, middlewares...)...)
}

func (group *RouterGroup) DELETE(path string, handler HandleFunc, middlewares ...Middleware) {
	group.group.DELETE(path, convert(handler, middlewares...)...)
}

func (group *RouterGroup) Group(path string, middlewares ...Middleware) *RouterGroup {
	return &RouterGroup{group.group.Group(path, convertMiddleware(middlewares...)...)}
}

func (group *RouterGroup) UseMiddle(middlewares ...Middleware) {
	group.group.Use(convertMiddleware(middlewares...)...)
}

func convert(handler HandleFunc, middlewares ...Middleware) []gin.HandlerFunc {
	h := convertHandleFunc(handler)
	ms := convertMiddleware(middlewares...)
	return append(ms, h)
}

func convertHandleFunc(handler HandleFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resp := handler(&Context{Context: ctx})
		ctx.JSON(http.StatusOK, resp)
	}
}

func convertMiddleware(middlewares ...Middleware) []gin.HandlerFunc {
	var ginMiddlewares = make([]gin.HandlerFunc, len(middlewares))
	for i := range middlewares {
		handler := middlewares[i]
		ginMiddlewares[i] = func(ctx *gin.Context) {
			handler(&Context{Context: ctx})
		}
	}
	return ginMiddlewares
}
