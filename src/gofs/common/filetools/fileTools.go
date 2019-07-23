package filetools

/**
 *@description 文件工具类
 *@author	wupeng364@outlook.com
*/
import (
	"fmt"
	"io"
	"os"
	"time"
	"encoding/json"
	"path/filepath"
	"strings"
)
// =========================== type 
type CopyCallback func(count int, s_src, s_dst string, err error)error
type MoveCallback func(count int, s_src, s_dst string, err error)error

// 获取文件信息对象
func GetFileInfo(path string) (os.FileInfo, error){
	return  os.Stat(path)
}
// 获取文件信息对象
func OpenFile(path string) (*os.File, error){
	return  os.Open(path)
}
// 获取文件信息对象
func GetWriter(path string) (*os.File, error){
	return os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
}
// 文件/夹是否存在
func IsExist(path string)bool{
	 _, _err := GetFileInfo(path)
    return _err == nil
}
// 是否是文件夹
func IsDir(path string)bool{
	_stat, _err := GetFileInfo(path)
    if _err == nil {
        return _stat.IsDir( )
    }
    return false
}
// 创建文件夹-多级
func MkdirAll(path string)error{
	 return os.MkdirAll(path, os.ModePerm)
}// 创建文件夹
func Mkdir(path string)error{
	 return os.Mkdir(path, os.ModePerm)
}
// 是否是文件
func IsFile(path string)bool{
	_stat, _err := GetFileInfo(path)
    if _err == nil  {
        return !_stat.IsDir()
    }
    return false
}
// 获取一级子目录(无序), 无路径, 只是文件|目录名字
func GetDirList(path string)[]string{
	f, err := OpenFile(path)
	defer	f.Close()
	if err != nil {
		return nil
	}
	
	list, err := f.Readdir(-1)
	if err != nil {
		return nil
	}
	result := make([]string, 0, len(list))
	for _, fi := range list{
		result = append(result, fi.Name() )
	}
	return result
}
// 获取文件大小
func GetFileSize(path string)int64{
	f, err := OpenFile(path)
	if err != nil {
		return 0
	}
	f_info, err_stat := f.Stat()
	f.Close( )
	if err_stat != nil {
		return 0
	}
	return f_info.Size( )
}
// 获取创建时间(时间戳/S)
func GetCreateTime(path string)time.Time{
	
	return GetModifyTime( path )
}
// 获取修改时间(时间戳/S)
func GetModifyTime(path string)time.Time{
	f, err := OpenFile(path)
	if err != nil {
		return time.Time{}
	}
	f_info, err_stat := f.Stat()
	f.Close()
	if err_stat != nil {
		return time.Time{}
	}
	return f_info.ModTime()
}
// 删除文件
func RemoveFile(file string) error{
   if !IsExist(file) {
      return PathError_NotExist("RemoveFile", file)
   }
   return os.Remove(file)
}

// 删除文件
func RemoveAll(file string) error{
   if !IsExist(file) {
      return PathError_NotExist("RemoveAll", file)
   }
   return os.RemoveAll(file)
}
// 重命名 
func Rename(old, newName string ) error{
	_path := filepath.Clean(old)
	if strings.Index(_path, "\\") > -1 {
	    return os.Rename(old, old[:strings.LastIndex(old, "\\")+1]+newName)
	}
	return os.Rename(old, old[:strings.LastIndex(old, "/")+1]+newName)
}
// 移动文件|文件夹 - 可跨分区移动
func MoveDir(src, dst string, replace, ignore bool, callback MoveCallback)error{
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)
	if src == dst && len(src) >0 {
		return nil
	}
	if !IsExist(src) {
	    return PathError_NotExist("MoveDir", src)
    }
	// 尝试本分区移动
	err := MoveFile(src, dst, replace, ignore, func(count int, s_src, s_dst string, mv_err error)error{
		if IsError_DifferentDiskDrive(mv_err) {
			return mv_err
		}
		return callback(count, s_src, s_dst, mv_err)
	})
	// 尝试跨分区移动
	if IsError_DifferentDiskDrive(err) {
		return MoveFileByCopying(src, dst, replace, ignore, callback)
	}else{
		return err
	}
}
// 移动文件|夹 - 如果存在的话就列表后逐一移动
func MoveFile(src, dst string, replace, ignore bool, callback MoveCallback)error{
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)
	if src == dst && len(src) >0 {
		return nil
	}
	if !IsExist(src) {
	    return PathError_NotExist("MoveFile", src)
    }
	if IsExist(dst) {
	 	// 如果是文件, 则需要判断是否覆盖|忽略
	 	if IsFile(dst) {
			if ignore {
				// 如果忽略则不返回错误
				return callback(1, src, dst, nil)
			} else if replace {
				// 如果覆盖, 则需要选择是否忽略错误继续
				err1 := os.Remove(dst)
				err  := callback(1, src, dst, err1)
				if err != nil {
					return err // 错误被正常返回, 说明是向终止操作
				}else if err1 != nil {
				 	return nil // 手动忽略了这个错误, 虽然不报错但是也不进行操作文件
			 	}
			}
			return callback(1, src, dst, os.Rename(src, dst))
			
		 }else{
		 	// 如果文件夹存在, 则处理里面的文件就行了
		 	list := GetDirList(src)
		 	for _, val := range list{
			 	err := MoveFile(src+"/"+val, dst+"/"+val, replace, ignore, callback)
			 	if err != nil {
					return err	// 返回终止信号
			 	}
		 	}
		 	if IsExist( src ){
			 	return callback(1, src, dst, os.Remove( src ))
		 	}
		 	
	 	}
	}else{ // 如果目标文件/夹不存在就直接移动
		return callback(1, src, dst, os.Rename(src, dst))
	}
	return nil
}
// 移动文件夹 - 跨分区-拷贝
func MoveFileByCopying(src, dst string, replace, ignore bool, callback MoveCallback) error{
	if src == dst && len(src) >0 {
		return nil
	}
	if IsFile(src) {
		var err1 error
		err := CopyFile(src, dst, replace, ignore)
		if err != nil {
			err1 = callback(1, src, dst, err )
		}
		if err1 != nil {
			return err1
		}else if err != nil {
			return nil
		}
		return os.Remove(src)
	}else{
		_err := CopyDir(src, dst, replace, ignore, func(count int, s_src, s_dst string, err error)error{
			if err == nil && IsFile(s_src){
				err = os.Remove(s_src)
			}
			return callback(count, s_src, s_dst, err )
		})
		// 最后的清理
		if _err == nil {
			_err = os.RemoveAll( src )
		}
		return _err
	}
}
// 复制文件
func CopyFile(src, dst string, replace, ignore bool)error{
	if src == dst && len(src) >0 {
		return nil
	}
    if !IsFile(src) {
	    return  PathError_NotExist("CopyFile", src)
    }
    if IsExist(dst) {
    	if replace {
	    	err := os.Remove(dst)
	    	if err != nil {
		    	return err
	    	}
    	}else if ignore {
	    	return nil
    	}else{
		    return PathError_Exist("CopyFile", dst)
    	}
    }
    r_src, err := OpenFile(src)
	defer r_src.Close( )
    if err != nil {
	    return err
    }
	w_dst, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, 0777)
	defer w_dst.Close()
	if err != nil{
		return err
	}
    _, err = io.Copy(w_dst, r_src)
    return err
}
// 复制文件夹
// 源路径, 目标路径, 是否覆盖, 是否跳过, 复制回调
// callback 内返回错误即可终止后续拷贝; 如果callback返回有错误, 该文件则为处理失败, 需要手动处理出错文件; 如果callback返回nil则继续往下拷贝
func CopyDir(src, dst string, replace, ignore bool, callback CopyCallback) error{
	if src == dst && len(src) >0 {
		return nil
	}
	countSuccess := 0
	if !IsExist(src) {
		return callback(countSuccess, src, dst, PathError_NotExist("CopyDir", src) )
	}
	if IsFile(dst) {
		return callback(countSuccess, src, dst, PathError_Exist("CopyDir", dst) )
	}
	dst = filepath.Clean(dst)
	src = filepath.Clean(src)
	len_src := len(src)
	return filepath.Walk(src, func(s string, f os.FileInfo, err error) error {
		d := dst+s[len_src:]
		if err == nil {
			if f.IsDir( ) {
				if !IsDir(d) {
					err = os.Mkdir(d, os.ModePerm)
				}
			}else{
				err = CopyFile(s, d, replace, ignore)
			}
		}
		countSuccess++
		return callback(countSuccess, s, d, err )
	})
}
// 读取Json文件
func ReadFileAsJson( path string, v interface{} ) error{
	if len(path) == 0 {
		return PathError_NotExist("ReadFile", "")
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
func WriteFileAsJson( path string, v interface{} ) error{
	if len(path) == 0 {
		return PathError_NotExist("WriteFile", "")
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

// ====================================== errors 
// 路径已经存在的错误
func PathError_Exist( op, dst string ) error {
	return &os.PathError{op, dst, os.ErrExist}
}
// 路径不存在的错误
func PathError_NotExist( op, src string ) error {
	return &os.PathError{op, src, os.ErrNotExist}
}
// 是否是目标位置已经存在的错误
func IsError_DestExist( err error ) bool{
	if nil == err{
		return false
	}
	var _err error
	switch err := err.(type) {
		case *os.PathError:
			_err = err.Err
	}
	if _err == nil {
		return false
	}
	return _err.Error( ) == os.ErrExist.Error( )
}
// 是否是目标位置不存在的错误
func IsError_DestNotExist( err error ) bool{
	var _err error
	switch err := err.(type) {
		case *os.PathError:
			_err = err.Err
	}
	if _err == nil {
		return false
	}
	return _err.Error( ) == os.ErrNotExist.Error( )
}
// 是否是跨磁盘错误
func IsError_DifferentDiskDrive( err error ) bool{
	var _err error
	switch err := err.(type) {
		case *os.LinkError:
			_err = err.Err
	}
	if _err == nil {
		return false
	}
	return _err.Error( ) == "The system cannot move the file to a different disk drive."
}

func say( ){
	fmt.Println("")
}