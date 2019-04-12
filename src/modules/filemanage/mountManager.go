package filemanage
/**
 *@description 挂载管理器
 *@author	wupeng364@outlook.com
*/
import(
	"fmt"
	"strings"
	"strconv"
	"time"
	"errors"
	"path"
	"common/filetools"
	"path/filepath"
)

// 常量
const(
	cfg_data_mount string = "data"
	cfg_data_type  string = "type"
	cfg_data_addr  string = "addr"

	mt_type_local string  = "LOCAL"
	mt_type_smb   string  = "SMB"
	mt_type_oss	  string  = "OSS"
	
	sys_dir = ".sys"
	locl_temp = sys_dir+"/.cache"
	locl_deleting = sys_dir+"/.deleting"
)
// 挂在的节点配置
type MtItems struct {
	MtPath string	// 挂载路径-虚拟路径
	MtType string	// 挂载类型
	MtAddr string	// 实际挂载路径
	Depth  int		// 深度
}
// 挂载管理器
type MountManager struct{
	Mts []MtItems
}

// 初始化挂载节点
func (mtm *MountManager) initMountItems( mounts map[string]interface{} ) *MountManager{
	if len(mounts) == 0 {
		panic("mounts is nil")
	}
	mts := make([]MtItems, len(mounts))
	_count := 0
	for _key, _val := range mounts{
		_val_new := _val.(map[string]interface{})
		mts[_count] = parseMtItems(MtItems{
			_key, 
			_val_new[cfg_data_type].(string), 
			_val_new[cfg_data_addr].(string),
			0,
		})
		// 初始化必要的文件夹
		_locl_temp := filepath.Clean( mts[_count].MtAddr+"/"+locl_temp )
		if !filetools.IsDir( _locl_temp ){
			if err := filetools.MkdirAll(_locl_temp); nil != err {
				panic("Create Folder Failed, Path: "+_locl_temp+", "+err.Error( ))
			}
		}
		_locl_deleting := filepath.Clean( mts[_count].MtAddr+"/"+locl_deleting )
		if !filetools.IsDir( _locl_deleting ){
			if err := filetools.MkdirAll(_locl_deleting); nil != err {
				panic("Create Folder Failed, Path: "+_locl_deleting+", "+err.Error( ))
			}
		}
		// 删除零时文件
		_dirs := filetools.GetDirList(_locl_temp)
		if nil != _dirs {
			for _, temp := range _dirs{
				err := filetools.RemoveAll( filepath.Clean(_locl_temp+"/"+temp) )
				if nil != err {
					panic("Clear temps Failed, Error: "+ err.Error( ))
				}
			}
		}
		_dirs = filetools.GetDirList(_locl_deleting)
		if nil != _dirs {
			for _, temp := range _dirs{
				err := filetools.RemoveAll( filepath.Clean(_locl_deleting+"/"+temp) )
				if nil != err {
					panic("Clear temps Failed, Error: "+ err.Error( ))
				}
			}
		}
		_count++
	}
	mtm.Mts = mts
	return mtm
}

// 根据相对路径获取对应驱动类
func (mtm *MountManager) getInterface( relativePath string) fmInterface {
	if len(strings.Replace(relativePath, " ", "", -1)) == 0 {
		relativePath = "/"
	}
	// 挂载节点
	recentMtItems := mtm.getMountItem(relativePath)
	// 解析 recentMtItems
	if recentMtItems.MtPath == "" {
		panic(errors.New("Mount path is not find"))
	}
	if recentMtItems.MtAddr == "" {
		panic(errors.New("Mount address is nil, at mount path: "+recentMtItems.MtPath))
	}
	if recentMtItems.MtType == "" {
		panic(errors.New("Mount Type is not find"))
	}
	// 
	switch recentMtItems.MtType {
		case mt_type_local:{ // 本地存储
			return &imp_local{ recentMtItems, mtm }
			break
		}
		case mt_type_smb:{ // Smb协议
			panic(errors.New("This type of partition mount type is not implemented: Smb"))
			break
		}
		case mt_type_oss:{ // oss对象存储
			panic(errors.New("This type of partition mount type is not implemented: Oss"))
			break
		}
		default:{ // 不支持的分区挂载类型
			panic(errors.New("Unsupported partition mount type: "+recentMtItems.MtType))
		}
	}
	return nil
}
// 查找相对路径下的分区挂载信息
func (mtm *MountManager) getMountItem( relativePath string ) MtItems{
	// 如果传入路径和挂载节点匹配, 则记录下来
	_PathLen := -1
	var recentMtItems MtItems
	for _, _val := range mtm.Mts{
		// 如果挂载路径再传入路径的头部, 则认为有效
		// "/"==>/A || /A==>/A || /A/==> /A/B/ 
		if "/" == _val.MtPath || 
			_val.MtPath == relativePath || 
			strings.HasPrefix(relativePath, _val.MtPath+"/") {
			// /A==>/A/B/C < /A/B==>/A/B/C 
			if _PathLen < len(_val.MtPath){
				_PathLen = len(_val.MtPath)
				recentMtItems = _val
			}
		}
		
	}
	return recentMtItems
}
// 查找符合当前路径下的子挂载分区路径 /==>/Mount
func (mtm *MountManager) findMountChild( relativePath string ) (res []string){
	if relativePath != "/" {
		return res
	}
	depth := len(strings.Split(relativePath, "/")) // 这个地方实质上+1了
	for _, _val := range mtm.Mts{
		if relativePath == "/" {
			// 如果为 / 则取挂载目录深度为 1 的 /==>/mount1 /mount2
			if  _val.Depth == 1 && _val.MtPath != "/" {
				res = append(res, _val.MtPath)
			}
		}else 
		// 其他目录则取当前目录深度加一目录&以他开头的 /ps==>/ps/mount1 /ps/mount2
		if _val.Depth == depth && _val.MtPath != "/"&&
			strings.HasPrefix(_val.MtPath, relativePath+"/") {
			res = append(res, _val.MtPath)
		}
	}
	return res
}
// ==================================
// 转换配置信息, 如: 相对路径转绝对路径
func parseMtItems( mi MtItems ) MtItems{
	
	// 需要统一挂载类型大消息
	mi.MtType = strings.ToUpper(mi.MtType)
	if	mi.MtType != mt_type_local &&
		mi.MtType != mt_type_smb && 
		mi.MtType != mt_type_oss {
			panic(errors.New("Unsupported partition mount type: "+mi.MtType))
	}
	// 本地挂载需要处理路径
	if mi.MtType == mt_type_local {
		if !filepath.IsAbs( mi.MtAddr ) {
			var err error
			mi.MtAddr, err = filepath.Abs( mi.MtAddr )
			if err != nil {
				panic( err )
			}
		}
		mi.MtAddr = filepath.Clean(mi.MtAddr)
	}
	// 需要注意挂载路径的结尾符号 /
	_LastIndex := strings.LastIndex(mi.MtPath, "/")
	if _LastIndex >0 && _LastIndex == len(mi.MtPath)-1{
		mi.MtPath = mi.MtPath[0:_LastIndex]
	}
	mi.Depth = len(strings.Split(mi.MtPath, "/"))-1
	fmt.Println("   > Mounting partition: ", mi)
	return mi
}

// 处理路径拼接
func getAbsolutePath( MountInfo MtItems, relativePath string )(abs string, rlPath string, err error){
	rlPath = relativePath
	if "/" != MountInfo.MtPath {
		rlPath = relativePath[len(MountInfo.MtPath):]
		if rlPath == "" {
			rlPath = "/"
		}
	}
	// /Mount/.sys/.cache=>/.sys/.cache
	if 	rlPath == sys_dir || 
		rlPath == "/"+sys_dir ||
		0 == strings.Index(rlPath, "/"+sys_dir+"/") {
		return abs, rlPath, errors.New("Does not allow access: "+ rlPath)
	}
	abs = filepath.Clean(MountInfo.MtAddr+rlPath)
	//fmt.Println( "getAbsolutePath: ", rlPath, abs )
	return
}
// 获取相对路径
func getRelativePath( mti MtItems, absolute string ) string{
	// fmt.Println("getRelativePath: ", mti.MtAddr, absolute)
	if strings.HasPrefix(absolute, mti.MtAddr) {
		return path.Clean( mti.MtPath+"/"+strings.Replace(absolute[len(mti.MtAddr):], "\\", "/", -1) )
	}
	return path.Clean( strings.Replace(absolute, "\\", "/", -1) ) 
}
// 获取该分区下的缓存目录
func getAbsoluteTempPath( MountInfo MtItems )string{
	return filepath.Clean(MountInfo.MtAddr+"/"+locl_temp+"/"+strconv.FormatInt(time.Now( ).UnixNano( ), 10))
}
// 获取一个放置删除文件的目录
func getAbsoluteDeletingPath( MountInfo MtItems )string{
	return filepath.Clean(MountInfo.MtAddr+"/"+locl_deleting+"/"+strconv.FormatInt(time.Now( ).UnixNano( ), 10))
}