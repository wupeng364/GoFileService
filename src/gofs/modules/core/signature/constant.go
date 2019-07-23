package signature
import(
	"errors"
	"encoding/json"
)
const(
	SingnatureType_API SingnatureType = 1	// 放置数据, 不删除就不过期
	SingnatureType_Web SingnatureType = 0   // 放置内存, 会过期
	SingnatureType_Web_DestroyTime = 30*60	// 默认15分钟销毁
	SingnatureType_Web_CacheType   = "SingnatureType_Web" // 标记
	
	Request_Header_AccessKey = "ack" // 用于验证签名的key
	Request_Header_Sign = "sign" 	   // 客户端签名结果的key
)

// SingnatureType 用于约束参数
type SingnatureType int
// 访问密钥和签名
type AccessToken struct{
	UserId    string    // 用户信息
	AccessKey string	// 访问密钥
	SecretKey string	// 加密签名
}
func (ack AccessToken)ToJson() string{
	bt, err := json.Marshal( ack )
	if nil != err {
		return err.Error( )
	}
	return string(bt)
}
// 放置到内存中的内容字段, 一次会话期间有效
type AccessBody struct {
	UserId 			string
	SecretKey 		string
	SessionAttrs	map[string]string
}

var Error_UserIdIsNil = errors.New("User ID cannot be empty")
var Error_NotSupport  = errors.New("This type is not supported")
var Error_ParamsNotEmpty  = errors.New("Property name or content cannot be empty")
var Error_SessionExpired  = errors.New("Session has expired")