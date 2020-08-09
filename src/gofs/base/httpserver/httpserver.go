// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// HTTP服务模块, 服务路由注册、请求拦截、端口侦听

package httpserver

import (
	"errors"
	"fmt"
	"gutils/hstool"
	"gutils/mloader"
	"gutils/reflecttool"
	"gutils/strtool"
	"net/http"
	"strings"
)

// HTTPServer HTTP服务器
// 配置参数(mloader.GetParam): DEBUG、httpserver.listen
type HTTPServer struct {
	serviceRouter *hstool.ServiceRouter
	mctx          *mloader.Loader
	listenAddr    string
}

// ModuleOpts 模块加载器接口实现, 返回模块信息&配置
func (httpserver *HTTPServer) ModuleOpts() mloader.Opts {
	return mloader.Opts{
		Name:        "HTTPServer",
		Version:     1.0,
		Description: "HTTP服务模块",
		OnReady: func(mctx *mloader.Loader) {
			httpserver.mctx = mctx
			httpserver.listenAddr = mctx.GetParam("httpserver.listen").ToString("127.0.0.1:8080")
			httpserver.serviceRouter = &hstool.ServiceRouter{}
			httpserver.serviceRouter.ClearHandlersMap()
			httpserver.serviceRouter.SetDebug(mctx.GetParam("DEBUG").ToBool(false))
		},
	}
}

// AddIgnoreFilter 注册忽略签名验证的路径
func (httpserver *HTTPServer) AddIgnoreFilter(url string) {
	httpserver.serviceRouter.AddURLFilter(url, func(w http.ResponseWriter, r *http.Request, next hstool.FilterNext) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if "OPTIONS" == strings.ToUpper(r.Method) {
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.Header().Set("Access-Control-Allow-Methods", "*")
			w.WriteHeader(http.StatusOK)
			return
		}
		next()
	})
}

// AddGlobalFilter 默认全局过滤器, 可用于校验签名和session
func (httpserver *HTTPServer) AddGlobalFilter(globalFilter hstool.FilterFunc) error {
	if globalFilter == nil {
		return errors.New("globalFilter is nil")
	}
	httpserver.serviceRouter.SetGlobalFilter(globalFilter)
	return nil
}

// AddURLFilter URL过滤器, 匹配方式: startWith url
func (httpserver *HTTPServer) AddURLFilter(url string, urlFilter hstool.FilterFunc) error {
	if len(url) == 0 {
		return errors.New("Filter url is nil")
	}
	if urlFilter == nil {
		return errors.New("globalFilter func is nil")
	}
	httpserver.serviceRouter.AddURLFilter(url, urlFilter)
	return nil
}

// DefaultHandler 注册默认处理器, 在无匹配时调用
func (httpserver *HTTPServer) DefaultHandler(df hstool.HandlersFunc) (err error) {
	if df == nil {
		return errors.New("http.Handler is nil")
	}
	httpserver.serviceRouter.SetDefaultHandler(df)
	return nil
}

// AddHandlerFunc 通过 HandlersFunc 注册
func (httpserver *HTTPServer) AddHandlerFunc(path string, hf hstool.HandlersFunc) (err error) {
	if len(path) == 0 {
		return errors.New("Path is nil")
	}
	if hf == nil {
		return errors.New("HandlersFunc is nil")
	}
	httpserver.serviceRouter.AddHandler(path, hf)
	return nil
}

// AddRegistrar 通过服务模板注册路由服务
func (httpserver *HTTPServer) AddRegistrar(rs Registrar) error {
	handlers, err := BuildHandlersMap(rs)
	if err != nil {
		return err
	}
	httpserver.serviceRouter.AddHandlers(handlers)
	return nil
}

// DoStartServer 启动服务
func (httpserver *HTTPServer) DoStartServer(server *http.Server) error {
	fmt.Println("   >Server listened in :" + httpserver.listenAddr)
	if len(httpserver.listenAddr) == 0 {
		return errors.New("addr is nil or empty")
	}
	s := server
	if nil == s {
		s = &http.Server{
			ReadTimeout:    0,
			WriteTimeout:   0,
			MaxHeaderBytes: 1 << 20,
		}
	}
	s.Addr = httpserver.listenAddr
	s.Handler = httpserver.serviceRouter

	return s.ListenAndServe()
}

// BuildHandlersMap 根据类型来构造请求路由
func BuildHandlersMap(rs Registrar) (map[string]hstool.HandlersFunc, error) {
	if rs == nil {
		return nil, errors.New("Not find Registrar")
	}
	srg := rs.RoutList()
	hf := srg.HandlerFunc
	if len(hf) == 0 {
		return nil, errors.New("Not find method of Registrar")
	}
	baseURL := strtool.Parse2UnixPath(srg.BasePath)
	if baseURL != "/" {
		baseURL += "/"
	}
	handlersMap := make(map[string]hstool.HandlersFunc, len(hf))
	for i := 0; i < len(hf); i++ {
		_fn := reflecttool.GetFunctionName(hf[i], '.')
		_fm := strings.Index(_fn, "-")
		if _fm > -1 {
			_fn = _fn[:_fm]
		}
		_fn = baseURL + _fn
		if srg.IsToLower {
			_fn = strings.ToLower(_fn)
		}
		handlersMap[_fn] = hf[i]
	}
	return handlersMap, nil
}
