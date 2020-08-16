// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// FPmsAPI 文件权限管理api

package fpmsapi

import (
	"encoding/json"
	"fmt"
	"gofs/base/httpserver"
	"gofs/base/signature"
	"gofs/data/filepermission"
	"gofs/data/usermanage"
	"gutils/hstool"
	"gutils/mloader"
	"net/http"
	"strconv"
)

// FPmsAPI 用户管理api
type FPmsAPI struct {
	um    *usermanage.UserManager
	fpmsr *filepermission.FPmsManager
	hs    *httpserver.HTTPServer
	sg    *signature.Signature
}

// ModuleOpts 模块加载器接口实现, 返回模块信息&配置
func (api *FPmsAPI) ModuleOpts() mloader.Opts {
	return mloader.Opts{
		Name:        "FPmsAPI",
		Version:     1.0,
		Description: "文件权限管理Api接口模块",
		OnReady: func(mctx *mloader.Loader) {
			api.um = mctx.GetModuleByTemplate(api.um).(*usermanage.UserManager)
			api.fpmsr = mctx.GetModuleByTemplate(api.fpmsr).(*filepermission.FPmsManager)
			api.hs = mctx.GetModuleByTemplate(api.hs).(*httpserver.HTTPServer)
			api.sg = mctx.GetModuleByTemplate(api.sg).(*signature.Signature)
		},
		OnInit: api.init,
	}
}

// 初始化
func (api *FPmsAPI) init() {
	api.hs.AddRegistrar(api)
	// 注册Api签名拦截器
	api.hs.AddURLFilter(baseurl+"/:"+`[\S]+`, api.sg.RestfulAPIFilter)
	fmt.Println("   > FPmsAPI http registered end")
}

// RoutList 向 Server Router 中注册下列处理器 , 实现接口 httpserver.Registrar
func (api *FPmsAPI) RoutList() httpserver.StructRegistrar {
	return httpserver.StructRegistrar{
		IsToLower: true,
		BasePath:  baseurl,
		HandlerFunc: []hstool.HandlersFunc{
			api.ListFPermissions,
			api.GetUserPermissionSum,
			api.ListUserFPermissions,
			api.AddFPermission,
			api.DelFPermission,
			api.UpdateFPermission,
		},
	}
}

// checkPermission 检查是否是管理员
func (api *FPmsAPI) checkPermission(w http.ResponseWriter, r *http.Request) bool {
	userID := api.sg.GetUserID4Request(r)
	if len(userID) > 0 {
		qUserID := r.FormValue("userid")
		if len(qUserID) > 0 && qUserID == userID {
			return true
		}
		if pms, err := api.um.QueryUser(userID); nil == err {
			if pms.UserType == usermanage.AdminRole {
				return true
			}
		}
	}
	w.WriteHeader(http.StatusForbidden)
	return false
}

// ListFPermissions 列出所有数据, 无分页
func (api *FPmsAPI) ListFPermissions(w http.ResponseWriter, r *http.Request) {
	if !api.checkPermission(w, r) {
		return
	}
	if pms, err := api.fpmsr.ListFPermissions(); nil == err {
		if tb, err := json.Marshal(pms); nil == err {
			httpserver.SendSuccess(w, string(tb))
		} else {
			httpserver.SendError(w, err)
		}
	} else {
		httpserver.SendError(w, err)
	}
}

// GetUserPermissionSum 根据用户ID查询权限值
func (api *FPmsAPI) GetUserPermissionSum(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("userid")
	var paths []string
	if len(userID) == 0 {
		httpserver.SendError(w, filepermission.ErrorUserIDIsNil)
		return
	}
	err := json.Unmarshal([]byte(r.FormValue("paths")), &paths)
	if nil != err || len(paths) == 0 {
		httpserver.SendError(w, filepermission.ErrorPermissionPathIsNil)
		return
	}
	if !api.checkPermission(w, r) {
		return
	}
	permissions := make(map[string]int64, len(paths))
	for i := 0; i < len(paths); i++ {
		permissions[paths[i]] = api.fpmsr.GetUserPermissionSum(userID, paths[i])
	}

	if tb, err := json.Marshal(permissions); nil == err {
		httpserver.SendSuccess(w, string(tb))
	} else {
		httpserver.SendError(w, err)
	}
}

// ListUserFPermissions 根据用户ID查询详细信息
func (api *FPmsAPI) ListUserFPermissions(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("userid")
	if len(userID) == 0 {
		httpserver.SendError(w, filepermission.ErrorUserIDIsNil)
		return
	}
	if !api.checkPermission(w, r) {
		return
	}
	if pms, err := api.fpmsr.ListUserFPermissions(userID); nil == err {
		if tb, err := json.Marshal(pms); nil == err {
			httpserver.SendSuccess(w, string(tb))
		} else {
			httpserver.SendError(w, err)
		}
	} else {
		httpserver.SendError(w, err)
	}
}

// AddFPermission 添加
func (api *FPmsAPI) AddFPermission(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("userid")
	path := r.FormValue("path")
	permissionStr := r.FormValue("permission")
	if len(userID) == 0 {
		httpserver.SendError(w, filepermission.ErrorUserIDIsNil)
		return
	}
	if len(path) == 0 {
		httpserver.SendError(w, filepermission.ErrorPermissionPathIsNil)
		return
	}
	if len(permissionStr) == 0 {
		httpserver.SendError(w, filepermission.ErrorPermissionIsNil)
		return
	}
	permission, err := strconv.Atoi(permissionStr)
	if nil != err {
		httpserver.SendError(w, err)
		return
	}
	if !api.checkPermission(w, r) {
		return
	}
	pmsInfo := filepermission.FPermissionInfo{
		UserID:     userID,
		Path:       path,
		Permission: int64(permission),
	}

	if err := api.fpmsr.AddFPermission(pmsInfo); nil == err {
		httpserver.SendSuccess(w, "")
	} else {
		httpserver.SendError(w, err)
	}
}

// UpdateFPermission 修改
func (api *FPmsAPI) UpdateFPermission(w http.ResponseWriter, r *http.Request) {
	permissionID := r.FormValue("permissionid")
	permissionStr := r.FormValue("permission")
	if len(permissionStr) == 0 {
		httpserver.SendError(w, filepermission.ErrorPermissionIDIsNil)
		return
	}
	permission, err := strconv.Atoi(permissionStr)
	if nil != err {
		httpserver.SendError(w, err)
		return
	}

	if !api.checkPermission(w, r) {
		return
	}
	pmsInfo := filepermission.FPermissionInfo{
		PermissionID: permissionID,
		Permission:   int64(permission),
	}
	if err := api.fpmsr.UpdateFPermission(pmsInfo); nil == err {
		httpserver.SendSuccess(w, "")
	} else {
		httpserver.SendError(w, err)
	}
}

// DelFPermission 根据ID删除
func (api *FPmsAPI) DelFPermission(w http.ResponseWriter, r *http.Request) {
	permissionID := r.FormValue("permissionid")
	if len(permissionID) == 0 {
		httpserver.SendError(w, filepermission.ErrorPermissionIDIsNil)
		return
	}
	if !api.checkPermission(w, r) {
		return
	}

	if err := api.fpmsr.DelFPermission(permissionID); nil == err {
		httpserver.SendSuccess(w, "")
	} else {
		httpserver.SendError(w, err)
	}
}
