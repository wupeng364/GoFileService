package main

/**
 *@description 文件管理器启动入口
 *@author	wupeng364@outlook.com
 */
import (
	"fmt"
	"gofs/common/moduleloader"
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
		mloader := moduleloader.New("gofs")
		// 加载基础模块
		mloader.Loads(&sqlite.SqliteModule{}, &httpserver.HttpServerModule{})
		// 加载业务模块
		mloader.Loads(&filemanage.FileManageModule{}, &usermanage.UserManageModule{}, &signature.SignatureModule{})
		// 加载Api网络模块
		mloader.Loads(&fileapi.FsApiModule{}, &userapi.UserApiModule{})
		// 加载拓展模块
		mloader.Load(&htmlpage.HtmlModule{})

		// 启动监听
		fmt.Println(mloader.Invoke("HttpServerModule", "DoStartServer")[0].Interface().(error))
	}
}
