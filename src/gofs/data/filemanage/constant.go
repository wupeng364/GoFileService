// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package filemanage

const (
	tokenExpiredSecond = 60 // 复制操作令牌失效时间
	// MountsConfKey 配置文件 - 目录挂载
	MountsConfKey = "filemanager.mounts"
)

// CopyCallback 复制回调
type CopyCallback func(src, dst string, err *CopyError) error

// MoveCallback 移动回调
type MoveCallback func(src, dst string, err *MoveError) error

// CopyError 复制错误结构
type CopyError struct {
	SrcIsExist  bool
	DstIsExist  bool
	ErrorString string
}

// MoveError 移动错误结构
type MoveError CopyError

// FsInfo 文件|夹基础属性
type FsInfo struct {
	Path     string
	CtTime   int64
	IsFile   bool
	FileSize int64
}
