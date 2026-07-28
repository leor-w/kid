package kid

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/lsjhtang/kid/utils"
	"github.com/spf13/cast"
)

const UserBindKey = "bind_user"

type Context struct {
	*gin.Context
}

// SetResponseType 指定响应格式类型
func (ctx *Context) SetResponseType(respType ResponseType) {
	ctx.Set(ResponseTypeKey, respType)
}

// GetResponseType 获取响应格式化类型
func (ctx *Context) GetResponseType() ResponseType {
	val, exist := ctx.Get(ResponseTypeKey)
	if !exist {
		return 0
	}
	respType, ok := val.(int)
	if ok {
		return ResponseType(respType)
	}
	return 0
}

// BindUser 绑定用户
func (ctx *Context) BindUser(val interface{}) {
	ctx.Set(UserBindKey, val)
}

// GetBindUser 获取绑定的用户
func (ctx *Context) GetBindUser(recipient interface{}) error {
	val, exist := ctx.Get(UserBindKey)
	if !exist {
		return fmt.Errorf("not found bind user")
	}
	if utils.IsNilPointer(recipient) {
		return fmt.Errorf("recipient cannot is nil pointer")
	}
	if err := copier.Copy(recipient, val); err != nil {
		return err
	}
	return nil
}

func (ctx *Context) GetValue(key string, recipient interface{}) error {
	val, exist := ctx.Get(key)
	if !exist {
		return fmt.Errorf("[%s] value not exist", key)
	}
	if utils.IsNilPointer(recipient) {
		return fmt.Errorf("recipient cannot is nil pointer")
	}
	if err := copier.Copy(recipient, val); err != nil {
		return err
	}
	return nil
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

func (ctx *Context) FindInt8(key string) int8 {
	return cast.ToInt8(ctx.find(key))
}

func (ctx *Context) FindDefaultInt8(key string, defaultVal int8) int8 {
	val := ctx.FindInt8(key)
	if val > 0 {
		return val
	}
	return defaultVal
}

func (ctx *Context) FindUint8(key string) uint8 {
	return cast.ToUint8(ctx.find(key))
}

func (ctx *Context) FindDefaultUint8(key string, defaultVal uint8) uint8 {
	val := ctx.FindUint8(key)
	if val > 0 {
		return val
	}
	return defaultVal
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

func (ctx *Context) FindFloat64(key string) float64 {
	return cast.ToFloat64(ctx.find(key))
}

func (ctx *Context) FindDefaultFloat64(key string, defaultValue float64) float64 {
	val := ctx.FindFloat64(key)
	if val > 0 {
		return val
	}
	return defaultValue
}

func (ctx *Context) FindFloat32(key string) float32 {
	return cast.ToFloat32(ctx.find(key))
}

func (ctx *Context) FindDefaultFloat32(key string, defaultValue float32) float32 {
	val := ctx.FindFloat32(key)
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

func (ctx *Context) FindBool(key string) bool {
	return cast.ToBool(ctx.find(key))
}

func (ctx *Context) FindDefaultBool(key string, defaultValue bool) bool {
	val := ctx.find(key)
	if len(val) > 0 {
		return cast.ToBool(val)
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
	val = ctx.GetString(key)
	if len(val) > 0 {
		return val
	}
	val = ctx.GetHeader(key)
	if len(val) > 0 {
		return val
	}
	val = ctx.PostForm(key)
	if len(val) > 0 {
		return val
	}
	return ""
}
