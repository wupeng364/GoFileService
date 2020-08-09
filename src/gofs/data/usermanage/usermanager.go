// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 用户管理模块

package usermanage

import (
	"gofs/base/conf"
	"gutils/mloader"
	"path/filepath"
)

// UserManager 用户管理模块
type UserManager struct {
	umi UserManage
}

// ModuleOpts 模块加载器接口实现, 返回模块信息&配置
func (umg *UserManager) ModuleOpts() mloader.Opts {
	return mloader.Opts{
		Name:        "UserManager",
		Version:     1.1,
		Description: "用户管理模块",
		OnReady: func(mctx *mloader.Loader) {
			umg.umi = umg.getUserManageInterfaceImp(mctx)
		},
		OnSetup: func() {
			// 执行建库、建表
			err := umg.umi.InitTables()
			if nil != err {
				panic(err)
			}
		},
		OnUpdate: func(hv float64) {
			doUpdateModule(hv, umg)
		},
	}
}

// ListAllUsers 列出所有用户数据, 无分页
func (umg *UserManager) ListAllUsers() (*[]UserInfo, error) {
	return umg.umi.ListAllUsers()
}

// QueryUser 根据用户ID查询详细信息
func (umg *UserManager) QueryUser(userID string) (*UserInfo, error) {
	return umg.umi.QueryUser(userID)
}

// AddUser 添加用户
func (umg *UserManager) AddUser(user *UserInfo) error {
	return umg.umi.AddUser(user)
}

// UpdateUserPwd 修改用户密码
func (umg *UserManager) UpdateUserPwd(userID, pwd string) error {
	userOld, err := umg.QueryUser(userID)
	if nil != err {
		return err
	}
	if nil == userOld {
		return ErrorUserNotExist
	}
	userOld.UserPWD = pwd
	return umg.umi.UpdatePWD(userOld)
}

// UpdateUserName 修改用户昵称
func (umg *UserManager) UpdateUserName(userID, userName string) error {
	if len(userName) == 0 {
		return ErrorUserNameIsNil
	}
	userOld, err := umg.QueryUser(userID)
	if nil != err {
		return err
	}
	if nil == userOld {
		return ErrorUserNotExist
	}
	userOld.UserName = userName
	return umg.umi.UpdateUser(userOld)
}

// DelUser 根据userID删除用户
func (umg *UserManager) DelUser(userID string) error {
	return umg.umi.DelUser(userID)
}

// CheckPwd 校验密码是否一致
func (umg *UserManager) CheckPwd(userID, pwd string) bool {
	return umg.umi.CheckPwd(userID, pwd)
}

// getUserManageInterfaceImp 获取启动类型, 并实例对象指针
func (umg *UserManager) getUserManageInterfaceImp(mctx *mloader.Loader) UserManage {
	// 默认使用sqlite驱动
	dbType := mctx.GetParam(conf.DataBasType).ToString("")
	if len(dbType) == 0 {
		panic(conf.DataBasType + " is empty")
	}
	dbSource := mctx.GetParam(conf.DataBaseSource).ToString("")
	if len(dbSource) == 0 {
		panic(conf.DataBaseSource + " is empty")
	}
	var umi UserManage
	switch dbType {
	case conf.DefaultDataBaseType:
		{
			if !filepath.IsAbs(dbSource) {
				dbSource, _ = filepath.Abs(dbSource)
			}
			umi = &impBySqlite{}
			err := umi.InitDriver(dbSource)
			if nil != err {
				panic(err)
			}
		}
		break
	}
	return umi
}
