package gomodule
/**
 *@description 配置工具
 *@author	wupeng364@outlook.com
*/
import (
	"os"
	"encoding/json"
	"strings"
	"errors"
	"path/filepath"
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
	if isFile(cfg) {
		gmc.configPath  = cfg
	}else{
		err = writeFileAsJson(cfg, make(map[string]interface{}))
		if err == nil {
			gmc.configPath  = cfg
		}
	}
	if err != nil { panic(err); }
	
	// Json to map
	gmc.moduleConfig = make( map[string]interface{} )
	err = readFileAsJson(cfg, &gmc.moduleConfig)
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
			err := writeFileAsJson(gmc.configPath, gmc.moduleConfig )
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

// 读取Json文件
func readFileAsJson( path string, v interface{} ) error{
	if len(path) == 0 {
		return &os.PathError{"ReadFile", "", os.ErrNotExist}
	}
	fp, err := os.OpenFile(path, os.O_RDONLY, 0755)
    defer fp.Close()
    
    if err == nil {
    	st, err_st := fp.Stat( )
    	if err == nil{
	        data := make([]byte, st.Size( ))
			_, err = fp.Read(data)
			if err == nil {
				return json.Unmarshal(data, v)
			}
    	}else{
    		err = err_st
    	}
    }
    return err
}
// 写入Json文件
func writeFileAsJson( path string, v interface{} ) error{
	if len(path) == 0 {
		return &os.PathError{"WriteFile", "", os.ErrNotExist}
	}
	fp, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
    defer fp.Close()
    
    if err == nil {
    	data, err := json.Marshal(v)
    	if err == nil {
    		_, err := fp.Write(data)
    		return err
    	}else{
    		return err
    	}
    }else{
	    return err
    }
}
// 是否是文件
func isFile(path string)bool{
	_stat, _err := os.Stat(path)
    if _err == nil  {
        return !_stat.IsDir()
    }
    return false
}