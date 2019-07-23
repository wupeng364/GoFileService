package fileapi
/**
 *@description 文件API接口模块
 *文件的新建、删除、移动、复制等操作, 实现fsApiHttpInterface接口
 *@author	wupeng364@outlook.com
*/
import (
	"gofs/common/gomodule"
	"gofs/modules/core/filemanage"
	"gofs/modules/core/signature"
	"gofs/modules/common/httpserver"
	"fmt"
)
type FsApiModule struct{
	fm *filemanage.FileManageModule
	hs *httpserver.HttpServerModule
	sg *signature.SignatureModule
}

// 返回模块信息
func (fa *FsApiModule)MInfo( )(*gomodule.ModuleInfo)	{
	return &gomodule.ModuleInfo{
		fa,
		"FsApiModule",
		1.0,
		"文件管理Api接口模块",
	}
}
// 模块安装, 一个模块只初始化一次
func (fa *FsApiModule)OnMSetup( ref gomodule.ReferenceModule ) {
	
}
// 模块升级, 一个版本执行一次
func (fa *FsApiModule)OnMUpdate( ref gomodule.ReferenceModule ) {
	
}

// 每次启动加载模块执行一次
func (fa *FsApiModule)OnMInit(ref gomodule.ReferenceModule) {
	fa.fm = ref(fa.fm).(*filemanage.FileManageModule)
	fa.hs = ref(fa.hs).(*httpserver.HttpServerModule)
	fa.sg = ref(fa.sg).(*signature.SignatureModule)
	
	// http 方式
	(&imp_http{fm:fa.fm, hs:fa.hs, sg:fa.sg,}).init( )
	
}
// 系统执行销毁时执行
func (fa *FsApiModule)OnMDestroy( ref gomodule.ReferenceModule ) {
	
}

// ==============================================================================================



func sayHello( ){	
	fmt.Println("..")
}