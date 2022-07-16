package kid

import (
	"github.com/gin-gonic/gin"
	"github.com/leor-w/kid/guard"
	"github.com/leor-w/kid/validation"
	"github.com/spf13/cast"
)

type Context struct {
	gin.Context
	Guard guard.Guard
}

func (ctx *Context) Valid(recipient interface{}) error {
	if err := ctx.ShouldBind(&recipient); err != nil {
		return err
	}
	if err := validation.Struct(recipient); err != nil {
		return err
	}
	return nil
}

func (ctx *Context) ValidFiled(rules validation.Rules) error {
	return validation.Verify(ctx, rules)
}

func (ctx *Context) GetInt(key string) int {
	return cast.ToInt(ctx.get(key))
}

func (ctx *Context) GetDefaultInt(key string, defaultValue int) int {
	val := ctx.GetInt(key)
	if val > 0 {
		return val
	}
	return defaultValue
}

func (ctx *Context) GetInt64(key string) int64 {
	return cast.ToInt64(ctx.get(key))
}

func (ctx *Context) GetDefaultInt64(key string, defaultValue int64) int64 {
	val := ctx.GetInt64(key)
	if val > 0 {
		return val
	}
	return defaultValue
}

func (ctx *Context) GetString(key string) string {
	return ctx.get(key)
}

func (ctx *Context) GetDefaultString(key, defaultValue string) string {
	val := ctx.get(key)
	if len(val) > 0 {
		return val
	}
	return defaultValue
}

func (ctx *Context) get(key string) string {
	var val string
	val = ctx.Param(key)
	if len(val) > 0 {
		return val
	}
	val = ctx.Query(key)
	if len(val) > 0 {
		return val
	}
	val = ctx.GetString(key)
	if len(val) > 0 {
		return val
	}
	val = ctx.GetHeader(key)
	if len(val) > 0 {
		return val
	}
	return ""
}
