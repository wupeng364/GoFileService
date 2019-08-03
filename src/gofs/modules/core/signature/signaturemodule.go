package signature
/**
 *@description 请求签名模块
 *对http请求的拦截, 验证参数完整性、身份合法性检测、用户session管理
 *@author	wupeng364@outlook.com
*/
import (
	"sort"
	"net/http"
	"gofs/common/moduletools"
	hst "gofs/common/httpservertools"
	"gofs/modules/common/httpserver"
	"gofs/modules/common/config"
)

type SignatureModule struct{
	cfg 		*config.ConfigModule
	httpserver  *httpserver.HttpServerModule
	sign		signature
}

// 返回模块信息
func (this *SignatureModule)MInfo( )(*moduletools.ModuleInfo)	{
	return &moduletools.ModuleInfo{
		this,
		"SignatureModule",
		1.0,
		"Api接口签名模块",
	}
}

// 模块安装, 一个模块只初始化一次
func (this *SignatureModule)OnMSetup( ref moduletools.ReferenceModule ) {
	
}
// 模块升级, 一个版本执行一次
func (this *SignatureModule)OnMUpdate( ref moduletools.ReferenceModule ) {
	
}

// 每次启动加载模块执行一次
func (this *SignatureModule)OnMInit( ref moduletools.ReferenceModule ) {
	this.cfg = ref(this.cfg).(*config.ConfigModule)
	this.httpserver = ref(this.httpserver).(*httpserver.HttpServerModule)
	
	// 这里暂时只实现单机、本地内存版本
	this.sign = &implement_Local{ }
	this.sign.SignatureInitial()
}

// 系统执行销毁时执行
func (this *SignatureModule)OnMDestroy( ref moduletools.ReferenceModule ) {
	
}

// ==============================================================================================
// 添加会话
func (this *SignatureModule)CreateWebSession( userId string, r *http.Request )(AccessToken, error) {
	return this.sign.GenerateAccessToken(userId, SingnatureType_Web)
}
// 获取会话信息
func (this *SignatureModule)GetUserIdByAccessKey( ack string)string{
	return this.sign.GetUserId(ack)
}
// Api签名拦截器
func (this *SignatureModule)ApiFilter_Http(w http.ResponseWriter, r *http.Request, next hst.FilterNext){
	//fmt.Println("ApiFilter: ", r.RemoteAddr, r.URL.Path)
	// 从请求中获取accessKey, 不能为空
	accessKey := r.Header.Get(Request_Header_AccessKey)
	// 从请求中获取signature, 不能为空
	sign := r.Header.Get(Request_Header_Sign)
	if len(accessKey) == 0 || len(sign) == 0{
		w.WriteHeader(http.StatusUnauthorized); return // 401
	}
	// 填充Form对象
	if nil == r.Form {
		err := r.ParseForm( )
		if nil != err {
			// 出现异常, 不继续处理
			w.WriteHeader(http.StatusInternalServerError); return
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
		if  _keysLen > 0 {
			sort.Strings(keys)
			for i, val := range keys {
				requestparameter += val+"="+r.Form[val][0]
				if i < _keysLen-1 {
					requestparameter += "&"
				}
			}
			
		}
	}
	// 校验参数合法性
	if !this.sign.SignatureVerification(accessKey, sign, requestparameter) {
		w.WriteHeader(http.StatusUnauthorized); return // 401
	}
	
	next( ) // 校验参数合法性-通过
}