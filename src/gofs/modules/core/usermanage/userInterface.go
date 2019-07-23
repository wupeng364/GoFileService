package usermanage

/**
 * 用户管理接口
 */
type userManageInterface interface{
	InitDriver(db interface{}) error	// 初始化驱动指针
	InitTables( ) error					// 创建初始表和数据
	ListAllUsers( )(*[]UserInfo, error)	// 列出所有用户数据, 无分页
	QueryUser(userId string)(*UserInfo, error)	// 根据用户ID查询详细信息
	AddUser(user *UserInfo) error		// 添加用户
	UpdateUser(user *UserInfo) error	// 修改用户
	DelUser(userId string) error		// 根据userId删除用户
	CheckPwd(userId, pwd string) bool	// 校验密码是否一致
	
}