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
		root.$fpmsApi = factory();
	}
}(this, function ( ){
	let api = { sync:{} };

	// listFPermissions
	api.listFPermissions = function( ){
		return $apitools.apiPost("/fpms/listfpermissions", { })
	};
	// listUserFPermissions
	api.listUserFPermissions = function( userid ){
		return $apitools.apiPost("/fpms/listuserfpermissions", {userid: userid?userid:''})
	};
	// addFPermission
	api.addFPermission = function( userid, path, permission){
		return $apitools.apiPost("/fpms/addfpermission", {
			'userid': userid?userid:'',
			'path': path?path:'',
			'permission': permission?permission:'',
		 });
	};
	// delFPermission
	api.delFPermission = function( permissionid){
		return $apitools.apiPost("/fpms/delfpermission", {
			'permissionid': permissionid?permissionid:'',
		 });
	};
	// updateFPermission
	api.updateFPermission = function( permissionid, permission){
		return $apitools.apiPost("/fpms/updatefpermission", {
			"permissionid": permissionid?permissionid:'',
			"permission": permission?permission:'',
		});
	};
	// delFPermission
	api.sync.delFPermission = function( permissionid ){
		return $apitools.apiPostSync("/fpms/delfpermission", {
			'permissionid': permissionid?permissionid:'',
		 });
	};
	return api;
}));
