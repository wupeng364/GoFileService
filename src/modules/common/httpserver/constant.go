package httpserver

import(
	"net/http"
	"encoding/json"
	"common/httpservertools"
)
const(
	cfg_http 	  = "http"
	cfg_http_addr = "addr"
	cfg_http_port = "port"
)

// struct 批量注册器
type StructRegistrar struct{
	IsToLower 	bool	        // 是否需要转小写访问
	BasePath 	string  		// 基础路径 /base/child....
	HandlerFunc []httpservertools.HandlersFunc	// 需要注册的 fuc
}
// struct 批量注册器接口
type Registrar interface{
	RoutList( ) StructRegistrar
}
// APi 接口返回规范
type ApiResponse struct{
	Code int
	Data string
}
// 返回成功结果
func SendSuccess(w http.ResponseWriter, msg string){
	w.WriteHeader(http.StatusOK)
	w.Header( ).Set("Content-type", "application/json;charset=utf-8")
	w.Write(parse2ApiJson(http.StatusOK, msg))
}
// 返回失败结果
func SendError( w http.ResponseWriter, err error ){
	w.WriteHeader(http.StatusBadRequest)
	w.Header( ).Set("Content-type", "application/json;charset=utf-8")
	w.Write( parse2ApiJson(http.StatusBadRequest, err.Error( )) )
}
func parse2ApiJson( code int, str string) []byte{
	bt, err := json.Marshal(ApiResponse{Code: code, Data: str})
	if nil != err {
		return []byte(err.Error())
	}
	return bt
}
