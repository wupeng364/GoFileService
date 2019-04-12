package htmlmodel

/**
 *@description 静态资源加载器
 *@author	wupeng364@outlook.com
*/

import (
	g "common/gomodule"
	"modules/configmodel"
	"modules/httpservermodel"
	"path/filepath"
	"net/http"
	"fmt"
)
type Htmlmodel struct{
	cfgModel *configmodel.ConfigModel
	hsvModel *httpservermodel.HttpServerModel
}

// 返回模块信息
func (h *Htmlmodel)MInfo( )(*g.ModelInfo)	{
	return &g.ModelInfo{
		h,
		"Htmlmodel",
		1.0,
		"网页/静态资源处理",
	}
}

// 模块安装, 一个模块只初始化一次
func (h *Htmlmodel)MSetup( ) {
	
}
// 模块升级, 一个版本执行一次
func (h *Htmlmodel)MUpdate( ) {
	
}

// 每次启动加载模块执行一次
func (h *Htmlmodel)OnMInit( getPointer func(m interface{})interface{} ) {
	h.cfgModel = getPointer(h.cfgModel).(*configmodel.ConfigModel)
	h.hsvModel = getPointer(h.hsvModel).(*httpservermodel.HttpServerModel)
	// config
	static := h.cfgModel.GetConfig(cfg_http_static)
	if !filepath.IsAbs(static) {
		temp, err := filepath.Abs(static)
		if err != nil {
			panic( err )
		}
		static = temp
	}
	
	// 首页重定向
	err := h.hsvModel.AddHandlerFunc("/", PageDispatch{}.Index)
	if err != nil {
		panic( err )
	}
	// 页面资源
	err = h.hsvModel.DefaultHandler( http.StripPrefix("/", http.FileServer(http.Dir( static ))) )
	if err != nil {
		panic( err )
	}
	fmt.Println("   > Htmlmodel http registered end")
}
// 系统执行销毁时执行
func (h *Htmlmodel)OnMDestroy( ) {
	
}

// ==============================================================================================



func sayHello( ){
	fmt.Println("..")
}