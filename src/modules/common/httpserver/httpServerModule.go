package httpserver
/**
 *@description HTTP服务模块
 *http服务路由注册、请求拦截、端口侦听
 *@author	wupeng364@outlook.com
*/

import (
	hst "common/httpservertools"
	"common/stringtools"
	"common/gomodule"
	"modules/common/config"
	"strings"
	"net/http"
	"errors"
	"fmt"
)
type HttpServerModule struct{
	serviceRouter  *hst.ServiceRouter
	cfgModule *config.ConfigModule
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
func (hs *HttpServerModule)OnMSetup( ref gomodule.ReferenceModule ) {
	
}
// 模块升级, 一个版本执行一次
func (hs *HttpServerModule)OnMUpdate( ref gomodule.ReferenceModule ) {
	
}

// 每次启动加载模块执行一次
func (hs *HttpServerModule)OnMInit( ref gomodule.ReferenceModule ) {
	hs.serviceRouter = &hst.ServiceRouter{}
	hs.serviceRouter.ClearHandlersMap( )
	hs.serviceRouter.SetDebug(true)
	
	hs.cfgModule = ref(hs.cfgModule).(*config.ConfigModule)
	cfgMap := hs.cfgModule.GetConfigs(cfg_http).(map[string]interface{})
	if len( cfgMap ) == 0 {
		panic("http server config is nil")
	}
	hs.listenAddr = cfgMap[cfg_http_addr].(string)+":"+cfgMap[cfg_http_port].(string)
}
// 系统执行销毁时执行
func (hs *HttpServerModule)OnMDestroy( ref gomodule.ReferenceModule ) {
	
}

// ==============================================================================================
func (hs *HttpServerModule)GetListenAddr( )string{
	return hs.listenAddr;
}
// 注册忽略签名验证的路径
func (hs *HttpServerModule)AddIgnoreFilter(url string){
	hs.serviceRouter.AddUrlFilter(url, func(w http.ResponseWriter, r *http.Request, next hst.FilterNext){
		next()
	});
}
// 默认全局过滤器, 可用于校验签名和session
func (hs *HttpServerModule)AddGlobalFilter(globalFilter hst.FilterFunc)error{
	if globalFilter == nil {
		return errors.New("globalFilter is nil")
	}
	hs.serviceRouter.SetGlobalFilter(globalFilter)
	return nil
}
// URL过滤器, 匹配方式: startWith url
func (hs *HttpServerModule)AddUrlFilter(url string, urlFilter hst.FilterFunc)error{
	if len(url) == 0 {
		return errors.New("Filter url is nil")
	}
	if urlFilter == nil {
		return errors.New("globalFilter func is nil")
	}
	hs.serviceRouter.AddUrlFilter(url, urlFilter)
	return nil
}
// 注册默认处理器
func (hs *HttpServerModule)DefaultHandler(df hst.HandlersFunc)(err error){
	if df == nil {
		return errors.New("http.Handler is nil")
	}
	hs.serviceRouter.SetDefaultHandler(df)
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
	hs.serviceRouter.AddHandler(path, hf)
	return nil
}
// 通过服务模板注册路由服务
func (hs *HttpServerModule)AddRegistrar(rs Registrar)error{
	handlers, err := BuildHandlersMap(rs)
	if err != nil {
		return err
	}
	hs.serviceRouter.AddHandlers(handlers)
	return nil
}
// 启动服务
func (hs *HttpServerModule)DoStartServer( ) error {
	fmt.Println("   >Server listened in :"+ hs.listenAddr)
	if len(hs.listenAddr) == 0 {
		return errors.New("addr is nil or empty")
	}
	
	server := &http.Server{
        Addr:           hs.listenAddr,
        Handler:        hs.serviceRouter,
        ReadTimeout:    0,
        WriteTimeout:   0,
        MaxHeaderBytes: 1 << 20,
    }
	return server.ListenAndServe( )
}

// =============================================================================================
// 根据类型来构造请求路由
func BuildHandlersMap(rs Registrar )(map[string]hst.HandlersFunc, error) {
	if rs == nil {
		return nil, errors.New("Not find Registrar")
	}
	srg := rs.RoutList( )
	hf  := srg.HandlerFunc
	if len(hf) == 0 {
		return nil, errors.New("Not find method of Registrar")
	}
	baseUrl := stringtools.UnixPathClear(srg.BasePath)
	if baseUrl != "/" {
		baseUrl += "/"
	}
	handlersMap := make(map[string]hst.HandlersFunc, len(hf))
	for _, fn := range hf{
		_fn := stringtools.GetFunctionName(fn, '.')
		_fm := strings.Index(_fn, "-")
		if _fm > -1 {
			_fn = _fn[:_fm]
		}
		_fn = baseUrl+_fn
		if srg.IsToLower {
			_fn = strings.ToLower( _fn )
		}
		handlersMap[_fn] = fn
	}
	return handlersMap, nil
}