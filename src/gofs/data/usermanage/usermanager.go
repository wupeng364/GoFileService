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

import "gutils/mloader"

// UserManager 用户管理模块
// 配置参数(mloader.GetParam): usermanager.dbtype
type UserManager struct {
	mctx *mloader.Loader
	umi  userManageInterface
}

// ModuleOpts 模块加载器接口实现, 返回模块信息&配置
func (umg *UserManager) ModuleOpts() mloader.Opts {
	return mloader.Opts{
		Name:        "UserManager",
		Version:     1.1,
		Description: "用户管理模块",
		OnReady: func(mctx *mloader.Loader) {
			umg.mctx = mctx
			dbType := umg.mctx.GetParam("usermanager.dbtype").ToString("sqlite")
			umg.umi = umg.getUserManageInterfaceImp(dbType)
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
func (umg *UserManager) getUserManageInterfaceImp(dbType string) userManageInterface {
	// 默认使用sqlite驱动
	if true {
		umi := &impBySqlite{}
		err := umi.InitDriver(umg.mctx.GetModuleByTemplate(umi.db))
		if nil != err {
			panic(err)
		}
		return umi
	}
	return nil
}
