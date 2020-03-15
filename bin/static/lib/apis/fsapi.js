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
		root.$fsApi = factory();
	}
}(this, function ( ){
	var api = {
		// 获取一个传输的Token
		GetTransferToken: function( path ){
			return $apitools.apiPost("/fsapi/transfertoken", {
				"path": path?path:""
			});
		},
		// 获取token里的信息
		BatchOperationTokenInfo:function( token ){
			return $apitools.apiGet("/fsapi/batchoperationtokenstauts", {"token": token?token:""});
		},
		// 操作token中的值
		SetBatchOperationToken:function( token, operation ){
			return $apitools.apiPost("/fsapi/batchoperationtokenstauts", {
				"token": token?token:"",
				"operation": operation?operation:""
			});
		},
		// 列表路径
		List:function( path ){
			return $apitools.apiGet("/fsapi/list", {"path": path?path:""});
		},
		// 获取一个下载的Url
		GetDownloadUrl: function( path ){
			return api.GetTransferToken( path ).then(function( data ){
				return "/fsapi/download/"+data
			});
		},
		// 获取一个打开的Url - 流
		GetSteamUrl: function( path ){
			return api.GetTransferToken( path ).then(function( data ){
				return "/fsapi/openfile/"+data+path.getSuffixed( );
			});
		},
		// 获取一个上载的Url
		GetUploadUrl: function( path ){
			return api.GetTransferToken( path ).then(function( data ){
				return "/fsapi/upload/"+data;
			}).catch(function( ){
				return "";
			});
		},
		// 复制文件|文件夹
		CopyAsync:function( src, dest, replaceExist, ignoreError ){
			return $apitools.apiPost("/fsapi/copyasync", {
				srcPath: src?src:"",
				dstPath: dest?dest:"",
				replace: replaceExist?replaceExist:false,
				ignore: ignoreError?ignoreError:false
			});
		},
		// 移动文件|文件夹
		MoveAsync:function( src, dest, replaceExist, ignoreError ){
			return $apitools.apiPost("/fsapi/moveasync", {
				srcPath: src?src:"",
				dstPath: dest?dest:"",
				replace: replaceExist?replaceExist:false,
				ignore: ignoreError?ignoreError:false
			});
		},
		// 移动文件|文件夹
		Delete: function( path ){
			return $apitools.apiPost("/fsapi/del", {
				path: path?path:"",
			});
		},
		// 重命名文件|文件夹
		Rename: function( path, name ){
			return $apitools.apiPost("/fsapi/rename", {
				path: path?path:"",
				name: name?name:"",
			});
		},
		// 新建文件夹
		NewFolder: function( path ){
			return $apitools.apiPost("/fsapi/newfolder", {
				path: path?path:"",
			});
		},
	};
	return api;
}));
