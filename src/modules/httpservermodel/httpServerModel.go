package httpservermodel
/**
 *@description HTTP服务模块
 *@author	wupeng364@outlook.com
*/

import (
	hst "common/httpservertools"
	"modules/configmodel"
	"common/gomodule"
	"net/http"
	"errors"
	"fmt"
)
type HttpServerModel struct{
	cfgModel *configmodel.ConfigModel
	defaultHandler http.Handler
	handlers map[string]hst.HandlersFunc
	listenAddr string
}

// 返回模块信息
func (hs *HttpServerModel)MInfo( )(*gomodule.ModelInfo)	{
	return &gomodule.ModelInfo{
		hs,
		"HttpServerModel",
		1.0,
		"HTTP服务模块",
	}
}

// 模块安装, 一个模块只初始化一次
func (hs *HttpServerModel)MSetup( ) {
	
}
// 模块升级, 一个版本执行一次
func (hs *HttpServerModel)MUpdate( ) {
	
}

// 每次启动加载模块执行一次
func (hs *HttpServerModel)OnMInit( getPointer func(m interface{})interface{}  ) {
	hs.cfgModel = getPointer(hs.cfgModel).(*configmodel.ConfigModel)
	hs.handlers = make(map[string]hst.HandlersFunc)
	
	cfgMap := hs.cfgModel.GetConfigs(cfg_http).(map[string]interface{})
	if len( cfgMap ) == 0 {
		panic("http server config is nil")
	}
	hs.listenAddr = cfgMap[cfg_http_addr].(string)+":"+cfgMap[cfg_http_port].(string)
}
// 系统执行销毁时执行
func (hs *HttpServerModel)OnMDestroy( ) {
	
}

// ==============================================================================================
func (hs *HttpServerModel)GetListenAddr( )string{
	return hs.listenAddr;
}
// 注册默认处理器
func (hs *HttpServerModel)DefaultHandler(df http.Handler)(err error){
	if df == nil {
		return errors.New("http.Handler is nil")
	}
	hs.defaultHandler = df
	return nil
}
// 通过 HandlersFunc 注册
func (hs *HttpServerModel)AddHandlerFunc(path string, hf hst.HandlersFunc)(err error){
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
func (hs *HttpServerModel)AddRegistrar(rs hst.Registrar)(err error){
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
func (hs *HttpServerModel)DoStartServer( ) error {
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