package signature
/**
 *@description 请求签名模块
 *对http请求的拦截, 验证参数完整性、身份合法性检测、用户session管理
 *@author	wupeng364@outlook.com
*/
import (
	"sort"
	"net/http"
	"gofs/common/moduleloader"
	hst "gofs/common/httpservertools"
)

type SignatureModule struct{
	sign		signature
}

// 返回模块信息
func (this *SignatureModule)ModuleOpts( )(moduleloader.Opts) {
	return moduleloader.Opts{
		Name: "SignatureModule",
		Version: 1.0,
		Description: "Api接口签名模块",
		OnReady: func (mctx *moduleloader.Loader) {
		},
		OnInit: this.onMInit,
	}
}

// 每次启动加载模块执行一次
func (this *SignatureModule)onMInit( ) {
	// 这里暂时只实现单机、本地内存版本
	this.sign = &implement_Local{ }
	this.sign.SignatureInitial()
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