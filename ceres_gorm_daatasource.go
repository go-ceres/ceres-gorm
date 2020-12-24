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
	"fmt"
	"net/url"
	"strings"
)

// 数据库连接结构体
type DataSource struct {
	Drive    string     // 驱动
	Conn     string     // 连接字符串
	User     string     // 用户名
	Password string     // 密码
	Addr     string     // 连接地址
	DbName   string     // 数据库名
	Params   url.Values // 连接参数
}

// 解析数据库连接
func (c *Config) ParseUrl() error {
	// 初始化
	c.dataSource = new(DataSource)
	// 使用url解析包
	data, err := url.Parse(c.Url)
	if err != nil {
		return fmt.Errorf("URL parsing to data source connection failed：%s", err.Error())
	}
	if data == nil {
		return fmt.Errorf("URL parsing to data source connection failed：%s", "parse "+c.Url+", result is nil")
	}
	if data.Scheme == "" {
		return fmt.Errorf("URL parsing to data source connection failed：%s", "Driver not resolved")
	}
	// 设置数据
	// 设置驱动
	c.dataSource.Drive = data.Scheme
	// 设置连接地址
	c.dataSource.Addr = data.Host
	// 设置用户名
	c.dataSource.User = data.User.Username()
	// 设置密码
	c.dataSource.Password, _ = data.User.Password()
	// 设置数据库名
	c.dataSource.DbName = strings.Replace(data.Path, "/", "", 1)
	// 设置额外连接参数
	params, e := url.ParseQuery(data.RawQuery)
	if e != nil {
		c.Logger.Panicf("URL parsing to data source connection failed：%s", e)
	}
	c.dataSource.Params = params
	// 设置连接字符串
	c.dataSource.Conn = strings.Replace(c.Url, c.dataSource.Drive+"://", "", 1)
	return nil
}
