// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package asynctask

import (
	"errors"
	"gutils/mloader"
	"net/http"
)

// AsyncTask 接口
type AsyncTask interface {
	Name() string                                  // 动作名字
	Init(mctx *mloader.Loader) AsyncTask           // 初始化对象
	Execute(r *http.Request) (string, error)       // 动作执行, 返回一个tooken
	Status(w http.ResponseWriter, r *http.Request) // 查询动作状态, 在内部返回数据
}

// ErrorDiscontinue ErrorDiscontinue
var ErrorDiscontinue = errors.New("Discontine")

// ErrorOprationExpires ErrorOprationExpires
var ErrorOprationExpires = errors.New("Opration expires")

// ErrorOprationFailed ErrorOprationFailed
var ErrorOprationFailed = errors.New("Opration failed")

// ErrorOprationUnknown ErrorOprationUnknown
var ErrorOprationUnknown = errors.New("Opration unknown")

// ErrorFileNotExist ErrorFileNotExist
var ErrorFileNotExist = errors.New("file does not exist")

// ErrorParentFolderNotExist ErrorParentFolderNotExist
var ErrorParentFolderNotExist = errors.New("parent folder does not exist")

// ErrorNewNameIsEmpty ErrorNewNameIsEmpty
var ErrorNewNameIsEmpty = errors.New("New name cannot be empty")

// ErrorPermissionInsufficient 权限不足
var ErrorPermissionInsufficient = errors.New("权限不足")
