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
		root.$preview = factory();
	}
}(this, function ( ){
	var api = {
		supported: {
			audio: ['mp3','flac'],
			video: ['mp4','mov'],
		},
		// 获取一个预览的Token
		askToken: function( path ){
			return $apitools.apiGet("/preview/asktoken", {
				"path": path?path:""
			});
		},
		// Status
		status: function( path ){
			return $apitools.apiGet("/preview/status", { });
		},
		// 获取token内的信息 token, type=[list||steam], path
		// 读取信息: token & type=list 
		// 获取流: token & type=stream || path 
		tokenDatas: function( token, type, path ){
			return $apitools.apiGet("/preview/tokendatas", {
				"token": token?token:"",
				"type": type?type:"",
				"path": path?path:""
			});
		},
		// 预览
		doPreview: function(path, suffix){
			return new Promise(function(resolve, reject){
				api.askToken( path ).then(function( data ){
					api.openPreview(data, suffix?suffix.toLowerCase():path.getSuffixed(false).toLowerCase());
					resolve();
				}).catch(reject);
			});
		},
		// 打开预览地址
		openPreview: function(token, suffix){
			let type = '';
			if(api.isSupport('audio', suffix)){
				type = 'audio';
			}else if(api.isSupport('video', suffix)){
				type = 'video';
			}else{
				throw "不支持预览该文件";
			}
			window.open("/pages/preview?token="+token+"&type="+type);
		},
		// 是否是支持的播放类型
		isSupport: function(type, suffix){
			return type && suffix && api.supported[type] && api.supported[type].indexOf(suffix) > -1;
		}
	};
	return api;
}));
