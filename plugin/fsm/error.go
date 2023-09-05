package fsm

type fsmError struct {
	currentState State
	event        Event
}

func (e fsmError) Error() string {
	return "状态机错误: 无法找到事件为 [" + string(e.event) + "] 且状态为 [" + string(e.currentState) + "] 的转换器"
}
