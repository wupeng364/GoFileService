package fileapi
/**
 *@description 文件API接口模块
 *文件的新建、删除、移动、复制等操作, 实现fsApiHttpInterface接口
 *@author	wupeng364@outlook.com
*/
import (
	"gofs/common/moduleloader"
	"gofs/modules/core/filemanage"
	"gofs/modules/core/signature"
	"gofs/modules/common/httpserver"
)
type FsApiModule struct{
	fm *filemanage.FileManageModule
	hs *httpserver.HttpServerModule
	sg *signature.SignatureModule
}

// 返回模块信息Name: 
func (this *FsApiModule)ModuleOpts( )(moduleloader.Opts)	{
	return moduleloader.Opts{
		Name: "FsApiModule",
		Version: 1.0,
		Description: "文件管理Api接口模块",
		OnReady: func (mctx *moduleloader.Loader) {
			this.fm = mctx.GetModuleByTemplate(this.fm).(*filemanage.FileManageModule)
			this.hs = mctx.GetModuleByTemplate(this.hs).(*httpserver.HttpServerModule)
			this.sg = mctx.GetModuleByTemplate(this.sg).(*signature.SignatureModule)
		},
		OnInit: func ( ) {
			// http 方式
			(&imp_http{fm:this.fm, hs:this.hs, sg:this.sg,}).init( )
		},
	}
}