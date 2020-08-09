// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 用sqlite3存放数据

package filepermission

import (
	"database/sql"
	_ "go-sqlite3"
	"gofs/data/usermanage"
	"gutils/strtool"
	"time"
)

// sqlite3实现
type impBySqlite struct {
	dbSource string
}

// InitDriver 初始化驱动
func (sqlti *impBySqlite) InitDriver(db interface{}) error {
	if nil == db {
		return ErrorConnIsNil
	}
	sqlti.dbSource = db.(string)
	return nil
}
func (sqlti *impBySqlite) Open() (*sql.DB, error) {
	return sql.Open("sqlite3", sqlti.dbSource)
}

// InitTables 初始化 filepermissions 表
// 添加默认权限 admin, 密码为空
func (sqlti *impBySqlite) InitTables() error {
	// 打开数据库
	db, err := sqlti.Open()
	if err != nil {
		return err
	}
	defer db.Close()

	// 开启事务
	ts, err := db.Begin()
	if err != nil {
		return err
	}
	// 建表 - userinfo
	_, err = ts.Exec(
		`CREATE TABLE IF NOT EXISTS filepermissions(
		    permissionid VARCHAR(64) PRIMARY KEY,
		    path TEXT(1000) NULL,
		    userid VARCHAR(64) NULL,
		    permission INTEGER(255) NULL,
		    cttime DATE NULL
		);`)
	if err != nil {
		ts.Rollback()
		return err
	}

	// 插入默认权限 - admin
	stmt, err := ts.Prepare("INSERT INTO filepermissions(permissionid,path,userid,permission,cttime) values(?,?,?,?,?)")
	if err != nil {
		ts.Rollback()
		return err
	}
	allPermission := (1 << VISIBLE) + (1 << READ) + (1 << WRITE)
	_, err = stmt.Exec(strtool.GetUUID(), "/", usermanage.AdminUserID, allPermission, time.Now())
	if err != nil {
		ts.Rollback()
		return err
	}

	// 提交更改
	err = ts.Commit()
	if err != err {
		ts.Rollback()
		return err
	}
	return nil
}

// ListFPermissions 列出所有权限数据, 无分页
func (sqlti *impBySqlite) ListFPermissions() ([]FPermissionInfo, error) {
	// 打开数据库
	db, err := sqlti.Open()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT permissionid, path, userid, permission, cttime FROM filepermissions")
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	//
	res := make([]FPermissionInfo, 0)
	for rows.Next() {
		fpms := FPermissionInfo{}
		err := rows.Scan(&fpms.PermissionID, &fpms.Path, &fpms.UserID, &fpms.Permission, &fpms.CtTime)
		if err != nil {
			return nil, err
		}
		res = append(res, fpms)
	}
	return res, nil
}

// ListUserFPermissions 列出用户所有权限数据, 无分页
func (sqlti *impBySqlite) ListUserFPermissions(userID string) ([]FPermissionInfo, error) {
	// 打开数据库
	db, err := sqlti.Open()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT permissionid, path, userid, permission, cttime FROM filepermissions WHERE userid = '" + userID + "'")
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	//
	res := make([]FPermissionInfo, 0)
	for rows.Next() {
		fpms := FPermissionInfo{}
		err := rows.Scan(&fpms.PermissionID, &fpms.Path, &fpms.UserID, &fpms.Permission, &fpms.CtTime)
		if err != nil {
			return nil, err
		}
		res = append(res, fpms)
	}
	return res, nil
}

// QueryFPermission 根据权限ID查询详细信息
func (sqlti *impBySqlite) QueryFPermission(permissionID string) (*FPermissionInfo, error) {
	// 打开数据库
	db, err := sqlti.Open()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT permissionid, path, userid, permission, cttime FROM filepermissions where permissionid='" + permissionID + "'")
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	//
	for rows.Next() {
		fpms := FPermissionInfo{}
		err := rows.Scan(&fpms.PermissionID, &fpms.Path, &fpms.UserID, &fpms.Permission, &fpms.CtTime)
		if err != nil {
			return nil, err
		}
		return &fpms, nil
	}
	return nil, nil
}

// AddFPermission 添加权限
func (sqlti *impBySqlite) AddFPermission(fpms FPermissionInfo) error {
	if len(fpms.UserID) == 0 {
		return ErrorUserIDIsNil
	}
	if len(fpms.Path) == 0 {
		return ErrorPermissionPathIsNil
	}
	if fpms.Permission <= 0 {
		return ErrorPermissionIsNil
	}
	// 打开数据库
	db, err := sqlti.Open()
	if err != nil {
		return err
	}
	defer db.Close()

	// 开启事务
	ts, err := db.Begin()
	if err != nil {
		return err
	}

	//
	stmt, err := ts.Prepare("INSERT INTO filepermissions(permissionid, path, userid, permission, cttime) values(?,?,?,?,?)")
	if err != nil {
		ts.Rollback()
		return err
	}
	_, err = stmt.Exec(strtool.GetUUID(), fpms.Path, fpms.UserID, fpms.Permission, time.Now())
	if err != nil {
		ts.Rollback()
		return err
	}

	// 提交
	err = ts.Commit()
	if nil != err {
		ts.Rollback()
		return err
	}
	return nil
}

// UpdateFPermission 修改权限
func (sqlti *impBySqlite) UpdateFPermission(fpms FPermissionInfo) error {
	// if len(fpms.UserID) == 0 {
	// 	return ErrorUserIDIsNil
	// }
	if len(fpms.PermissionID) == 0 {
		return ErrorPermissionIDIsNil
	}
	// if len(fpms.Path) == 0 {
	// 	return ErrorPermissionPathIsNil
	// }
	if fpms.Permission <= 0 {
		return ErrorPermissionIsNil
	}

	// 打开数据库
	db, err := sqlti.Open()
	if err != nil {
		return err
	}
	defer db.Close()

	// 开启事务
	ts, err := db.Begin()
	if err != nil {
		return err
	}

	// permissionid, path, userid, permission, cttime
	stmt, err := ts.Prepare("UPDATE filepermissions SET permission=? WHERE permissionid=?")
	if err != nil {
		ts.Rollback()
		return err
	}
	_, err = stmt.Exec(fpms.Permission, fpms.PermissionID)
	if err != nil {
		ts.Rollback()
		return err
	}

	// 提交
	err = ts.Commit()
	if nil != err {
		ts.Rollback()
		return err
	}
	return nil
}

// DelFPermission 删除权限
func (sqlti *impBySqlite) DelFPermission(permissionID string) error {
	if len(permissionID) == 0 {
		return ErrorPermissionIDIsNil
	}

	// 打开数据库
	db, err := sqlti.Open()
	if err != nil {
		return err
	}
	defer db.Close()

	// 开启事务
	ts, err := db.Begin()
	if err != nil {
		return err
	}

	//
	stmt, err := ts.Prepare("DELETE FROM filepermissions WHERE permissionid = ?")
	if err != nil {
		ts.Rollback()
		return err
	}
	_, err = stmt.Exec(permissionID)
	if err != nil {
		ts.Rollback()
		return err
	}

	// 提交
	err = ts.Commit()
	if nil != err {
		ts.Rollback()
		return err
	}
	return nil
}
