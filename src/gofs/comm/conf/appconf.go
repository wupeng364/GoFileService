// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 程序配置模块

package conf

import (
	"gutils/conftool"
	"gutils/mloader"
	"path/filepath"
)

// AppConf 程序配置模块
// 配置参数(mloader.GetParam):
type AppConf struct {
	conftool.JSONCFG
	mctx *mloader.Loader
}

// ModuleOpts 模块加载器接口实现, 返回模块信息&配置
func (appconf *AppConf) ModuleOpts() mloader.Opts {
	return mloader.Opts{
		Name:        "AppConf",
		Version:     1.0,
		Description: "App配置模块",
		OnReady: func(mctx *mloader.Loader) {
			appconf.mctx = mctx
		},
		OnInit: func() {
			confpath := appconf.mctx.GetParam(AppconfDir).ToString("./conf/" + appconf.mctx.GetParam("app.name").ToString("app") + ".json")
			if !filepath.IsAbs(confpath) {
				confpath, _ = filepath.Abs(confpath)
			}
			err := appconf.InitConfig(confpath)
			if nil != err {
				panic(err)
			}
			// 读取配置信息
			appconf.mctx.SetParam("htmlpage.static", appconf.GetConfig(HTMLpageStaticDir).ToString("./static"))
			appconf.mctx.SetParam("httpserver.listen", appconf.GetConfig(HTTPServerListen).ToString("0.0.0.0:8080"))
		},
	}
}
