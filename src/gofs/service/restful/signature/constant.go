// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package signature

import (
	"encoding/json"
	"errors"
)

const (
	// SingnatureTypeAsAPI 放置数据, 不删除就不过期
	SingnatureTypeAsAPI SingnatureType = 1
	// SingnatureTypeAsWeb 放置内存, 会过期
	SingnatureTypeAsWeb SingnatureType = 0
	// SingnatureTypeAsWebDestroyTime 默认30分钟销毁
	SingnatureTypeAsWebDestroyTime = 30 * 60
	// SingnatureTypeAsWebCacheType 标记
	SingnatureTypeAsWebCacheType = "SingnatureType_Web"
	// RequestHeaderAccessKey 用于验证签名的key
	RequestHeaderAccessKey = "ack"
	// RequestHeaderSign 客户端签名结果的key
	RequestHeaderSign = "sign"
)

// SingnatureType 用于约束参数
type SingnatureType int

// AccessToken 访问密钥和签名
type AccessToken struct {
	UserID    string // 用户信息
	AccessKey string // 访问密钥
	SecretKey string // 加密签名
}

// ToJSON to JSON
func (ack AccessToken) ToJSON() string {
	bt, err := json.Marshal(ack)
	if nil != err {
		return err.Error()
	}
	return string(bt)
}

// AccessBody 放置到内存中的内容字段
type AccessBody struct {
	UserID       string
	SecretKey    string
	SessionAttrs map[string]string
}

// ErrorUserIDIsNil ErrorUserIDIsNil
var ErrorUserIDIsNil = errors.New("User ID cannot be empty")

// ErrorNotSupport ErrorNotSupport
var ErrorNotSupport = errors.New("This type is not supported")

// ErrorParamsNotEmpty ErrorParamsNotEmpty
var ErrorParamsNotEmpty = errors.New("Property name or content cannot be empty")

// ErrorSessionExpired ErrorSessionExpired
var ErrorSessionExpired = errors.New("Session has expired")
