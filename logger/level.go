package logger

type Level int8

const (
	// PanicLevel 该级别的日志会记录并调用 log.Panic, 为最高级别的错误
	PanicLevel Level = iota
	// FatalLevel 该级别的日志会记录并调用 log.Exit(1), 为最高级别的错误
	FatalLevel
	// ErrorLevel 错误级别 通常是非常值得关注的错误信息
	ErrorLevel
	// WarnLevel 警告级别 通常输出值得关注的条目
	WarnLevel
	// InfoLevel 默认的日志级别, 通常输出程序正常执行期间的日志
	InfoLevel
	// DebugLevel debug 调试信息, 内容相对更详细
	DebugLevel
	// TraceLevel 指定比 debug 更小粒度的日志信息
	TraceLevel
)

func (l Level) String() string {
	switch l {
	case TraceLevel:
		return "trace"
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	}
	return ""
}

// Enabled 如果给定的 Level 大于当前 Level 则返回 true
func (l Level) Enabled(lvl Level) bool {
	return lvl >= l
}
