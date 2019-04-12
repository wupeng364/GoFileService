package gomodule
/**
 *@description 配置工具
 *@author	wupeng364@outlook.com
*/
import (
	"os"
	"strings"
	"errors"
	"path/filepath"
	"common/filetools"
)
const(
	cfgPath = "/config/GoModuleConfig.json"
)

type GoModuleConfig struct{
	moduleConfig map[string]interface{}
	appWorkPath string
	configPath string
}

func (gmc *GoModuleConfig) InitConfig( ){
	var appwd, cfg string
	var err error

	// app运行路径
	appwd, err = os.Getwd( )
	if( err == nil ){
		if len(appwd) > 0{
			gmc.appWorkPath = appwd
		}else{
			err = errors.New("Getwd is nil or empty")
		}
	}
	if err != nil { panic(err); }
	
	// 配置文件位置
	cfg = filepath.Join(appwd, cfgPath)
	if filetools.IsFile(cfg) {
		gmc.configPath  = cfg
	}else{
		err = filetools.WriteFileAsJson(cfg, make(map[string]interface{}))
		if err == nil {
			gmc.configPath  = cfg
		}
	}
	if err != nil { panic(err); }
	
	// Json to map
	gmc.moduleConfig = make( map[string]interface{} )
	err = filetools.ReadFileAsJson(cfg, &gmc.moduleConfig)
	if err != nil { panic(err); }
	// fmt.Printf("%+v\r", gmc.moduleConfig)
}

func (gmc *GoModuleConfig) GetConfig(key string)(res string){
	if len( key ) == 0 || gmc.moduleConfig == nil || len(gmc.moduleConfig) == 0 {
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
				if tp, ok := gmc.moduleConfig[keys[i]]; ok {
					res = tp.(string)
				}
			}else if temp != nil {
				if tp, ok := temp.(map[string]interface{})[keys[i]]; ok {
					res = tp.(string)
				}
			}
			return
		}
		
		// 
		var _temp interface{}
		if temp == nil { // first
			if tp, ok := gmc.moduleConfig[keys[i]]; ok {
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

func (gmc *GoModuleConfig) SetConfig(key string, value string) error{
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
				gmc.moduleConfig[keys[i]] = value
			}else if temp != nil {
				temp.(map[string]interface{})[keys[i]] = value
			}
			// fmt.Println( gmc.moduleConfig )
			err := filetools.WriteFileAsJson(gmc.configPath, gmc.moduleConfig )
			return err
		}
		
		// 
		var _temp interface{}
		if temp == nil { // first
			if tp, ok := gmc.moduleConfig[keys[i]]; ok {
				_temp =  tp
			}else{
				_temp = make(map[string]interface{})
				gmc.moduleConfig[keys[i]] = _temp
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