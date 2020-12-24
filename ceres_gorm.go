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
	"github.com/jinzhu/gorm"
)

type (
	DB       = gorm.DB
	CallBack = gorm.Callback
)

// Open 打开数据库
func Open(c *Config) (*DB, error) {
	// 创建DB
	db, err := gorm.Open(c.Drive, c.Url)
	if err != nil {
		c.Logger.Panicd("open "+c.Drive+" error", CeresLogger.FieldPkg("ceres-gorm"), CeresLogger.FieldString("url", c.Url), CeresLogger.FieldErr(err))
	}
	// 日志模式
	db.LogMode(c.Debug)
	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	if c.MaxIdleConns != 0 {
		db.DB().SetMaxIdleConns(c.MaxIdleConns)
	}
	// SetMaxOpenCons 设置数据库的最大连接数量。
	if c.MaxOpenConns != 0 {
		db.DB().SetMaxOpenConns(c.MaxOpenConns)
	}
	// SetConnMaxLifetiment 设置连接的最大可复用时间。
	if c.ConnMaxLifetime != 0 {
		db.DB().SetConnMaxLifetime(c.ConnMaxLifetime)
	}
	return db, nil
}
