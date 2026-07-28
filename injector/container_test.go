package injector

import (
	"context"
	"fmt"
	"testing"
)

type A struct {
	C1 *C `inject:"scope:scope1"`
	C2 *C `inject:"alias:zhangsan,scope:scope1"`
	C3 *C `inject:"alias:wuwei,scope:scope1"`
	B  *B `inject:"alias:lisi,opts:NR,scope:scope1"`
}

func (a *A) Provide(context.Context) interface{} {
	return a
}

type B struct {
	Name string
	Tab  string
}

func (b *B) Provide(context.Context) interface{} {
	return b
}

type C struct {
	Name string
	Age  int
}

func (c *C) Provide(context.Context) interface{} {
	return c
}

func TestInject(t *testing.T) {
	container := New()
	scope1 := container.scope.Scope("scope1")
	if err := scope1.Provide(&C{
		Name: "lisi",
		Age:  22,
	}); err != nil {
		t.Fatal(err.Error())
		return
	}
	if err := scope1.Provide(&C{
		Name: "zhangsan",
		Age:  28,
	}, WithAlias("zhangsan")); err != nil {
		t.Fatal(err.Error())
		return
	}
	if err := scope1.Provide(&C{
		Name: "wuwei",
		Age:  30,
	}, WithAlias("wuwei")); err != nil {
		t.Fatal(err.Error())
		return
	}
	if err := scope1.Provide(&B{
		Name: "b",
		Tab:  "b",
	}, WithAlias("lisi")); err != nil {
		t.Fatal(err.Error())
		return
	}
	scope2 := container.Scope("scope2")
	a := &A{}
	if err := scope2.Provide(a); err != nil {
		t.Fatal(err.Error())
		return
	}
	if err := container.Populate(); err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(a.C1, a.C2, a.C3, a.B)
}

func TestInvoke(t *testing.T) {
	container := New()
	scope1 := container.scope.Scope("scope1")
	if err := scope1.Provide(&C{
		Name: "lisi",
		Age:  22,
	}); err != nil {
		t.Fatal(err.Error())
		return
	}
	if err := scope1.Provide(&C{
		Name: "zhangsan",
		Age:  28,
	}, WithAlias("zhangsan")); err != nil {
		t.Fatal(err.Error())
		return
	}
	if err := scope1.Provide(&C{
		Name: "wuwei",
		Age:  30,
	}, WithAlias("wuwei")); err != nil {
		t.Fatal(err.Error())
		return
	}
	if err := scope1.Provide(&B{
		Name: "b",
		Tab:  "b",
	}, WithAlias("lisi")); err != nil {
		t.Fatal(err.Error())
		return
	}
	scope2 := container.Scope("scope2")
	a := &A{}
	if err := scope2.Provide(a); err != nil {
		t.Fatal(err.Error())
		return
	}
	if err := container.Populate(); err != nil {
		t.Fatal(err.Error())
	}
	if _, err := container.Invoke(func(a *A) {
		fmt.Println(a.C1, a.C2, a.C3, a.B)
	}, WithInvokeInfo(
		NewInvokeInfo(new(A), "scope2", "", ""),
	)); err != nil {
		t.Fatal(err.Error())
	}
}
