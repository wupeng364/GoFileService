// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 文件权限管理模块

package filepermission

import (
	"errors"
	"gofs/base/conf"
	"gofs/data/usermanage"
	"gutils/mloader"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

// FPmsManager 用户管理模块
type FPmsManager struct {
	fpmsmg FPmsManage
	lock   *sync.RWMutex
	upms   map[string]map[string]int64
}

// ModuleOpts 模块加载器接口实现, 返回模块信息&配置
func (fpmsmgr *FPmsManager) ModuleOpts() mloader.Opts {
	return mloader.Opts{
		Name:        "FPmsManager",
		Version:     1.0,
		Description: "文件权限模块",
		OnReady: func(mctx *mloader.Loader) {
			fpmsmgr.lock = new(sync.RWMutex)
			fpmsmgr.upms = make(map[string]map[string]int64, 0)
			fpmsmgr.fpmsmg = fpmsmgr.getFPmsInterfaceImp(mctx)
		},
		OnSetup: func() {
			// 执行建库、建表
			err := fpmsmgr.fpmsmg.InitTables()
			if nil != err {
				panic(err)
			}
		},
		OnInit: func() {
			fpmsmgr.loadPermission2Memory()
		},
	}
}

// HashPermission 是否拥有某个权限
func (fpmsmgr *FPmsManager) HashPermission(userID, path string, permission int64) bool {
	userPms := fpmsmgr.getUserPermissionInMemory(userID)
	if len(userPms) == 0 {
		return false
	}
	val, ok := userPms[path]
	if ok { // 是个文件夹&有记录
		if permission == VISIBLECHILD && val >= VISIBLECHILD {
			return true
		}
		if 1<<permission == val&(1<<permission) {
			return true
		}
	} else { // 可能只是没有权限记录, 或者是个文件, 我们则需要尝试上级目录是否拥有>VISIBLECHILD的权限
		parent := path
		for {
			parent = parent[:strings.LastIndex(parent, "/")]
			if len(parent) == 0 {
				parent = "/"
			}
			// 如果从上级找到了权限, 则需要是>VISIBLECHILD的权限
			// 组最近的权限设定, 即便上级再上级有权限也不管
			if val, ok := userPms[parent]; ok && val > VISIBLECHILD {
				// 如果上级有明确的权限, 则下级默认可见
				if permission <= VISIBLE {
					return true
				}
				// 计算
				if 1<<permission == val&(1<<permission) {
					return true
				}
				return false
			}
			// 如果没有找到, 则继续网上找, 如果到/还没有, 则无权限
			if parent == "/" {
				return false
			}
		}
	}

	return false
}

// ListFPermissions 列出所有权限数据, 无分页
func (fpmsmgr *FPmsManager) ListFPermissions() ([]FPermissionInfo, error) {
	return fpmsmgr.fpmsmg.ListFPermissions()
}

// ListUserFPermissions 列出用户所有权限数据, 无分页
func (fpmsmgr *FPmsManager) ListUserFPermissions(userID string) ([]FPermissionInfo, error) {
	return fpmsmgr.fpmsmg.ListUserFPermissions(userID)
}

// QueryFPermission 根据权限ID查询详细信息
func (fpmsmgr *FPmsManager) QueryFPermission(permissionID string) (*FPermissionInfo, error) {
	return fpmsmgr.fpmsmg.QueryFPermission(permissionID)
}

// AddFPermission 添加权限
func (fpmsmgr *FPmsManager) AddFPermission(fpm FPermissionInfo) error {
	if fpm.UserID == usermanage.AdminUserID {
		return errors.New("不能增加管理员权限")
	}
	err := fpmsmgr.fpmsmg.AddFPermission(fpm)
	fpmsmgr.loadPermission2Memory()
	return err
}

// UpdateFPermission 修改权限
func (fpmsmgr *FPmsManager) UpdateFPermission(fpm FPermissionInfo) error {
	pms, err := fpmsmgr.fpmsmg.QueryFPermission(fpm.PermissionID)
	if nil != err {
		return err
	}
	if nil != pms {
		if pms.UserID == usermanage.AdminUserID {
			return errors.New("不能修改管理员权限")
		}
	}
	err = fpmsmgr.fpmsmg.UpdateFPermission(fpm)
	fpmsmgr.loadPermission2Memory()
	return err
}

// DelFPermission 根据permissionID删除权限
func (fpmsmgr *FPmsManager) DelFPermission(permissionID string) error {
	pms, err := fpmsmgr.fpmsmg.QueryFPermission(permissionID)
	if nil != err {
		return err
	}
	if nil != pms {
		if pms.UserID == usermanage.AdminUserID {
			return errors.New("不能删除管理员权限")
		}
	}
	err = fpmsmgr.fpmsmg.DelFPermission(permissionID)
	fpmsmgr.loadPermission2Memory()
	return err
}

// loadPermission2Memory 加载权限结构到内存
func (fpmsmgr *FPmsManager) loadPermission2Memory() {
	list, err := fpmsmgr.ListFPermissions()
	if nil != err {
		panic(err)
	}
	fpmsmgr.lock.Lock()
	defer fpmsmgr.lock.Unlock()
	fpmsmgr.upms = make(map[string]map[string]int64, 0)
	// 初始化数据
	for i := 0; i < len(list); i++ {
		val, ok := fpmsmgr.upms[list[i].UserID]
		if !ok {
			fpmsmgr.upms[list[i].UserID] = make(map[string]int64, 0)
			val, _ = fpmsmgr.upms[list[i].UserID]
		}
		// 当前用户&当前路径
		list[i].Path = path.Clean(list[i].Path)
		val[list[i].Path] = list[i].Permission
		// 检查上级目录, 生成 VISIBLECHILD 权限
		parent := list[i].Path
		for {
			parent = parent[:strings.LastIndex(parent, "/")]
			if len(parent) == 0 {
				if _, ok := val["/"]; !ok {
					val["/"] = VISIBLECHILD
				}
				break
			}
			_, ok := val[parent]
			if !ok {
				val[parent] = VISIBLECHILD
			}
		}
	}
}

// getUserPermissionsInMemory 获取用户权限结构
func (fpmsmgr *FPmsManager) getUserPermissionInMemory(userID string) map[string]int64 {
	fpmsmgr.lock.RLock()
	defer fpmsmgr.lock.RUnlock()
	val, ok := fpmsmgr.upms[userID]
	if ok {
		return val
	}
	return nil
}

// getFPmsInterfaceImp 获取启动类型, 并实例对象指针
func (fpmsmgr *FPmsManager) getFPmsInterfaceImp(mctx *mloader.Loader) FPmsManage {
	// 默认使用sqlite驱动
	dbType := mctx.GetParam(conf.DataBasType).ToString("")
	if len(dbType) == 0 {
		panic(conf.DataBasType + " is empty")
	}
	dbSource := mctx.GetParam(conf.DataBaseSource).ToString("")
	if len(dbSource) == 0 {
		panic(conf.DataBaseSource + " is empty")
	}
	var fpmsmg FPmsManage
	switch dbType {
	case conf.DefaultDataBaseType:
		{
			if !filepath.IsAbs(dbSource) {
				dbSource, _ = filepath.Abs(dbSource)
			}
			fpmsmg = &impBySqlite{}
			err := fpmsmg.InitDriver(dbSource)
			if nil != err {
				panic(err)
			}
		}
		break
	}
	return fpmsmg
}
