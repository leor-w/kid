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

func TestPhoneNumberValidate(t *testing.T) {
	tests := []struct {
		phoneNumber       string
		expectedIsValid   bool
		expectedFormatted string
	}{
		{
			phoneNumber:       "+1 (650) 643-8195", // 美国的有效号码
			expectedIsValid:   true,
			expectedFormatted: "+16506438195",
		},
		{
			phoneNumber:       "+16506438195", // 美国的有效号码
			expectedIsValid:   true,
			expectedFormatted: "+16506438195",
		},
		{
			phoneNumber:       "+165064381951", // 美国的无效号码
			expectedIsValid:   false,
			expectedFormatted: "",
		},
		{
			phoneNumber:       "+85244230025", // 香港的有效号码
			expectedIsValid:   true,
			expectedFormatted: "+85244230025",
		},
		{
			phoneNumber:       "+8524423002", // 香港的有效号码
			expectedIsValid:   false,
			expectedFormatted: "",
		},
		{
			phoneNumber:       "+8612345678901", // 中国大陆的有效号码
			expectedIsValid:   true,
			expectedFormatted: "+8612345678901",
		},
		{
			phoneNumber:       "+11234567890", // 错误美国电话
			expectedIsValid:   false,
			expectedFormatted: "+11234567890",
		},
	}

	for _, tt := range tests {
		t.Run(tt.phoneNumber, func(t *testing.T) {
			valid, formatted, err := ValidatePhoneNumber(tt.phoneNumber)
			if err != nil {
				t.Errorf("验证失败: %v", err)
			}
			if valid != tt.expectedIsValid {
				t.Errorf("预期有效性为 %v，但得到 %v", tt.expectedIsValid, valid)
			}
			if formatted != tt.expectedFormatted {
				t.Errorf("预期格式化号码为 %v，但得到 %v", tt.expectedFormatted, formatted)
			}
		})
	}
}

// 测试无效的电话号码
func TestValidatePhoneNumber_Invalid(t *testing.T) {
	tests := []struct {
		phoneNumber string
		expectedErr string
	}{
		{
			phoneNumber: "+1 650-555-000", // 假设是无效号码
			expectedErr: "无效电话号码",
		},
		{
			phoneNumber: "12345", // 纯数字，格式错误
			expectedErr: "错误解析电话号码",
		},
		{
			phoneNumber: "+999 000-000-0000", // 不存在的国家代码
			expectedErr: "无效电话号码",
		},
	}

	for _, tt := range tests {
		t.Run(tt.phoneNumber, func(t *testing.T) {
			valid, _, err := ValidatePhoneNumber(tt.phoneNumber)
			if valid {
				t.Errorf("预期无效号码，但结果有效: %v", tt.phoneNumber)
			}
			if err == nil || !contains(err.Error(), tt.expectedErr) {
				t.Errorf("预期错误包含 %v，但得到 %v", tt.expectedErr, err)
			}
		})
	}
}

// 检查错误消息中是否包含期望的文本
func contains(str, substr string) bool {
	return len(str) >= len(substr) && str[:len(substr)] == substr
}
