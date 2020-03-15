// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 本地内存存储的方式实现签名校验

package signature

import (
	"fmt"
	"gutils/strtool"
	"gutils/tokentool"
)

// impByLocalCache 本机缓存实现
type impByLocalCache struct {
	cache *tokentool.TokenManager
}

// 初始化模块, 这里会初始化Web和API的信息到内存中
func (signature *impByLocalCache) SignatureInitial() error {
	signature.cache = (&tokentool.TokenManager{}).Init()
	return nil
}

// GenerateAccessToken 生成访问令牌, 返回AccessToken
// accessBody={UserID:"admin", SecretKey:Guid2, SessionAttrs:{"key":"val",}}
func (signature *impByLocalCache) GenerateAccessToken(userID string, singnatureType SingnatureType) (AccessToken, error) {
	accessToken := AccessToken{}
	if len(userID) == 0 {
		return accessToken, ErrorUserIDIsNil
	}
	accessToken.UserID = userID
	accessToken.SecretKey = strtool.GetUUID() // 放到accessBody中, 后续可以根据AccessKey取出作为校验

	accessBody := AccessBody{}
	accessBody.SecretKey = accessToken.SecretKey
	accessBody.SessionAttrs = make(map[string]string) // 在本次会话中有效, 和AccessKey生命周期一致
	accessBody.UserID = userID                        // 编辑当前用户ID
	if singnatureType == SingnatureTypeAsWeb {
		// 注册AccessKey到内存, 并放置accessBody
		accessToken.AccessKey = signature.cache.AskToken(&accessBody, SingnatureTypeAsWebDestroyTime)
	} else if singnatureType == SingnatureTypeAsAPI {
		// 注册到数据库和持久内存中
		// accessKey := strtool.GetUUID( )

	} else {
		return AccessToken{}, ErrorNotSupport
	}
	return accessToken, nil
}

// SignatureVerification 验证签名是否有效, 通过accessKey查找SecretKey然后校验参数
// Todd 可以尝试绑定IP
func (signature *impByLocalCache) SignatureVerification(accessKey, sign string, requestparameter string) bool {
	if len(accessKey) == 0 || len(requestparameter) == 0 || len(sign) == 0 {
		return false
	}
	tokenBody, exist := signature.cache.GetTokenBody(accessKey)
	if !exist {
		return false
	}
	accessBody := tokenBody.(*AccessBody)
	calcSign := strtool.GetMD5(requestparameter + accessBody.SecretKey)
	fmt.Println("SignatureVerification: ", accessKey, requestparameter+accessBody.SecretKey, calcSign, sign)
	if calcSign == sign {
		signature.cache.RefreshToken(accessKey) // 刷新过期时间
		return true
	}
	return false
}

// SignatureDestroy 销毁签名, 使其无效
func (signature *impByLocalCache) SignatureDestroy(accessKey string) error {
	signature.cache.DestroyToken(accessKey)
	return nil
}

// GetUserID 获取用户ID
func (signature *impByLocalCache) GetUserID(accessKey string) string {
	if len(accessKey) == 0 {
		return ""
	}
	tokenBody, exist := signature.cache.GetTokenBody(accessKey)
	if !exist {
		return ""
	}
	accessBody := tokenBody.(*AccessBody)
	return accessBody.UserID
}

// SetSessionAttr 设置属性到session里面, 会话过期自动删除
func (signature *impByLocalCache) SetSessionAttr(accessKey, key, val string) error {
	if len(accessKey) == 0 || len(key) == 0 || len(val) == 0 {
		return ErrorParamsNotEmpty
	}
	tokenBody, exist := signature.cache.GetTokenBody(accessKey)
	if !exist {
		return ErrorSessionExpired
	}
	accessBody := tokenBody.(*AccessBody)
	accessBody.SessionAttrs[key] = val
	return nil
}

// GetSessionAttr 获取用户放在session里面的属性
func (signature *impByLocalCache) GetSessionAttr(accessKey, key string) (string, error) {
	if len(accessKey) == 0 || len(key) == 0 {
		return "", ErrorParamsNotEmpty
	}
	tokenBody, exist := signature.cache.GetTokenBody(accessKey)
	if !exist {
		return "", ErrorSessionExpired
	}
	accessBody := tokenBody.(*AccessBody)
	val, exist := accessBody.SessionAttrs[key]
	if exist {
		return val, nil
	}
	return "", nil
}
