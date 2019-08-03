package userapi
/**
 *@description 
 *@author	wupeng364@outlook.com
*/
import (
	"gofs/common/moduletools"
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
func (this *UserApiModule)MInfo( )(*moduletools.ModuleInfo)	{
	return &moduletools.ModuleInfo{
		this,
		"UserApiModule",
		1.0,
		"用户管理Api接口模块",
	}
}

// 模块安装, 一个模块只初始化一次
func (this *UserApiModule)OnMSetup( ref moduletools.ReferenceModule ) {
	
}
// 模块升级, 一个版本执行一次
func (this *UserApiModule)OnMUpdate( ref moduletools.ReferenceModule ) {
	
}

// 每次启动加载模块执行一次
func (this *UserApiModule)OnMInit( ref moduletools.ReferenceModule ) {
	this.um = ref(this.um).(*usermanage.UserManageModule)
	this.hs = ref(this.hs).(*httpserver.HttpServerModule)
	this.sg = ref(this.sg).(*signature.SignatureModule)
	
	// http 方式
	(&imp_http{ um:this.um, hs: this.hs, sg: this.sg, }).init( )
}

// 系统执行销毁时执行
func (this *UserApiModule)OnMDestroy( ref moduletools.ReferenceModule ) {
	
}

// ==============================================================================================