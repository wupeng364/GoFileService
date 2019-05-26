package httpserver
/**
 *@description HTTP服务模块
 *@author	wupeng364@outlook.com
*/

import (
	hst "common/httpservertools"
	"modules/config"
	"common/gomodule"
	"net/http"
	"errors"
	"fmt"
)
type HttpServerModule struct{
	cfgModule *config.ConfigModule
	defaultHandler http.Handler
	handlers map[string]hst.HandlersFunc
	listenAddr string
}

// 返回模块信息
func (hs *HttpServerModule)MInfo( )(*gomodule.ModuleInfo)	{
	return &gomodule.ModuleInfo{
		hs,
		"HttpServerModule",
		1.0,
		"HTTP服务模块",
	}
}

// 模块安装, 一个模块只初始化一次
func (hs *HttpServerModule)MSetup( ) {
	
}
// 模块升级, 一个版本执行一次
func (hs *HttpServerModule)MUpdate( ) {
	
}

// 每次启动加载模块执行一次
func (hs *HttpServerModule)OnMInit( ref gomodule.ReferenceModule ) {
	hs.cfgModule = ref(hs.cfgModule).(*config.ConfigModule)
	hs.handlers = make(map[string]hst.HandlersFunc)
	
	cfgMap := hs.cfgModule.GetConfigs(cfg_http).(map[string]interface{})
	if len( cfgMap ) == 0 {
		panic("http server config is nil")
	}
	hs.listenAddr = cfgMap[cfg_http_addr].(string)+":"+cfgMap[cfg_http_port].(string)
}
// 系统执行销毁时执行
func (hs *HttpServerModule)OnMDestroy( ) {
	
}

// ==============================================================================================
func (hs *HttpServerModule)GetListenAddr( )string{
	return hs.listenAddr;
}
// 注册默认处理器
func (hs *HttpServerModule)DefaultHandler(df http.Handler)(err error){
	if df == nil {
		return errors.New("http.Handler is nil")
	}
	hs.defaultHandler = df
	return nil
}
// 通过 HandlersFunc 注册
func (hs *HttpServerModule)AddHandlerFunc(path string, hf hst.HandlersFunc)(err error){
	if len(path) == 0 {
		return errors.New("Path is nil")
	}
	if hf == nil {
		return errors.New("HandlersFunc is nil")
	}
	hs.handlers[path] = hf
	return nil
}
// 通过服务模板注册路由服务
func (hs *HttpServerModule)AddRegistrar(rs hst.Registrar)(err error){
	len_hds := len(hs.handlers)
	if len_hds == 0 {
		hs.handlers, err = hst.BuildHandlersMap(rs)
	}else{
		temp, err := hst.BuildHandlersMap(rs)
		if err != nil {
			return err
		}
		len_temp := len(temp)
		if len_temp == 0 {
			return nil
		}
		for key, val := range temp {
			hs.handlers[key] = val
		}
	}
	return err
}
// 启动服务
func (hs *HttpServerModule)DoStartServer( ) error {
	fmt.Println("   >Server listened in :"+ hs.listenAddr)	
	err := hst.RegistService(hs.listenAddr, hs.handlers, hs.defaultHandler)
	if err != nil {
		return err 
	}
	return nil
}

func sayHello( ){
	fmt.Println("..")
}