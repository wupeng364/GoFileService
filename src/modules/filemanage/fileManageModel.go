package filemanage
/**
 *@description 文件管理模块
 *@author	wupeng364@outlook.com
*/
import (
	"common/gomodule"
	"common/tokenmanager"
	"modules/configmodel"
	"io"
	"fmt"
)
type FileManageModel struct{
	cfgModel *configmodel.ConfigModel
	mt 		 *MountManager
	tk		 *tokenmanager.TokenManager
}

// 返回模块信息
func (fm *FileManageModel)MInfo( )(*gomodule.ModelInfo)	{
	return &gomodule.ModelInfo{
		fm,
		"FileManageModel",
		1.0,
		"文件管理模块",
	}
}
// 模块安装, 一个模块只初始化一次
func (fm *FileManageModel)MSetup( ) {
	
}
// 模块升级, 一个版本执行一次
func (fm *FileManageModel)MUpdate( ) {
	
}

// 每次启动加载模块执行一次
func (fm *FileManageModel)OnMInit( getPointer func(m interface{})interface{}  ) {
	fm.cfgModel = getPointer(fm.cfgModel).(*configmodel.ConfigModel)
	fm.mt = (&MountManager{}).initMountItems( fm.cfgModel.GetConfigs(cfg_data_mount).(map[string]interface{}) )
	fm.tk = (&tokenmanager.TokenManager{}).Init( )
}
// 系统执行销毁时执行
func (fm *FileManageModel)OnMDestroy( ) {
	
}

// ==============================================================================================
// 申请一个Token用于跟踪和控制操作
// 复制, 移动 等出现重复或者异常后, 需要返回 跳过/重试 控制权限
// 后端的操作逻辑根据对象中的值进行跳过/重试操作, 如果客户端超过60s没有响应则放弃操作
func (fm *FileManageModel) AskToken(operationType string, tokenBody interface{})string{
	return fm.tk.AskToken(&tokenmanager.TokenObject{
			TypeStr: operationType,
			Second: tokenExpired_Second,
			TokenBody: tokenBody, 
	})
}
// 查询Token的内容
func (fm *FileManageModel) GetToken(token string) *tokenmanager.TokenObject{
	tokenobject, ok := fm.tk.GetTokenInfo(token)
	if ok {
		return tokenobject
	}
	return nil
}
func (fm *FileManageModel) RefreshToken(token string){
	fm.tk.RefreshToken(token)
}
func (fm *FileManageModel) RemoveToken(token string){
	fm.tk.DestroyToken(token)
}
// newName
func (fm *FileManageModel) DoRename(relativePath, newName string) error{
	fs := fm.mt.getInterface(relativePath)
	return fs.DoRename(relativePath, newName)
}
// DoNewFolder
func (fm *FileManageModel) DoNewFolder(relativePath string) error{
	fs := fm.mt.getInterface(relativePath)
	return fs.DoNewFolder(relativePath)
}
// 删除文件|文件夹
func (fm *FileManageModel) DoDelete(relativePath string) error{
	fs := fm.mt.getInterface(relativePath)
	return fs.DoDelete(relativePath)
}
// 移动文件|文件夹
func (fm *FileManageModel) DoMove(src, dest string, replace, ignore bool, callback MoveCallback)error{
	fs := fm.mt.getInterface(src)
	return fs.DoMove(src, dest, replace, ignore, callback)
}
// 复制文件|夹
func (fm *FileManageModel) DoCopy(src, dest string, replace, ignore bool, callback CopyCallback)error{
	fs := fm.mt.getInterface(src)
	return fs.DoCopy(src, dest, replace, ignore, callback)
}
// 写入文件
func (fm *FileManageModel)DoWrite(relativePath string, ioReader io.Reader)error{
	fs := fm.mt.getInterface(relativePath)
	return fs.DoWrite(relativePath, ioReader)
}
// 读取文件
func (fm *FileManageModel)DoRead(relativePath string)(io.Reader, error){
	fs := fm.mt.getInterface(relativePath)
	return fs.DoRead(relativePath)
}
// 是否是文件, 如果路径不对或者驱动不对则为 false
func (fm *FileManageModel) IsFile( relativePath string ) bool {
	fs := fm.mt.getInterface(relativePath)
	ok,_ := fs.IsFile(relativePath)
	return ok
}
// 是否存在, 如果路径不对或者驱动不对则为 false
func (fm *FileManageModel) IsExist( relativePath string ) bool {
	fs := fm.mt.getInterface(relativePath)
	ok,_ := fs.IsExist(relativePath)
	return ok
}
// 获取文件大小
func (fm *FileManageModel) GetFileSize(relativePath string) (int64, error){
	fs := fm.mt.getInterface(relativePath)
	return fs.GetFileSize(relativePath)
}
// 获取文件夹列表
func (fm *FileManageModel) GetDirList(relativePath string) ([]string, error){
	fs := fm.mt.getInterface(relativePath)
	return fs.GetDirList( relativePath )
}
// 获取文件夹下文件的基本信息
func (fm *FileManageModel) GetDirListInfo(relativePath string) ([]F_BaseInfo, error){
	fs := fm.mt.getInterface(relativePath)
	ls, err := fs.GetDirList( relativePath )
	_len_ls := len(ls)
	
	f_bi_file   := make([]F_BaseInfo, 0)
	f_bi_folder := make([]F_BaseInfo, 0)
	if err == nil && _len_ls > 0{
		for _, _p := range ls {
			_childPath := "/"+_p
			if relativePath != "/" {
				_childPath = relativePath+_childPath
			}
			// fmt.Println("_childPath: ", _childPath)
			isFile, _:=fs.IsFile(_childPath)
			fbi := F_BaseInfo{
				_childPath,
				(func( )int64{res,_:=fs.GetModifyTime(_childPath);return res;})( ),
				isFile,
				(func( )int64{if !isFile {return 0;}; res,_:=fs.GetFileSize(_childPath);return res;})( ),
			}
			if isFile {
				f_bi_file = append(f_bi_file, fbi)
			}else{
				f_bi_folder = append(f_bi_folder, fbi)
			}
			// fmt.Println(f_bi[i])
		}
	}
	m_ls := fm.mt.findMountChild(relativePath)
	if len(m_ls) > 0 {
		f_bi_folder = BaseInfoMerge(f_bi_folder, m_ls)
	}
	// 把文件夹排到前面去
	len_fbi_file   := len(f_bi_file)
	len_fbi_folder := len(f_bi_folder)
	res := make([]F_BaseInfo, len_fbi_file + len_fbi_folder )
	if len_fbi_folder > 0 {
		for i, val := range f_bi_folder {
			res[i] = val
		}
	}
	if len_fbi_file > 0 {
		for i, val := range f_bi_file {
			res[i+len_fbi_folder] = val
		}
	}
	return res, err
}

// ==============================
// 合并挂载路径到返回结果中去
func BaseInfoMerge(x []F_BaseInfo, y []string) []F_BaseInfo { 
	xlen := len(x)	//x数组的长度 
	z := make([]F_BaseInfo, xlen)
	// x
	for i, val := range x {
		z[i] = val
	}
	// y
	for _, val := range y {
		has := false
		for _, val1 := range x {
			if val1.Path == val {
				has = true; break
			} 
		}
		if !has {
			z = append(z, F_BaseInfo{ val, 0, false, 0 })
		}
	}
	return z 
}
func sayHello( ){
	fmt.Println("..")
}
