package usermanage
/**
 *@description 用户管理模块
 *
 *@author	wupeng364@outlook.com
*/
import (
	//"fmt"
	"common/gomodule"
	"modules/common/config"
)

type UserManageModule struct{
	cfg *config.ConfigModule
	umi  userManageInterface
}

// 返回模块信息
func (um *UserManageModule)MInfo( )(*gomodule.ModuleInfo) {
	return &gomodule.ModuleInfo{
		um,
		"UserManageModule",
		1.0,
		"用户管理模块",
	}
}

// 模块安装, 一个模块只初始化一次
func (um *UserManageModule)OnMSetup( ref gomodule.ReferenceModule ) {
	cfg    := ref(um.cfg).(*config.ConfigModule)
	dbType := cfg.GetConfig(cfc_db_type)
	umi    := getUserManageInterfaceImp(dbType, ref)
	
	// 执行建库、建表
	err := umi.InitTables( )
	if nil != err {
		panic( err )
	}
}
// 模块升级, 一个版本执行一次
func (um *UserManageModule)OnMUpdate( ref gomodule.ReferenceModule ) {
	
}

// 每次启动加载模块执行一次
func (um *UserManageModule)OnMInit( ref gomodule.ReferenceModule ) {
	um.cfg = ref(um.cfg).(*config.ConfigModule)
	dbType := um.cfg.GetConfig(cfc_db_type)
	um.umi = getUserManageInterfaceImp(dbType, ref)
}

// 系统执行销毁时执行
func (um *UserManageModule)OnMDestroy( ref gomodule.ReferenceModule ) {
	
}

// 列出所有用户数据, 无分页
func (um *UserManageModule)ListAllUsers( )(*[]UserInfo, error){
	return um.umi.ListAllUsers( )
}
// 根据用户ID查询详细信息
func (um *UserManageModule)QueryUser(userId string)(*UserInfo, error){
	return um.umi.QueryUser( userId )
}
// 添加用户
func (um *UserManageModule)AddUser(user *UserInfo) error{
	return um.umi.AddUser( user )
}
// 修改用户密码
func (um *UserManageModule)UpdateUserPwd(userId, pwd string) error{
	user_old, err := um.QueryUser(userId)
	if nil != err {
		return err
	}
	if nil == user_old {
		return Error_UserNotExist
	}
	user_old.Userpwd = pwd
	return um.umi.UpdateUser( user_old)
}
// 修改用户昵称
func (um *UserManageModule)UpdateUserName(userId, userName string) error{
	if len(userName) == 0 {
		return Error_UserNameIsNil
	}
	user_old, err := um.QueryUser(userId)
	if nil != err {
		return err
	}
	if nil == user_old {
		return Error_UserNotExist
	}
	user_old.Username = userName
	return um.umi.UpdateUser( user_old)
}
// 根据userId删除用户
func (um *UserManageModule)DelUser(userId string) error{
	return um.umi.DelUser( userId )
}
// 校验密码是否一致
func (um *UserManageModule)CheckPwd(userId, pwd string) bool{
	return um.umi.CheckPwd( userId, pwd )
}
// ==============================================================================================
// 获取启动类型, 并实例对象指针
func getUserManageInterfaceImp( dbType string, ref gomodule.ReferenceModule )userManageInterface{
	// 默认使用sqlite驱动
	if true {
		umi := &imp_sqlite{}
		err := umi.InitDriver(ref(umi.db))
		if nil != err {
			panic( err )
		}
		return umi
	}
	return nil
}