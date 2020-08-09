// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package filepermission

// FPmsManage 文件权限管理接口
type FPmsManage interface {
	InitDriver(db interface{}) error                                // 初始化驱动指针
	InitTables() error                                              // 创建初始表和数据
	ListFPermissions() ([]FPermissionInfo, error)                   // 列出所有权限数据, 无分页
	ListUserFPermissions(userID string) ([]FPermissionInfo, error)  // 列出用户所有权限数据, 无分页
	QueryFPermission(permissionID string) (*FPermissionInfo, error) // 根据权限ID查询详细信息
	AddFPermission(fpm FPermissionInfo) error                       // 添加权限
	UpdateFPermission(fpm FPermissionInfo) error                    // 修改权限
	DelFPermission(permissionID string) error                       // 根据permissionID删除权限
}

// FPmsManageCheck 文件权限校验
type FPmsManageCheck interface {
	HashPermission(path string, permission int64)
}
