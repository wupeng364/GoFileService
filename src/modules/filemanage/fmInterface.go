package filemanage

import(
	"io"
)

/**
 * 文件基础操作接口
 */
type fmInterface interface{
	
	IsExist( relativePath string )(bool, error)
	IsDir( relativePath string )(bool, error)
	IsFile( relativePath string )(bool, error)
	
	GetDirList( relativePath string )([]string, error)
	GetFileSize( relativePath string )(int64, error)
	GetCreateTime( relativePath string )(int64, error)
	GetModifyTime( relativePath string )(int64, error)
	GetCreateUser( relativePath string )(string, error)
	GetModifyUser( relativePath string )(string, error)
	GetFileLatestVersion( relativePath string )(string, error)
	GetFileVersionList( relativePath string )([]string, error)
	
	DoNewFolder(path string)error
	DoRename(src string, dest string)error
	DoMove(src, dest string, repalce, ignore bool, callback MoveCallback)error
	DoDelete( relativePath string )error
	
	DoRead( relativePath string )(io.Reader, error)
	DoWrite( relativePath string, ioReader io.Reader )error
	DoCopy(src, dst string, replace, ignore bool, callback CopyCallback)error
	
	// 传输Token
	AskUploadToken( relativePath string )(string, error)
	SaveTokenFile(token string, src io.Reader)(bool)
	SubmitToken(token string, isCreateVer bool)(bool, error)
}