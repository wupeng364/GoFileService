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

	var apitools = {
		
		/**
		 * Api自动签名请求
		 */
		apiRequest: function(opt){
			var session = apitools.getSession( );
			opt.session = session;
			if(!session||!session.SecretKey||!session.AccessKey){
				throw new Error("Signature is empty");
			}
			var request = $utils.AjaxRequest( opt );
			var signPayload = request.payload.payload;
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
			var result = {
				Code: 0,
				Data: '',
			};
			// 请求错误&没有响应数据
			if( xhr.status !== 200 && ""==xhr.responseText ){
				result.Code = xhr.status;
				result.Data = xhr.statusText;
			}else{
				var _obj;
				if( undefined != xhr.responseText ){
					if( $utils.isString(xhr.responseText) ){
						_obj = JSON.parse( xhr.responseText );
					}else{
						_obj = xht.responseText;
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
			}
			apitools.apiResponseStautsHandler(result);
			// console.log("apiResponseFormat: ", result, xhr);
			return result;
		},
		// Api 状态返回翻译和处理
		apiResponseStautsHandler: function(res){
			if(res){
				if(res.Code == 401){
					res.Data = "登陆过期,请刷新页面";
					window.location.href = "/";
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
						var res = apitools.apiResponseFormat(xhr);
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
						var res = apitools.apiResponseFormat(xhr);
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
		 * 保存会话到cookie中
		 */
		saveSession: function( accessObj ){
			var _key = "access_object";
			// 在cookie中保存一天
			localStorage.setItem(_key, "");
			if(accessObj && accessObj.AccessKey){
				localStorage.setItem(_key, JSON.stringify(accessObj));
				$utils.setCookie("ack", accessObj.AccessKey);
			}
		},
		/**
		 * 获取会话信息
		 * {UserId, AccessKey, SecretKey}
		 */
		getSession: function(){
			var _key = "access_object";
			try{
				return JSON.parse(localStorage.getItem(_key));
			}catch(err){
				return null;
			}
		},
	};

	/*apitools.saveSession({
		UserId: "system",
		AccessKey: "0000000000000000000000000000",
		SecretKey: "0000000000000000000000000000",
	});*/
	return apitools;
}));