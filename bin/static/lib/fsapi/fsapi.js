"use strict";
/**
 *@description FsApi
 *@author	wupeng364@outlook.com
*/
;(function (root, factory) {
	if (typeof exports === "object") {
		module.exports = exports = factory();
	}else {
		// Global (browser)
		root.$fsApi = factory();
	}
}(this, function ( ){
	var _fsApi = {
		// 获取一个传输的Token
		GetTransferToken: function( path ){
			return $fhttp.post("/fsapi/transfertoken", {
				"path": path?path:""
			});
		},
		// 获取token里的信息
		BatchOperationTokenInfo:function( token ){
			return $fhttp.get("/fsapi/batchoperationtokenstauts", {"token": token?token:""});
		},
		// 操作token中的值
		SetBatchOperationToken:function( token, operation ){
			return $fhttp.post("/fsapi/batchoperationtokenstauts", {
				"token": token?token:"",
				"operation": operation?operation:""
			});
		},
		// 列表路径
		List:function( path ){
			return $fhttp.get("/fsapi/list", {"path": path?path:""});
		},
		// 获取一个下载的Url
		GetDownloadUrl: function( path ){
			return _fsApi.GetTransferToken( path ).then(function( data ){
				return "/fsapi/download/"+data
			});
		},
		// 获取一个打开的Url - 流
		GetSteamUrl: function( path ){
			return _fsApi.GetTransferToken( path ).then(function( data ){
				return "/fsapi/openfile/"+data
			});
		},
		// 获取一个上载的Url
		GetUploadUrl: function( path ){
			return _fsApi.GetTransferToken( path ).then(function( data ){
				return "/fsapi/upload/"+data
			}).catch(function( ){
				return "";
			});
		},
		// 复制文件|文件夹
		CopyAsync:function( src, dest, replaceExist, ignoreError ){
			return $fhttp.post("/fsapi/copyasync", {
				srcPath: src?src:"",
				dstPath: dest?dest:"",
				replace: replaceExist?replaceExist:false,
				ignore: ignoreError?ignoreError:false
			});
		},
		// 移动文件|文件夹
		MoveAsync:function( src, dest, replaceExist, ignoreError ){
			return $fhttp.post("/fsapi/moveasync", {
				srcPath: src?src:"",
				dstPath: dest?dest:"",
				replace: replaceExist?replaceExist:false,
				ignore: ignoreError?ignoreError:false
			});
		},
		// 移动文件|文件夹
		Delete: function( path ){
			return $fhttp.post("/fsapi/del", {
				path: path?path:"",
			});
		},
		// 重命名文件|文件夹
		Rename: function( path, name ){
			return $fhttp.post("/fsapi/rename", {
				path: path?path:"",
				name: name?name:"",
			});
		},
		// 新建文件夹
		NewFolder: function( path ){
			return $fhttp.post("/fsapi/newfolder", {
				path: path?path:"",
			});
		},
	};
	return _fsApi;
}));
