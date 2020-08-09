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
		// 异步执行一个动作
		AsyncExec: function(func, params ){
			if(!params){
				params = {};
			}
			params.func = func;
			return $apitools.apiPost("/file/asyncexec", params);
		},
		// 异步执行Token查询 CopyFile, MoveFile
		AsyncExecToken: function(func, token, params ){
			if(!params){
				params = {};
			}
			params.func = func;
			params.token = token;
			return $apitools.apiGet("/file/asyncexectoken", params);
		},
		// 获取一个Token, data可以存放值
		GetStreamToken: function( type, datas ){
			return $apitools.apiPost("/file/streamtoken", {
				"type": type?type:'',
				"data": datas?datas:''
			});
		},
		// 列表路径
		List:function( path ){
			return $apitools.apiGet("/file/list", {"path": path?path:''});
		},
		// 获取一个下载的Url
		GetDownloadUrl: function( path ){
			return api.GetStreamToken('download', path).then(function( data ){
				return $apitools.buildAPIURL($apitools.getSignAPIURL("/file/stream", {
					token: data,
				}));
			});
		},
		// 获取一个打开的Url - 流
		GetSteamUrl: function( path ){
			return api.GetStreamToken('stream', path).then(function( data ){
				return encodeURI( $apitools.buildAPIURL(
							$apitools.getSignAPIURL("/file/stream/"+(path.getName()), { token: data, }
						)
				));
			});
		},
		// 获取一个上载的Url
		GetUploadUrl: function( path ){
			return api.GetStreamToken('upload', path).then(function( data ){
				return $apitools.buildAPIURL($apitools.getSignAPIURL("/file/stream", {
					token: data,
				}));
			});
		},
		// 复制文件|文件夹
		CopyAsync:function( src, dest, replaceExist, ignoreError ){
			return this.AsyncExec("CopyFile", {
				srcPath: src?src:'',
				dstPath: dest?dest:'',
				replace: replaceExist?replaceExist:false,
				ignore: ignoreError?ignoreError:false
			});
		},
		// 移动文件|文件夹
		MoveAsync:function( src, dest, replaceExist, ignoreError ){
			return this.AsyncExec("MoveFile", {
				srcPath: src?src:'',
				dstPath: dest?dest:'',
				replace: replaceExist?replaceExist:false,
				ignore: ignoreError?ignoreError:false
			});
		},
		// 移动文件|文件夹
		Delete: function( path ){
			return $apitools.apiPost("/file/del", {
				path: path?path:'',
			});
		},
		// 重命名文件|文件夹
		Rename: function( path, name ){
			return $apitools.apiPost("/file/rename", {
				path: path?path:'',
				name: name?name:'',
			});
		},
		// 新建文件夹
		NewFolder: function( path ){
			return $apitools.apiPost("/file/newfolder", {
				path: path?path:'',
			});
		},
	};
	return api;
}));
