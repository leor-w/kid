package logger

type Logger interface {
	Init(...Option) error
	Options() *Options
	// WithFields 添加输出固定字段
	WithFields(map[string]interface{})
	// Hook 添加钩子到日志
	Hook(Hook)
	// Log 输出日志
	Log(Level, ...interface{})
	// Logf 输出带格式化的日志
	Logf(Level, string, ...interface{})
}

type Option func(*Options)
