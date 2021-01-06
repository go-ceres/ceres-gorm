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
	CeresConfig "github.com/go-ceres/ceres-config"
	"github.com/go-ceres/ceres-logger"
	"gorm.io/gorm"
	log "gorm.io/gorm/logger"
	"time"
)

// 定义配置信息
type Config struct {
	// 驱动
	Drive string
	// 连接字符串
	DNS string
	// 是否开启debug
	Debug bool
	// 最大空闲连接数
	MaxIdleConns int
	// 最大活动连接数
	MaxOpenConns int
	// 连接的最大存活时间
	ConnMaxLifetime time.Duration
	// 日志库
	Logger CeresLogger.Logger
	// 驱动适配器
	Dialector Dialector
	// 日志配置
	LogConfig
	// gorm的配置
	GormConfig
}

// 日志配置
type LogConfig struct {
	// 慢日志阈值
	SlowThreshold time.Duration
	// 是否开启日志颜色区别
	Colorful bool
	// 日志等级
	LogLevel string
}

// gorm的配置
type GormConfig gorm.Config

// NewDefaultConfig 创建一个默认的配置
func newDefaultConfig() *Config {
	return &Config{
		Drive:           "mysql",
		DNS:             "",
		Debug:           false,
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: time.Hour,
		Dialector:       drivers["mysql"](""),
		Logger:          CeresLogger.FrameLogger.With(CeresLogger.FieldMod("ceres-gorm")).AddCallerSkip(-1),
		LogConfig: LogConfig{
			SlowThreshold: time.Second,
			Colorful:      false,
			LogLevel:      "",
		},
	}
}

// newDefaultLogConf 创建一个默认的日志配置
func newDefaultLogConf() log.Config {
	return log.Config{
		SlowThreshold: time.Second, // 慢 SQL 阈值
		LogLevel:      log.Silent,  // Log level
		Colorful:      false,       // 禁用彩色打印
	}
}

// RawConfig 根据完整key解析配置信息
func RawConfig(key string) *Config {
	// 默认的配置
	conf := newDefaultConfig()
	// 解析配置信息
	if err := CeresConfig.Get(key).Scan(conf); err != nil {
		CeresLogger.FrameLogger.Panicd("scan config", CeresLogger.FieldMod("ceres-gorm"), CeresLogger.FieldErr(err))
	}
	return conf
}

// ScanConfig 根据配置名解析配置
func ScanConfig(name string) *Config {
	return RawConfig("ceres.database." + name)
}

// WithLogger 设置日志输出库
func (c *Config) WithLogger(l CeresLogger.Logger) *Config {
	c.Logger = l
	return c
}

// WithDriver 设置驱动实例
func (c *Config) WithDriver(dialect Dialector) *Config {
	c.Dialector = dialect
	return c
}

// Build 构建数据库
func (c *Config) Build() *DB {
	// 创建驱动
	if driver, ok := drivers[c.Drive]; !ok {
		c.Logger.Panicf("%s driver is not set", driver)
	} else {
		c.Dialector = driver(c.DNS)
	}
	// 初始化日志
	c.initLogger()
	// 数据库
	db, err := Open(c.Dialector, c)
	if err != nil {
		c.Logger.Panicd("open gorm", CeresLogger.FieldErr(err), CeresLogger.FieldAny("value", c))
	}
	return db
}

// initLogger 初始化日志
func (c *Config) initLogger() {
	// 默认日志配置
	logConf := newDefaultLogConf()
	// 转换等级
	if c.LogLevel != "" {
		logConf.LogLevel = ConvertLevel(c.LogLevel)
	}
	logConf.Colorful = c.Colorful
	logConf.SlowThreshold = c.SlowThreshold
	dbLog := newLog(c.Logger, logConf)
	if c.Debug {
		dbLog = dbLog.LogMode(log.Info)
	}
	// gorm的配置信息
	c.GormConfig.Logger = dbLog
}
