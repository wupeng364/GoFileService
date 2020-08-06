// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 文件管理模块, 文件操作(新建、删除、移动、复制等)、虚拟分区挂载

package filemanage

import (
	"errors"
	"gofs/comm/conf"
	"gutils/mloader"
	"gutils/strtool"
	"gutils/tokentool"
	"io"
	"strings"
)

// FileManager 文件管理
type FileManager struct {
	mctx *mloader.Loader
	mt   *mountManager
	tk   *tokentool.TokenManager
	conf *conf.AppConf
}

// ModuleOpts 模块加载器接口实现, 返回模块信息&配置
func (fmg *FileManager) ModuleOpts() mloader.Opts {
	return mloader.Opts{
		Name:        "FileManager",
		Version:     1.0,
		Description: "文件管理模块",
		OnReady: func(mctx *mloader.Loader) {
			fmg.mctx = mctx
			fmg.conf = mctx.GetModuleByTemplate(fmg.conf).(*conf.AppConf)
		},
		OnSetup: func() {
			fmg.conf.SetConfig(MountsConfKey+"./."+mountTypeKey, "LOCAL")
			fmg.conf.SetConfig(MountsConfKey+"./."+mountAddrKey, "./datas")
		},
		OnInit: func() {
			mounts := fmg.conf.GetConfig(MountsConfKey).ToStrMap(make(map[string]interface{}))
			fmg.mt = (&mountManager{}).initMountItems(mounts)
			fmg.tk = (&tokentool.TokenManager{}).Init()
		},
	}
}

// AskToken 申请一个Token用于跟踪和控制操作
// 复制, 移动 等出现重复或者异常后, 需要返回 跳过/重试 控制权限
// 后端的操作逻辑根据对象中的值进行跳过/重试操作, 如果客户端超过60s没有响应则放弃操作
func (fmg *FileManager) AskToken(operationType string, tokenBody interface{}) string {
	return fmg.tk.AskToken(tokenBody, tokenExpiredSecond)
}

// GetToken 查询Token的内容
func (fmg *FileManager) GetToken(token string) interface{} {
	tokenobject, ok := fmg.tk.GetTokenBody(token)
	if ok {
		return tokenobject
	}
	return nil
}

// RefreshToken RefreshToken
func (fmg *FileManager) RefreshToken(token string) {
	fmg.tk.RefreshToken(token)
}

// RemoveToken RemoveToken
func (fmg *FileManager) RemoveToken(token string) {
	fmg.tk.DestroyToken(token)
}

// DoRename DoRename
func (fmg *FileManager) DoRename(relativePath, newName string) error {
	relativePath, err := checkPathSafety(relativePath)
	if nil != err {
		return err
	}
	fs := fmg.mt.getInterface(relativePath)
	return fs.DoRename(relativePath, newName)
}

// DoNewFolder DoNewFolder
func (fmg *FileManager) DoNewFolder(relativePath string) error {
	relativePath, err := checkPathSafety(relativePath)
	if nil != err {
		return err
	}
	fs := fmg.mt.getInterface(relativePath)
	return fs.DoNewFolder(relativePath)
}

// DoDelete 删除文件|文件夹
func (fmg *FileManager) DoDelete(relativePath string) error {
	relativePath, err := checkPathSafety(relativePath)
	if nil != err {
		return err
	}
	fs := fmg.mt.getInterface(relativePath)
	return fs.DoDelete(relativePath)
}

// DoMove 移动文件|文件夹
func (fmg *FileManager) DoMove(src, dest string, replace, ignore bool, callback MoveCallback) error {
	src, err := checkPathSafety(src)
	if nil != err {
		return err
	}
	dest, err = checkPathSafety(dest)
	if nil != err {
		return err
	}
	fs := fmg.mt.getInterface(src)
	return fs.DoMove(src, dest, replace, ignore, callback)
}

// DoCopy 复制文件|夹
func (fmg *FileManager) DoCopy(src, dest string, replace, ignore bool, callback CopyCallback) error {
	src, err := checkPathSafety(src)
	if nil != err {
		return err
	}
	dest, err = checkPathSafety(dest)
	if nil != err {
		return err
	}
	fs := fmg.mt.getInterface(src)
	return fs.DoCopy(src, dest, replace, ignore, callback)
}

// DoWrite 写入文件
func (fmg *FileManager) DoWrite(relativePath string, ioReader io.Reader) error {
	relativePath, err := checkPathSafety(relativePath)
	if nil != err {
		return err
	}
	fs := fmg.mt.getInterface(relativePath)
	return fs.DoWrite(relativePath, ioReader)
}

// DoRead 读取文件
func (fmg *FileManager) DoRead(relativePath string, offset int64) (Reader, error) {
	relativePath, err := checkPathSafety(relativePath)
	if nil != err {
		return nil, err
	}
	fs := fmg.mt.getInterface(relativePath)
	return fs.DoRead(relativePath, offset)
}

// IsFile 是否是文件, 如果路径不对或者驱动不对则为 false
func (fmg *FileManager) IsFile(relativePath string) bool {
	relativePath, err := checkPathSafety(relativePath)
	if nil != err {
		return false
	}
	fs := fmg.mt.getInterface(relativePath)
	ok, _ := fs.IsFile(relativePath)
	return ok
}

// IsExist 是否存在, 如果路径不对或者驱动不对则为 false
func (fmg *FileManager) IsExist(relativePath string) bool {
	relativePath, err := checkPathSafety(relativePath)
	if nil != err {
		return false
	}
	fs := fmg.mt.getInterface(relativePath)
	ok, _ := fs.IsExist(relativePath)
	return ok
}

// GetFileSize 获取文件大小
func (fmg *FileManager) GetFileSize(relativePath string) (int64, error) {
	relativePath, err := checkPathSafety(relativePath)
	if nil != err {
		return -1, err
	}
	fs := fmg.mt.getInterface(relativePath)
	return fs.GetFileSize(relativePath)
}

// GetDirList 获取文件夹列表
func (fmg *FileManager) GetDirList(relativePath string) ([]string, error) {
	relativePath, err := checkPathSafety(relativePath)
	if nil != err {
		return make([]string, 0), err
	}
	fs := fmg.mt.getInterface(relativePath)
	return fs.GetDirList(relativePath)
}

// GetDirListInfo 获取文件夹下文件的基本信息
func (fmg *FileManager) GetDirListInfo(relativePath string) ([]FsInfo, error) {
	relativePath, err := checkPathSafety(relativePath)
	if nil != err {
		return make([]FsInfo, 0), err
	}
	fs := fmg.mt.getInterface(relativePath)
	ls, err := fs.GetDirList(relativePath)
	lenLS := len(ls)

	files := make([]FsInfo, 0)
	folders := make([]FsInfo, 0)
	if err == nil && lenLS > 0 {
		for _, p := range ls {
			childPath := "/" + p
			if relativePath != "/" {
				childPath = relativePath + childPath
			}
			// fmt.Println("childPath: ", childPath)
			isFile, _ := fs.IsFile(childPath)
			fbi := FsInfo{
				childPath,
				(func() int64 { res, _ := fs.GetModifyTime(childPath); return res })(),
				isFile,
				(func() int64 {
					if !isFile {
						return 0
					}
					res, _ := fs.GetFileSize(childPath)
					return res
				})(),
			}
			if isFile {
				files = append(files, fbi)
			} else {
				folders = append(folders, fbi)
			}
			// fmt.Println(f_bi[i])
		}
	}
	mLS := fmg.mt.findMountChild(relativePath)
	if len(mLS) > 0 {
		for _, val := range mLS {
			existedIndex := -1
			for lsi, lsVal := range ls {
				if lsVal == val {
					existedIndex = lsi
				}
			}
			if existedIndex == -1 {
				childPath := val
				if relativePath != "/" {
					childPath = relativePath + childPath
				}
				mtfs := fmg.mt.getInterface(childPath)
				folders = append(folders, FsInfo{
					childPath,
					(func() int64 { res, _ := mtfs.GetModifyTime(childPath); return res })(),
					false,
					0,
				})
			}
		}
	}
	// 把文件夹排到前面去
	lenFiles := len(files)
	lenFolders := len(folders)
	res := make([]FsInfo, lenFiles+lenFolders)
	if lenFolders > 0 {
		for i, val := range folders {
			res[i] = val
		}
	}
	if lenFiles > 0 {
		for i, val := range files {
			res[i+lenFolders] = val
		}
	}
	return res, err
}

// checkPathSafety 路径合规检查, 避免 ../ ./之类的路径
func checkPathSafety(path string) (string, error) {
	if len(path) == 0 {
		return "/", nil
	}
	if strings.Index(path, "../") > -1 {
		return "", errors.New("Unsupported path format '../'")
	}
	if strings.Index(path, "./") > -1 {
		return "", errors.New("Unsupported path format './'")
	}
	return strtool.Parse2UnixPath(path), nil
}
