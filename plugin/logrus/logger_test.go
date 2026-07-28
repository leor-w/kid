package logrus

import (
	"github.com/lsjhtang/kid/logger"
	"github.com/lsjhtang/kid/plugin/logrus/rotates/simple"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	logger.Init(NewLogger())
	// 日志文件分割 hook
	rotateHook, _ := simple.NewSimpleRotate()
	logger.UserHook(rotateHook)
	// 日志上传 elasticsearch hook

	logger.Error("error log test")
	time.Sleep(1 * time.Second)
}
