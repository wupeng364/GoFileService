package userapi
/**
 *@description 用户管理接口-http
 *@author	wupeng364@outlook.com
*/
import(
	"fmt"
	"net/http"
	"encoding/json"
	"gofs/modules/core/usermanage"
	"gofs/modules/core/signature"
	"gofs/modules/common/httpserver"
	hst "gofs/common/httpservertools"
)
const(
	base_url = "/userapi"
)
/**
 * 用户管理操作网络接口
 */
type imp_http struct{
	um *usermanage.UserManageModule
	hs *httpserver.HttpServerModule
	sg *signature.SignatureModule
}
// 初始化
func (ih *imp_http) init( ){
	ih.hs.AddRegistrar(ih)
	
	// 用户密码校验不需要会话
	ih.hs.AddIgnoreFilter(base_url+"/checkpwd");
	// 注册Api签名拦截器
	ih.hs.AddUrlFilter(base_url+"/:"+`[\S]+`, ih.sg.ApiFilter_Http)
	fmt.Println("   > UserApiModule http registered end")
}
// 向 Server Router 中注册下列处理器 , 实现接口 httpserver.Registrar
func (fa *imp_http) RoutList( ) httpserver.StructRegistrar{
	return httpserver.StructRegistrar{
		true,
		base_url,
		[]hst.HandlersFunc{
			fa.ListAllUsers,
			fa.QueryUser,
			fa.AddUser,
			fa.DelUser,
			fa.UpdateUserName,
			fa.UpdateUserPwd,
			fa.CheckPwd,
		},
	}
}

// 列出所有用户数据, 无分页
func (ih *imp_http)ListAllUsers(w http.ResponseWriter, r *http.Request){
	if users, err := ih.um.ListAllUsers(); nil == err {
		if tb, err := json.Marshal( users ); nil == err {
			httpserver.SendSuccess(w, string(tb))	
		}else{
			httpserver.SendError(w, err)
		}
	}else{
		httpserver.SendError(w, err)
	}
}
// 根据用户ID查询详细信息
func (ih *imp_http)QueryUser(w http.ResponseWriter, r *http.Request){
	userId   := r.FormValue("userId")
	if len(userId) == 0 {
		httpserver.SendError(w, Error_UserIdIsNil); return
	}
	
	if user, err := ih.um.QueryUser(userId); nil == err {
		if tb, err := json.Marshal( user ); nil == err {
			httpserver.SendSuccess(w, string(tb))	
		}else{
			httpserver.SendError(w, err)
		}
	}else{
		httpserver.SendError(w, err)
	}
}
// 添加用户
func (ih *imp_http)AddUser(w http.ResponseWriter, r *http.Request){
	userId   := r.FormValue("userId")
	userName := r.FormValue("userName")
	userPwd  := r.FormValue("userPwd")
	if len(userId) == 0 {
		httpserver.SendError(w, Error_UserIdIsNil); return
	}
	if len(userName) == 0 {
		httpserver.SendError(w, Error_UserNameIsNil); return
	}
	uinfo := usermanage.UserInfo{
		UserId: userId,
		Username: userName,
		Userpwd: userPwd,
	}
	
	if err := ih.um.AddUser(&uinfo); nil == err {
		httpserver.SendSuccess(w, "")
	}else{
		httpserver.SendError(w, err);
	}
}
// 修改用户密码
func (ih *imp_http)UpdateUserPwd(w http.ResponseWriter, r *http.Request){
	userId := r.FormValue("userId")
	userPwd := r.FormValue("userPwd")
	if len(userId) == 0 {
		httpserver.SendError(w, Error_UserIdIsNil); return
	}
	if err := ih.um.UpdateUserPwd(userId, userPwd); nil == err {
		httpserver.SendSuccess(w, "")
	}else{
		httpserver.SendError(w, err);
	}
}
// 修改用户昵称
func (ih *imp_http)UpdateUserName(w http.ResponseWriter, r *http.Request){
	userId := r.FormValue("userId")
	userName := r.FormValue("userName")
	if len(userId) == 0 {
		httpserver.SendError(w, Error_UserIdIsNil); return
	}
	if len(userName) == 0 {
		httpserver.SendError(w, Error_UserNameIsNil); return
	}
	if err := ih.um.UpdateUserName(userId, userName); nil == err {
		httpserver.SendSuccess(w, "")
	}else{
		httpserver.SendError(w, err);
	}
}
// 根据userId删除用户
func (ih *imp_http)DelUser(w http.ResponseWriter, r *http.Request){
	userId := r.FormValue("userId")
	if len(userId) == 0 {
		httpserver.SendError(w, Error_UserIdIsNil); return
	}
	if err := ih.um.DelUser(userId); nil == err {
		httpserver.SendSuccess(w, "")
	}else{
		httpserver.SendError(w, err);
	}
}
// 校验密码是否一致
// 校验成功返回session
func (ih *imp_http)CheckPwd(w http.ResponseWriter, r *http.Request){
	userId := r.FormValue("userId")
	pwd  := r.FormValue("pwd")
	if len(userId) == 0 {
		httpserver.SendError(w, Error_UserIdIsNil); return
	}
	// 检查密码是否正确, 如果正确需要返回签名信息
	if ih.um.CheckPwd(userId, pwd) {
		ack, err := ih.sg.CreateWebSession(userId, r)
		if nil != err {
			httpserver.SendError(w, err); return
		}
		httpserver.SendSuccess(w, ack.ToJson( ))
	}else{
		httpserver.SendError(w, Error_PwdIsError)
	}
}