// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package filepermission

import (
	"errors"
	"time"
)

const (
	// CHILDVISIBLE 子节点可见导致的父节点可见, 父节点其实并没有权限
	CHILDVISIBLE = iota
	// VISIBLE 可见
	VISIBLE
	// READ 只读
	READ
	// WRITE 可写
	WRITE
)

// ErrorConnIsNil 空连接
var ErrorConnIsNil = errors.New("The data source is empty")

// ErrorUserIDIsNil ErrorUserIDIsNil
var ErrorUserIDIsNil = errors.New("The permissionid is empty")

// ErrorPermissionIDIsNil ErrorPermissionIDIsNil
var ErrorPermissionIDIsNil = errors.New("The permissionid is empty")

// ErrorPermissionPathIsNil ErrorPermissionPathIsNil
var ErrorPermissionPathIsNil = errors.New("The permission path is empty")

// ErrorPermissionIsNil ErrorPermissionIsNil
var ErrorPermissionIsNil = errors.New("The permission is empty")

// FPermissionInfo 文件权限表存储的结构
type FPermissionInfo struct {
	PermissionID string    // 权限路径
	Path         string    // 文件路径
	UserID       string    // 用户ID
	Permission   int64     // 权限值
	CtTime       time.Time // 插入时间
}
