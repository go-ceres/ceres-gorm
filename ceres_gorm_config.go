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
	"github.com/go-ceres/ceres-logger"
	"time"
)

// 定义配置信息
type Config struct {
	// 连接字符串
	Url string
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
	// 数据源
	dataSource *DataSource
	// 回调函数
	callback *CallbackManager
}

// NewDefaultConfig 创建一个默认的配置
func NewDefaultConfig() *Config {
	return &Config{
		Url:             "",
		Debug:           false,
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: time.Hour,
		Logger:          CeresLogger.FrameLogger,
		callback:        newDefaultCallbackManager(),
	}
}

// WithLogger 设置日志输出库
func (c *Config) WithLogger(l CeresLogger.Logger) *Config {
	c.Logger = l
	return c
}

// RegisterCallback 注册回调函数
func (c *Config) RegisterCallback(hook string, op string, cb Callback) error {
	return c.callback.Register(hook, op, cb)
}

// Build 构建数据库
func (c *Config) Build() *DB {
	// 解析连接字符串
	if err := c.ParseUrl(); err != nil {
		c.Logger.Panic(err)
	}

	db, err := Open(c)
	if err != nil {
		c.Logger.Panic(err)
	}

	return db
}
