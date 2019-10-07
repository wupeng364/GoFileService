package userapi
/**
 *@description 
 *@author	wupeng364@outlook.com
*/
import (
	"gofs/common/moduleloader"
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
func (this *UserApiModule)ModuleOpts( )(moduleloader.Opts)	{
	return moduleloader.Opts{
		Name: "UserApiModule",
		Version: 1.0,
		Description: "用户管理Api接口模块",
		OnReady: func (mctx *moduleloader.Loader) {
			this.um = mctx.GetModuleByTemplate(this.um).(*usermanage.UserManageModule)
			this.hs = mctx.GetModuleByTemplate(this.hs).(*httpserver.HttpServerModule)
			this.sg = mctx.GetModuleByTemplate(this.sg).(*signature.SignatureModule)
		},
		OnInit: func ( ) {
			// http 方式
			(&imp_http{ um:this.um, hs: this.hs, sg: this.sg, }).init( )
		},
	}
}