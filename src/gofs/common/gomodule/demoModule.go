package gomodule
/**
 *@description 示例模板
 *@author	wupeng364@outlook.com
*/
import (
	"fmt"
	"strconv"
)
type DemoModule struct{}

// 返回模块信息
func (dm *DemoModule)MInfo( )(*ModuleInfo)	{
	return &ModuleInfo{ dm, "DemoModule", 1.0, "测试模板" }
}
// 模块安装, 一个模块只初始化一次
func (dm *DemoModule)OnMSetup( ref ReferenceModule ) {
	
}
// 模块升级, 一个版本执行一次
func (dm *DemoModule)OnMUpdate( ref ReferenceModule ) {
	
}

// 每次启动加载模块执行一次
func (dm *DemoModule)OnMInit( ref ReferenceModule ) {
	fmt.Println("DemoModule init start")
}
// 系统执行销毁时执行
func (dm *DemoModule)OnMDestroy( ref ReferenceModule ) {
	
}

// ==============================================================================================
func (dm *DemoModule) SayHole( text string, count int) interface{}{
	for i:=0; i<count; i++ { 
		fmt.Printf("%d-"+text+"\n\r", i)
	}
	return map[string]string{"text":text, "count": strconv.Itoa(count), }
}