package kid

import (
	"github.com/gin-gonic/gin"
	"github.com/leor-w/kid/guard"
	"github.com/spf13/cast"
)

type Context struct {
	*gin.Context
	Guard guard.Guard
}

func (ctx *Context) Valid(recipient interface{}) error {
	if err := ctx.ShouldBind(&recipient); err != nil {
		return err
	}
	if err := Struct(recipient); err != nil {
		return err
	}
	return nil
}

func (ctx *Context) ValidFiled(rules Rules) error {
	return Verify(ctx, rules)
}

func (ctx *Context) FindInt(key string) int {
	return cast.ToInt(ctx.find(key))
}

func (ctx *Context) FindDefaultInt(key string, defaultValue int) int {
	val := ctx.FindInt(key)
	if val > 0 {
		return val
	}
	return defaultValue
}

func (ctx *Context) FindInt64(key string) int64 {
	return cast.ToInt64(ctx.find(key))
}

func (ctx *Context) FindDefaultInt64(key string, defaultValue int64) int64 {
	val := ctx.FindInt64(key)
	if val > 0 {
		return val
	}
	return defaultValue
}

func (ctx *Context) FindString(key string) string {
	return ctx.find(key)
}

func (ctx *Context) FindDefaultString(key, defaultValue string) string {
	val := ctx.find(key)
	if len(val) > 0 {
		return val
	}
	return defaultValue
}

func (ctx *Context) find(key string) string {
	var val string
	val = ctx.Param(key)
	if len(val) > 0 {
		return val
	}
	val = ctx.Query(key)
	if len(val) > 0 {
		return val
	}
	val = ctx.FindString(key)
	if len(val) > 0 {
		return val
	}
	val = ctx.GetHeader(key)
	if len(val) > 0 {
		return val
	}
	return ""
}
