package formatters

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type SimpleFormatter struct {
}

const timeFormat = "2006-01-02 15:04:05"

func (s *SimpleFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var (
		t      = time.Now()
		format = timeFormat
	)
	msg := fmt.Sprintf("[%s] [%s] %s\n",
		t.Format(format),
		strings.ToUpper(entry.Level.String()),
		entry.Message,
	)
	return []byte(msg), nil
}

func NewSimpleFormatter() *SimpleFormatter {
	formatter := &SimpleFormatter{}
	return formatter
}
