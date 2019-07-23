package filemanage

const(
	tokenExpired_Second = 60	// 复制操作令牌失效时间
)
// 复制回调
type CopyCallback func(count int, src, dst string, err *CopyError)error
type MoveCallback func(count int, src, dst string, err *MoveError)error
type CopyError struct{
	IsSrcExist	bool
	IsDstExist	bool
	ErrorString string
}
type MoveError CopyError
// 文件基础属性
type F_BaseInfo struct{
	Path   string
	CtTime int64
	IsFile bool
	FileSize	int64
}