package fileapi
/**
 *@description 文件API接口模块
 *文件的新建、删除、移动、复制等操作, 实现fsApiHttpInterface接口
 *@author	wupeng364@outlook.com
*/
import (
	"gofs/common/moduletools"
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
func (this *FsApiModule)MInfo( )(*moduletools.ModuleInfo)	{
	return &moduletools.ModuleInfo{
		this,
		"FsApiModule",
		1.0,
		"文件管理Api接口模块",
	}
}
// 模块安装, 一个模块只初始化一次
func (this *FsApiModule)OnMSetup( ref moduletools.ReferenceModule ) {
	
}
// 模块升级, 一个版本执行一次
func (this *FsApiModule)OnMUpdate( ref moduletools.ReferenceModule ) {
	
}

// 每次启动加载模块执行一次
func (this *FsApiModule)OnMInit(ref moduletools.ReferenceModule) {
	this.fm = ref(this.fm).(*filemanage.FileManageModule)
	this.hs = ref(this.hs).(*httpserver.HttpServerModule)
	this.sg = ref(this.sg).(*signature.SignatureModule)
	
	// http 方式
	(&imp_http{fm:this.fm, hs:this.hs, sg:this.sg,}).init( )
	
}
// 系统执行销毁时执行
func (this *FsApiModule)OnMDestroy( ref moduletools.ReferenceModule ) {
	
}

// ==============================================================================================



func sayHello( ){	
	fmt.Println("..")
}