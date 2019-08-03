package usermanage
/**
 *@description 用户管理模块
 *
 *@author	wupeng364@outlook.com
*/
import (
	//"fmt"
	"gofs/common/moduletools"
	"gofs/modules/common/config"
)

type UserManageModule struct{
	cfg *config.ConfigModule
	umi  userManageInterface
}

// 返回模块信息
func (this *UserManageModule)MInfo( )(*moduletools.ModuleInfo) {
	return &moduletools.ModuleInfo{
		this,
		"UserManageModule",
		1.0,
		"用户管理模块",
	}
}

// 模块安装, 一个模块只初始化一次
func (this *UserManageModule)OnMSetup( ref moduletools.ReferenceModule ) {
	cfg    := ref(this.cfg).(*config.ConfigModule)
	dbType := cfg.GetConfig(cfc_db_type)
	umi    := getUserManageInterfaceImp(dbType, ref)
	
	// 执行建库、建表
	err := umi.InitTables( )
	if nil != err {
		panic( err )
	}
}
// 模块升级, 一个版本执行一次
func (this *UserManageModule)OnMUpdate( ref moduletools.ReferenceModule ) {
	
}

// 每次启动加载模块执行一次
func (this *UserManageModule)OnMInit( ref moduletools.ReferenceModule ) {
	this.cfg = ref(this.cfg).(*config.ConfigModule)
	dbType := this.cfg.GetConfig(cfc_db_type)
	this.umi = getUserManageInterfaceImp(dbType, ref)
}

// 系统执行销毁时执行
func (this *UserManageModule)OnMDestroy( ref moduletools.ReferenceModule ) {
	
}

// 列出所有用户数据, 无分页
func (this *UserManageModule)ListAllUsers( )(*[]UserInfo, error){
	return this.umi.ListAllUsers( )
}
// 根据用户ID查询详细信息
func (this *UserManageModule)QueryUser(userId string)(*UserInfo, error){
	return this.umi.QueryUser( userId )
}
// 添加用户
func (this *UserManageModule)AddUser(user *UserInfo) error{
	return this.umi.AddUser( user )
}
// 修改用户密码
func (this *UserManageModule)UpdateUserPwd(userId, pwd string) error{
	user_old, err := this.QueryUser(userId)
	if nil != err {
		return err
	}
	if nil == user_old {
		return Error_UserNotExist
	}
	user_old.Userpwd = pwd
	return this.umi.UpdateUser( user_old)
}
// 修改用户昵称
func (this *UserManageModule)UpdateUserName(userId, userName string) error{
	if len(userName) == 0 {
		return Error_UserNameIsNil
	}
	user_old, err := this.QueryUser(userId)
	if nil != err {
		return err
	}
	if nil == user_old {
		return Error_UserNotExist
	}
	user_old.Username = userName
	return this.umi.UpdateUser( user_old)
}
// 根据userId删除用户
func (this *UserManageModule)DelUser(userId string) error{
	return this.umi.DelUser( userId )
}
// 校验密码是否一致
func (this *UserManageModule)CheckPwd(userId, pwd string) bool{
	return this.umi.CheckPwd( userId, pwd )
}
// ==============================================================================================
// 获取启动类型, 并实例对象指针
func getUserManageInterfaceImp( dbType string, ref moduletools.ReferenceModule )userManageInterface{
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