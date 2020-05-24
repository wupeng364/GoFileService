// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

"use strict";
;(function (root, factory) {
	if (typeof exports === "object") {
		module.exports = exports = factory();
	}else {
		// Global (browser)
		root.$userApi = factory();
	}
}(this, function ( ){
	var api = { sync:{} };

	/**
	 * 登陆
	 * return {UserID, AccessKey, SecretKey}
	 */
	api.login = function(user, pwd){
		return new Promise(function(resolve, reject){
			$utils.AjaxRequest({
				uri: "/userapi/checkpwd",
				datas: {
					"userid": user,
					"pwd": pwd,
				},				
			}).do(function(xhr, opt){
				if(xhr.readyState === 4){
					var res = $apitools.apiResponseFormat(xhr);
					if( res.Code === 200 ){
						resolve( res.Data );
					}else{
						reject( res.Data );
					}
				}
			});
		});
	};
	// logout
	api.logout = function( ){
		return $apitools.apiPost("/userapi/logout")
	};
	// QueryUser
	api.queryuser = function( userid ){
		return $apitools.apiPost("/userapi/queryuser", {
			"userid": userid
		})
	};
	// UpdateUserName
	api.updateUserName = function( userid, username){
		return $apitools.apiPost("/userapi/updateusername", {
			"userid": userid,
			"username": username,
		})
	};
	// UpdateUserPwd
	api.updateUserPwd = function( userid, userpwd){
		return $apitools.apiPost("/userapi/updateuserpwd", {
			"userid": userid,
			"userpwd": userpwd,
		})
	};
	// ListAllUsers
	api.listAllUsers = function( userid, userpwd){
		return $apitools.apiPost("/userapi/listallusers", { })
	};
	// AddUser
	api.addUser = function( userid, username, userpwd){
		return $apitools.apiPost("/userapi/adduser", {
			'userid': userid,
			'username': username,
			'userpwd': userpwd?userpwd:'',
		 })
	};
	// DelUser
	api.delUser = function( userid){
		return $apitools.apiPost("/userapi/deluser", {
			'userid': userid,
		 })
	};
	// UpdateUserName
	api.sync.updateUserName = function( userid, username){
		return $apitools.apiPostSync("/userapi/updateusername", {
			"userid": userid,
			"username": username,
		})
	};
	// UpdateUserPwd
	api.sync.updateUserPwd = function( userid, userpwd){
		return $apitools.apiPostSync("/userapi/updateuserpwd", {
			"userid": userid,
			"userpwd": userpwd,
		})
	};
	// DelUser
	api.sync.delUser = function( userid){
		return $apitools.apiPostSync("/userapi/deluser", {
			'userid': userid,
		 })
	};
	return api;
}));
