package htmlpage

/**
 *@description 静态资源加载器
 *@author	wupeng364@outlook.com
*/

import (
	"common/gomodule"
	"modules/config"
	"modules/httpserver"
	"path/filepath"
	"net/http"
	"fmt"
)
type HtmlModule struct{
	cfgModule *config.ConfigModule
	hsvModule *httpserver.HttpServerModule
}

// 返回模块信息
func (h *HtmlModule)MInfo( )(*gomodule.ModuleInfo)	{
	return &gomodule.ModuleInfo{
		h,
		"HtmlModule",
		1.0,
		"网页/静态资源处理",
	}
}

// 模块安装, 一个模块只初始化一次
func (h *HtmlModule)MSetup( ) {
	
}
// 模块升级, 一个版本执行一次
func (h *HtmlModule)MUpdate( ) {
	
}

// 每次启动加载模块执行一次
func (h *HtmlModule)OnMInit( ref gomodule.ReferenceModule ) {
	h.cfgModule = ref(h.cfgModule).(*config.ConfigModule)
	h.hsvModule = ref(h.hsvModule).(*httpserver.HttpServerModule)
	// config
	static := h.cfgModule.GetConfig(cfg_http_static)
	if !filepath.IsAbs(static) {
		temp, err := filepath.Abs(static)
		if err != nil {
			panic( err )
		}
		static = temp
	}
	
	// 首页重定向
	err := h.hsvModule.AddHandlerFunc("/", PageDispatch{}.Index)
	if err != nil {
		panic( err )
	}
	// 页面资源
	err = h.hsvModule.DefaultHandler( http.StripPrefix("/", http.FileServer(http.Dir( static ))) )
	if err != nil {
		panic( err )
	}
	fmt.Println("   > Htmlpage Module http registered end")
}
// 系统执行销毁时执行
func (h *HtmlModule)OnMDestroy( ) {
	
}

// ==============================================================================================



func sayHello( ){
	fmt.Println("..")
}