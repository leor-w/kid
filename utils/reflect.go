package utils

import (
	"fmt"
	"reflect"
)

// ConvertPointer 将一个类型转换为指针类型
func ConvertPointer(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		return t.Elem()
	}
	return reflect.PtrTo(t)
}

// IsNilPointer 判断一个值是否是一个空指针
func IsNilPointer(i interface{}) bool {
	iv := reflect.ValueOf(i)
	return iv.IsNil() || iv.Kind() != reflect.Ptr
}

// IsFunc 判断一个值是否是一个函数
func IsFunc(i interface{}) bool {
	iv := reflect.ValueOf(i)
	return !iv.IsNil() && iv.Kind() == reflect.Func
}

// RemoveTypePtr 移除一个类型的指针
func RemoveTypePtr(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// RemoveValuePtr 移除一个值的指针
func RemoveValuePtr(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

// StructToMap 将一个结构体转换为 map[string]interface{}
func StructToMap(in interface{}, omits ...string) (map[string]interface{}, error) {
	var out = make(map[string]interface{})
	v := reflect.ValueOf(in)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("in value must be struct or struct pointer : %T", v)
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := t.Field(i)
		if f.Anonymous {
			fv := v.Field(i)
			if !fv.CanInterface() {
				continue
			}
			return StructToMap(fv.Interface(), omits...)
		}
		if ContainString(omits, CamelToSnake(f.Name)) {
			continue
		}
		fv := v.Field(i)
		if !fv.CanInterface() {
			continue
		}
		out[CamelToSnake(f.Name)] = fv.Interface()
	}
	return out, nil
}
