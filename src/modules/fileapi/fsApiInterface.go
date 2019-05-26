package fileapi

import(
	"net/http"
)
/**
 * 文件基础操作网络接口
 */
type fsApiHttpInterface interface{
	init( )
	List(http.ResponseWriter, *http.Request)
	Del(http.ResponseWriter, *http.Request)
	DelVer(http.ResponseWriter, *http.Request)
	ReName(http.ResponseWriter, *http.Request)
	Copy(http.ResponseWriter, *http.Request)
	Info(http.ResponseWriter, *http.Request)
	NameSearch(http.ResponseWriter, *http.Request)
	Upload(http.ResponseWriter, *http.Request)
	Download(http.ResponseWriter, *http.Request)
}