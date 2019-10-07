package htmlpage

/**
 *@description 静态资源加载器
 *静态资源处理
 *@author	wupeng364@outlook.com
*/

import (
	"strings"
	"gofs/common/moduleloader"
	// "common/filetools"
	hst "gofs/common/httpservertools"
	"gofs/modules/common/httpserver"
	"gofs/modules/core/signature"
	"path/filepath"
	"net/http"
	"fmt"
)
type HtmlModule struct{
	mctx *moduleloader.Loader
	hsvModule *httpserver.HttpServerModule
	sign 	  *signature.SignatureModule
}

// 返回模块信息
func (this *HtmlModule)ModuleOpts( )(moduleloader.Opts){
	return moduleloader.Opts{
		Name: "HtmlModule",
		Version: 1.0,
		Description: "静态资源处理",
		OnReady: func (mctx *moduleloader.Loader) {
			this.mctx = mctx
			this.hsvModule = mctx.GetModuleByTemplate(this.hsvModule).(*httpserver.HttpServerModule)
			this.sign = mctx.GetModuleByTemplate(this.sign).(*signature.SignatureModule)
		},
		OnInit: this.onMInit,
	}
}
// 每次启动加载模块执行一次
func (this *HtmlModule)onMInit( ) {	
	// config
	static := this.mctx.GetConfig(cfg_http_static)
	if !filepath.IsAbs(static) {
		temp, err := filepath.Abs(static)
		if err != nil {
			panic( err )
		}
		static = temp
	}
	// 全局过滤器
	err := this.hsvModule.AddGlobalFilter(this.staticFilter)
	
	// 首页重定向
	err = this.hsvModule.AddHandlerFunc("/", PageDispatch{}.Index)
	if err != nil {
		panic( err )
	}
	
	// 页面资源
	err = this.hsvModule.DefaultHandler( http.StripPrefix(StaticSource_BasePath, http.FileServer(http.Dir( static ))).ServeHTTP )
	if err != nil {
		panic( err )
	}
	
	// 全局默认过滤器, 如果没有特殊设定的url都会到这里来
	this.hsvModule.AddGlobalFilter(this.staticFilter)
	
	fmt.Println("   > Htmlpage Module http registered end")
}
// ==============================================================================================
// 静态资源过滤器
func (this *HtmlModule)staticFilter(w http.ResponseWriter, r *http.Request, next hst.FilterNext){
	// fmt.Println("Static Path: ", r.URL.Path)
	if strings.HasSuffix(r.URL.Path, ".html"){
		if ack, err := r.Cookie("ack"); nil == err{
			if len(ack.Value) > 0{
				// 是否只要 cookie 里面有合法 accessKey 就可以了
				userId := this.sign.GetUserIdByAccessKey(ack.Value)
				//fmt.Println("session: ", ack.Value, userId)
				if len(userId) > 0{
					next( ); return;
				}
			}
		}
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
	next( )
}
// 会话过期处理器
func (this *HtmlModule)sessionExpiredHandler(w http.ResponseWriter, r *http.Request){
	fmt.Println("Session expired: ", r.URL.Path)
}