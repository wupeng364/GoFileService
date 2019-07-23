package filemanage

/**
 *@description 本地磁盘驱动
 *@author	wupeng364@outlook.com
*/
import (
	"path/filepath"
	"gofs/common/filetools"
	"errors"
	"strings"
	"io"
	"fmt"
)

type imp_local struct{
	MountInfo MtItems
	mtm *MountManager
}


/**
 * 文件是否存在
**/
func (il imp_local) IsExist( relativePath string ) (bool, error){
	_abs_p, _, err := getAbsolutePath(il.MountInfo, relativePath)
    return filetools.IsExist( _abs_p ), err
}
func (il imp_local) IsDir( relativePath string ) (bool, error){
	_abs_p, _, err := getAbsolutePath(il.MountInfo, relativePath)
    return filetools.IsDir( _abs_p ), err
}
func (il imp_local) IsFile( relativePath string ) (bool, error){
	_abs_p, _, err := getAbsolutePath(il.MountInfo, relativePath)
    return filetools.IsFile( _abs_p ), err
}
func (il imp_local) GetDirList( relativePath string ) ([]string, error){
	_abs_p, _rl_p, err := getAbsolutePath(il.MountInfo, relativePath)
	if err != nil {
		return make([]string, 0), err
	}
	ls := filetools.GetDirList( _abs_p )
	// 如果是挂载目录根目录, 需要处理 缓存目录
	if _rl_p == "/" {
		if ls != nil && len(ls) > 0 {
			res := make([]string, 0)
			for _, p := range ls {
				// 如果是挂载目录根目录, 忽略系统目录
				if sys_dir == p {
					continue
				}
				res = append(res, p)
			}
			return res, nil
		}
	}
	return ls, nil
}
func (il imp_local) GetFileSize( relativePath string ) (int64, error){
	_abs_p, _, err := getAbsolutePath(il.MountInfo, relativePath)
    return filetools.GetFileSize( _abs_p ), err
}
func (il imp_local) GetCreateTime( relativePath string ) (int64, error){
	_abs_p, _, err := getAbsolutePath(il.MountInfo, relativePath)
    return filetools.GetCreateTime( _abs_p ).UnixNano( ) / 1e6, err
}
func (il imp_local) GetModifyTime( relativePath string ) (int64, error){
	_abs_p, _, err := getAbsolutePath(il.MountInfo, relativePath)
    return filetools.GetModifyTime( _abs_p ).UnixNano( ) / 1e6, err
}
func (il imp_local) GetCreateUser( relativePath string ) (string, error){
	return "", nil
}
func (il imp_local) GetModifyUser( relativePath string ) (string, error){
	return "", nil
}
func (il imp_local) GetFileLatestVersion( relativePath string ) (string, error){
	return "", nil
}
func (il imp_local) GetFileVersionList( relativePath string ) ([]string, error){
	return []string{ }, nil
}
// 移动文件|夹 
func (il imp_local) DoMove(src string, dst string, replace, ignore bool, callback MoveCallback)error{
	if il.MountInfo.MtPath == src {
		return errors.New("Does not allow access: "+ src)
	}
	_abs_src, _, _err := getAbsolutePath(il.MountInfo, src)
	if nil != _err {
		return _err
	}
	if filepath.Clean(il.MountInfo.MtAddr) == _abs_src {
		return errors.New(src+" is mount root, cannot move")
	}
	// 目标位置驱动接口
	dstMountItem := il.mtm.getMountItem(dst)
	_abs_dst, _, err := getAbsolutePath(dstMountItem, dst)
	if nil != err {
		return err
	}
	switch dstMountItem.MtType {
		case mt_type_local:{ // 本地存储
			return filetools.MoveDir(_abs_src, _abs_dst, replace, ignore, func(count int, s_src, s_dst string, err error)error{
				_r_src := getRelativePath(il.MountInfo, s_src)
				_r_dst := getRelativePath(dstMountItem, s_dst)
				if nil != err {
					// 是否是目标文件家已存在, 如是是则跳过处理
					if filetools.IsDir(s_src) && filetools.IsError_DestExist(err) {
						return callback(count, _r_src, _r_dst, nil)
					}
					// 出现错误
					return callback(count, _r_src, _r_dst, &MoveError{
								IsSrcExist:  filetools.IsExist(s_src),
								IsDstExist:  filetools.IsExist(s_dst),
								ErrorString: parseErrorString(il.MountInfo.MtAddr, dstMountItem.MtAddr, err),
							})
				}
				return callback(count, _r_src, _r_dst, nil)
			})
		}
		case mt_type_oss:{ // oss对象存储
			return errors.New("This type of partition mount type is not implemented: Oss")
		}
		default:{ // 不支持的分区挂载类型
			return errors.New("Unsupported partition mount type: "+dstMountItem.MtType)
		}
	}
}
// 重命名文件|文件夹
func (il imp_local) DoRename(relativePath string, newName string)error{
	if il.MountInfo.MtPath == relativePath {
		return errors.New("Does not allow access: "+ relativePath)
	}
	_abs_src, _, _err := getAbsolutePath(il.MountInfo, relativePath)
	if nil != _err {
		return _err
	}
	if len(newName) == 0 {
		return nil
	}
	return filetools.Rename(_abs_src, newName)
}
// 新建文件夹
func (il imp_local) DoNewFolder(relativePath string)error{
	if il.MountInfo.MtPath == relativePath {
		return errors.New("Does not allow access: "+ relativePath)
	}
	_abs_src, _, _err := getAbsolutePath(il.MountInfo, relativePath)
	if nil != _err {
		return _err
	}
	return filetools.Mkdir(_abs_src)
}
// 删除文件|文件夹
func (il imp_local) DoDelete( relativePath string )error{
	if il.MountInfo.MtPath == relativePath {
		return errors.New("Does not allow access: "+ relativePath)
	}
	_abs_src, _, _err := getAbsolutePath(il.MountInfo, relativePath)
	if nil != _err {
		return _err
	}
	_deleting_path := getAbsoluteDeletingPath(il.MountInfo)
	// 移动到删除零时目录, 如果存在则覆盖
	// 通过这种方式可以减少函数等待时间, 但是如果线程删除失败则可能导致文件无法删除
	// 所以再启动或者周期性的检擦删除零时目录, 进行清空
	mv_err := filetools.MoveDir(_abs_src, _deleting_path, true, false, func(count int, s_src, s_dst string, err error)error{
		return err
	})
	// 开一个线程去移除它, 移除可能需要更多的时间
	if nil == mv_err {
		go il.DoClearDeletings( )
	}
	return mv_err
}
// 删除各个分区内的'临时删除文件'
func (il imp_local) DoClearDeletings( ){
	for _, _val := range il.mtm.Mts{
		dirs := filetools.GetDirList(_val.MtAddr+"/"+locl_deleting)
		if nil == dirs { continue; }
		for _, temp := range dirs{
			err := filetools.RemoveAll( filepath.Clean(_val.MtAddr+"/"+locl_deleting+"/"+temp) )
			if nil != err {
				fmt.Println("DoClearDeletings", err)
			}
		}
	}
}
// 拷贝文件
func (il imp_local) DoCopy(src, dst string, replace, ignore bool, callback CopyCallback)error{
	_abs_src, _, _err := getAbsolutePath(il.MountInfo, src)
	if nil != _err {
		return _err
	}
	// 目标位置驱动接口
	dstMountItem := il.mtm.getMountItem(dst)
	_abs_dst, _, err := getAbsolutePath(dstMountItem, dst)
	if nil != err {
		return err
	}
	switch dstMountItem.MtType {
		case mt_type_local:{ // 本地存储
			if filetools.IsFile(_abs_src){
				_r_src := getRelativePath(il.MountInfo, _abs_src)
				_r_dst := getRelativePath(dstMountItem, _abs_dst)
				_err   := filetools.CopyFile(_abs_src, _abs_dst, replace, ignore )
				if nil != err {
					return callback(1, _r_src, _r_dst, &CopyError{
								IsSrcExist:  filetools.IsExist(_abs_src),
								IsDstExist:  filetools.IsExist(_abs_dst),
								ErrorString: parseErrorString(il.MountInfo.MtAddr, dstMountItem.MtAddr, _err),
							})
				}
				return callback(1, _r_src, _r_dst, nil)
			}else{
				return filetools.CopyDir(_abs_src, _abs_dst, replace, ignore, func(count int, s_src, s_dst string, err error)error{
					_r_src := getRelativePath(il.MountInfo, s_src)
					_r_dst := getRelativePath(dstMountItem, s_dst)
					if nil != err {
						// 是否是目标文件家已存在, 如是是则跳过处理
						if filetools.IsDir(s_src) && filetools.IsError_DestExist(err) {
							return callback(count, _r_src, _r_dst, nil)
						}
						// 出现错误
						return callback(count, _r_src, _r_dst, &CopyError{
									IsSrcExist:  filetools.IsExist(s_src),
									IsDstExist:  filetools.IsExist(s_dst),
									ErrorString: parseErrorString(il.MountInfo.MtAddr, dstMountItem.MtAddr, err),
								})
					}
					return callback(count, _r_src, _r_dst, nil)
				})
				
			}
		}
		case mt_type_oss:{ // oss对象存储
			return errors.New("This type of partition mount type is not implemented: Oss")
		}
		default:{ // 不支持的分区挂载类型
			return errors.New("Unsupported partition mount type: "+dstMountItem.MtType)
		}
	}
}
func (il imp_local) DoCreat( )(bool, error){
	return true, nil
}
// 读取文件
func (il imp_local) DoRead( relativePath string )(io.Reader, error){
	abs_dst, _, gp_err := getAbsolutePath(il.MountInfo, relativePath)
	if nil != gp_err {
		return nil, gp_err
	}
	return filetools.OpenFile(abs_dst)
}
// 写入文件
// 先写入临时位置, 然后移动到正确位置
func (il imp_local) DoWrite( relativePath string, ioReader io.Reader ) error{
	if ioReader == nil {
		return errors.New("IO Reader is nil")
	}
	abs_dst, _, gp_err := getAbsolutePath(il.MountInfo, relativePath)
	if nil != gp_err {
		return  gp_err
	}
	tempPath := getAbsoluteTempPath(il.MountInfo)
	fs, w_err := filetools.GetWriter( tempPath )
	if w_err != nil {
		return w_err
	}
	_, cp_err := io.Copy(fs, ioReader)
	if nil == cp_err {
		fsClose_err := fs.Close( )
		if fsClose_err == nil{
			return filetools.MoveFile(tempPath, abs_dst, true, false, func(count int, s_src, s_dst string, err error)error{
				return err
			})
		}
		return fsClose_err
	}else{
		fsClose_err := fs.Close( )
		if nil != fsClose_err {
			return fsClose_err
		}
		rm_err := filetools.RemoveFile(tempPath)
		if rm_err != nil {
			return rm_err
		}
	}
	return cp_err
}
func (il imp_local) AskUploadToken( relativePath string )(string, error){
	return "", nil
}
func (il imp_local) SaveTokenFile(token string, src io.Reader)(bool){
	return true
}
func (il imp_local) SubmitToken(token string, isCreateVer bool)(bool, error){
	return true, nil
}

// ============================
// 去除具体位置信息
func parseErrorString( src, dsc string, err error ) string{
	if nil != err{
		errorString := err.Error( )
		if strings.Index(errorString, src) > -1 {
			return strings.Replace(errorString, src, "", -1)
		} else if strings.Index(errorString, dsc) > -1 {
			return strings.Replace(errorString, dsc, "", -1)
		}
		src = filepath.Clean(src)
		dsc = filepath.Clean(dsc)
		if strings.Index(errorString, src) > -1 {
			return strings.Replace(errorString, src, "", -1)
		} else if strings.Index(errorString, dsc) > -1 {
			return strings.Replace(errorString, dsc, "", -1)
		}
		return errorString
	}
	return ""
}

func test( ){
	fmt.Println("...")
}