// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// API请求包装
;(function (root, factory) {
	if (typeof exports === "object") {
		module.exports = exports = factory();
	}else {
		// Global (browser)
		root.$apitools = factory();
	}
}(this, function ( ){
	"use strict";
	let apitools = {
		_Const: {
			ack: "access_object",
			host: "access_host",
			hosturi: localStorage.getItem("access_host"),
		},
		/**
		 * 设置API主机地址
		 */
		setAPIHost: function( uri ){
			if(uri && uri.length > 0){
				localStorage.setItem(this._Const.host, uri);
				this._Const.hosturi = uri;
			}else{
				localStorage.removeItem(this._Const.host);
			}
		},
		getAPIHost: function( ){
			return this._Const.hosturi?this._Const.hosturi:'';
		},
		/**
		 * 生成一个完整的url
		 * 如: /a/b/b --> https://127.0.0.1/a/b/c
		 */
		buildAPIURL: function( url ){
			return this._Const.hosturi?this._Const.hosturi+url:window.location.origin+url;
		},
		// 构建签名url
		getSignAPIURL: function (url, params) {
			if(!params){ params = {}; }
			let session = apitools.getSession( );
			if(!session||!session.SecretKey||!session.AccessKey){
				this.apiResponseStautsHandler({Code: 401});return;
			}
			// 去除无效字段
			let paramsmap = new Map( );
			for(let key in params){
				if (params[key] == undefined || params[key] == null || params[key].length == 0) {
					continue;
				}
				paramsmap.set(key, params[key]);
			}
			// 构建请求负载
			let payloads = []; let payloads_encode = [];
			if(paramsmap.size > 0 ){
				let keys = paramsmap.keys( ).sort( );
				for (let i = 0; i < keys.length; i++) {
					let val = paramsmap.get(keys[i]);
					payloads.push(keys[i]+"="+val);
					payloads_encode.push(keys[i]+"="+encodeURIComponent(val));
				}
			}
			if(payloads.length >0){
				return url +'?'+payloads_encode.join("&")+'&ack='+session.AccessKey+"&sign="+md5(payloads.join("&") + session.SecretKey)
			}else{
				return url +'?ack='+session.AccessKey+'&sign='+md5(session.SecretKey)
			}
		},
		/**
		 * Api自动签名请求
		 */
		apiRequest: function(opt){
			let session = apitools.getSession( );
			opt.session = session;
			if(!session||!session.SecretKey||!session.AccessKey){
				this.apiResponseStautsHandler({Code: 401});return;
			}
			if(this._Const.hosturi && this._Const.hosturi.length > 0){
				if(!opt.uri.startWith('http://') && !opt.uri.startWith('https://')){
					opt.uri = this._Const.hosturi + opt.uri;
				}
			}
			let request = $utils.AjaxRequest( opt );
			let signPayload = request.payload.payload;
			if(signPayload&&signPayload.length>0){
				signPayload += opt.session.SecretKey;
			}else{
				signPayload = opt.session.SecretKey;
			}
			request.setHeader("ack", session.AccessKey);
			request.setHeader("sign", md5(signPayload));
			return request;
		},
		/**
		 * 相应结构格式化
		 */
		apiResponseFormat: function(xhr){
			let result = {
				Code: 0,
				Data: '',
			};
			// 请求错误&没有响应数据
			// if( xhr.status !== 200 ){
			// 	result.Code = xhr.status;
			// 	result.Data = xhr.statusText;
			// }else{
				let _obj = {};
				if( undefined != xhr.responseText ){
					if( $utils.isString(xhr.responseText) ){
						try{
							_obj = JSON.parse( xhr.responseText );
						}catch(err){
							_obj = xhr.responseText;
						}
					}else{
						_obj = xhr.responseText;
					}
				}
				// 没有结构化的数据返回
				if( undefined == _obj.Code && undefined == _obj.Data){
					result.Code = xhr.status;
					result.Data = _obj;
				}else{
					result.Code = undefined==_obj.Code?xhr.status:_obj.Code;
					result.Data = undefined==_obj.Data?'':_obj.Data;
				}
			// }
			apitools.apiResponseStautsHandler(result);
			// console.log("apiResponseFormat: ", result, xhr);
			return result;
		},
		// Api 状态返回翻译和处理
		apiResponseStautsHandler: function(res){
			if(res){
				if(res.Code == 401){
					res.Data = "登录过期,请刷新页面";
					(top.location||window.location).href = "/";
				}
				else if(res.Code == 403){
					res.Data = "请求被禁止,可能是权限不足";
				}
			}
		},
		/**
		 * APi-Get请求
		 */
		apiGet: function(uri, datas){
			return new Promise(function(resolve, reject){
				apitools.apiRequest({
					uri: uri,
					datas: datas,				
				}).do(function(xhr, opt){
					if(xhr.readyState === 4){
						let res = apitools.apiResponseFormat(xhr);
						if( res.Code === 200 ){
							resolve( res.Data );
						}else{
							reject( res.Data );
						}
					}
				});
			})
		},
		/**
		 * APi-Post请求
		 */
		apiPost: function(uri, datas){
			return new Promise(function(resolve, reject){
				apitools.apiRequest({
					uri: uri,
					method:"POST",
					datas: datas,				
				}).do(function(xhr, opt){
					if(xhr.readyState === 4){
						let res = apitools.apiResponseFormat(xhr);
						if( res.Code === 200 ){
							resolve( res.Data );
						}else{
							reject( res.Data );
						}
					}
				});
			})
		},
		apiPostSync: function(uri, datas){
			return apitools.apiRequest({
				uri: uri,
				method:"POST",
				datas: datas,	
				async: false,			
			}).do(function(xhr, opt){
				if(xhr.readyState === 4){
					let res = apitools.apiResponseFormat(xhr);
					if( res.Code === 200 ){
						return res.Data;
					}else{
						throw res.Data;
					}
				}
			});
		},
		/**
		 * 保存会话到cookie中
		 */
		saveSession: function( accessObj ){
			localStorage.setItem(apitools._Const.ack, "");
			if(accessObj && accessObj.AccessKey){
				localStorage.setItem(apitools._Const.ack, JSON.stringify(accessObj));
				$utils.setCookie("ack", accessObj.AccessKey);
			}
		},
		/**
		 * 获取会话信息
		 * {UserId, AccessKey, SecretKey}
		 */
		getSession: function(){
			try{
				return JSON.parse(localStorage.getItem(apitools._Const.ack));
			}catch(err){
				return {};
			}
		},
		/**
		 * 注销会话
		 */
		destroySession: function(){
			localStorage.setItem(apitools._Const.ack, '');
			$utils.setCookie("ack", '');
		},
	};

	/*apitools.saveSession({
		UserId: "system",
		AccessKey: "0000000000000000000000000000",
		SecretKey: "0000000000000000000000000000",
	});*/
	return apitools;
}));