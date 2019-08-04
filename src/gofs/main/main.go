package main

/**
 *@description 文件管理器启动入口
 *@author	wupeng364@outlook.com
 */
import (
	"fmt"
	"gofs/common/moduletools"
	"gofs/modules/common/config"
	"gofs/modules/common/httpserver"
	"gofs/modules/common/sqlite"
	"gofs/modules/apis/fileapi"
	"gofs/modules/apis/userapi"
	"gofs/modules/core/signature"
	"gofs/modules/core/filemanage"
	"gofs/modules/core/usermanage"
	"gofs/modules/extend/htmlpage"
)

func main() {
	// 加载模块&监听端口
	{
		// 加载基础模块
		moduletools.LoadModule(&config.ConfigModule{})
		moduletools.LoadModule(&sqlite.SqliteModule{})
		moduletools.LoadModule(&httpserver.HttpServerModule{})
		// 加载业务模块
		moduletools.LoadModule(&filemanage.FileManageModule{})
		moduletools.LoadModule(&usermanage.UserManageModule{})
		moduletools.LoadModule(&signature.SignatureModule{})
		// 加载Api网络模块
		moduletools.LoadModule(&fileapi.FsApiModule{})
		moduletools.LoadModule(&userapi.UserApiModule{})
		// 加载拓展模块
		moduletools.LoadModule(&htmlpage.HtmlModule{})

		// 启动监听
		fmt.Println(moduletools.Invoke("HttpServerModule", "DoStartServer")[0].Interface().(error))
	}
}
