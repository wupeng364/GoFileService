package main

/**
 *@description 文件管理器启动入口
 *@author	wupeng364@outlook.com
 */
import (
	GModule "common/gomodule"
	"fmt"
	"modules/apis/fileapi"
	"modules/apis/userapi"
	"modules/common/config"
	"modules/common/httpserver"
	"modules/common/sqlite"
	"modules/core/signature"
	"modules/core/filemanage"
	"modules/core/usermanage"
	"modules/extend/htmlpage"
)

func main() {
	// 加载模块&监听端口
	{
		// 加载基础模块
		GModule.LoadModule(&config.ConfigModule{})
		GModule.LoadModule(&sqlite.SqliteModule{})
		GModule.LoadModule(&httpserver.HttpServerModule{})
		// 加载业务模块
		GModule.LoadModule(&filemanage.FileManageModule{})
		GModule.LoadModule(&usermanage.UserManageModule{})
		GModule.LoadModule(&signature.SignatureModule{})
		// 加载Api网络模块
		GModule.LoadModule(&fileapi.FsApiModule{})
		GModule.LoadModule(&userapi.UserApiModule{})
		// 加载拓展模块
		GModule.LoadModule(&htmlpage.HtmlModule{})

		// 启动监听
		fmt.Println(GModule.Invoke("HttpServerModule", "DoStartServer")[0].Interface().(error))
	}
}
