package fileapimodel
/**
 *@description 文件API接口模块
 *@author	wupeng364@outlook.com
*/
import (
	"common/gomodule"
	"modules/filemanage"
	"modules/httpservermodel"
	"fmt"
)
type FsApimodel struct{
	fm *filemanage.FileManageModel
	hs *httpservermodel.HttpServerModel
}

// 返回模块信息
func (fa *FsApimodel)MInfo( )(*gomodule.ModelInfo)	{
	return &gomodule.ModelInfo{
		fa,
		"FsApimodel",
		1.0,
		"文件管理对外API接口模块",
	}
}
// 模块安装, 一个模块只初始化一次
func (fa *FsApimodel)MSetup( ) {
	
}
// 模块升级, 一个版本执行一次
func (fa *FsApimodel)MUpdate( ) {
	
}

// 每次启动加载模块执行一次
func (fa *FsApimodel)OnMInit( getPointer func(m interface{})interface{} ) {
	fa.fm = getPointer(fa.fm).(*filemanage.FileManageModel)
	fa.hs = getPointer(fa.hs).(*httpservermodel.HttpServerModel)
	imp_http{}.init(fa.fm, fa.hs )
	
}
// 系统执行销毁时执行
func (fa *FsApimodel)OnMDestroy( ) {
	
}

// ==============================================================================================



func sayHello( ){	
	fmt.Println("..")
}