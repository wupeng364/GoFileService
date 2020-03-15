// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// sqlite模块, 提供系统数据库库管理

package sqlite

import (
	"database/sql"
	"fmt"
	_ "go-sqlite3"
	"gutils/mloader"
	"path/filepath"
	"strings"
)

// SqliteConn 系统数据库
// 配置参数(mloader.GetParam): app.name
type SqliteConn struct {
	dbSource string
	mctx     *mloader.Loader
}

// ModuleOpts 模块加载器接口实现, 返回模块信息&配置
func (sqliteconn *SqliteConn) ModuleOpts() mloader.Opts {
	return mloader.Opts{
		Name:        "SqliteConn",
		Version:     1.0,
		Description: "Sqlite模块",
		OnReady: func(mctx *mloader.Loader) {
			sqliteconn.mctx = mctx
		},
		OnInit: func() {
			path, err := filepath.Abs(strings.Replace("./conf/{name}.db", "{name}", sqliteconn.mctx.GetParam("app.name").ToString("app"), -1))
			if nil != err {
				panic(err)
			}
			sqliteconn.dbSource = path

			fmt.Println("   > SqliteConn dbSource=" + sqliteconn.dbSource)
		},
	}
}

// Open 打开一个数据库连接
func (sqliteconn *SqliteConn) Open() (*sql.DB, error) {
	return sql.Open("sqlite3", sqliteconn.dbSource)
}
