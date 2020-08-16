// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 本地文件挂载操作驱动

package filemanage

import (
	"errors"
	"fmt"
	"gutils/fstool"
	"gutils/strtool"
	"io"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// localDriver 本地文件挂载操作驱动
type localDriver struct {
	mountNode mountNode
	mtm       *mountManager
}

// IsExist 文件是否存在
func (locl localDriver) IsExist(relativePath string) (bool, error) {
	absPath, _, err := locl.getAbsolutePath(locl.mountNode, relativePath)
	return fstool.IsExist(absPath), locl.wrapError(err)
}

// IsDir IsDir
func (locl localDriver) IsDir(relativePath string) (bool, error) {
	absPath, _, err := locl.getAbsolutePath(locl.mountNode, relativePath)
	return fstool.IsDir(absPath), locl.wrapError(err)
}

// IsFile IsFile
func (locl localDriver) IsFile(relativePath string) (bool, error) {
	absPath, _, err := locl.getAbsolutePath(locl.mountNode, relativePath)
	return fstool.IsFile(absPath), locl.wrapError(err)
}

// GetDirList 获取路径列表, 返回相对路径
func (locl localDriver) GetDirList(relativePath string) ([]string, error) {
	absPath, mRelativePath, err := locl.getAbsolutePath(locl.mountNode, relativePath)
	if err != nil {
		return make([]string, 0), err
	}
	ls, err := fstool.GetDirList(absPath)
	if nil != err {
		return make([]string, 0), locl.wrapError(err)
	}
	// 如果是挂载目录根目录, 需要处理 缓存目录
	if mRelativePath == "/" {
		if ls != nil && len(ls) > 0 {
			res := make([]string, 0)
			for i := 0; i < len(ls); i++ {
				// 如果是挂载目录根目录, 忽略系统目录
				if sysDir == ls[i] {
					continue
				}
				res = append(res, ls[i])
			}
			return res, nil
		}
	}
	return ls, nil
}

// GetFileSize GetFileSize
func (locl localDriver) GetFileSize(relativePath string) (int64, error) {
	absPath, _, err := locl.getAbsolutePath(locl.mountNode, relativePath)
	if nil != err {
		return -1, err
	}
	size, err := fstool.GetFileSize(absPath)
	return size, locl.wrapError(err)
}

// GetModifyTime GetModifyTime
func (locl localDriver) GetModifyTime(relativePath string) (int64, error) {
	absPath, _, err := locl.getAbsolutePath(locl.mountNode, relativePath)
	if nil != err {
		return -1, err
	}
	time, err := fstool.GetModifyTime(absPath)
	if nil != err {
		return -1, err
	}
	return time.UnixNano() / 1e6, nil
}

// DoMove 移动文件|夹
func (locl localDriver) DoMove(src string, dst string, replace, ignore bool, callback MoveCallback) error {
	if locl.mountNode.mtPath == src {
		return errors.New("Does not allow access: " + src)
	}
	absSrc, _, err := locl.getAbsolutePath(locl.mountNode, src)
	if nil != err {
		return err
	}
	if locl.mountNode.mtAddr == absSrc {
		return errors.New(src + " is mount root, cannot move")
	}
	// 目标位置驱动接口
	dstMountItem := locl.mtm.getMountItem(dst)
	absDst, _, err := locl.getAbsolutePath(dstMountItem, dst)
	if nil != err {
		return err
	}
	switch dstMountItem.mtType {
	case localTypeKey:
		{ // 本地存储
			return fstool.MoveFiles(absSrc, absDst, replace, ignore, func(srcPath, dstPath string, err error) error {
				rSrc := locl.getRelativePath(locl.mountNode, srcPath)
				rDst := locl.getRelativePath(dstMountItem, dstPath)
				if nil != err {
					// 出现错误
					return callback(rSrc, rDst, &MoveError{
						SrcIsExist:  fstool.IsExist(srcPath),
						DstIsExist:  fstool.IsExist(dstPath),
						ErrorString: locl.clearMountAddr(locl.mountNode, dstMountItem, err),
					})
				}
				return callback(rSrc, rDst, nil)
			})
		}
	case ossTypeKey:
		{ // oss对象存储
			return errors.New("locl type of partition mount type is not implemented: Oss")
		}
	default:
		{ // 不支持的分区挂载类型
			return errors.New("Unsupported partition mount type: " + dstMountItem.mtType)
		}
	}
}

// 重命名文件|文件夹
func (locl localDriver) DoRename(relativePath string, newName string) error {
	if locl.mountNode.mtPath == relativePath {
		return errors.New("Does not allow access: " + relativePath)
	}
	absSrc, _, err := locl.getAbsolutePath(locl.mountNode, relativePath)
	if nil != err {
		return err
	}
	if len(newName) == 0 {
		return nil
	}
	return locl.wrapError(fstool.Rename(absSrc, newName))
}

// 新建文件夹
func (locl localDriver) DoNewFolder(relativePath string) error {
	if locl.mountNode.mtPath == relativePath {
		return errors.New("Does not allow access: " + relativePath)
	}
	absSrc, _, err := locl.getAbsolutePath(locl.mountNode, relativePath)
	if nil != err {
		return err
	}
	return locl.wrapError(fstool.Mkdir(absSrc))
}

// DoDelete 删除文件|文件夹
func (locl localDriver) DoDelete(relativePath string) error {
	if locl.mountNode.mtPath == relativePath {
		return errors.New("Does not allow access: " + relativePath)
	}
	absSrc, _, err := locl.getAbsolutePath(locl.mountNode, relativePath)
	if nil != err {
		return err
	}
	deletingPath := locl.getAbsoluteDeletingPath(locl.mountNode)
	// 移动到删除零时目录, 如果存在则覆盖
	// 通过这种方式可以减少函数等待时间, 但是如果线程删除失败则可能导致文件无法删除
	// 所以再启动或者周期性的检擦删除零时目录, 进行清空
	mvErr := fstool.MoveFiles(absSrc, deletingPath, true, false, func(srcPath, dstPath string, err error) error {
		return err
	})
	// 开一个线程去移除它, 移除可能需要更多的时间
	if nil == mvErr {
		go locl.DoClearDeletings()
	}
	return locl.wrapError(mvErr)
}

// DoClearDeletings 删除各个分区内的'临时删除文件'
func (locl localDriver) DoClearDeletings() {
	for i := 0; i < len(locl.mtm.mtnds); i++ {
		if locl.mtm.mtnds[i].mtType == localTypeKey {
			dirs, _ := fstool.GetDirList(filepath.Clean(locl.mtm.mtnds[i].mtAddr + "/" + deletingDir))
			if nil == dirs {
				continue
			}
			for j := 0; j < len(dirs); j++ {
				err := fstool.RemoveAll(filepath.Clean(locl.mtm.mtnds[i].mtAddr + "/" + deletingDir + "/" + dirs[j]))
				if nil != err {
					fmt.Println("DoClearDeletings", err)
				}
			}
		}
	}
}

// DoCopy 拷贝文件
func (locl localDriver) DoCopy(src, dst string, replace, ignore bool, callback CopyCallback) error {
	absSrc, _, err := locl.getAbsolutePath(locl.mountNode, src)
	if nil != err {
		return err
	}
	// 目标位置驱动接口
	dstMountItem := locl.mtm.getMountItem(dst)
	absDst, _, err := locl.getAbsolutePath(dstMountItem, dst)
	if nil != err {
		return err
	}
	switch dstMountItem.mtType {
	case localTypeKey:
		{ // 本地存储
			if fstool.IsFile(absSrc) {
				rSrc := locl.getRelativePath(locl.mountNode, absSrc)
				rDst := locl.getRelativePath(dstMountItem, absDst)
				err = fstool.CopyFile(absSrc, absDst, replace, ignore)
				if nil != err {
					return callback(rSrc, rDst, &CopyError{
						SrcIsExist:  fstool.IsExist(absSrc),
						DstIsExist:  fstool.IsExist(absDst),
						ErrorString: locl.clearMountAddr(locl.mountNode, dstMountItem, err),
					})
				}
				return callback(rSrc, rDst, nil)
			}
			return fstool.CopyFiles(absSrc, absDst, replace, ignore, func(srcPath, dstPath string, err error) error {
				rSrc := locl.getRelativePath(locl.mountNode, srcPath)
				rDst := locl.getRelativePath(dstMountItem, dstPath)
				if nil != err {
					// 出现错误
					return callback(rSrc, rDst, &CopyError{
						SrcIsExist:  fstool.IsExist(srcPath),
						DstIsExist:  fstool.IsExist(dstPath),
						ErrorString: locl.clearMountAddr(locl.mountNode, dstMountItem, err),
					})
				}
				return callback(rSrc, rDst, nil)
			})
		}
	case ossTypeKey:
		{ // oss对象存储
			return errors.New("locl type of partition mount type is not implemented: Oss")
		}
	default:
		{ // 不支持的分区挂载类型
			return errors.New("Unsupported partition mount type: " + dstMountItem.mtType)
		}
	}
}
func (locl localDriver) DoCreat() (bool, error) {
	return true, nil
}

// DoRead 读取文件, 需要手动关闭流
func (locl localDriver) DoRead(relativePath string, offset int64) (Reader, error) {
	absDst, _, gpErr := locl.getAbsolutePath(locl.mountNode, relativePath)
	if nil != gpErr {
		return nil, gpErr
	}
	fs, err := fstool.OpenFile(absDst)
	if nil != err {
		return nil, locl.wrapError(err)
	}
	_, err = fs.Seek(offset, io.SeekStart)
	if nil != err {
		return nil, locl.wrapError(err)
	}
	return fs, nil
}

// DoWrite 写入文件， 先写入临时位置, 然后移动到正确位置
func (locl localDriver) DoWrite(relativePath string, ioReader io.Reader) error {
	if ioReader == nil {
		return errors.New("IO Reader is nil")
	}
	absDst, _, gpErr := locl.getAbsolutePath(locl.mountNode, relativePath)
	if nil != gpErr {
		return gpErr
	}
	tempPath := locl.getAbsoluteTempPath(locl.mountNode)
	fs, wErr := fstool.GetWriter(tempPath)
	if wErr != nil {
		return locl.wrapError(wErr)
	}
	_, cpErr := io.Copy(fs, ioReader)
	if nil == cpErr {
		fsCloseErr := fs.Close()
		if fsCloseErr == nil {
			return fstool.MoveFiles(tempPath, absDst, true, false, func(srcPath, dstPath string, err error) error {
				return locl.wrapError(err)
			})
		}
		return locl.wrapError(fsCloseErr)
	}
	fsCloseErr := fs.Close()
	if nil != fsCloseErr {
		return locl.wrapError(fsCloseErr)
	}
	rmErr := fstool.RemoveFile(tempPath)
	if rmErr != nil {
		return locl.wrapError(rmErr)
	}
	return locl.wrapError(cpErr)
}

// getAbsolutePath 处理路径拼接
func (locl localDriver) getAbsolutePath(mountNode mountNode, relativePath string) (abs string, rlPath string, err error) {
	rlPath = relativePath
	if "/" != mountNode.mtPath {
		rlPath = relativePath[len(mountNode.mtPath):]
		if rlPath == "" {
			rlPath = "/"
		}
	}
	// /Mount/.sys/.cache=>/.sys/.cache
	if rlPath == sysDir ||
		rlPath == "/"+sysDir ||
		0 == strings.Index(rlPath, "/"+sysDir+"/") {
		return abs, rlPath, errors.New("Does not allow access: " + rlPath)
	}
	abs = filepath.Clean(mountNode.mtAddr + rlPath)
	//fmt.Println( "getAbsolutePath: ", rlPath, abs )
	return
}

// getRelativePath 获取相对路径
func (locl localDriver) getRelativePath(mti mountNode, absolute string) string {
	// fmt.Println("locl.getRelativePath: ", mti.mtAddr, absolute)
	absolute = filepath.Clean(absolute)
	if filepath.IsAbs(absolute) {
		if strings.HasPrefix(absolute, mti.mtAddr) {
			return strtool.Parse2UnixPath(mti.mtPath + "/" + absolute[len(mti.mtAddr):])
		} else if strings.HasPrefix(mti.mtAddr, absolute) {
			return mti.mtPath
		}
	}
	return absolute
}

// getAbsoluteTempPath 获取该分区下的缓存目录
func (locl localDriver) getAbsoluteTempPath(mountNode mountNode) string {
	return filepath.Clean(mountNode.mtAddr + "/" + tempDir + "/" + strtool.GetUUID())
}

// getAbsoluteDeletingPath 获取一个放置删除文件的目录
func (locl localDriver) getAbsoluteDeletingPath(mountNode mountNode) string {
	return filepath.Clean(mountNode.mtAddr + "/" + deletingDir + "/" + strconv.FormatInt(time.Now().UnixNano(), 10))
}

// wrapLocalError 重新包装本地驱动错误信息, 避免真实路径暴露
func (locl localDriver) wrapError(err error) error {
	if nil != err {
		errStr := err.Error()
		if len(errStr) > 0 {
			rStr := locl.mountNode.mtPath
			if locl.mountNode.mtPath == "/" {
				rStr = ""
			}
			errStr = strings.Replace(errStr, "\\", "/", -1)
			errStr = strings.Replace(errStr, strtool.Parse2UnixPath(locl.mountNode.mtAddr), rStr, -1)
			return errors.New(errStr)
		}
	}
	return err
}

// clearMountAddr 去除挂载目录的位置信息
func (locl localDriver) clearMountAddr(srcMount, destMount mountNode, err error) string {
	if nil != err {
		errStr := err.Error()
		if len(errStr) > 0 {
			errStr = strings.Replace(errStr, "\\", "/", -1)
			srcMtAddr := strtool.Parse2UnixPath(srcMount.mtAddr)
			destMtAddr := strtool.Parse2UnixPath(destMount.mtAddr)
			if strings.Index(errStr, srcMtAddr) > -1 {
				// /root/datas/a/b -> /a/b/a/b
				rStr := srcMount.mtPath
				if srcMount.mtPath == "/" {
					rStr = ""
				}
				errStr = strings.Replace(errStr, srcMtAddr, rStr, -1)
			}
			if srcMtAddr != destMtAddr && strings.Index(errStr, destMtAddr) > -1 {
				rStr := destMount.mtPath
				if destMount.mtPath == "/" {
					rStr = ""
				}
				errStr = strings.Replace(errStr, destMtAddr, rStr, -1)
			}
			return errStr
		}
	}
	return ""
}
