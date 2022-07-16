package kid

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"reflect"
)

type Rules map[string]interface{}

var Validator = validator.New()

func Struct(value interface{}) error {
	return Validator.Struct(value)
}

func Verify(ctx *Context, rules Rules) (err error) {
	var val = make(map[string]interface{})
	for key, _ := range rules {
		v, exist := ctx.Get(key)
		if !exist {
			return &validator.InvalidValidationError{Type: reflect.TypeOf(v)}
		}
		val[key] = v
	}
	result := Validator.ValidateMap(val, rules)
	if len(result) > 0 {
		for k, v := range result {
			err = errors.Wrap(err, fmt.Sprintf("%s:%v", k, v))
		}
		err = fmt.Errorf("verify failed: %w", err)
	}
	return
}
