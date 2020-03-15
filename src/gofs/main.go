package main

import (
	"flag"
	"fmt"
	"gofs/comm/conf"
	"gofs/comm/httpserver"
	"gofs/comm/sqlite"
	"gofs/data/filemanage"
	"gofs/data/usermanage"
	"gofs/extend/htmlpage"
	"gofs/service/restful/fileapi"
	"gofs/service/restful/signature"
	"gofs/service/restful/userapi"
	"gutils/mloader"
	"gutils/strtool"
	"net/http"
	"path/filepath"
)

func main() {
	// 读取配置
	debug := flag.String("debug", "true", "The name of the configuration to load, default gofs")
	name := flag.String("name", "gofs", "The name of the configuration to load, default gofs")
	flag.Parse()

	savePath, _ := filepath.Abs("./conf/modules")
	mloader, err := mloader.NewAsJSONRecorder(savePath)
	if nil != err {
		panic(err)
	}
	// 设置全局配置
	mloader.SetParam("DEBUG", strtool.String2Bool(*debug))
	mloader.SetParam("app.name", *name)
	// 加载模块
	mloader.Loads(new(conf.AppConf), new(httpserver.HTTPServer), new(sqlite.SqliteConn))
	mloader.Loads(new(filemanage.FileManager), new(usermanage.UserManager))
	mloader.Loads(new(signature.Signature), new(userapi.UserAPI), new(fileapi.FileAPI))
	mloader.Loads(new(htmlpage.HTMLPage))
	// 启动服务
	res, err := mloader.Invoke("HTTPServer", "DoStartServer", &http.Server{
		ReadTimeout:    0,
		WriteTimeout:   0,
		MaxHeaderBytes: 1 << 20,
	})
	fmt.Println(res, err)
}
