package utils

import (
	"fmt"
	"testing"
)

type Name struct {
	AB string
	B  int
	C  int64
	D  bool
	e  string
}

func TestStructToMap(t *testing.T) {
	n := &Name{
		AB: "leor",
		B:  8,
		C:  3,
		D:  true,
		e:  "1112",
	}
	m, err := StructToMap(n)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(m)
}

func TestCamelToSnake(t *testing.T) {
	v := "CamelSnake"
	fmt.Println(CamelToSnake(v))
}

func TestSnakeToCamel(t *testing.T) {
	v := "camel_snake"
	fmt.Println(SnakeToCamel(v))
}
