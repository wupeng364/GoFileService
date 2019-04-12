/**
 *@description 一些常量的结构
 *@author	wupeng364@outlook.com
*/
package gomodule

import(
	"reflect"
)

const(
	moduleVersionPrec = 2
)

var models  = make(map[string]interface{}) // 模块Map表
var configs = &GoModuleConfig{ }           // 配置模块

// 函数执行后的返回值, 暂时不封装
type Returns []reflect.Value

// 模块模板
type ModuleTemplate interface{
	MInfo( )(*ModelInfo)			// 返回模块信息
	MSetup( )					    // 模块安装, 一个模块只初始化一次
	MUpdate( )					    // 模块升级, 一个版本执行一次
	
	OnMInit( func(m interface{})interface{} )   // 每次启动加载模块执行一次
	OnMDestroy( )								// 系统执行销毁时执行
}

// 模块的描述
type ModelInfo struct{
	Pointer interface{}
	Name    string
	Version float64
	Description string
}

// init
func init( ){
	configs.InitConfig( )    // module config
}