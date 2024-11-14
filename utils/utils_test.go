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

type V struct {
	Name
}

func TestStructToMap(t *testing.T) {
	n := &V{
		Name{
			AB: "leor",
			B:  8,
			C:  3,
			D:  true,
			e:  "1112",
		},
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

func TestReplaceFileDir(t *testing.T) {
	fmt.Println(ReplaceFileDir("material/allinpbr.com_6820521013_8192_Blender CycleEevee_Alpha.png", "allinpbr.com_6820521013_8192_Blender CycleEevee_Alpha"))
}

func TestFileName(t *testing.T) {
	fmt.Println(FileName("material/allinpbr.com_6820521013_8192_Blender CycleEevee_Alpha.png"))
}

func TestCheckPhone(t *testing.T) {
	phone := "16637196891"
	fmt.Println(RegexpMatchPhone(phone))
}

func TestGenSMSCode(t *testing.T) {
	fmt.Println(RandomSMSCode(6))
}

func TestRecursiveURLDecode(t *testing.T) {
	value := "https://baidu.com/?baidu=baidu"
	//encodeVal := url.QueryEscape(url.QueryEscape(value))
	encodeVal := value
	t.Logf("Url qurey escape: %s", encodeVal)
	decodeVal, err := RecursiveURLDecode(encodeVal)
	if err != nil {
		t.Errorf("RecursiveURLDecode error: %v", err)
	}
	t.Logf("RecursiveURLDecode: %s", decodeVal)
}
