package htmlpage
/**
 *@description 静态资源路由
 *@author	wupeng364@outlook.com
*/
import(
	"net/http"
)
/**
 * 文件基础操作网络接口
 */
type PageDispatch struct{
}
// 首页重定向
func (this PageDispatch)Index( w http.ResponseWriter, r *http.Request ){
	 http.Redirect(w, r, "/pages", http.StatusTemporaryRedirect)
}
// =========================================================
func sendError( w http.ResponseWriter, err error ){
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err.Error( )))
}
func setJson( w http.ResponseWriter ){
	w.Header( ).Set("Content-type", "application/json")
}
