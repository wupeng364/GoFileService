package moduleloader
/**
 *@description 配置类, json格式存储
 *@author	wupeng364@outlook.com
*/
import (
	"os"
	"errors"
	"strings"
	"path/filepath"
	"encoding/json"
)

const(
	cfgPath = "/config.json"
)
type jsoncfg struct{
	jsonConfig map[string]interface{}
	configPath string
}

func (this *jsoncfg) InitConfig( configpath string ) error{
	var cfg string
	var err error

	// 配置文件位置
	if len(configpath) == 0 {
		// app运行路径
		appwd, err := os.Getwd( )
		if nil != err {
			return err
		}
		cfg = filepath.Join(appwd, cfgPath)
	}else{
		cfg = configpath
	}
	
	if isFile(cfg) {
		this.configPath  = cfg
	}else{
		err = writeFileAsJson(cfg, make(map[string]interface{}))
		if nil != err {
			return err
		}
		if err == nil {
			this.configPath  = cfg
		}
	}
	
	// Json to map
	this.jsonConfig = make( map[string]interface{} )
	return readFileAsJson(cfg, &this.jsonConfig)
}
func (this *jsoncfg) GetConfig(key string)(res string){
	res, ok := this.GetConfigs(key).(string)
	if ok {
		return res
	}
	return ""
}

func (this *jsoncfg) GetConfigs(key string)(res interface{}){
	if len( key ) == 0 || this.jsonConfig == nil || len(this.jsonConfig) == 0 {
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
				if tp, ok := this.jsonConfig[keys[i]]; ok {
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
			if tp, ok := this.jsonConfig[keys[i]]; ok {
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

func (this *jsoncfg) SetConfig(key string, value string) error{
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
				this.jsonConfig[keys[i]] = value
			}else if temp != nil {
				temp.(map[string]interface{})[keys[i]] = value
			}
			// fmt.Println( this.jsonConfig )
			err := writeFileAsJson(this.configPath, this.jsonConfig )
			return err
		}
		
		// 
		var _temp interface{}
		if temp == nil { // first
			if tp, ok := this.jsonConfig[keys[i]]; ok {
				_temp =  tp
			}else{
				_temp = make(map[string]interface{})
				this.jsonConfig[keys[i]] = _temp
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
	parent := getParentPath(path)
	if !isDir( parent ) {
		os.MkdirAll(parent, 0755)
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
// 是否是文件夹
func isDir(path string)bool{
	_stat, _err := os.Stat(path)
    if _err == nil  {
        return _stat.IsDir()
    }
    return false
}
// 范围最后一个'/'前的文字
func getParentPath( s_path string )string{
	if strings.Index(s_path, "\\") > -1 {
		return s_path[:strings.LastIndex(s_path, "\\")]
	}else{
		return s_path[:strings.LastIndex(s_path, "/")]
	}
}