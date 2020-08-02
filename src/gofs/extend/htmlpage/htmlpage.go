// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 静态资源加载器

package htmlpage

import (
	"gofs/comm/httpserver"
	"gofs/service/restful/signature"
	"gutils/hstool"
	"gutils/mloader"
	"strings"

	"fmt"
	"net/http"
	"path/filepath"
)

// HTMLPage 静态资源加载器
// 配置参数(mloader.GetParam): htmlpage.static
type HTMLPage struct {
	mctx      *mloader.Loader
	hsvModule *httpserver.HTTPServer
	sign      *signature.Signature
}

// ModuleOpts 模块加载器接口实现, 返回模块信息&配置
func (html *HTMLPage) ModuleOpts() mloader.Opts {
	return mloader.Opts{
		Name:        "HTMLPage",
		Version:     1.0,
		Description: "静态资源处理",
		OnReady: func(mctx *mloader.Loader) {
			html.mctx = mctx
			html.hsvModule = mctx.GetModuleByTemplate(html.hsvModule).(*httpserver.HTTPServer)
			html.sign = mctx.GetModuleByTemplate(html.sign).(*signature.Signature)
		},
		OnInit: html.onInit,
	}
}

// 每次启动加载模块执行一次
func (html *HTMLPage) onInit() {
	// config
	static := html.mctx.GetParam("htmlpage.static").ToString("./static")
	if !filepath.IsAbs(static) {
		temp, err := filepath.Abs(static)
		if err != nil {
			panic(err)
		}
		static = temp
	}
	// 全局过滤器
	err := html.hsvModule.AddGlobalFilter(html.staticFilter)

	// 首页重定向
	err = html.hsvModule.AddHandlerFunc("/", PageDispatch{}.Index)
	if err != nil {
		panic(err)
	}

	// 页面资源
	err = html.hsvModule.DefaultHandler(http.StripPrefix("/", http.FileServer(http.Dir(static))).ServeHTTP)
	if err != nil {
		panic(err)
	}

	// 全局默认过滤器, 如果没有特殊设定的url都会到这里来
	html.hsvModule.AddGlobalFilter(html.staticFilter)

	fmt.Println("   > Htmlpage Module http registered end")
}

// staticFilter 静态资源过滤器
func (html *HTMLPage) staticFilter(w http.ResponseWriter, r *http.Request, next hstool.FilterNext) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if "OPTIONS" == strings.ToUpper(r.Method) {
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.WriteHeader(http.StatusOK)
		return
	}
	// fmt.Println("Static Path: ", r.URL.Path)
	if strings.HasSuffix(r.URL.Path, ".html") {
		if ack, err := r.Cookie("ack"); nil == err {
			if len(ack.Value) > 0 {
				// 是否只要 cookie 里面有合法 accessKey 就可以了
				userID := html.sign.GetUserIDByAccessKey(ack.Value)
				//fmt.Println("session: ", ack.Value, userID)
				if len(userID) > 0 {
					next()
					return
				}
			}
		}
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
	next()
}

// sessionExpiredHandler 会话过期处理器
func (html *HTMLPage) sessionExpiredHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Session expired: ", r.URL.Path)
}
