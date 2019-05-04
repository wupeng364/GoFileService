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

	axios.defaults.timeout = 60000; 
	axios.defaults.withCredentials = true;
	axios.defaults.headers.post['Content-Type'] = 'application/x-www-form-urlencoded;charset=UTF-8';
	var _http = {
		session:function( ){
			return {accesskey:undefined};
		},
		buildHeader:function(ak, method, uri, datas ){
			return {};
		},
		fetchSync: function(method, uri, datas){
			if(datas!=undefined){ 
				for(var key in datas){
					if(datas[key]==undefined){
						delete datas[key];
					}
				}
			}
			var headers = {};		
			var session = _http.session( );
			var ak = (session==null || session.accesskey==null || (typeof session.accesskey == 'undefined')) ? null : session.accesskey;
			if(ak!=null){
				headers = _http.buildHeader(ak, method.toUpperCase(), uri, datas);
			}

			var url = uri;
			var payload = null;
			if(datas!=null){
				var query = "";
				for(var key in datas){
					if(query.length>0)
						query += "&";
					query += key+"="+encodeURIComponent(datas[key]);
				}

				if(method=='GET'){				
					url += "?"+query;
				}else if(method=='POST'){
					payload = query;
				}
			}

			var xhr = null;
			if (window.ActiveXObject) {  
				xhr = new ActiveXObject("Microsoft.XMLHTTP");  
			} else if (window.XMLHttpRequest) {  
				xhr = new XMLHttpRequest();  
			}

			xhr.open(method, url, false);  
			xhr.setRequestHeader('If-Modified-Since', '0');
			if(method.toUpperCase() =="POST"){
				xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded;charset=UTF-8');
			}
			if(headers!=null){
				for (var key in headers){
					xhr.setRequestHeader(key, headers[key]);
				}
			}
			try{
				xhr.send(payload);  		
			}catch(e){
				var resp = new Object();
				resp.Code = 2;
				resp.Data = "访问服务器失败";
				return resp;
			}
			var xresp = xhr.response;
			var resp = _parseJsonDatas( xresp );
			return resp;
		},
		getSync:function(uri, datas){
			var response = _http.fetchSync('GET', uri, datas);	
			if(response.Code == 0){
				return response.data;
			}else{
				throw response.data;
			}
			
		},
		postSync:function(uri, datas){
			var response = _http.fetchSync('POST', uri, datas);	
			if(response.Code == 0){
				return response.data;
			}else{
				throw response.data;
			}
		},
		get:function(uri, datas){
			if(datas==null || typeof datas=='undefined'){
			}else{
				uri = uri + "?" + json2query(datas);
			}
			return new Promise(function(resolve, reject){
				_http.fetch( uri, null,  "get").then( function(res) {
					var Res = _http.formatResponse( res );	
					if (Res.Code === 200) {
						resolve(Res.Data);
					} else {
						reject(Res.Data)
					}
				}).catch(function(err){
					reject(err)
				})
			})
		},
		post:function (uri, datas) {
			return new Promise( function(resolve, reject) {
				_http.fetch( uri, datas, "post").then(function(res) {
					var Res = _http.formatResponse( res );				
					if (Res.Code === 200) {
						resolve(Res.Data)
					} else {
						reject(Res.Data)
					}
				}).catch( function(err){
					reject(err)
				})
			})
		},
		put : function( uri, datas )  {
			return new Promise(function(resolve, reject){
				_http.fetch( uri, datas, "put").then(function(res){
					var Res = _http.formatResponse( res );
					if (Res.Code === 200) {
						resolve(Res.Data)
					} else {
						reject(Res.Data)
					}
				}).catch( function(err){
					reject(err)
				})
			})
		},
		delete : function( uri, datas )  {
			return new Promise(function(resolve, reject){
				_http.fetch( uri, datas, "delete").then(function(res){
					var Res = _http.formatResponse( res );
					if (Res.Code === 200) {
						resolve(Res.Data)
					} else {
						reject(Res.Data)
					}
				}).catch( function(err){
					reject(err)
				})
			})
		},
		fetch : function(uri, params, method) {
			//reqConfig.data=FormData(表单传值) 	/  reqConfig.params=Query String(url带参)
			//POST		支持reqConfig.data传递参数	 FormData模式，会在request拦截器做处理(axios的bug)，将json转querystring格式
			//GET		支持reqConfig.params传递参数 querystring模式
			//DELETE    支持reqConfig.params传递参数 querystring模式
			//PUT       支持reqConfig.params传递参数 querystring模式
			//DELETE:设置reqConfig.data 服务器端无法接收到参数，需要设置为reqConfig.params
			if(params){
				for(var key in params){
					if(params[key]==undefined){
						delete params[key];
					}
				}
			}
			var reqConfig = {url: uri,method:method,params:params};
			var headers = {};
			var session = _http.session();
			var ak = (session==null || session.accesskey==null || (typeof session.accesskey == 'undefined')) ? null : session.accesskey;
			if(ak!=null){
				headers = _http.buildHeader(ak, method.toUpperCase(), uri, params);							
			}
			//此头信息对post有效
			if(method.toUpperCase() =="POST"){
				delete reqConfig["params"];
				reqConfig.data=json2query(params); 
				headers["Content-Type"] = "application/x-www-form-urlencoded;charset=UTF-8";
			}
			headers["X-Requested-With"] = "XMLHttpRequest";
			reqConfig.headers=headers;
			return axios(reqConfig);		
		},
		
		multipart:function(uri, datas, progress ){
			return new Promise( function(resolve, reject){
				_http.httpUpload( uri, datas, progress,"post",true).then(function(res){
					var Res = _http.formatResponse( res );
					if (Res.Code === 200) {
						resolve(Res.Data);
					} else {
						reject(Res.msg);
					}
				}).catch( function(err){
					reject(err);
				})
			});
		},
		
		upload:function(uri, datas, callbaclUploadProgress,method,usingDefaultAuth,optionsHeaders){
			if( method == undefined || method=="undefined" ){
				method="post";
			}
			var formdata = new FormData();
			for(var key in datas){
				formdata.append(key,datas[key]);
			}
			var headers = {};		
			if (usingDefaultAuth !== undefined && usingDefaultAuth === false) {
				
			}else{
				headers = HttpSignature(method.toUpperCase(), uri, null);
			}
			headers["Content-Type"] = "multipart/form-data";
			headers["X-Requested-With"] = "XMLHttpRequest";		
			if(optionsHeaders){
				for(var key in optionsHeaders){
					headers[key] = optionsHeaders[key];
				}			
			}
			return axios({
					url: uri,
					method: method,				
					data: formdata,
					onUploadProgress:function(progressEvent){
						if(progressEvent.lengthComputable){
							if(callbaclUploadProgress){
								callbaclUploadProgress(progressEvent);
							}
						}
					}
			});		
		},
		formatResponse: function( response ) {
			var resp = response;
			if(resp==undefined || resp == null){
				return null;
			}
			try{
				if( (typeof resp)=="string" ){
					resp = _parseJsonDatas(resp);	
				}
			}catch(e){
				return resp;
			}
			//获取服务器端返回的真实数据...
			var data = {};
			if(_http.hasKey(resp,"data")){
				data = resp.data;
				if( (typeof data)=="string" ){
					data = _parseJsonDatas(data);
				}
				// 如果数据格式不正确, 构造一个正确的数据
				if( response.status === 200 &&
					!_http.hasKey(data, "Code") && 
					!_http.hasKey(data, "Data") ){
					data = {
						Code: response.status,
						Data: data
					};
				}
			}else{
				data.isError = true;
				data.msg = "unknow error";
				return data;
			}
			//将服务器端的真实数据返回在数据结构中，可能有些调用需要response相关信息
			data.source=response;
			if(data.Code == 200 || data.Code=="200"){
				data.isError = false;
			}else{
				data.isError = true;
			}		
			return data;
		},
		hasKey:function(obj,key){
			if(obj==undefined || obj == null){
				return false;
			}
			return Object.prototype.hasOwnProperty.call(obj,key);
		}
	};

	function _parseJsonDatas(baseStr) {
	    if (!baseStr || typeof baseStr != 'string') return;
	    var jsonData = null;
	    try {
	        jsonData = JSON.parse(baseStr);
	    } catch (err){
	        return null;
	    }
	    var needReplaceStrs = [];
	    _loopFindArrOrObj(jsonData,needReplaceStrs);
	    needReplaceStrs.forEach(function (replaceInfo) {
	        var matchArr = baseStr.match(eval('/"'+ replaceInfo.key + '":[0-9]{15,}/'));
	        if (matchArr) {
	            var str = matchArr[0];
	            var replaceStr = str.replace('"' + replaceInfo.key + '":','"' + replaceInfo.key + '":"');
	            replaceStr += '"';
	            baseStr = baseStr.replace(str,replaceStr);
	        }
	    });
	    var returnJson = null;
	    try {
	        returnJson = JSON.parse(baseStr);
	    }catch (err){
	        return null;
	    }
	    return returnJson;
	};
	//遍历对象类型的
	function _getNeedRpStrByObj(obj,needReplaceStrs) {
	    for (var key in obj) {
	        var value = obj[key];
	        if (typeof value == 'number' && value > 9007199254740992){
	            needReplaceStrs.push({key:key});
	        }
	        _loopFindArrOrObj(value,needReplaceStrs);
	    }
	};
	//遍历数组类型的
	function _getNeedRpStrByArr(arr,needReplaceStrs) {
	    for(var i=0; i<arr.length; i++){
	        var value = arr[i];
	        _loopFindArrOrObj(value,needReplaceStrs);
	    }
	};
	//递归遍历
	function _loopFindArrOrObj(value,needRpStrArr) {
	    var valueTypeof = Object.prototype.toString.call(value);
	    if (valueTypeof == '[object Object]') {
	        needRpStrArr.concat(_getNeedRpStrByObj(value,needRpStrArr));
	    }
	    if (valueTypeof == '[object Array]') {
	        needRpStrArr.concat(_getNeedRpStrByArr(value,needRpStrArr));
	    }
	};

	function json2query(obj, prefix) {
	  prefix = prefix || '';
	  var pairs = [];
	  var has = Object.prototype.hasOwnProperty;
	  //
	  // Optionally prefix with a '?' if needed
	  //
	  if ('string' !== typeof prefix) prefix = '?';
	  for (var key in obj) {
	    if (has.call(obj, key)) {
	      pairs.push(key +'='+ encodeURIComponent(obj[key]));
	    }
	  }
	  return pairs.length ? prefix + pairs.join('&') : '';
	};

	return _http;
}));