package httpservertools

/**
 *@description http服务端工具类
 *@author	wupeng364@outlook.com
*/
import (
	"strings"
	"common/stringtools"
    "net/http"
    "errors"
    "regexp"
    "fmt"
)
// routers
type router struct {
	isDebug		   bool
	defaultHandler http.Handler
	urlHandlersMap *map[string]HandlersFunc
	regexpHandlersMap *map[string]HandlersFunc
}
// 实现 Handler 的 ServeHTTP 接口
func (rt *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if rt.isDebug {
		fmt.Println("URL.Path: ", r.URL.Path)
	}
	// strings.ToLower( r.URL.Path )
	// stringtools.UnixPathClear(r.URL.Path)
    if h, ok := (*rt.urlHandlersMap)[r.URL.Path]; ok {
        h(w, r)
    }else{
    	// 检查正则
    	for key, h := range (*rt.regexpHandlersMap){
    		_SymbolIndex := strings.Index(key, ":")
    		if _SymbolIndex == -1 {
	    		continue
    		}
	    	ok, _ := regexp.MatchString(key[:_SymbolIndex]+key[_SymbolIndex+1:], r.URL.Path)
    		if rt.isDebug {
				fmt.Println("URL.Regexp: ", key[:_SymbolIndex]+key[_SymbolIndex+1:], ok)
			}
	    	if ok {
		    	 h(w, r); return
	    	}
    	}
    	// 默认处理器
    	if rt.defaultHandler != nil {
			rt.defaultHandler.ServeHTTP(w, r)	
    	}else{
	    	w.WriteHeader(http.StatusNotFound)
    	}
    }
}
// 定义请求处理器
type HandlersFunc func(http.ResponseWriter, *http.Request)

// struct 注册器接口
type Registrar interface{
	RoutList( ) StructRegistrar
}
// struct 注册器
type StructRegistrar struct{
	IsToLower 	bool	        // 是否需要转小写访问
	BasePath 	string  		// 基础路径 /base/child....
	HandlerFunc []HandlersFunc	// 需要注册的 fuc
}


// 注册路由, 支持正则表达式(以':'符号开始, 如: /upload/:\S+)
func RegistService( addr string, handlersMap map[string]HandlersFunc, defaultHandler http.Handler ) error{
	if handlersMap == nil || len(handlersMap) == 0 {
		return errors.New("handlersMap is nil")
	}
	if addr == "" {
		return errors.New("addr is nil or empty")
	}
	_router := &router{}
	_router.isDebug = true
	_router.urlHandlersMap = getUrlHandlersMap(handlersMap)
	_router.regexpHandlersMap = getRegexpHandlersMap(handlersMap)
	_router.defaultHandler = defaultHandler
	server := &http.Server{
        Addr:           addr,
        Handler:        _router,
        ReadTimeout:    0,
        WriteTimeout:   0,
        MaxHeaderBytes: 1 << 20,
    }
	return server.ListenAndServe()
}
func getUrlHandlersMap( handlersMap map[string]HandlersFunc )*map[string]HandlersFunc{
	if len(handlersMap) == 0 {
		return nil
	}
	temp := make(map[string]HandlersFunc)
	for key, val := range handlersMap {
		if strings.Index(key, ":") == -1 {
			temp[key] = val
		}
	}
	return &temp
}
func getRegexpHandlersMap( handlersMap map[string]HandlersFunc )*map[string]HandlersFunc{
	if len(handlersMap) == 0 {
		return nil
	}
	temp := make(map[string]HandlersFunc)
	for key, val := range handlersMap {
		if strings.Index(key, ":") > -1 {
			temp[key] = val
		}
	}
	return &temp
}
// 根据类型来构造请求路由
func BuildHandlersMap(rs Registrar )(map[string]HandlersFunc, error) {
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
	handlersMap := make(map[string]HandlersFunc, len(hf))
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