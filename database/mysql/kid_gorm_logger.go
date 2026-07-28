package mysql

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lsjhtang/kid/logger"
	logger2 "github.com/lsjhtang/kid/logger"
	gLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type KidLogger struct {
	logger.Logger
	*gLogger.Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

func NewKidLogger(config *gLogger.Config) gLogger.Interface {
	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)

	if config.Colorful {
		infoStr = gLogger.Green + "%s\n" + gLogger.Reset + gLogger.Green + "[info] " + gLogger.Reset
		warnStr = gLogger.BlueBold + "%s\n" + gLogger.Reset + gLogger.Magenta + "[warn] " + gLogger.Reset
		errStr = gLogger.Magenta + "%s\n" + gLogger.Reset + gLogger.Red + "[error] " + gLogger.Reset
		traceStr = gLogger.Green + "%s\n" + gLogger.Reset + gLogger.Yellow + "[%.3fms] " + gLogger.BlueBold + "[rows:%v]" + gLogger.Reset + " %s"
		traceWarnStr = gLogger.Green + "%s " + gLogger.Yellow + "%s\n" + gLogger.Reset + gLogger.RedBold + "[%.3fms] " + gLogger.Yellow + "[rows:%v]" + gLogger.Magenta + " %s" + gLogger.Reset
		traceErrStr = gLogger.RedBold + "%s " + gLogger.MagentaBold + "%s\n" + gLogger.Reset + gLogger.Yellow + "[%.3fms] " + gLogger.BlueBold + "[rows:%v]" + gLogger.Reset + " %s"
	}

	return &KidLogger{
		Logger:       logger.Default(),
		Config:       config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

func (log *KidLogger) LogMode(level gLogger.LogLevel) gLogger.Interface {
	newLogger := *log
	newLogger.LogLevel = level
	return &newLogger
}

func (log *KidLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if log.LogLevel >= gLogger.Info {
		log.Logf(logger2.InfoLevel, log.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

func (log *KidLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if log.LogLevel >= gLogger.Warn {
		log.Logf(logger2.WarnLevel, log.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

func (log *KidLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if log.LogLevel >= gLogger.Error {
		log.Logf(logger2.ErrorLevel, log.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

func (log *KidLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if log.LogLevel <= gLogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && log.LogLevel >= gLogger.Error && (!errors.Is(err, gLogger.ErrRecordNotFound) || !log.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			log.Logf(logger2.ErrorLevel, log.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			log.Logf(logger2.ErrorLevel, log.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > log.SlowThreshold && log.SlowThreshold != 0 && log.LogLevel >= gLogger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", log.SlowThreshold)
		if rows == -1 {
			log.Logf(logger2.WarnLevel, log.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			log.Logf(logger2.WarnLevel, log.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case log.LogLevel == gLogger.Info:
		sql, rows := fc()
		if rows == -1 {
			log.Logf(logger2.InfoLevel, log.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			log.Logf(logger2.InfoLevel, log.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
