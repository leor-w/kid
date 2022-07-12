package kid

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type Context struct {
	*gin.Context
}

func (ctx *Context) GetInt(key string) int {
	return cast.ToInt(ctx.GetString(key))
}

func (ctx *Context) GetDefaultInt(key string, defaultValue int) int {
	val := ctx.GetInt(key)
	if val > 0 {
		return val
	}
	return defaultValue
}

func (ctx *Context) GetInt64(key string) int64 {
	return cast.ToInt64(ctx.GetString(key))
}

func (ctx *Context) GetDefaultInt64(key string, defaultValue int64) int64 {
	val := ctx.GetInt64(key)
	if val > 0 {
		return val
	}
	return defaultValue
}

func (ctx *Context) GetDefaultString(key, defaultValue string) string {
	val := ctx.GetString(key)
	if len(val) > 0 {
		return val
	}
	return defaultValue
}

func (ctx *Context) GetString(key string) string {
	var val string
	val = ctx.Param(key)
	if len(val) > 0 {
		return val
	}
	val = ctx.Query(key)
	if len(val) > 0 {
		return val
	}
	return ""
}
