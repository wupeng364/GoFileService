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
	"io"
	"path/filepath"
	"strings"
)

// localDriver 本地文件挂载操作驱动
type localDriver struct {
	mountNode mountNodes
	mtm       *mountManager
}

// IsExist 文件是否存在
func (locl localDriver) IsExist(relativePath string) (bool, error) {
	absPath, _, err := getAbsolutePath(locl.mountNode, relativePath)
	return fstool.IsExist(absPath), err
}

// IsDir IsDir
func (locl localDriver) IsDir(relativePath string) (bool, error) {
	absPath, _, err := getAbsolutePath(locl.mountNode, relativePath)
	return fstool.IsDir(absPath), err
}

// IsFile IsFile
func (locl localDriver) IsFile(relativePath string) (bool, error) {
	absPath, _, err := getAbsolutePath(locl.mountNode, relativePath)
	return fstool.IsFile(absPath), err
}

// GetDirList 获取路径列表, 返回相对路径
func (locl localDriver) GetDirList(relativePath string) ([]string, error) {
	absPath, relativePath, err := getAbsolutePath(locl.mountNode, relativePath)
	if err != nil {
		return make([]string, 0), err
	}
	ls, err := fstool.GetDirList(absPath)
	if nil != err {
		return make([]string, 0), err
	}
	// 如果是挂载目录根目录, 需要处理 缓存目录
	if relativePath == "/" {
		if ls != nil && len(ls) > 0 {
			res := make([]string, 0)
			for _, p := range ls {
				// 如果是挂载目录根目录, 忽略系统目录
				if sysDir == p {
					continue
				}
				res = append(res, p)
			}
			return res, nil
		}
	}
	return ls, nil
}

// GetFileSize GetFileSize
func (locl localDriver) GetFileSize(relativePath string) (int64, error) {
	absPath, _, err := getAbsolutePath(locl.mountNode, relativePath)
	if nil != err {
		return -1, err
	}
	return fstool.GetFileSize(absPath)
}

// GetModifyTime GetModifyTime
func (locl localDriver) GetModifyTime(relativePath string) (int64, error) {
	absPath, _, err := getAbsolutePath(locl.mountNode, relativePath)
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
	absSrc, _, err := getAbsolutePath(locl.mountNode, src)
	if nil != err {
		return err
	}
	if filepath.Clean(locl.mountNode.mtAddr) == absSrc {
		return errors.New(src + " is mount root, cannot move")
	}
	// 目标位置驱动接口
	dstMountItem := locl.mtm.getMountItem(dst)
	absDst, _, err := getAbsolutePath(dstMountItem, dst)
	if nil != err {
		return err
	}
	switch dstMountItem.mtType {
	case localTypeKey:
		{ // 本地存储
			return fstool.MoveFiles(absSrc, absDst, replace, ignore, func(srcPath, dstPath string, err error) error {
				rSrc := getRelativePath(locl.mountNode, srcPath)
				rDst := getRelativePath(dstMountItem, dstPath)
				if nil != err {
					// 出现错误
					return callback(rSrc, rDst, &MoveError{
						SrcIsExist:  fstool.IsExist(srcPath),
						DstIsExist:  fstool.IsExist(dstPath),
						ErrorString: parseErrorString(locl.mountNode.mtAddr, dstMountItem.mtAddr, err),
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
	absSrc, _, err := getAbsolutePath(locl.mountNode, relativePath)
	if nil != err {
		return err
	}
	if len(newName) == 0 {
		return nil
	}
	return fstool.Rename(absSrc, newName)
}

// 新建文件夹
func (locl localDriver) DoNewFolder(relativePath string) error {
	if locl.mountNode.mtPath == relativePath {
		return errors.New("Does not allow access: " + relativePath)
	}
	absSrc, _, err := getAbsolutePath(locl.mountNode, relativePath)
	if nil != err {
		return err
	}
	return fstool.Mkdir(absSrc)
}

// DoDelete 删除文件|文件夹
func (locl localDriver) DoDelete(relativePath string) error {
	if locl.mountNode.mtPath == relativePath {
		return errors.New("Does not allow access: " + relativePath)
	}
	absSrc, _, err := getAbsolutePath(locl.mountNode, relativePath)
	if nil != err {
		return err
	}
	deletingPath := getAbsoluteDeletingPath(locl.mountNode)
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
	return mvErr
}

// DoClearDeletings 删除各个分区内的'临时删除文件'
func (locl localDriver) DoClearDeletings() {
	for _, val := range locl.mtm.mtnds {
		if val.mtType == localTypeKey {
			dirs, _ := fstool.GetDirList(val.mtAddr + "/" + deletingDir)
			if nil == dirs {
				continue
			}
			for _, temp := range dirs {
				err := fstool.RemoveAll(filepath.Clean(val.mtAddr + "/" + deletingDir + "/" + temp))
				if nil != err {
					fmt.Println("DoClearDeletings", err)
				}
			}
		}
	}
}

// DoCopy 拷贝文件
func (locl localDriver) DoCopy(src, dst string, replace, ignore bool, callback CopyCallback) error {
	absSrc, _, err := getAbsolutePath(locl.mountNode, src)
	if nil != err {
		return err
	}
	// 目标位置驱动接口
	dstMountItem := locl.mtm.getMountItem(dst)
	absDst, _, err := getAbsolutePath(dstMountItem, dst)
	if nil != err {
		return err
	}
	switch dstMountItem.mtType {
	case localTypeKey:
		{ // 本地存储
			if fstool.IsFile(absSrc) {
				rSrc := getRelativePath(locl.mountNode, absSrc)
				rDst := getRelativePath(dstMountItem, absDst)
				err = fstool.CopyFile(absSrc, absDst, replace, ignore)
				if nil != err {
					return callback(rSrc, rDst, &CopyError{
						SrcIsExist:  fstool.IsExist(absSrc),
						DstIsExist:  fstool.IsExist(absDst),
						ErrorString: parseErrorString(locl.mountNode.mtAddr, dstMountItem.mtAddr, err),
					})
				}
				return callback(rSrc, rDst, nil)
			}
			return fstool.CopyFiles(absSrc, absDst, replace, ignore, func(srcPath, dstPath string, err error) error {
				rSrc := getRelativePath(locl.mountNode, srcPath)
				rDst := getRelativePath(dstMountItem, dstPath)
				if nil != err {
					// 出现错误
					return callback(rSrc, rDst, &CopyError{
						SrcIsExist:  fstool.IsExist(srcPath),
						DstIsExist:  fstool.IsExist(dstPath),
						ErrorString: parseErrorString(locl.mountNode.mtAddr, dstMountItem.mtAddr, err),
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
	absDst, _, gpErr := getAbsolutePath(locl.mountNode, relativePath)
	if nil != gpErr {
		return nil, gpErr
	}
	fs, err := fstool.OpenFile(absDst)
	if nil != err {
		return nil, err
	}
	_, err = fs.Seek(offset, io.SeekStart)
	if nil != err {
		return nil, err
	}
	return fs, nil
}

// DoWrite 写入文件， 先写入临时位置, 然后移动到正确位置
func (locl localDriver) DoWrite(relativePath string, ioReader io.Reader) error {
	if ioReader == nil {
		return errors.New("IO Reader is nil")
	}
	absDst, _, gpErr := getAbsolutePath(locl.mountNode, relativePath)
	if nil != gpErr {
		return gpErr
	}
	tempPath := getAbsoluteTempPath(locl.mountNode)
	fs, wErr := fstool.GetWriter(tempPath)
	if wErr != nil {
		return wErr
	}
	_, cpErr := io.Copy(fs, ioReader)
	if nil == cpErr {
		fsCloseErr := fs.Close()
		if fsCloseErr == nil {
			return fstool.MoveFiles(tempPath, absDst, true, false, func(srcPath, dstPath string, err error) error {
				return err
			})
		}
		return fsCloseErr
	}
	fsCloseErr := fs.Close()
	if nil != fsCloseErr {
		return fsCloseErr
	}
	rmErr := fstool.RemoveFile(tempPath)
	if rmErr != nil {
		return rmErr
	}
	return cpErr
}

// parseErrorString 去除具体位置信息
func parseErrorString(src, dsc string, err error) string {
	if nil != err {
		src = filepath.Clean(src)
		dsc = filepath.Clean(dsc)
		errorString := err.Error()
		if strings.Index(errorString, src) > -1 {
			errorString = strings.Replace(errorString, src, "", -1)
		}
		if strings.Index(errorString, dsc) > -1 {
			errorString = strings.Replace(errorString, dsc, "", -1)
		}
		return errorString
	}
	return ""
}
