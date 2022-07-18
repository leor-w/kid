package kid

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
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
		v := ctx.find(key)
		if len(v) <= 0 {
			return &validator.InvalidValidationError{Type: reflect.TypeOf(v)}
		}
		val[key] = v
	}
	result := Validator.ValidateMap(val, rules)
	if len(result) > 0 {
		for k, v := range result {
			if err != nil {
				err = fmt.Errorf("%w : %s: %v", err, k, v)
				continue
			}
			err = errors.New(fmt.Sprintf("%s:%v", k, v))
		}
		err = fmt.Errorf("verify failed: %w", err)
	}
	return
}
