package httpservertools

/**
 *@description http服务端工具类
 *@author	wupeng364@outlook.com
*/
import (
	"strings"
    "net/http"
    "errors"
    "regexp"
    "fmt"
)
// 定义请求处理器
type HandlersFunc func(http.ResponseWriter, *http.Request)
// http请求过滤器, next函数用于触发下一步操作, 不执行就不继续处理请求
type FilterNext func( )
type FilterFunc func(http.ResponseWriter, *http.Request, FilterNext)

// routers一个实现了ServeHTTP的Handler对象
// 提供简单的url路由
type ServiceRouter struct {
	isDebug		      	bool			// 调试模式可以打印信息
	defaultHandler    	HandlersFunc  // 默认的url处理, 可以用于处理静态资源
	urlHandlersMap    	map[string]HandlersFunc // url路径全匹配路由表
	regexpHandlersMap 	map[string]HandlersFunc // url路径正则配路由表
	regexpHandlersIndex []string				// url路径正则配路由表-索引(用于保存顺序)
	urlFiltersMap 	  	map[string]FilterFunc   // url路径过滤器
	regexpFiltersMap  	map[string]FilterFunc   // url路径正则匹配过滤器
	regexpFiltersIndex  []string			    // url路径正则匹配过滤器-索引(用于保存顺序)
	filterWhitelist	  	map[string]string
	globalFileter	  	FilterFunc
}

// 根据注册的路由表调用对应的函数
// 优先匹配全url > 正则url > 默认处理器 > 404
func (this *ServiceRouter)doHandle(w http.ResponseWriter, r *http.Request) {
	// 如果是url全匹配, 则直接执行hand函数
    if h, ok := this.urlHandlersMap[r.URL.Path]; ok {
    	if this.isDebug {
			fmt.Println("URL.Handler: ", r.URL.Path)
		}
        h(w, r); return
    }else{
    	// 如果是url正则检查, 则需要检查正则, 正则为':'后面的字符
    	for _, key := range this.regexpHandlersIndex{
    		_SymbolIndex := strings.Index(key, ":")
    		if _SymbolIndex == -1 {
	    		continue
    		}
    		_BaseUrl := key[:_SymbolIndex]
    		if !strings.HasPrefix(r.URL.Path, _BaseUrl) {
	    		continue
    		}
	    	if ok, _ := regexp.MatchString(_BaseUrl+key[_SymbolIndex+1:], r.URL.Path); ok {
	    		if this.isDebug {
					fmt.Println("URL.Handler.Regexp: ", key)
				}
		    	this.regexpHandlersMap[key](w, r); 
		    	return
	    	}
    	
    	}
    }
	// 没有注册的地址, 使用默认处理器
	if this.defaultHandler != nil {
		this.defaultHandler(w, r)	
	}else{
    	w.WriteHeader(http.StatusNotFound)
	}
}
func (this *ServiceRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 处理前进行过滤处理
	if this.isDebug {
		fmt.Println("URL.Path: ", r.URL.Path)
	}
	// 1.1 检擦是否有指定路径的全路径匹配过滤器设定, 优先处理
	if nil != this.urlFiltersMap {
		if h, exist := this.urlFiltersMap[r.URL.Path]; exist {
			if this.isDebug {
				fmt.Println("URL.Filter: ", r.URL.Path)
			}
			h(w, r, func( ){
				this.doHandle(w, r)
			}); return
		}
	}
	// 1.2 检擦是否有指定路径的正则匹配过滤器设定, 优先处理
	if nil != this.regexpFiltersIndex && len(this.regexpFiltersIndex) > 0 {
		for _, key := range this.regexpFiltersIndex{
			_SymbolIndex := strings.Index(key, ":")
    		if _SymbolIndex == -1 {
	    		continue
    		}
    		_BaseUrl := key[:_SymbolIndex]
    		if !strings.HasPrefix(r.URL.Path, _BaseUrl) {
	    		continue
    		}
			
	    	if ok, _ := regexp.MatchString(key[:_SymbolIndex]+key[_SymbolIndex+1:], r.URL.Path); ok {
	    		if this.isDebug {
					fmt.Println("URL.Filter.Regexp: ", key)
				}
	    		this.regexpFiltersMap[key](w, r, func( ){
					this.doHandle(w, r)
				}); return
	    	}
		}
	}
	// 2. 检擦是否有全局过滤器存在, 如果有则执行它
	if nil != this.globalFileter {
		this.globalFileter(w, r, func( ){
			this.doHandle(w, r)
		}); return
	}
	// 3. 啥也没有设定
	this.doHandle(w, r)
}
// 清空路由表
func (this *ServiceRouter) ClearHandlersMap( ){
	this.urlHandlersMap    = make(map[string]HandlersFunc)
	this.regexpHandlersMap = make(map[string]HandlersFunc)
	this.regexpHandlersIndex = make([]string, 0, 0)
}
// 是否输出url
func (this *ServiceRouter) SetDebug( isDebug bool ){
	this.isDebug = isDebug
}
// 设置默认相应函数, 当无匹配时触发
func (this *ServiceRouter) SetDefaultHandler( defaultHandler HandlersFunc ){
	this.defaultHandler = defaultHandler
}
// 设置全局过滤器, 设置后, 如果不调用next函数则不进行下一步处理
// type FilterFunc func(http.ResponseWriter, *http.Request, func( ))
func (this *ServiceRouter)SetGlobalFilter( globalFilter FilterFunc ){
	this.globalFileter = globalFilter
}
// 设置url过滤器, 设置后, 如果不调用next函数则不进行下一步处理
// 过滤器有优先调用权, 正则匹配路径有先后顺序
// type FilterFunc func(http.ResponseWriter, *http.Request, func( ))
func (this *ServiceRouter)AddUrlFilter( url string, filter FilterFunc ){
	if len(url) == 0 {
		return
	}
	if nil == this.urlFiltersMap {
		this.urlFiltersMap = make(map[string]FilterFunc)
	}
	if nil == this.regexpFiltersMap {
		this.regexpFiltersMap = make(map[string]FilterFunc)
	}
	if nil == this.regexpFiltersIndex {
		this.regexpFiltersIndex = make([]string,0,0)
	}
	if strings.Index(url, ":") > -1 {
		this.regexpFiltersMap[url] = filter
		this.regexpFiltersIndex = append(this.regexpFiltersIndex, url)
	}else{
		this.urlFiltersMap[url] = filter
	}
}
// 删除filter索引
func (this *ServiceRouter) removeFilterIndex(url string){
	if len(url) > 0 {
		for i, key := range this.regexpFiltersIndex{
			if key == url {
				this.regexpFiltersIndex = append(this.regexpFiltersIndex[:i], this.regexpFiltersIndex[i+i:]...)
				break
			}
		}
	}
}
// 删除一个过滤器
func (this *ServiceRouter) RemoveFilter(url string){
	if len(url) == 0 {
		return
	}
	if nil != this.regexpHandlersMap {
		if _, ok := this.regexpFiltersMap[url]; ok {
			delete(this.regexpFiltersMap, url)
			this.removeFilterIndex(url)
		}
	}
	if nil != this.urlFiltersMap {
		if _, ok := this.urlFiltersMap[url]; ok {
			delete(this.urlFiltersMap, url)
		}
	}
}
// 构建urlmap, 全匹配和正则匹配分开存放, 正则表达式以':'符号开始, 如: /upload/:\S+
func (this *ServiceRouter) AddHandlers( handlersMap map[string]HandlersFunc ){
	if len(handlersMap) == 0 {
		return
	}
	if nil == this.regexpHandlersMap {
		this.regexpHandlersMap = make(map[string]HandlersFunc)
		this.regexpHandlersIndex = make([]string, 0, 0)
	}
	if nil == this.urlHandlersMap {
		this.urlHandlersMap = make(map[string]HandlersFunc)
	}	
	for key, val := range handlersMap {
		if strings.Index(key, ":") > -1 {
			this.regexpHandlersMap[key] = val
			this.regexpHandlersIndex = append(this.regexpHandlersIndex, key)
		}else{
			this.urlHandlersMap[key] = val
		}
	}
}
// 构建urlmap, 全匹配和正则匹配分开存放, 正则表达式以':'符号开始, 如: /upload/:\S+
func (this *ServiceRouter) AddHandler(url string, handler HandlersFunc ){
	if len(url) == 0 {
		return
	}
	if nil == this.regexpHandlersMap {
		this.regexpHandlersMap = make(map[string]HandlersFunc)
		this.regexpHandlersIndex = make([]string, 0, 0)
	}
	if nil == this.urlHandlersMap {
		this.urlHandlersMap = make(map[string]HandlersFunc)
	}
	if strings.Index(url, ":") > -1 {
		this.regexpHandlersMap[url] = handler
		this.regexpHandlersIndex = append(this.regexpHandlersIndex, url)
	}else{
		this.urlHandlersMap[url] = handler
	}
}
// 删除handler索引
func (this *ServiceRouter) removeHandlerIndex(url string){
	if len(url) > 0 {
		for i, key := range this.regexpHandlersIndex{
			if key == url {
				this.regexpHandlersIndex = append(this.regexpHandlersIndex[:i], this.regexpHandlersIndex[i+i:]...)
				break
			}
		}
	}
}
// 删除一个路由表
func (this *ServiceRouter) RemoveHandler(url string){
	if len(url) == 0 {
		return
	}
	if nil != this.regexpHandlersMap {
		if _, ok := this.regexpHandlersMap[url]; ok {
			delete(this.regexpHandlersMap, url)
			this.removeHandlerIndex(url)
		}
	}
	if nil != this.urlHandlersMap {
		if _, ok := this.urlHandlersMap[url]; ok {
			delete(this.urlHandlersMap, url)
		}
	}
}
// 注册路由,并启动服务 
// 此函数为ServiceRouter简版的使用方式, 可以根据ServiceRouter自己实现
func StartService(addr string, handlersMap map[string]HandlersFunc, defaultHandler HandlersFunc) error{
	if handlersMap == nil || len(handlersMap) == 0 {
		return errors.New("handlersMap is nil")
	}
	if addr == "" {
		return errors.New("addr is nil or empty")
	}
	_router := &ServiceRouter{}
	_router.isDebug = true
	_router.defaultHandler = defaultHandler
	_router.AddHandlers(handlersMap)
	
	server := &http.Server{
        Addr:           addr,
        Handler:        _router,
        ReadTimeout:    0,
        WriteTimeout:   0,
        MaxHeaderBytes: 1 << 20,
    }
	return server.ListenAndServe()
}