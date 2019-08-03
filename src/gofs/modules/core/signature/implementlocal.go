package signature
import(
	"gofs/common/tokenmanager"
	"gofs/common/stringtools"
	"fmt"
)
/**
 * 本地内存存储的方式实现类
 */
type implement_Local struct{
	cache *tokenmanager.TokenManager
}

// 初始化模块, 这里会初始化Web和API的信息到内存中
func (this *implement_Local)SignatureInitial( ) error{
	this.cache = (&tokenmanager.TokenManager{}).Init( )
	return nil
}
// 生成访问令牌, 返回AccessToken
// accessBody={UserId:"admin", SecretKey:Guid2, SessionAttrs:{"key":"val",}}
func (this *implement_Local)GenerateAccessToken(userId string, singnatureType SingnatureType)(AccessToken, error){
	accessToken := AccessToken{}
	if len(userId) == 0 {
		return accessToken, Error_UserIdIsNil
	}
	accessToken.UserId		= userId
	accessToken.SecretKey   = stringtools.GetUUID( )  // 放到accessBody中, 后续可以根据AccessKey取出作为校验
	
	accessBody  := AccessBody{}
	accessBody.SecretKey    = accessToken.SecretKey
	accessBody.SessionAttrs = make(map[string]string) // 在本次会话中有效, 和AccessKey生命周期一致
	accessBody.UserId     = userId					  // 编辑当前用户ID
	if singnatureType == SingnatureType_Web {
		// 注册AccessKey到内存, 并放置accessBody
		accessToken.AccessKey = this.cache.AskToken(&tokenmanager.TokenObject{
			Second: SingnatureType_Web_DestroyTime,
			TypeStr: SingnatureType_Web_CacheType,
			TokenBody: accessBody,
		})
	}else if singnatureType == SingnatureType_API {
		// 注册到数据库和持久内存中
		// accessKey := stringtools.GetUUID( )
		
	}else{
		return AccessToken{}, Error_NotSupport 
	}
	return accessToken, nil
}
// 验证签名是否有效, 通过accessKey查找SecretKey然后校验参数
// Todd 可以尝试绑定IP
func (this *implement_Local)SignatureVerification(accessKey, sign string, requestparameter string)bool{
	if len(accessKey) == 0 || len(requestparameter) == 0 || len(sign) == 0{
		return false
	}
	tokenObject, exist := this.cache.GetTokenInfo(accessKey)
	if !exist {
		return false
	}
	accessBody := tokenObject.TokenBody.(AccessBody)
	calcSign := stringtools.String2MD5(requestparameter+accessBody.SecretKey)
	fmt.Println("SignatureVerification: ", accessKey, requestparameter+accessBody.SecretKey, calcSign, sign )
	if calcSign == sign {
		this.cache.RefreshToken(accessKey) // 刷新过期时间
		return true
	}
	return false
}
// 销毁签名, 使其无效
func (this *implement_Local)SignatureDestroy(accessKey string)error{
	this.cache.DestroyToken(accessKey)
	return nil
}
// 获取用户ID
func (this *implement_Local)GetUserId(accessKey string)string{
	if len(accessKey) == 0 {
		return ""
	}
	tokenObject, exist := this.cache.GetTokenInfo(accessKey)
	if !exist {
		return ""
	}
	accessBody := tokenObject.TokenBody.(AccessBody)
	return accessBody.UserId
}
// 设置属性到session里面, 会话过期自动删除
func (this *implement_Local)SetSessionAttr(accessKey, key, val string) error{
	if len(accessKey) == 0 || len(key) == 0 || len(val) == 0 {
		return Error_ParamsNotEmpty
	}
	tokenObject, exist := this.cache.GetTokenInfo(accessKey)
	if !exist {
		return Error_SessionExpired
	}
	accessBody := tokenObject.TokenBody.(AccessBody)
	accessBody.SessionAttrs[key] = val
	return nil
}
// 获取用户放在session里面的属性
func (this *implement_Local)GetSessionAttr(accessKey, key string)(string, error){
	if len(accessKey) == 0 || len(key) == 0 {
		return "", Error_ParamsNotEmpty
	}
	tokenObject, exist := this.cache.GetTokenInfo(accessKey)
	if !exist {
		return "", Error_SessionExpired
	}
	accessBody := tokenObject.TokenBody.(AccessBody)
	val, exist := accessBody.SessionAttrs[key]
	if exist {
		return val, nil
	}
	return "", nil
}
