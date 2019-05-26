/**
 *@description httpUtils
 *@author	wupeng364@outlook.com
*/
;(function (root, factory) {
	if( !$tools ){
		throw new Error("need $tools object, please require tools.js");
	}
	if( !axios ){
		throw new Error("need axios object, please require axios.XXX.js");
	}
	if (typeof exports === "object") {
		module.exports = exports = factory();
	}else {
		// Global (browser)
		root.$fhttp = factory();
	}
}(this, function ( ){

	axios.defaults.timeout = 60000; 
	axios.defaults.withCredentials = true;
	/**
	 * 格式化服务器响应的数据, 规范格式
	 * 在此处可以自定义一些错误描述
	 */
	function formatResponse( response ){
		//console.log("response", response )
		var res = {
			IsError: false,
			Code: response.status,
			Data: "",
			Msg: "",
		};
		res.IsError = res.Code != 200;
		if($tools.hasKey(response, "data")){
			var res_data = response.data;
			if( res_data ){
				if( $tools.isString(res_data) ){
					res_data = JSON.parse( res_data );
				}

				// 写入数据
				res.Data = res_data.Data?res_data.Data:res_data;
				res.Msg = res_data.Msg?res_data.Msg:"";
			}
		}
		// 
		res.toString = function( ){
			if( res.IsError ){
				return res.Msg?res.Msg:res.Data;
			}
			return JSON.stringify(res);
		}
		return res;
	};

	var _http = {
		// 通过 axios 发起请求
		fetch: function(uri, method, params) {
			// reqConfig.data=FormData(表单传值) 	/  reqConfig.params=Query String(url带参)
			// POST		支持reqConfig.data传递参数	 FormData模式, 会在request拦截器做处理, 将json转querystring格式
			// GET		支持reqConfig.params传递参数 querystring模式
			// DELETE    支持reqConfig.params传递参数 querystring模式
			// PUT       支持reqConfig.params传递参数 querystring模式
			// DELETE: 设置reqConfig.data 服务器端无法接收到参数, 需要设置为reqConfig.params
			if(params){
				for(var key in params){
					if(params[key]==undefined){
						delete params[key];
					}
				}
			}
			var methodUppercase = method.toUpperCase( );
			var reqConfig = {url: uri, method: method, headers: {}};
			if(methodUppercase == "POST" ){
				//此头信息对post有效
				reqConfig["data"] = $tools.json2query(params); 
				reqConfig.headers["Content-Type"] = "application/x-www-form-urlencoded;charset=UTF-8";
			}else{
				reqConfig["params"] = params;
			}
			reqConfig.headers["X-Requested-With"] = "XMLHttpRequest";
			return axios(reqConfig);		
		},
				// Get 请求 - 参数在url上
		get: function(uri, datas){
			if(datas==null || typeof datas=='undefined'){
			}else{
				uri = uri + "?" + $tools.json2query(datas);
			}
			return new Promise(function(resolve, reject){
				_http.fetch(uri, "get").then( function(res) {
					var Res = formatResponse( res );	
					if (!Res.IsError) {
						resolve(Res.Data);
					} else {
						reject(Res.Msg);
					}
				}).catch(function( err ){
					reject(err.response?formatResponse(err.response):err);
				});
			})
		},
		// Post 请求 - 参数在FormData上
		post: function (uri, datas) {
			return new Promise( function(resolve, reject) {
				_http.fetch( uri, "post", datas).then(function(res) {
					var Res = formatResponse( res );				
					if (!Res.IsError) {
						resolve(Res.Data)
					} else {
						reject(Res.Msg);
					}
				}).catch(function( err ){
					reject(err.response?formatResponse(err.response):err);
				});
			})
		},
		// Put 请求 - 参数在url上
		put: function( uri, datas )  {
			return new Promise(function(resolve, reject){
				_http.fetch( uri, datas, "put").then(function(res){
					var Res = formatResponse( res );
					if (!Res.IsError) {
						resolve(Res.Data)
					} else {
						reject(Res.Msg);
					}
				}).catch(function( err ){
					reject(err.response?formatResponse(err.response):err);
				});
			})
		},
		// Del 请求 - 参数在url上
		delete: function( uri, datas )  {
			return new Promise(function(resolve, reject){
				_http.fetch( uri, datas, "delete").then(function(res){
					var Res = formatResponse( res );
					if (!Res.IsError) {
						resolve(Res.Data)
					} else {
						reject(Res.Msg);
					}
				}).catch(function( err ){
					reject(err.response?formatResponse(err.response):err);
				});
			})
		},

		session: function( ){
			return {accesskey:undefined};
		},
		buildHeader: function(ak, method, uri, datas ){
			return { };
		},
	};

	return _http;
}));