package logger

var (
	logger Logger
)

func Init(log Logger) {
	logger = log
}

func Default() Logger {
	return logger
}

func UserHook(hook Hook) {
	logger.Hook(hook)
}

func Trace(args ...interface{}) {
	logger.Log(TraceLevel, args...)
}

func Debug(args ...interface{}) {
	logger.Log(DebugLevel, args...)
}

func Info(args ...interface{}) {
	logger.Log(InfoLevel, args...)
}

func Warn(args ...interface{}) {
	logger.Log(WarnLevel, args...)
}

func Error(args ...interface{}) {
	logger.Log(ErrorLevel, args...)
}

func Fatal(args ...interface{}) {
	logger.Log(FatalLevel, args...)
}

func Tracef(fmt string, args ...interface{}) {
	logger.Logf(TraceLevel, fmt, args...)
}

func Debugf(fmt string, args ...interface{}) {
	logger.Logf(DebugLevel, fmt, args...)
}

func Infof(fmt string, args ...interface{}) {
	logger.Logf(InfoLevel, fmt, args...)
}

func Warnf(fmt string, args ...interface{}) {
	logger.Logf(WarnLevel, fmt, args...)
}

func Errorf(fmt string, args ...interface{}) {
	logger.Logf(ErrorLevel, fmt, args...)
}

func Fatalf(fmt string, args ...interface{}) {
	logger.Logf(FatalLevel, fmt, args...)
}
