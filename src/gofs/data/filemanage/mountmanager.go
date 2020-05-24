// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 挂载管理器, 文件虚拟目录解析&路径驱动获取

package filemanage

import (
	"errors"
	"fmt"
	"gutils/fstool"
	"gutils/strtool"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// 常量
const (
	mountTypeKey string = "type" // 配置文件-挂载类别
	mountAddrKey string = "addr" // 配置文件-挂载地址

	localTypeKey string = "LOCAL" // 本地挂载方式
	ossTypeKey   string = "OSS"   // oss挂载方式

	sysDir      = ".sys"                // 系统文件存放位置
	tempDir     = sysDir + "/.cache"    // 系统文件-缓存位置
	deletingDir = sysDir + "/.deleting" // 系统文件-待删除文件位置
)

// 挂在的节点配置
type mountNodes struct {
	mtPath string // 挂载路径-虚拟路径
	mtType string // 挂载类型
	mtAddr string // 实际挂载路径
	depth  int    // 深度
}

// 挂载管理器
type mountManager struct {
	mtnds []mountNodes
}

// initMountItems 初始化挂载节点
func (mtg *mountManager) initMountItems(mounts map[string]interface{}) *mountManager {
	if len(mounts) == 0 {
		panic("mounts is nil")
	}
	mtnds := make([]mountNodes, len(mounts))
	count := 0
	for key, val := range mounts {
		newVal := val.(map[string]interface{})
		mtnds[count] = parseMountNodes(mountNodes{
			mtPath: key,
			mtType: newVal[mountTypeKey].(string),
			mtAddr: newVal[mountAddrKey].(string),
			depth:  0,
		})
		// 本地驱动
		if mtnds[count].mtType == localTypeKey {
			// 初始化必要的文件夹
			loclTemp := filepath.Clean(mtnds[count].mtAddr + "/" + tempDir)
			if !fstool.IsDir(loclTemp) {
				if err := fstool.MkdirAll(loclTemp); nil != err {
					panic("Create Folder Failed, Path: " + loclTemp + ", " + err.Error())
				}
			}
			loclDeleting := filepath.Clean(mtnds[count].mtAddr + "/" + deletingDir)
			if !fstool.IsDir(loclDeleting) {
				if err := fstool.MkdirAll(loclDeleting); nil != err {
					panic("Create Folder Failed, Path: " + loclDeleting + ", " + err.Error())
				}
			}
			// 删除零时文件
			dirs, err := fstool.GetDirList(loclTemp)
			if nil != err {
				panic(err.Error())
			}
			if nil != dirs {
				for _, temp := range dirs {
					err = fstool.RemoveAll(filepath.Clean(loclTemp + "/" + temp))
					if nil != err {
						panic("Clear temps Failed, Error: " + err.Error())
					}
				}
			}
			dirs, err = fstool.GetDirList(loclDeleting)
			if nil != err {
				panic(err.Error())
			}
			if nil != dirs {
				for _, temp := range dirs {
					err := fstool.RemoveAll(filepath.Clean(loclDeleting + "/" + temp))
					if nil != err {
						panic("Clear temps Failed, Error: " + err.Error())
					}
				}
			}
		} else {
			panic("Unsupported mount type: " + mtnds[count].mtType)
		}
		count++
	}
	mtg.mtnds = mtnds
	return mtg
}

// getInterface 根据相对路径获取对应驱动类
func (mtg *mountManager) getInterface(relativePath string) fmInterface {
	if len(strings.Replace(relativePath, " ", "", -1)) == 0 {
		relativePath = "/"
	}
	// 挂载节点
	recentMountNodes := mtg.getMountItem(relativePath)
	// 解析 recentMountNodes
	if recentMountNodes.mtPath == "" {
		panic(errors.New("Mount path is not find"))
	}
	if recentMountNodes.mtAddr == "" {
		panic(errors.New("Mount address is nil, at mount path: " + recentMountNodes.mtPath))
	}
	if recentMountNodes.mtType == "" {
		panic(errors.New("Mount Type is not find"))
	}
	//
	switch recentMountNodes.mtType {
	case localTypeKey:
		// 本地存储
		return &localDriver{recentMountNodes, mtg}
	case ossTypeKey:
		// oss对象存储
		panic(errors.New("mtg type of partition mount type is not implemented: Oss"))
	default:
		// 不支持的分区挂载类型
		panic(errors.New("Unsupported partition mount type: " + recentMountNodes.mtType))
	}
}

// getMountItem 查找相对路径下的分区挂载信息
func (mtg *mountManager) getMountItem(relativePath string) mountNodes {
	// 如果传入路径和挂载节点匹配, 则记录下来
	pathLen := -1
	var recentMountNodes mountNodes
	for _, val := range mtg.mtnds {
		// 如果挂载路径再传入路径的头部, 则认为有效
		// "/"==>/A || /A==>/A || /A/==> /A/B/
		if "/" == val.mtPath ||
			val.mtPath == relativePath ||
			strings.HasPrefix(relativePath, val.mtPath+"/") {
			// /A==>/A/B/C < /A/B==>/A/B/C
			if pathLen < len(val.mtPath) {
				pathLen = len(val.mtPath)
				recentMountNodes = val
			}
		}

	}
	return recentMountNodes
}

// findMountChild 查找符合当前路径下的子挂载分区路径 /==>/Mount
func (mtg *mountManager) findMountChild(relativePath string) (res []string) {
	if relativePath != "/" {
		return res
	}
	depth := len(strings.Split(relativePath, "/")) // 这个地方实质上+1了
	for _, val := range mtg.mtnds {
		if relativePath == "/" {
			// 如果为 / 则取挂载目录深度为 1 的 /==>/mount1 /mount2
			if val.depth == 1 && val.mtPath != "/" {
				res = append(res, val.mtPath)
			}
		} else
		// 其他目录则取当前目录深度加一目录&以他开头的 /ps==>/ps/mount1 /ps/mount2
		if val.depth == depth && val.mtPath != "/" &&
			strings.HasPrefix(val.mtPath, relativePath+"/") {
			res = append(res, val.mtPath)
		}
	}
	return res
}

// parseMountNodes 转换配置信息, 如: 相对路径转绝对路径
func parseMountNodes(mi mountNodes) mountNodes {

	// 需要统一挂载类型大消息
	mi.mtType = strings.ToUpper(mi.mtType)
	if mi.mtType != localTypeKey &&
		mi.mtType != ossTypeKey {
		panic(errors.New("Unsupported partition mount type: " + mi.mtType))
	}
	// 本地挂载需要处理路径
	if mi.mtType == localTypeKey {
		if !filepath.IsAbs(mi.mtAddr) {
			var err error
			mi.mtAddr, err = filepath.Abs(mi.mtAddr)
			if err != nil {
				panic(err)
			}
		}
		mi.mtAddr = filepath.Clean(mi.mtAddr)
	}
	// 需要注意挂载路径的结尾符号 /
	lastIndex := strings.LastIndex(mi.mtPath, "/")
	if lastIndex > 0 && lastIndex == len(mi.mtPath)-1 {
		mi.mtPath = mi.mtPath[0:lastIndex]
	}
	mi.depth = len(strings.Split(mi.mtPath, "/")) - 1
	fmt.Println("   > Mounting partition: ", mi)
	return mi
}

// getAbsolutePath 处理路径拼接
func getAbsolutePath(mountNode mountNodes, relativePath string) (abs string, rlPath string, err error) {
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
func getRelativePath(mti mountNodes, absolute string) string {
	// fmt.Println("getRelativePath: ", mti.mtAddr, absolute)
	if strings.HasPrefix(absolute, mti.mtAddr) {
		return path.Clean(mti.mtPath + "/" + strings.Replace(absolute[len(mti.mtAddr):], "\\", "/", -1))
	}
	return path.Clean(strings.Replace(absolute, "\\", "/", -1))
}

// getAbsoluteTempPath 获取该分区下的缓存目录
func getAbsoluteTempPath(mountNode mountNodes) string {
	return filepath.Clean(mountNode.mtAddr + "/" + tempDir + "/" + strtool.GetUUID())
}

// getAbsoluteDeletingPath 获取一个放置删除文件的目录
func getAbsoluteDeletingPath(mountNode mountNodes) string {
	return filepath.Clean(mountNode.mtAddr + "/" + deletingDir + "/" + strconv.FormatInt(time.Now().UnixNano(), 10))
}
