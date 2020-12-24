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
	CeresError "github.com/go-ceres/ceres-error"
	"github.com/jinzhu/gorm"
	"sync"
)

// 全局的回调（即）
var defaultCallbackManager = newDefaultCallbackManager()

// 回调函数
type CallbackManager struct {
	rw        *sync.RWMutex
	callbacks map[string]map[string]Callback
}

// newDefaultCallbackManager 创建一个回调管理器
func newDefaultCallbackManager() *CallbackManager {
	manager := &CallbackManager{
		rw:        &sync.RWMutex{},
		callbacks: make(map[string]map[string]Callback),
	}
	manager.callbacks["query"] = make(map[string]Callback)
	manager.callbacks["create"] = make(map[string]Callback)
	manager.callbacks["update"] = make(map[string]Callback)
	manager.callbacks["delete"] = make(map[string]Callback)
	return manager
}

// Register 注册插件
func (m *CallbackManager) Register(hook string, op string, cb Callback) error {
	m.rw.Lock()
	defer m.rw.Unlock()
	if _, ok := m.callbacks[hook]; !ok {
		return CeresError.New("ceres-gorm", "hook does not exist")
	} else {
		m.callbacks[hook][op] = cb
	}
	return nil
}

// 定义回调函数
type Callback func(scope *gorm.Scope)
