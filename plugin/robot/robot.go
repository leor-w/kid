package robot

import "fmt"

var (
	ErrWrongParamType = fmt.Errorf("错误的请求参数类型")
)

type Robot interface {
	SendMessage(params interface{}) error
	WithdrawMessage(params interface{}) error
}
