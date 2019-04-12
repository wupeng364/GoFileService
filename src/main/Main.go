package main
/**
 *@description 文件管理器启动入口
 *@author	wupeng364@outlook.com
*/
import (
	"fmt"
	GModel "common/gomodule"
	"modules/configmodel"
	"modules/httpservermodel"
	"modules/filemanage"
	"modules/fileapimodel"
	"modules/htmlmodel"
)

func main(){
	// 加载模块&监听端口
	go func( ){
		// 加载模块
	    GModel.LoadModel( &configmodel.ConfigModel{} )
	    GModel.LoadModel( &httpservermodel.HttpServerModel{} )
	    GModel.LoadModel( &filemanage.FileManageModel{} )
	    GModel.LoadModel( &htmlmodel.Htmlmodel{} )
	    GModel.LoadModel( &fileapimodel.FsApimodel{} )
	    fmt.Println("\r\n\r\n")
	    
	    // 启动监听
		fmt.Println( GModel.Invoke("HttpServerModel", "DoStartServer")[0].Interface( ).(error) )
	}( )
	
	var sc string
	for{
		if sc == "exit" {
			break
		}
		fmt.Scan(&sc)
		fmt.Println("Input 'exit' to exit this program.")
	}
}