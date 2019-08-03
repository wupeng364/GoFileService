package usermanage

/**
 *@description 用sqlite3存放数据
 *@author	wupeng364@outlook.com
*/
import (
	"fmt"
	"time"
	"gofs/common/stringtools"
	"gofs/modules/common/sqlite"
)

type imp_sqlite struct{
	db *sqlite.SqliteModule
}

// 初始化驱动
func (this *imp_sqlite)InitDriver(db interface{}) error{
	if nil == db {
		return Error_ConnIsNil
	}
	this.db = db.(*sqlite.SqliteModule)
	return nil
}

// 初始化 users 表
// 添加默认用户 admin, 密码为空
func (this *imp_sqlite)InitTables( )error{
	// 打开数据库
	db, err := this.db.Open( )
	defer db.Close( )
	if err != nil { return err }
	
	// 开启事务
	ts, err := db.Begin( )
	if err != nil { return err }
	
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
    	ts.Rollback( )
    	return err 
    }
    
    // 插入默认用户 - admin
    stmt, err := ts.Prepare("INSERT INTO users(userid,username,usertype,userpwd,cttime) values(?,?,?,?,?)")
    if err != nil {
        ts.Rollback( )
    	return err
    }
    _, err = stmt.Exec("admin", "管理员", UserType_Admin, stringtools.String2MD5(""), time.Now( ))
    if err != nil {
        ts.Rollback( )
    	return err
    }
	
	// 提交更改
    err = ts.Commit()
    if err != err {
	    ts.Rollback( )
    	return err
    }
    return nil
}
// 列出所有用户数据, 无分页
func (this *imp_sqlite)ListAllUsers( )(*[]UserInfo, error) {
	// 打开数据库
	db, err := this.db.Open( )
	defer db.Close( )
	if err != nil { return nil, err }
	
	rows, err := db.Query("SELECT userid,username,usertype,cttime FROM users")
    defer rows.Close( )
    if err != nil {  return nil, err }

	// 
	res := make([]UserInfo, 0)
    for rows.Next( ) {
        user := UserInfo{}
        err := rows.Scan(&user.UserId, &user.Username, &user.UserType, &user.Cttime)
        if err != nil {  return nil, err }
        res = append(res, user)
    }
	return &res, nil
}
// 根据用户ID查询详细信息
func (this *imp_sqlite)QueryUser(userId string)(*UserInfo, error) {
	// 打开数据库
	db, err := this.db.Open( )
	defer db.Close( )
	if err != nil { return nil, err }
	
	rows, err := db.Query("SELECT userid,username,usertype,cttime FROM users where userid="+userId)
    defer rows.Close( )
    if err != nil {  return nil, err }

	// 
    for rows.Next( ) {
		user := UserInfo{}
        err := rows.Scan(&user.UserId, &user.Username, &user.UserType, &user.Cttime)
        if err != nil {  return nil, err }
        return &user, nil
    }
	return nil, nil
}
// 添加用户
func (this *imp_sqlite)AddUser(user *UserInfo) error {
	if len(user.UserId) == 0 {
		return Error_UserIdIsNil
	}
	if len(user.Username) == 0 {
		return Error_UserNameIsNil
	}
	// 打开数据库
	db, err := this.db.Open( )
	defer db.Close( )
	if err != nil { return err }
	
	// 开启事务
	ts, err := db.Begin( )
	if err != nil { return err }
	
    // 
    stmt, err := ts.Prepare("INSERT INTO users(userid,username,usertype,userpwd,cttime) values(?,?,?,?,?)")
    if err != nil {
        ts.Rollback( )
    	return err
    }
    _, err = stmt.Exec(user.UserId, user.Username, user.UserType, stringtools.String2MD5(user.Userpwd), time.Now( ))
    if err != nil {
        ts.Rollback( )
    	return err
    }
    
    // 提交
    err = ts.Commit( )
    if nil != err {
	    ts.Rollback( )
	    return err
    }
	return nil
}
// 修改用户
func (this *imp_sqlite)UpdateUser(user *UserInfo) error {
	if len(user.UserId) == 0 {
		return Error_UserIdIsNil
	}
	
	// 查询旧数据 - 不校验
//	user_old, err := this.QueryUser(user.UserId)
//	if nil != err { return err }
//	if nil == user_old { return Error_UserNotExist }
	
	// 打开数据库
	db, err := this.db.Open( )
	defer db.Close( )
	if err != nil { return err }
	
	// 开启事务
	ts, err := db.Begin( )
	if err != nil { return err }
	
    // 
    stmt, err := ts.Prepare("UPDATE users SET username=?, usertype=?, userpwd=? WHERE userid = ?")
    if err != nil {
        ts.Rollback( )
    	return err
    }
    _, err = stmt.Exec(user.Username, user.UserType, stringtools.String2MD5(user.Userpwd), user.UserId)
    if err != nil {
        ts.Rollback( )
    	return err
    }
    
    // 提交
    err = ts.Commit( )
    if nil != err {
	    ts.Rollback( )
	    return err
    }
	return nil
}
// 根据userId删除用户
func (this *imp_sqlite)DelUser(userId string) error {
	if len(userId) == 0 {
		return Error_UserIdIsNil
	}
	
	// 打开数据库
	db, err := this.db.Open( )
	defer db.Close( )
	if err != nil { return err }
	
	// 开启事务
	ts, err := db.Begin( )
	if err != nil { return err }
	
    // 
    stmt, err := ts.Prepare("DELETE FROM users WHERE userid = ?")
    if err != nil {
        ts.Rollback( )
    	return err
    }
    _, err = stmt.Exec(userId)
    if err != nil {
        ts.Rollback( )
    	return err
    }
    
    // 提交
    err = ts.Commit( )
    if nil != err {
	    ts.Rollback( )
	    return err
    }
	return nil
}
// 校验密码是否一致
func (this *imp_sqlite)CheckPwd(userId, pwd string) bool {
	if len(userId) == 0 {
		return false
	}
	
	// 打开数据库
	db, err := this.db.Open( )
	defer db.Close( )
	if err != nil { return false }
	
	// 
	rows, err := db.Query("SELECT userid FROM users where userid='"+userId+"' and userpwd='"+stringtools.String2MD5(pwd)+"'")
    defer rows.Close( )
    if nil == rows {
		fmt.Println(err)
	    return false
    }
	// 
    if rows.Next( ) {
		return true
    }
	return false
}