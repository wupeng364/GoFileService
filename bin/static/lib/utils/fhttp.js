/**
 *@description httpUtils
 *@author	wupeng364@outlook.com
*/
;(function (root, factory) {
	if (typeof exports === "object") {
		module.exports = exports = factory();
	}else {
		// Global (browser)
		root.$fhttp = factory();
	}
}(this, function ( ){

	var _http = {
		/**
		 * Ajax请求
		 * opt={uri, method, datas, headers, async}
		 */
		AjaxRequest: function( opt ){
			opt.xhrPayload = "";
			opt.async  = undefined==opt.async?true:opt.async;
			opt.datas  = opt.datas?opt.datas:{};
			opt.method = opt.method?opt.method.toUpperCase():"GET";
			opt.headers= opt.headers?opt.headers:{};

			{
				// 构建请求参数
				opt.payload = _http.buildPayload(opt.datas);
				opt.xhrPayload = opt.payload.payload_encode;
				// GET 通过Url传递参数
				if(opt.xhrPayload&&opt.xhrPayload.length>0){
					if("GET"==opt.method){
						opt.uri = opt.uri+"?"+opt.xhrPayload;
						opt.xhrPayload = "";
					}
				}
				// Post 通过form传递参数
				if("POST" == opt.method){
					if( undefined == opt.headers["Content-Type"] ){
						opt.headers["Content-Type"] = "application/x-www-form-urlencoded;charset=UTF-8";
					}
				}
			}
			// 
			opt.xhr = new XMLHttpRequest();
			opt.xhr.open(opt.method, opt.uri, opt.async);
			opt.xhr.onreadystatechange = function( ){
				if(opt.onreadystatechange){
					opt.onreadystatechange( opt.xhr, opt );
				}
			};
			{
				if(opt.timeout){
					opt.xhr.timeout = parseInt(opt.timeout);
				}
				for(var key in opt.headers){
					opt.xhr.setRequestHeader(key, opt.headers[key]);
				}
			}
			// 
			opt.setHeader = function(key, val){
				opt.headers[key.toString()] = val;
				opt.xhr.setRequestHeader(key, val);
			};
			opt.do = function( callback ){
				if(callback){
					opt.onreadystatechange = callback;
				}
				opt.xhr.send(opt.xhrPayload);
			};
			return opt;
		},
		/**
		 * Api自动签名请求
		 */
		apiRequest: function(opt){
			var session = _http.getSession( );
			opt.session = session;
			if(!session||!session.SecretKey||!session.AccessKey){
				throw new Error("Signature is empty");
			}
			var request = _http.AjaxRequest( opt );
			var signPayload = request.payload.payload;
			if(signPayload&&signPayload.length>0){
				signPayload += opt.session.SecretKey;
			}else{
				signPayload = opt.session.SecretKey;
			}
			request.setHeader("ack", session.AccessKey);
			request.setHeader("sign", md5(signPayload));
			//console.log(request);
			return request;
		},
		/**
		* HTTP请求负载构建, 去除空字段并排序字段
		*/
		buildPayload: function (params) {
			if(!params){ return ""; }
			// 去除无效字段
			var _paramsmap = new Map( );
			for(var key in params){
				if (params[key] == undefined || params[key] == null || params[key].length == 0) {
					continue;
				}
				_paramsmap.set(key, params[key]);
			}
			// 构建请求负载
			var _payload, _payload_encode = "";
			if(_paramsmap.size > 0 ){
				var _keys = _paramsmap.keys( ).sort( );
				var _payloads = []; var _payloads_encode = [];
				for (var i = 0; i < _keys.length; i++) {
					var _val = _paramsmap.get(_keys[i]);
					_payloads.push(_keys[i]+"="+_val);
					_payloads_encode.push(_keys[i]+"="+encodeURIComponent(_val));
				}
				_payload = _payloads.join("&");
				_payload_encode = _payloads_encode.join("&");
			}
			return {
				payload: _payload,
				payload_encode: _payload_encode,
			};
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
					if( $tools.isString(xhr.responseText) ){
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
			_http.apiResponseStautsHandler(result);
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
				_http.apiRequest({
					uri: uri,
					datas: datas,				
				}).do(function(xhr, opt){
					if(xhr.readyState === 4){
						var res = _http.apiResponseFormat(xhr);
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
				_http.apiRequest({
					uri: uri,
					method:"POST",
					datas: datas,				
				}).do(function(xhr, opt){
					if(xhr.readyState === 4){
						var res = _http.apiResponseFormat(xhr);
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
			sessionStorage.setItem(_key, "");
			if(accessObj && accessObj.AccessKey){
				sessionStorage.setItem(_key, JSON.stringify(accessObj));
				$tools.setCookie("ack", accessObj.AccessKey);
			}
		},
		/**
		 * 获取会话信息
		 * {UserId, AccessKey, SecretKey}
		 */
		getSession: function(){
			var _key = "access_object";
			try{
				return JSON.parse(sessionStorage.getItem(_key));
			}catch(err){
				return null;
			}
		},
	};

	/*_http.saveSession({
		UserId: "system",
		AccessKey: "0000000000000000000000000000",
		SecretKey: "0000000000000000000000000000",
	});*/
	return _http;
}));