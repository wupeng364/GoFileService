// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// UserAPI 用户管理api

package userapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"gofs/comm/httpserver"
	"gofs/data/usermanage"
	"gofs/service/restful/signature"
	"gutils/hstool"
	"gutils/mloader"
	"net/http"
)

// UserAPI 用户管理api
type UserAPI struct {
	um *usermanage.UserManager
	hs *httpserver.HTTPServer
	sg *signature.Signature
}

// ModuleOpts 模块加载器接口实现, 返回模块信息&配置
func (userapi *UserAPI) ModuleOpts() mloader.Opts {
	return mloader.Opts{
		Name:        "UserAPI",
		Version:     1.0,
		Description: "用户管理Api接口模块",
		OnReady: func(mctx *mloader.Loader) {
			userapi.um = mctx.GetModuleByTemplate(userapi.um).(*usermanage.UserManager)
			userapi.hs = mctx.GetModuleByTemplate(userapi.hs).(*httpserver.HTTPServer)
			userapi.sg = mctx.GetModuleByTemplate(userapi.sg).(*signature.Signature)
		},
		OnInit: userapi.init,
	}
}

// 初始化
func (userapi *UserAPI) init() {
	userapi.hs.AddRegistrar(userapi)

	// 用户密码校验不需要会话
	userapi.hs.AddIgnoreFilter(baseurl + "/checkpwd")
	// 注册Api签名拦截器
	userapi.hs.AddURLFilter(baseurl+"/:"+`[\S]+`, userapi.sg.RestfulAPIFilter)
	fmt.Println("   > UserApiModule http registered end")
}

// RoutList 向 Server Router 中注册下列处理器 , 实现接口 httpserver.Registrar
func (userapi *UserAPI) RoutList() httpserver.StructRegistrar {
	return httpserver.StructRegistrar{
		IsToLower: true,
		BasePath:  baseurl,
		HandlerFunc: []hstool.HandlersFunc{
			userapi.ListAllUsers,
			userapi.QueryUser,
			userapi.AddUser,
			userapi.DelUser,
			userapi.UpdateUserName,
			userapi.UpdateUserPwd,
			userapi.CheckPwd,
			userapi.Logout,
		},
	}
}

// ListAllUsers 列出所有用户数据, 无分页
func (userapi *UserAPI) ListAllUsers(w http.ResponseWriter, r *http.Request) {
	if users, err := userapi.um.ListAllUsers(); nil == err {
		if tb, err := json.Marshal(*users); nil == err {
			httpserver.SendSuccess(w, string(tb))
		} else {
			httpserver.SendError(w, err)
		}
	} else {
		httpserver.SendError(w, err)
	}
}

// QueryUser 根据用户ID查询详细信息
func (userapi *UserAPI) QueryUser(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("userid")
	if len(userID) == 0 {
		httpserver.SendError(w, ErrorUserIDIsNil)
		return
	}

	if user, err := userapi.um.QueryUser(userID); nil == err {
		if tb, err := json.Marshal(*user); nil == err {
			httpserver.SendSuccess(w, string(tb))
		} else {
			httpserver.SendError(w, err)
		}
	} else {
		httpserver.SendError(w, err)
	}
}

// AddUser 添加用户
func (userapi *UserAPI) AddUser(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("userid")
	userName := r.FormValue("username")
	userPwd := r.FormValue("userpwd")
	if len(userID) == 0 {
		httpserver.SendError(w, ErrorUserIDIsNil)
		return
	}
	if len(userName) == 0 {
		httpserver.SendError(w, ErrorUserNameIsNil)
		return
	}
	uinfo := usermanage.UserInfo{
		UserID:   userID,
		UserName: userName,
		UserPWD:  userPwd,
		UserType: usermanage.NormalRole,
	}

	if err := userapi.um.AddUser(&uinfo); nil == err {
		httpserver.SendSuccess(w, "")
	} else {
		httpserver.SendError(w, err)
	}
}

// UpdateUserPwd 修改用户密码
func (userapi *UserAPI) UpdateUserPwd(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("userid")
	userPwd := r.FormValue("userpwd")
	if len(userID) == 0 {
		httpserver.SendError(w, ErrorUserIDIsNil)
		return
	}
	if err := userapi.um.UpdateUserPwd(userID, userPwd); nil == err {
		httpserver.SendSuccess(w, "")
	} else {
		httpserver.SendError(w, err)
	}
}

// UpdateUserName 修改用户昵称
func (userapi *UserAPI) UpdateUserName(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("userid")
	userName := r.FormValue("username")
	if len(userID) == 0 {
		httpserver.SendError(w, ErrorUserIDIsNil)
		return
	}
	if len(userName) == 0 {
		httpserver.SendError(w, ErrorUserNameIsNil)
		return
	}
	if err := userapi.um.UpdateUserName(userID, userName); nil == err {
		httpserver.SendSuccess(w, "")
	} else {
		httpserver.SendError(w, err)
	}
}

// DelUser 根据userID删除用户
func (userapi *UserAPI) DelUser(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("userid")
	if len(userID) == 0 {
		httpserver.SendError(w, ErrorUserIDIsNil)
		return
	}
	users, _ := userapi.um.ListAllUsers()
	if nil == users || len(*users) == 1 {
		httpserver.SendError(w, errors.New("Cannot delete the last user"))
		return
	}
	{
		count := 0
		lauid := ""
		for _, val := range *users {
			if val.UserType == usermanage.AdminRole {
				count++
				lauid = val.UserID
			}
		}
		if count <= 1 && (userID == lauid || len(lauid) == 0) {
			httpserver.SendError(w, errors.New("Cannot delete the last admin user"))
			return
		}
	}
	if err := userapi.um.DelUser(userID); nil == err {
		httpserver.SendSuccess(w, "")
	} else {
		httpserver.SendError(w, err)
	}
}

//CheckPwd 校验密码是否一致,  校验成功返回session
func (userapi *UserAPI) CheckPwd(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("userid")
	pwd := r.FormValue("pwd")
	if len(userID) == 0 {
		httpserver.SendError(w, ErrorUserIDIsNil)
		return
	}
	// 检查密码是否正确, 如果正确需要返回签名信息
	if userapi.um.CheckPwd(userID, pwd) {
		ack, err := userapi.sg.CreateWebSession(userID, r)
		if nil != err {
			httpserver.SendError(w, err)
			return
		}
		httpserver.SendSuccess(w, ack.ToJSON())
	} else {
		httpserver.SendError(w, ErrorPwdIsError)
	}
}

//Logout 注销会话
func (userapi *UserAPI) Logout(w http.ResponseWriter, r *http.Request) {

	err := userapi.sg.DestroySignature4HTTP(r)
	if nil != err {
		httpserver.SendError(w, err)
		return
	}
	httpserver.SendSuccess(w, "")
}
