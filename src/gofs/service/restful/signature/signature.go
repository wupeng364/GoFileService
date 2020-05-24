// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 请求签名模块, 请求的拦截、验证参数完整性、身份合法性检测、用户session管理

package signature

import (
	"errors"
	"gutils/hstool"
	"gutils/mloader"
	"net/http"
	"sort"
	"strconv"
)

// Signature 签名校验
type Signature struct {
	sign signature
}

// ModuleOpts 模块加载器接口实现, 返回模块信息&配置
func (signature *Signature) ModuleOpts() mloader.Opts {
	return mloader.Opts{
		Name:        "Signature",
		Version:     1.0,
		Description: "Api接口签名模块",
		OnReady: func(mctx *mloader.Loader) {
		},
		OnInit: func() {
			// 这里暂时只实现单机、本地内存版本
			signature.sign = &impByLocalCache{}
			signature.sign.SignatureInitial()
		},
	}
}

// CreateWebSession 添加会话
func (signature *Signature) CreateWebSession(userID string, r *http.Request) (AccessToken, error) {
	return signature.sign.GenerateAccessToken(userID, SingnatureTypeAsWeb)
}

// DestroySignature 注销会话
func (signature *Signature) DestroySignature(accessKey string) error {
	return signature.sign.DestroySignature(accessKey)
}

// DestroySignature4HTTP 注销会话-传入请求
func (signature *Signature) DestroySignature4HTTP(r *http.Request) error {
	// 从请求中获取accessKey, 不能为空
	accessKey := r.Header.Get(RequestHeaderAccessKey)
	if len(accessKey) == 0 {
		if ack, err := r.Cookie("ack"); nil == err {
			accessKey = ack.Value
		}
	}
	if len(accessKey) == 0 {
		return errors.New(strconv.Itoa(http.StatusUnauthorized))
	}
	return signature.DestroySignature(accessKey)
}

// GetUserIDByAccessKey 获取用户ID
func (signature *Signature) GetUserIDByAccessKey(ack string) string {
	return signature.sign.GetUserID(ack)
}

// SetSessionAttr 设置属性
func (signature *Signature) SetSessionAttr(accessKey, key, val string) error {
	return signature.sign.SetSessionAttr(accessKey, key, val)
}

// SetSessionAttr4HTTP 设置属性-传入请求
func (signature *Signature) SetSessionAttr4HTTP(r *http.Request, key, val string) error {
	// 从请求中获取accessKey, 不能为空
	accessKey := r.Header.Get(RequestHeaderAccessKey)
	if len(accessKey) == 0 {
		if ack, err := r.Cookie("ack"); nil == err {
			accessKey = ack.Value
		}
	}
	if len(accessKey) == 0 {
		return errors.New(strconv.Itoa(http.StatusUnauthorized))
	}
	return signature.SetSessionAttr(accessKey, key, val)
}

// GetSessionAttr 读取属性
func (signature *Signature) GetSessionAttr(accessKey, key string) (string, error) {
	return signature.sign.GetSessionAttr(accessKey, key)
}

// GetSessionAttr4HTTP 读取属性-传入请求
func (signature *Signature) GetSessionAttr4HTTP(r *http.Request, key string) (string, error) {
	// 从请求中获取accessKey, 不能为空
	accessKey := r.Header.Get(RequestHeaderAccessKey)
	if len(accessKey) == 0 {
		if ack, err := r.Cookie("ack"); nil == err {
			accessKey = ack.Value
		}
	}
	if len(accessKey) == 0 {
		return "", errors.New(strconv.Itoa(http.StatusUnauthorized))
	}
	return signature.GetSessionAttr(accessKey, key)
}

// RestfulAPIFilter Api签名拦截器
func (signature *Signature) RestfulAPIFilter(w http.ResponseWriter, r *http.Request, next hstool.FilterNext) {
	//fmt.Println("ApiFilter: ", r.RemoteAddr, r.URL.Path)
	// 从请求中获取accessKey, 不能为空
	accessKey := r.Header.Get(RequestHeaderAccessKey)
	// 从请求中获取signature, 不能为空
	sign := r.Header.Get(RequestHeaderSign)
	if len(accessKey) == 0 || len(sign) == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return // 401
	}
	// 填充Form对象
	if nil == r.Form {
		err := r.ParseForm()
		if nil != err {
			// 出现异常, 不继续处理
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	// 构建请求参数
	requestparameter := ""
	if nil != r.Form && len(r.Form) > 0 {
		keys := make([]string, 0) // 去掉参数为空的传值
		for key, val := range r.Form {
			if len(val) > 0 {
				keys = append(keys, key)
			}
		}

		_keysLen := len(keys)
		if _keysLen > 0 {
			sort.Strings(keys)
			for i, val := range keys {
				requestparameter += val + "=" + r.Form[val][0]
				if i < _keysLen-1 {
					requestparameter += "&"
				}
			}

		}
	}
	// 校验参数合法性
	if !signature.sign.VerificationSignature(accessKey, sign, requestparameter) {
		w.WriteHeader(http.StatusUnauthorized)
		return // 401
	}

	next() // 校验参数合法性-通过
}
