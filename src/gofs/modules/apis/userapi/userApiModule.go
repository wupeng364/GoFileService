package userapi
/**
 *@description 
 *@author	wupeng364@outlook.com
*/
import (
	"gofs/common/gomodule"
	"gofs/modules/core/usermanage"
	"gofs/modules/common/httpserver"
	"gofs/modules/core/signature"
)

type UserApiModule struct{
	um *usermanage.UserManageModule
	hs *httpserver.HttpServerModule
	sg *signature.SignatureModule
}

// 返回模块信息
func (ua *UserApiModule)MInfo( )(*gomodule.ModuleInfo)	{
	return &gomodule.ModuleInfo{
		ua,
		"UserApiModule",
		1.0,
		"用户管理Api接口模块",
	}
}

// 模块安装, 一个模块只初始化一次
func (ua *UserApiModule)OnMSetup( ref gomodule.ReferenceModule ) {
	
}
// 模块升级, 一个版本执行一次
func (ua *UserApiModule)OnMUpdate( ref gomodule.ReferenceModule ) {
	
}

// 每次启动加载模块执行一次
func (ua *UserApiModule)OnMInit( ref gomodule.ReferenceModule ) {
	ua.um = ref(ua.um).(*usermanage.UserManageModule)
	ua.hs = ref(ua.hs).(*httpserver.HttpServerModule)
	ua.sg = ref(ua.sg).(*signature.SignatureModule)
	
	// http 方式
	(&imp_http{ um:ua.um, hs: ua.hs, sg: ua.sg, }).init( )
}

// 系统执行销毁时执行
func (ua *UserApiModule)OnMDestroy( ref gomodule.ReferenceModule ) {
	
}

// ==============================================================================================