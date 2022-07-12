package container

import (
	"fmt"
	"testing"
)

type A struct {
	C1 *C `inject:""`
	C2 *C `inject:"zhangsan"`
	C3 *C `inject:"wuwei"`
	B  *B `inject:""`
}

type B struct {
	Name string
	Tab  string
}

type C struct {
	Name string
	Age  int
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
	}, WithName("zhangsan"))
	container.Provide(&C{
		Name: "wuwei",
		Age:  30,
	}, WithName("wuwei"))
	a := &A{}
	container.Provide(a)
	if err := container.Populate(); err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(a.C1, a.C2, a.C3, a.B)
}
