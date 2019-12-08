package moduleloader
/**
 *@description 模块加载器 [模块装载, 指针管理, 反射调用]
 *@author	wupeng364@outlook.com
*/
import (
	"fmt"
	"time"
	"strconv"
	"reflect"
	"errors"
	"strings"
	"path/filepath"
)

const(
	moduleVersionPrec  = 2     // 版本信息保留小数位数
	defaultCfgId       = "app" // 默认配置文件名称
	defaultCfgPath     = "./config/{name}.json" // 默认配置保存位置
)

// 模块配置项
type Opts struct{
	Name    string                  // 模块ID
	Version float64                 // 模块版本
	Description string              // 模块描述
	OnReady   func(mctx *Loader)    // 每次加载模块开始之前执行
	OnSetup   func( )               // 模块安装, 一个模块只初始化一次
	OnUpdate  func( )               // 模块升级, 一个版本执行一次
	OnInit    func( )               // 每次模块安装、升级后执行一次
	OnDestroy func( )               // 系统执行销毁时执行
}
// 模块模板, 实现这个接口便可加载
type Template interface{
	ModuleOpts( )( Opts )
}

// 配置接口
type cfgInterface interface{
	InitConfig(configpath string) error
	GetConfig(key string)string
	GetConfigs(key string)interface{}
	SetConfig(key string, value string)error
}
// 模块加载器对象
type Loader struct{
	name string                      // 加载器实例名称
	modules map[string]interface{}   // 模块Map表
	configs map[string]cfgInterface  // 配置模块
}
// 函数执行后的返回值, 暂时不封装
type Returns []reflect.Value

// 实例一个加载器对象
func New( name string ) *Loader{
	
	res := &Loader{
		configs: make(map[string]cfgInterface),
		modules: make(map[string]interface{}),
	}
	if 0 == len(name) {
		res.name = defaultCfgId
	}else{
		res.name = name
	}
	// 初始化模块版本配置
	res.RegConfig(res.getDefaultCfgName(), nil)
	res.RegConfig(res.getLoaderCfgName(), nil)
	
	return res
}

// 初始化模块 - DO Setup -> Check Ver -> Do Init
func (this *Loader)Loads( mts ...Template ){
	for _, mt := range mts{
		this.Load(mt)
	}
}

// 初始化模块 - DO Setup -> Check Ver -> Do Init
func (this *Loader)Load( mt Template ){
	opts := mt.ModuleOpts( )
	fmt.Printf(">Loading %s(%s)[%p] start \r\n", opts.Name, opts.Description, mt)
	// DO Ready
	this.doReady(opts)
	// DO Setup
	this.doSetup(opts)
	// Check Ver
	this.doCheckVersion(opts)
	// Do Init
	this.doInit(opts)
	// Load End
	this.doEnd(opts, mt )
	fmt.Printf(">Loading %s complete \r\n", opts.Name)
	
}

// >---------------- public --------------------<
// 获取实例的名字
func (this *Loader)GetLoaderName( ) string{
	return this.name
}
// 模块调用, 返回 []reflect.Value 
// 返回值暂时无法处理
func (this *Loader)Invoke(mId string, method string, params ...interface{} )Returns{
	if module, ok := this.modules[mId]; ok {
		val := reflect.ValueOf(module)
		fun := val.MethodByName(method)
		fmt.Printf( "> Invoke: "+mId+"."+method+", %v, %+v \r\n", fun, &fun )
		args := make([]reflect.Value, len(params))
		for i, temp := range params{
			args[i] = reflect.ValueOf(temp)
		}
		return fun.Call(args)
	}else{
		panic(errors.New("module not find: "+ mId))
	}
}
//  根据模块ID获取模块指针记录, 可以获取一个已经实例化的模块
func (this *Loader)GetModuleById( mId string )(val interface{}, ok bool){
	if v, ok := this.modules[mId]; ok {
		return v, true
	}
	return nil, false
}

// 根据模板对象获取模块指针记录, 可以获取一个已经实例化的模块
func (this *Loader)GetModuleByTemplate( mt Template ) interface{}{
	mopts := mt.ModuleOpts( )
	if val, ok := this.GetModuleById( mopts.Name ); ok {
		return val
	}
	panic(errors.New("module not find: "+ mopts.Name+"["+mopts.Description+"]"))
}
// 获取配置-默认配置文件
func (this *Loader)GetConfig( key string ) string{
	if cfg, ok := this.configs[this.getDefaultCfgName()]; ok {
		return cfg.GetConfig(key)
	}
	return ""
}
// 获取配置-默认配置文件
func (this *Loader)GetConfigs( key string ) interface{}{
	if cfg, ok := this.configs[this.getDefaultCfgName()]; ok {
		return cfg.GetConfigs(key)
	}
	return ""
}
// 设置配置-默认配置文件
func (this *Loader)SetConfig( key, val string ) error{
	if cfg, ok := this.configs[this.getDefaultCfgName()]; ok {
		return cfg.SetConfig(key, val)
	}
	return errors.New("config id not reg: "+this.getDefaultCfgName())
}
// 获取配置-根据配置名
func (this *Loader)GetConfigByName( id, key string ) string{
	if cfg, ok := this.configs[id]; ok {
		return cfg.GetConfig(key)
	}
	return ""
}
// 获取配置-根据配置名
func (this *Loader)GetConfigsByName( id, key string ) interface{}{
	if cfg, ok := this.configs[id]; ok {
		return cfg.GetConfigs(key)
	}
	return ""
}
// 设置配置-根据配置名
func (this *Loader)SetConfigByName( id, key, val string ) error{
	if cfg, ok := this.configs[id]; ok {
		return cfg.SetConfig(key, val)
	}
	return errors.New("config id not reg: "+id)
}
// 注册一个配置器
func (this *Loader)RegConfig( id string, cfg cfgInterface){
	if nil != cfg {
		this.configs[id] = cfg
	}else{
		this.configs[id] = this.GetDefaultCfgInterface( id )
	}
}
// 获取默认的Json配置器
func (this *Loader)GetDefaultCfgInterface( id string ) cfgInterface{
	if 0 == len(id) {
		panic(errors.New("config id is empty"))
	}
	cfgPath, err := filepath.Abs(strings.Replace(defaultCfgPath, "{name}", id, -1))
	if nil != err {
		panic( err )
	}
	cfg := &(jsoncfg{})
	err = cfg.InitConfig( cfgPath )
	if nil != err {
		panic( err )
	}
	return cfg
}
// >---------------- private --------------------<
// 记录模块信息的配置文件名
func (this *Loader)getLoaderCfgName( ) string{
	return this.name+".mloader"
}
// 获取默认配置文件名
func (this *Loader)getDefaultCfgName( ) string{
	return this.name+".default"
}
// 获取模块版本号
func (this *Loader)getInstalledVersion( opts Opts )(string){
	return this.GetConfigByName(this.getLoaderCfgName(), opts.Name+".SetupVer")
}
// 设置模块版本号 - 模块保留小数两位
func (this *Loader)setVersion( opts Opts ){
	this.SetConfigByName(this.getLoaderCfgName(), opts.Name+".SetupVer", formatVersion(opts) )
	this.SetConfigByName(this.getLoaderCfgName(), opts.Name+".SetupDate", strconv.FormatInt(time.Now( ).UnixNano( ), 10))
}
// 模块准备
func (this *Loader)doReady(opts Opts){
	if nil != opts.OnReady {
		fmt.Printf("  > On ready load \r\n")
		opts.OnReady(this)
	}
}
// 模块安装
func (this *Loader)doSetup( opts Opts ){
	if len(this.getInstalledVersion(opts)) == 0 {
		if nil != opts.OnSetup {
			fmt.Printf("  > On setup module \r\n")
			opts.OnSetup( )
		}
		
		this.setVersion(opts)
	}
}
// 模块升级
func (this *Loader)doCheckVersion( opts Opts ){
	setupVerStr := formatVersion(opts)
	_historyVer := this.getInstalledVersion(opts)
	if _historyVer != setupVerStr {
		if nil != opts.OnUpdate {
			fmt.Printf("  > On update version \r\n")
			opts.OnUpdate( )
		}
		
		this.setVersion(opts)
	}
}
// 模块初始化
func (this *Loader)doInit( opts Opts ){
	if nil != opts.OnInit {
		fmt.Printf("  > On init module \r\n")
		opts.OnInit( )
	}
}
// 模块加载结束
func (this *Loader)doEnd( opts Opts, mt Template){
	this.modules[opts.Name] = mt
}

// >-------------------- else ---------------------<
// 格式模块版本号 float64 => string
func formatVersion( opts Opts )(string){
	return strconv.FormatFloat(opts.Version, 'f', 2, 64)
}