// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package usermanage

// 用户管理接口
type userManageInterface interface {
	InitDriver(db interface{}) error            // 初始化驱动指针
	InitTables() error                          // 创建初始表和数据
	ListAllUsers() (*[]UserInfo, error)         // 列出所有用户数据, 无分页
	QueryUser(userID string) (*UserInfo, error) // 根据用户ID查询详细信息
	AddUser(user *UserInfo) error               // 添加用户
	UpdateUser(user *UserInfo) error            // 修改用户
	DelUser(userID string) error                // 根据userID删除用户
	CheckPwd(userID, pwd string) bool           // 校验密码是否一致
	UpdatePWD(user *UserInfo) error             // 修改用户密码
}
