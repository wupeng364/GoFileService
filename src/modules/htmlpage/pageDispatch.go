package htmlpage
/**
 *@description 静态资源路由
 *@author	wupeng364@outlook.com
*/
import(
	"net/http"
	hst "common/httpservertools"
)
/**
 * 文件基础操作网络接口
 */
type PageDispatch struct{
}


// 向 Server Router 中注册下列处理器 , 实现接口 httpservertools.Registrar
func (pd PageDispatch) RoutList( ) hst.StructRegistrar{
	return hst.StructRegistrar{
		true,
		"",
		[]hst.HandlersFunc{
			pd.Index,
		},
	}
}
// 首页重定向
func (pd PageDispatch)Index( w http.ResponseWriter, r *http.Request ){
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
