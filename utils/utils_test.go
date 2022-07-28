package utils

import (
	"fmt"
	"testing"
)

type Name struct {
	A string
	B int
	C int64
	D bool
	e string
}

func TestStructToMap(t *testing.T) {
	n := &Name{
		A: "leor",
		B: 8,
		C: 3,
		D: true,
		e: "1112",
	}
	m, err := StructToMap(n)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(m)
}
