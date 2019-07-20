package signature

/**
 * 会话管理接口
 * 会话存储在内存中, 或者Redis中
 */
type signature interface{
	// 初始化模块
	SignatureInitial( ) error
	// 生成访问令牌, 返回AccessToken
	GenerateAccessToken(userId string, singnatureType SingnatureType)(AccessToken, error)
	// 验证签名是否有效, 通过accessKey查找SecretKey, 通过MD5(SecretKey+requestparameter)==sign校验参数
	SignatureVerification(accessKey, sign string, requestparameter string)bool
	// 销毁签名, 使其无效
	SignatureDestroy(accessKey string)error
	// 获取用户ID
	GetUserId(accessKey string)string
	// 设置属性到session里面, 会话过期自动删除
	SetSessionAttr(accessKey, key, val string) error
	// 获取用户放在session里面的属性
	GetSessionAttr(accessKey, key string)(string, error)
	
}