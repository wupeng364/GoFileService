package fileapimodel

import(
	"errors"
)
// 错误处理支持的操作
// 忽略, 忽略全部
type errorOperationStruct struct{
	ignore 		string
	ignoreall 	string
	replace 	string
	replaceall  string
	discontinue string
}
// 申请Token的类型
type tokenTypeStruct struct{
	copyDir 	string
	moveDir		string
	download	string
}
// 复制文件Token保存对象
type FileBatchOperationTokenObject struct{
	CountIndex		int64
	ErrorString		string
	Src				string
	Dst				string
	IsSrcExist		bool
	IsDstExist		bool
	IsReplace		bool
	IsReplaceAll	bool
	IsIgnore		bool
	IsIgnoreAll		bool
	IsComplete		bool
	IsDiscontinue	bool
}
// 上传下载文件零时保存的数据
type FileTransferToken struct{
	FilePath        string
}

// 定义好的错误可选项目
var ErrorOperation = &errorOperationStruct{
	ignore: "ignore",
	ignoreall: "ignoreall",
	replace: "replace",
	replaceall: "replaceall", 
	discontinue: "discontinue",
}
// 定义好的Token类型
var TokenType = &tokenTypeStruct{
	copyDir: "copy_dir",
	moveDir: "move_dir",
	download: "download",
}
var Error_Discontinue = errors.New("Discontine")
var Error_OprationExpires = errors.New("Opration expires")
var Error_OprationFailed = errors.New("Opration failed")
var Error_OprationUnknown = errors.New("Opration unknown")
var Error_FileNotExist   = errors.New("file does not exist")
var Error_ParentFolderNotExist   = errors.New("parent folder does not exist")
var Error_NewNameIsEmpty = errors.New("New name cannot be empty")
