package container

import (
	"context"
	"fmt"
	"testing"
)

type A struct {
	C1 *C `inject:""`
	C2 *C `inject:"alias:zhangsan"`
	C3 *C `inject:"alias:wuwei"`
	B  *B `inject:"alias:lisi,opts:NR"`
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
	container.Provide(&C{
		Name: "lisi",
		Age:  22,
	})
	container.Provide(&C{
		Name: "zhangsan",
		Age:  28,
	}, WithAlias("zhangsan"))
	container.Provide(&C{
		Name: "wuwei",
		Age:  30,
	}, WithAlias("wuwei"))
	container.Provide(&B{
		Name: "b",
		Tab:  "b",
	})
	a := &A{}
	//container.Provide(a)
	//if err := container.Populate(); err != nil {
	//	t.Fatalf(err.Error())
	//}
	if _, err := container.PopulateSingle(a); err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(a.C1, a.C2, a.C3, a.B)
	//fmt.Println(a.C2, a.B)
}
