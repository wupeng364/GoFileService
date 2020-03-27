// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 文件基础操作接口

package filemanage

import (
	"io"
)

// Reader 文件流Reader, 必须包含Close操作
type Reader interface {
	io.Reader
	Close() error
}

// fmInterface 文件管理接口
type fmInterface interface {
	IsExist(relativePath string) (bool, error)
	IsDir(relativePath string) (bool, error)
	IsFile(relativePath string) (bool, error)

	GetDirList(relativePath string) ([]string, error)
	GetFileSize(relativePath string) (int64, error)
	GetCreateTime(relativePath string) (int64, error)
	GetModifyTime(relativePath string) (int64, error)
	GetCreateUser(relativePath string) (string, error)
	GetModifyUser(relativePath string) (string, error)
	GetFileLatestVersion(relativePath string) (string, error)
	GetFileVersionList(relativePath string) ([]string, error)

	DoNewFolder(path string) error
	DoRename(src string, dest string) error
	DoMove(src, dest string, repalce, ignore bool, callback MoveCallback) error
	DoDelete(relativePath string) error

	DoRead(relativePath string) (Reader, error)
	DoWrite(relativePath string, ioReader io.Reader) error
	DoCopy(src, dst string, replace, ignore bool, callback CopyCallback) error

	// 传输Token
	AskUploadToken(relativePath string) (string, error)
	SaveTokenFile(token string, src io.Reader) bool
	SubmitToken(token string, isCreateVer bool) (bool, error)
}
