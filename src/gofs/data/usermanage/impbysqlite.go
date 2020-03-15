// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 用sqlite3存放数据

package usermanage

import (
	"fmt"
	"gofs/comm/sqlite"
	"gutils/strtool"
	"time"
)

// sqlite3实现
type impBySqlite struct {
	db *sqlite.SqliteConn
}

// InitDriver 初始化驱动
func (sqlti *impBySqlite) InitDriver(db interface{}) error {
	if nil == db {
		return ErrorConnIsNil
	}
	sqlti.db = db.(*sqlite.SqliteConn)
	return nil
}

// InitTables 初始化 users 表
// 添加默认用户 admin, 密码为空
func (sqlti *impBySqlite) InitTables() error {
	// 打开数据库
	db, err := sqlti.db.Open()
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
		`CREATE TABLE IF NOT EXISTS users(
		    userid VARCHAR(64) PRIMARY KEY,
		    userpwd VARCHAR(255) NULL,
		    usertype VARCHAR(255) NULL,
		    username VARCHAR(255) NULL,
		    cttime DATE NULL
		);`)
	if err != nil {
		ts.Rollback()
		return err
	}

	// 插入默认用户 - admin
	stmt, err := ts.Prepare("INSERT INTO users(userid,username,usertype,userpwd,cttime) values(?,?,?,?,?)")
	if err != nil {
		ts.Rollback()
		return err
	}
	_, err = stmt.Exec("admin", "管理员", AdminRole, strtool.GetMD5(""), time.Now())
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

// ListAllUsers 列出所有用户数据, 无分页
func (sqlti *impBySqlite) ListAllUsers() (*[]UserInfo, error) {
	// 打开数据库
	db, err := sqlti.db.Open()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT userid,username,usertype,cttime FROM users")
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	//
	res := make([]UserInfo, 0)
	for rows.Next() {
		user := UserInfo{}
		err := rows.Scan(&user.UserID, &user.UserName, &user.UserType, &user.CtTime)
		if err != nil {
			return nil, err
		}
		res = append(res, user)
	}
	return &res, nil
}

// QueryUser 根据用户ID查询详细信息
func (sqlti *impBySqlite) QueryUser(userID string) (*UserInfo, error) {
	// 打开数据库
	db, err := sqlti.db.Open()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT userid,username,usertype,cttime FROM users where userid=" + userID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	//
	for rows.Next() {
		user := UserInfo{}
		err := rows.Scan(&user.UserID, &user.UserName, &user.UserType, &user.CtTime)
		if err != nil {
			return nil, err
		}
		return &user, nil
	}
	return nil, nil
}

// AddUser 添加用户
func (sqlti *impBySqlite) AddUser(user *UserInfo) error {
	if len(user.UserID) == 0 {
		return ErrorUserIDIsNil
	}
	if len(user.UserName) == 0 {
		return ErrorUserNameIsNil
	}
	// 打开数据库
	db, err := sqlti.db.Open()
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
	stmt, err := ts.Prepare("INSERT INTO users(userid,username,usertype,userpwd,cttime) values(?,?,?,?,?)")
	if err != nil {
		ts.Rollback()
		return err
	}
	_, err = stmt.Exec(user.UserID, user.UserName, user.UserType, strtool.GetMD5(user.UserPWD), time.Now())
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

// UpdateUser 修改用户
func (sqlti *impBySqlite) UpdateUser(user *UserInfo) error {
	if len(user.UserID) == 0 {
		return ErrorUserIDIsNil
	}

	// 查询旧数据 - 不校验
	//	user_old, err := sqlti.QueryUser(user.UserID)
	//	if nil != err { return err }
	//	if nil == user_old { return ErrorUserNotExist }

	// 打开数据库
	db, err := sqlti.db.Open()
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
	stmt, err := ts.Prepare("UPDATE users SET username=?, usertype=? WHERE userid=?")
	if err != nil {
		ts.Rollback()
		return err
	}
	_, err = stmt.Exec(user.UserName, user.UserType, user.UserID)
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

// DelUser 根据userId删除用户
func (sqlti *impBySqlite) DelUser(userID string) error {
	if len(userID) == 0 {
		return ErrorUserIDIsNil
	}

	// 打开数据库
	db, err := sqlti.db.Open()
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
	stmt, err := ts.Prepare("DELETE FROM users WHERE userid = ?")
	if err != nil {
		ts.Rollback()
		return err
	}
	_, err = stmt.Exec(userID)
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

// CheckPwd 校验密码是否一致
func (sqlti *impBySqlite) CheckPwd(userID, pwd string) bool {
	if len(userID) == 0 {
		return false
	}

	// 打开数据库
	db, err := sqlti.db.Open()
	if err != nil {
		return false
	}
	defer db.Close()

	//
	rows, err := db.Query("SELECT userid FROM users where userid='" + userID + "' and userpwd='" + strtool.GetMD5(pwd) + "'")
	defer rows.Close()
	if nil == rows {
		fmt.Println(err)
		return false
	}
	//
	if rows.Next() {
		return true
	}
	return false
}

// UpdateUser 修改用户密码
func (sqlti *impBySqlite) UpdatePWD(user *UserInfo) error {
	if len(user.UserID) == 0 {
		return ErrorUserIDIsNil
	}

	// 打开数据库
	db, err := sqlti.db.Open()
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
	stmt, err := ts.Prepare("UPDATE users SET uuserpwd=? WHERE userid=?")
	if err != nil {
		ts.Rollback()
		return err
	}
	_, err = stmt.Exec(strtool.GetMD5(user.UserPWD), user.UserID)
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
