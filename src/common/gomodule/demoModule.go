package gomodule
/**
 *@description 示例模板
 *@author	wupeng364@outlook.com
*/
import (
	"fmt"
	"strconv"
)
type DemoModel struct{}

// 返回模块信息
func (dm *DemoModel)MInfo( )(*ModelInfo)	{
	return &ModelInfo{ dm, "DemoModel", 1.0, "测试模板" }
}
// 模块安装, 一个模块只初始化一次
func (dm *DemoModel)MSetup( ) {
	
}
// 模块升级, 一个版本执行一次
func (dm *DemoModel)MUpdate( ) {
	
}

// 每次启动加载模块执行一次
func (dm *DemoModel)OnMInit( getPointer func(m interface{})interface{} ) {
	fmt.Println("DemoModel init start")
}
// 系统执行销毁时执行
func (dm *DemoModel)OnMDestroy( ) {
	
}

// ==============================================================================================
func (dm *DemoModel) SayHole( text string, count int) interface{}{
	for i:=0; i<count; i++ { 
		fmt.Printf("%d-"+text+"\n\r", i)
	}
	return map[string]string{"text":text, "count": strconv.Itoa(count), }
}