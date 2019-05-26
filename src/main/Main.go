package main
/**
 *@description 文件管理器启动入口
 *@author	wupeng364@outlook.com
*/
import (
	"fmt"
	GModule "common/gomodule"
	"modules/config"
	"modules/httpserver"
	"modules/filemanage"
	"modules/fileapi"
	"modules/htmlpage"
)

func main(){
	// 加载模块&监听端口
	{
		// 加载模块
	    GModule.LoadModule( &config.ConfigModule{} )
	    GModule.LoadModule( &httpserver.HttpServerModule{} )
	    GModule.LoadModule( &filemanage.FileManageModule{} )
	    GModule.LoadModule( &htmlpage.HtmlModule{} )
	    GModule.LoadModule( &fileapi.FsApiModule{} )
	    fmt.Println("\r\n\r\n")
	    
	    // 启动监听
		fmt.Println( GModule.Invoke("HttpServerModule", "DoStartServer")[0].Interface( ).(error) )
	}
}