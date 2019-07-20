package config
/**
 *@description 配置文件模块
 *系统的配置模块, 实现cfInterface接口
 *@author	wupeng364@outlook.com
*/
import (
	//"fmt"
	"os"
	"strings"
	"errors"
	"path/filepath"
	"common/gomodule"
	"common/filetools"
)
const(
	cfgPath = "/config/config.json"
)

type ConfigModule struct{
	appConfig map[string]interface{}
	appWorkPath string
	configPath string
}

// 返回模块信息
func (cm *ConfigModule)MInfo( )(*gomodule.ModuleInfo)	{
	return &gomodule.ModuleInfo{
		cm,
		"ConfigModule",
		1.0,
		"系统配置模块",
	}
}

// 模块安装, 一个模块只初始化一次
func (cm *ConfigModule)OnMSetup( ref gomodule.ReferenceModule ) {
	
}
// 模块升级, 一个版本执行一次
func (cm *ConfigModule)OnMUpdate( ref gomodule.ReferenceModule ) {
	
}

// 每次启动加载模块执行一次
func (cm *ConfigModule)OnMInit( ref gomodule.ReferenceModule ) {
	cm.InitConfig( )
	// fmt.Print("appConfig: ", cm.appConfig)
}

// 系统执行销毁时执行
func (cm *ConfigModule)OnMDestroy( ref gomodule.ReferenceModule ) {
	
}

// ==============================================================================================
func (cm *ConfigModule) InitConfig( ) error{
	var appwd, cfg string
	var err error

	// app运行路径
	appwd, err = os.Getwd( )
	if( err == nil ){
		if len(appwd) > 0{
			cm.appWorkPath = appwd
		}else{
			err = errors.New("Getwd is nil or empty")
		}
	}
	if err != nil { panic(err); }
	
	// 配置文件位置
	cfg = filepath.Join(appwd, cfgPath)
	if filetools.IsFile(cfg) {
		cm.configPath  = cfg
	}else{
		err = filetools.WriteFileAsJson(cfg, make(map[string]interface{}))
		if err == nil {
			cm.configPath  = cfg
		}
	}
	if err != nil { panic(err); }
	
	// Json to map
	cm.appConfig = make( map[string]interface{} )
	err = filetools.ReadFileAsJson(cfg, &cm.appConfig)
	if err != nil { panic(err); }
	//fmt.Println(cm.appConfig)
	return nil
}
func (cm *ConfigModule) GetConfig(key string)(res string){
	return cm.GetConfigs(key).(string)
}

func (cm *ConfigModule) GetConfigs(key string)(res interface{}){
	// fmt.Printf("\r==>%p", cm)
	if len( key ) == 0 || cm.appConfig == nil || len(cm.appConfig) == 0 {
		return
	}
	keys := strings.Split(key, ".")
	if keys == nil {
		return
	}
	var temp interface{}
	keyLength := len(keys)
	for i :=0; i<keyLength; i++ {
		// last key
		if i == keyLength-1 {
			if i == 0 {
				if tp, ok := cm.appConfig[keys[i]]; ok {
					res = tp
				}
			}else if temp != nil {
				if tp, ok := temp.(map[string]interface{})[keys[i]]; ok {
					res = tp
				}
			}
			return
		}
		
		// 
		var _temp interface{}
		if temp == nil { // first
			if tp, ok := cm.appConfig[keys[i]]; ok {
				_temp =  tp
			}
		}else{ // 
			if tp, ok := temp.(map[string]interface{})[keys[i]]; ok {
				_temp =  tp
			}
		}
		
		// find
		if _temp != nil {
			temp = _temp;
		}else{
			return 
		}
	}
	return
}
func (cm *ConfigModule) SetConfig(key string, value string)error{
	if len(key) ==0 || len(value)==0{
		return errors.New("key or value is empty")
	}
	keys := strings.Split(key, ".")
	keyLength := len(keys)
	var temp interface{}
	for i :=0; i<keyLength; i++ {
		// last key
		if i == keyLength-1 {
			if i == 0 {
				cm.appConfig[keys[i]] = value
			}else if temp != nil {
				temp.(map[string]interface{})[keys[i]] = value
			}
			// fmt.Println( cm.moduleConfig )
			return filetools.WriteFileAsJson(cm.configPath, cm.appConfig )
		}
		
		// 
		var _temp interface{}
		if temp == nil { // first
			if tp, ok := cm.appConfig[keys[i]]; ok {
				_temp =  tp
			}else{
				_temp = make(map[string]interface{})
				cm.appConfig[keys[i]] = _temp
			}
		}else{ // 
			if tp, ok := temp.(map[string]interface{})[keys[i]]; ok {
				_temp =  tp
			}else{
				_temp = make(map[string]interface{})
				temp.(map[string]interface{})[keys[i]] = _temp
			}
		}
		
		// find
		if _temp != nil {
			temp = _temp;
		}
	}
	return nil
}