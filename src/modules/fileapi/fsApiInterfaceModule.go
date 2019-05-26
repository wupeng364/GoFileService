package fileapi
/**
 *@description 文件API接口模块
 *@author	wupeng364@outlook.com
*/
import (
	"common/gomodule"
	"modules/filemanage"
	"modules/httpserver"
	"fmt"
)
type FsApiModule struct{
	fm *filemanage.FileManageModule
	hs *httpserver.HttpServerModule
}

// 返回模块信息
func (fa *FsApiModule)MInfo( )(*gomodule.ModuleInfo)	{
	return &gomodule.ModuleInfo{
		fa,
		"FsApiModule",
		1.0,
		"文件管理对外API接口模块",
	}
}
// 模块安装, 一个模块只初始化一次
func (fa *FsApiModule)MSetup( ) {
	
}
// 模块升级, 一个版本执行一次
func (fa *FsApiModule)MUpdate( ) {
	
}

// 每次启动加载模块执行一次
func (fa *FsApiModule)OnMInit(ref gomodule.ReferenceModule) {
	fa.fm = ref(fa.fm).(*filemanage.FileManageModule)
	fa.hs = ref(fa.hs).(*httpserver.HttpServerModule)
	imp_http{}.init(fa.fm, fa.hs )
	
}
// 系统执行销毁时执行
func (fa *FsApiModule)OnMDestroy( ) {
	
}

// ==============================================================================================



func sayHello( ){	
	fmt.Println("..")
}