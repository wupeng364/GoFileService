// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package conf

const (
	// AppconfDir 程序配置文件位置
	AppconfDir = "appconf.dir"
	// HTMLpageStaticDir 静态文件存放位置
	HTMLpageStaticDir = "htmlpage.static"
	// HTTPServerListen http服务监听地址
	HTTPServerListen = "httpserver.listen"
	// FileManagerMounts 文件挂载信息
	FileManagerMounts = "filemanager.mounts"
	// DataBaseSource 数据库位置|链接
	DataBaseSource = "database.dbsource"
	// DataBasType 数据库类型
	DataBasType = "database.dbtype"

	// DefaultDataBaseType 默认数据库类型
	DefaultDataBaseType = "sqllite"
)
