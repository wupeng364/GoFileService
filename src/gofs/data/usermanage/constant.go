// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package usermanage

import (
	"errors"
	"time"
)

const (
	// AdminRole 管理员角色标识
	AdminRole = 1
	// NormalRole 普通用户角色标识
	NormalRole = 0
)

// ErrorConnIsNil 空连接
var ErrorConnIsNil = errors.New("The data source is empty")

// ErrorUserIDIsNil ErrorUserIDIsNil
var ErrorUserIDIsNil = errors.New("The userID is empty")

// ErrorUserNotExist ErrorUserNotExist
var ErrorUserNotExist = errors.New("User does not exist")

// ErrorUserNameIsNil ErrorUserNameIsNil
var ErrorUserNameIsNil = errors.New("The userName is empty")

// UserInfo 用户表存储的结构
type UserInfo struct {
	UserType int
	UserID   string
	UserName string
	UserPWD  string
	CtTime   time.Time
}
