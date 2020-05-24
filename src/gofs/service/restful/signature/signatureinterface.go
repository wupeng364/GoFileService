// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package signature

// 会话管理接口, 会话可存储在内存中, 或者Redis中
type signature interface {
	// 初始化模块
	SignatureInitial() error
	// 生成访问令牌, 返回AccessToken
	GenerateAccessToken(userID string, singnatureType SingnatureType) (AccessToken, error)
	// 验证签名是否有效, 通过accessKey查找SecretKey, 通过MD5(SecretKey+requestparameter)==sign校验参数
	VerificationSignature(accessKey, sign string, requestparameter string) bool
	// 销毁签名, 使其无效
	DestroySignature(accessKey string) error
	// 获取用户ID
	GetUserID(accessKey string) string
	// 设置属性到session里面, 会话过期自动删除
	SetSessionAttr(accessKey, key, val string) error
	// 获取用户放在session里面的属性
	GetSessionAttr(accessKey, key string) (string, error)
}
