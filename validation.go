package kid

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

type Rules map[string]interface{}

var Validator = NewValidator()

func NewValidator() *validator.Validate {
	zh := zh.New()
	uni := ut.New(zh)
	trans, _ := uni.GetTranslator("zh")
	validate := validator.New()
	if err := zh_translations.RegisterDefaultTranslations(validate, trans); err != nil {
		panic(fmt.Errorf("注册参数验证翻译器失败: %s", err.Error()))
	}
	return validate
}

func Struct(value interface{}) error {
	return Validator.Struct(value)
}

func Verify(ctx *Context, rules Rules) error {
	var val = make(map[string]interface{})
	for key, _ := range rules {
		v := ctx.find(key)
		if len(v) <= 0 {
			return &validator.InvalidValidationError{Type: reflect.TypeOf(v)}
		}
		val[key] = v
	}
	result := Validator.ValidateMap(val, rules)
	var errBuild strings.Builder
	if len(result) > 0 {
		for k, v := range result {
			validErr, ok := v.(validator.ValidationErrors)
			if !ok {
				return errors.New("验证参数错误")
			}
			for _, e := range validErr {
				errBuild.WriteString(fmt.Sprintf("参数: [%s] 错误,参数要求为: [%s];", k, e.Tag()))
			}
		}
		return errors.New(errBuild.String())
	}
	return nil
}
