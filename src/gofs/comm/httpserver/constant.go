// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package httpserver

import (
	"encoding/json"
	"gutils/hstool"
	"net/http"
)

// StructRegistrar 批量注册器
type StructRegistrar struct {
	IsToLower   bool                  // 是否需要转小写访问
	BasePath    string                // 基础路径 /base/child....
	HandlerFunc []hstool.HandlersFunc // 需要注册的 fuc
}

// Registrar 批量注册器接口
type Registrar interface {
	RoutList() StructRegistrar
}

// APIResponse 接口返回格式约束
type APIResponse struct {
	Code int
	Data string
}

// SendSuccess 返回成功结果
func SendSuccess(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(parse2ApiJSON(http.StatusOK, msg))
}

// SendError 返回失败结果
func SendError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(parse2ApiJSON(http.StatusBadRequest, err.Error()))
}

func parse2ApiJSON(code int, str string) []byte {
	bt, err := json.Marshal(APIResponse{Code: code, Data: str})
	if nil != err {
		return []byte(err.Error())
	}
	return bt
}
