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
	"gorm.io/gorm"
)

type (
	DB         = gorm.DB
	Dialector  = gorm.Dialector
	DriverFunc = func(dns string) Dialector
)

// Open 打开链接
func Open(dialect gorm.Dialector, c *Config) (*DB, error) {

	inner, err := gorm.Open(dialect, (*gorm.Config)(&c.GormConfig))
	if err != nil {
		return nil, err
	}
	// 设置通用配置
	sqlDb, err := inner.DB()
	if err != nil {
		return nil, err
	}
	// 设置连接数
	sqlDb.SetMaxIdleConns(c.MaxIdleConns)
	sqlDb.SetMaxOpenConns(c.MaxOpenConns)
	// 设置连接存活时长
	if c.ConnMaxLifetime != 0 {
		sqlDb.SetConnMaxLifetime(c.ConnMaxLifetime)
	}
	// 测试是否能连通
	if err := sqlDb.Ping(); err != nil {
		return nil, err
	}
	return inner, nil
}
