package gomodule
/**
 *@description 模块加载器 [模块装载, 指针管理, 反射调用]
 *@author	wupeng364@outlook.com
*/
import (
	"reflect"
	"errors"
	"fmt"
	"time"
	"strconv"
)



// 模块加载
func LoadModule( mt ModuleTemplate ){
	
	// load mothod
	mInfo := mt.MInfo()
	
	doRecordModule(mInfo, mt)
	fmt.Printf("> Loading %s(%s)[%p] Start \r\n", mInfo.Name, mInfo.Description, mt)
	
	fmt.Printf("  > Do Check Setup \r\n")
	doSetup( mInfo, mt.OnMSetup )
	
	fmt.Printf("  > Do Check Update \r\n")
	doUpdate( mInfo, mt.OnMSetup )
	
	fmt.Printf("  > Do Init function \r\n")
	doInit( mInfo, mt.OnMInit )
	fmt.Println("> Loading Complate")
	// fmt.Println("moduleMethods: ", moduleMethods )
}

// 获取某个模块
func GetModule( mId string )(val interface{}, ok bool){
	if v, ok := modules[mId]; ok {
		return v, true
	}
	return nil, false
}

// 获取模块指针记录, 可以获取一个已经实例化的模块
func GetModuleReference( mt ModuleTemplate ) interface{}{
	if val, ok := GetModule( mt.(ModuleTemplate).MInfo( ).Name ); ok {
		return val
	}
	mInfo := mt.MInfo()
	panic(errors.New("module not find: "+ mInfo.Name+"["+mInfo.Description+"]"))
}

// 模块调用, 返回值暂时无法处理
func Invoke(mId string, method string, params ...interface{} )Returns{
	if module, ok := modules[mId]; ok {
		val := reflect.ValueOf(module)
		fun := val.MethodByName(method)
		fmt.Printf( "   > Invoke: "+mId+"."+method+", %v, %+v \r\n", fun, &fun )
		args := make([]reflect.Value, len(params))
		for i, temp := range params{
			args[i] = reflect.ValueOf(temp)
		}
		return fun.Call(args)
	}else{
		panic(errors.New("module not find: "+ mId))
	}
}
// ==================================extends==================================<
// 反射接口返回的接口内是否是空的
func ValsIsNil( res Returns, index int, doErr func(err error) ){
	// fmt.Println("ValsIsNil",res[index].Type( ).String( ), res[index].String( ), res[index].IsValid())
	if res == nil || !res[index].IsValid( ) {
		err := errors.New("The value of the return value subscript "+strconv.Itoa(index)+" is nil")
		if doErr == nil{
			panic( err )
		}else{
			doErr( err )
		}
	}
}
// 判断低X个参数是不是不为空的错误类型
func ValsIsErr( res Returns, index int, doErr func(err error) ) bool{
	if res != nil && !res[index].IsNil( ){
		if res[index].Type( ).String( ) == "error" {
			err := res[index].Interface( ).(error)
			if doErr == nil{
				panic( err )
			}else{
				doErr( err );
			}
			return true
		}
	}
	return false
}

// ==================================private==================================<
// 记录模块地址
func doRecordModule( mi *ModuleInfo, mt ModuleTemplate ){
	modules[mi.Name] = mt
	// fmt.Println(modules)
}
// 模块安装
func doSetup( mi *ModuleInfo, mst func(ReferenceModule) ){
	setupVerKey := "modules."+mi.Name+".SetupVer"
	if len(configs.GetConfig(setupVerKey)) == 0 {
		mst( GetModuleReference )
		configs.SetConfig("modules."+mi.Name+".SetupDate", strconv.FormatInt(time.Now( ).UnixNano( ), 10))
		configs.SetConfig(setupVerKey, strconv.FormatFloat(mi.Version, 'f', moduleVersionPrec, 64) )
	}
}
// 模块升级
func doUpdate( mi *ModuleInfo, mst func(ReferenceModule) ){
	setupVerKey := "modules."+mi.Name+".SetupVer"
	setupVerStr := strconv.FormatFloat(mi.Version, 'f', moduleVersionPrec, 64)
	_historyVer := configs.GetConfig(setupVerKey)
	if _historyVer != setupVerStr {
		mst( GetModuleReference )
		configs.SetConfig(setupVerKey, setupVerStr )
	}
}
// 模块初始化
func doInit( mi *ModuleInfo, mst func(ReferenceModule) ){
	mst( GetModuleReference )
}