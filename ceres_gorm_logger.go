//  Copyright 2020 Go-Ceres
//  Author https://github.com/go-ceres/go-ceres
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package CeresGorm

import (
	"context"
	"fmt"
	CeresLogger "github.com/go-ceres/ceres-logger"
	log "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

var (
	infoStr      = "%s\n[info] "
	warnStr      = "%s\n[warn] "
	errStr       = "%s\n[error] "
	traceStr     = "%s\n[%.3fms] [rows:%v] %s"
	traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
	traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
)

type logger struct {
	log.Config
	log                                 CeresLogger.Logger
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

// 字符串转等级数值
func ConvertLevel(level string) log.LogLevel {
	switch level {
	case "info", "INFO":
		return 4
	case "warn", "WARN":
		return 3
	case "error", "ERROR":
		return 2
	}
	return 1
}

// 创建日志实例
func newLog(l CeresLogger.Logger, c log.Config) log.Interface {
	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)

	if c.Colorful {
		infoStr = log.Green + "%s\n" + log.Reset + log.Green + "[info] " + log.Reset
		warnStr = log.BlueBold + "%s\n" + log.Reset + log.Magenta + "[warn] " + log.Reset
		errStr = log.Magenta + "%s\n" + log.Reset + log.Red + "[error] " + log.Reset
		traceStr = log.Green + "%s\n" + log.Reset + log.Yellow + "[%.3fms] " + log.BlueBold + "[rows:%v]" + log.Reset + " %s"
		traceWarnStr = log.Green + "%s " + log.Yellow + "%s\n" + log.Reset + log.RedBold + "[%.3fms] " + log.Yellow + "[rows:%v]" + log.Magenta + " %s" + log.Reset
		traceErrStr = log.RedBold + "%s " + log.MagentaBold + "%s\n" + log.Reset + log.Yellow + "[%.3fms] " + log.BlueBold + "[rows:%v]" + log.Reset + " %s"
	}

	return &logger{
		log:          l,
		Config:       c,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

// LogMode log mode
func (l *logger) LogMode(level log.LogLevel) log.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info print info
func (l logger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= log.Info {
		l.log.Infof(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Warn print warn messages
func (l logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= log.Warn {
		l.log.Warnf(l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Error print error messages
func (l logger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= log.Error {
		l.log.Errorf(l.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Trace print sql message
func (l logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel > 0 {
		elapsed := time.Since(begin)
		switch {
		case err != nil && l.LogLevel >= log.Error:
			sql, rows := fc()
			if rows == -1 {
				l.log.Errorf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				l.log.Errorf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= log.Warn:
			sql, rows := fc()
			slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
			if rows == -1 {
				l.log.Warnf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				l.log.Warnf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		case l.LogLevel >= log.Info:
			sql, rows := fc()
			if rows == -1 {
				l.log.Infof(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
			} else {
				l.log.Infof(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		}
	}
}
