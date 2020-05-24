// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 通用工具集: 对象合并, 字符转换|生成, 时间转换, DOM操作, HTTP请求, 文件上传


;(function (root, factory) {
	if (typeof exports === "object") {
		module.exports = exports = factory();
	}else {
		// Global (browser)
		root.$utils = factory();
	}
}(this, function ( ){

	// utils 导出对象, 每个分类通过 utils.extendAttrs 函数拓展, 可以根据需要来加载部分类别函数
	var utils = {
		// 拓展|覆盖对象属性, B_IsCopy 指定是否使用拷贝的方式给原有对象添加属性, 默认使用拷贝方式
		extendAttrs: function( OB_Src, OB_Add, B_IsCopy ){
			if( !OB_Src ){
				return OB_Add;
			}
			if(!utils.isType( OB_Src, "object" )){
				OB_Src = {};
			}
			for( var Tx_Key in OB_Add){
				if(utils.isType( OB_Add[ Tx_Key ], "object" ) && B_IsCopy !== false ){
					OB_Src[ Tx_Key ] = utils.extendAttrs(OB_Src[Tx_Key], OB_Add[ Tx_Key ]);
				}else{
					OB_Src[Tx_Key] = OB_Add[ Tx_Key ];
				}
			} 
			return OB_Src; 
		},
		// 获取对象类型
		getType: function( obj ){
			try{
				var testType = Object.prototype.toString.call(obj).slice(8,-1).toLowerCase( );
				return testType.toLowerCase( ); 
			}catch( e ){
				return e;
			}
		},
		// 判断是否是 XX 类型
		isType: function( obj, type ){
			try{
				var testType = Object.prototype.toString.call(obj).slice(8,-1).toLowerCase( );
				return (testType === type.toLowerCase( )); 
			}catch( e ){
				return false;
			}
		},
		// 是否是数组
		isArray: function(o){
			if(o==undefined){
				return false;
			}
			return utils.isType(o, "Array");
		},
		// 是否是字符
		isString: function(o){
			if(o==undefined){
				return false;
			}
			return utils.isType(o, "string");;
		},
		// 是否具有某个属性
		hasKey: function(obj, key){
			if( !obj ){
				return false;
			}
			return Object.prototype.hasOwnProperty.call(obj, key);
		}
	};

	// 字符转换|生成
	utils.extendAttrs(utils, {
		regExp:{
			CHINESE_CHARACTER: /[\u4e00-\u9fa5]/,
			NAME: /^[a-zA-Z\u4e00-\u9fa5]+$/,
			HTTP_ALL: /http(s?):\/\/[A-Za-z0-9]+\.[A-Za-z0-9]+[\/=\?%\-&_~`@[\]\’:+!]*([^<>\"\"])*/g,
			HTTP_STRICT: /((http|ftp|https):\/\/[\w\-_]+(\.[\w\-_]+)+([\w\-\.,@?^=%&amp;:\/~\+#]*[\w\-\@?^=%&amp;\/~\+#])?|www+(\.[\w\-_]+)+([\w\-\.,@?^=%&amp;:\/~\+#]*[\w\-\@?^=%&amp;\/~\+#])?)/gi,
			HTTP: /^http(s?):\/\/[A-Za-z0-9]+\.[A-Za-z0-9]+[\/=\?%\-&_~`@[\]\’:+!]*([^<>\"\"])*$/,
			DOMAIN : /^(http|https|ws):\/\/([^/:]+)(:\d*)?\//,
			EMAIL: /^((([a-z]|\d|[!#\$%&'\*\+\-\/=\?\^_`{\|}~]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])+(\.([a-z]|\d|[!#\$%&'\*\+\-\/=\?\^_`{\|}~]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])+)*)|((\x22)((((\x20|\x09)*(\x0d\x0a))?(\x20|\x09)+)?(([\x01-\x08\x0b\x0c\x0e-\x1f\x7f]|\x21|[\x23-\x5b]|[\x5d-\x7e]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(\\([\x01-\x09\x0b\x0c\x0d-\x7f]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF]))))*(((\x20|\x09)*(\x0d\x0a))?(\x20|\x09)+)?(\x22)))@((([a-z]|\d|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(([a-z]|\d|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])([a-z]|\d|-|\.|_|~|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])*([a-z]|\d|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])))\.)+(([a-z]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])|(([a-z]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])([a-z]|\d|-|\.|_|~|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])*([a-z]|[\u00A0-\uD7FF\uF900-\uFDCF\uFDF0-\uFFEF])))\.?$/i,
			NUMBER_AND_LETTER: /^([A-Z]|[a-z]|[\d])*$/,
			POSITIVE_NUMBER: /^[1-9]\d*$/,
			NON_NEGATIVE_NUMBER: /^(0|[1-9]\d*)$/,
			IP: /^((1?\d?\d|(2([0-4]\d|5[0-5])))\.){3}(1?\d?\d|(2([0-4]\d|5[0-5])))$/,
			URL: /^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,6}$/,
			PHONE: /^((0\d{2,3})-)?(\d{7,8})(-(\d{3,}))?$|^(13|14|15|17|18)[0-9]{9}$/,
			TELEPHONE: /^(13|14|15|17|18)[0-9]{9}$/,
			QQ: /^\d{1,10}$/,
			DATE: /^((?!0000)[0-9]{4}-((0[1-9]|1[0-2])-(0[1-9]|1[0-9]|2[0-8])|(0[13-9]|1[0-2])-(29|30)|(0[13578]|1[02])-31)|([0-9]{2}(0[48]|[2468][048]|[13579][26])|(0[48]|[2468][048]|[13579][26])00)-02-29)$/,
			POUND_TOPIC: /^#([^\/|\\|\:|\*|\?|\"|<|>|\|]+?)#/,
			AT: /(@[a-zA-Z0-9_\u4e00-\u9fa5（）()]+)(\W|$)/gi,
			DIR_NAME: /[\\/:*?\"<>|]/,
			EmplName: /[\\\/:*?\"'<>\{};#!|]/,
			PWD : /((?=.*\d)(?=.*\D)|(?=.*[a-zA-Z])(?=.*[^a-zA-Z]))^.{8,32}$/
		},	
		// 随机数生成
		random: function(len){
			var _chars = 'ABCDEFGHJKMNPQRSTWXYZabcdefhijkmnprstwxyz2345678'
			len = len || 32
			var maxPos = _chars.length
			var str = ''
			for (var i = 0; i < len; i++) {
				str += _chars.charAt(Math.floor(Math.random() * maxPos))
			}
			return str
		},
		// 唯一标识符生成
		guid: (function( ){
	        var counter = 0;
	        return function( prefix ) {
	            var guid = (+new Date()).toString( 32 ),
	                i = 0;
	            for ( ; i < 5; i++ ) {
	                guid += Math.floor( Math.random() * 65535 ).toString( 32 );
	            }
	            return (prefix || '') + guid + (counter++).toString( 32 );
	        };
	    })( ),
		// 格式化文件大小
		formatSize:function(bytes){
			try{
				var sOutput = bytes + " bytes";
				for (var aMultiples = ["KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"], nMultiple = 0, nApprox = bytes / 1024; nApprox > 1; nApprox /= 1024, nMultiple++) {
					sOutput = nApprox.toFixed(3) + " " + aMultiples[nMultiple];
				}
				return sOutput;
			}catch(e){
				return bytes;
			}
		},
		// 进度百分比
		toPercent:function(num1, num2){
			return num1<=0 || num2<=0 ? 0 : (Math.round(num2 / num1 * 100));
		},
	});

	// 时间转换
	utils.extendAttrs(utils, {
		// 时间戳转 yyyy-MM-dd HH:mm 格式时间
		long2Time: function( time ){
			if( !time || time == 0 ){
				return "";
			}
			return utils.dateFormat(new Date(time), 'yyyy-MM-dd HH:mm')
		},
		// 时间戳|时间对象转某种格式的时间, 默认 yyyyMMddHHmmss
		dateFormat: function(time, format){
			var t = Object.prototype.toString.call(time) == "[object Date]"?time:new Date(time);
			if(!format){ format = "yyyyMMddHHmmss"; }
			var tf = function(i){return (i < 10 ? '0' : '') + i};
			return format.replace(/yyyy|MM|dd|HH|mm|ss/g, function(a){
				switch(a){
				case 'yyyy':
					return tf(t.getFullYear( ));
					break;
				case 'MM':
					return tf(t.getMonth( ) + 1);
					break;
				case 'dd':
					return tf(t.getDate( ));
					break;
				case 'HH':
					return tf(t.getHours( ));
					break;
				case 'mm':
					return tf(t.getMinutes( ));
					break;
				case 'ss':
					return tf(t.getSeconds( ));
					break;
				}
			})
		},
	});

	// HTML | DOM操作
	utils.extendAttrs(utils, {
		// 获取url里面的查询值
		getQueryParam: function(url, name){
			var reg = new RegExp("(^|&)"+ name +"=([^&]*)(&|$)");
			var r = url.substr(1).match(reg);
			if(r!=null){
				return  decodeURIComponent(r[2]); 
			}
			return null;
		},
		// 设置cookie
		setCookie:function (cname, cvalue, exdays) {
			var expires ="";
			if(exdays!=undefined && exdays>0){
				var d = new Date();
				d.setTime(d.getTime() + (exdays * 24 * 60 * 60 * 1000));
				var expires = "expires=" + d.toUTCString() + ";";			
			}
			//$cookie.set("access","1");
			document.cookie = cname + "=" + cvalue + ";path=/;" + expires;
		},
		// 读取cookie
		getCookie:function (cname) {
			var arr, reg = new RegExp("(^| )" + cname + "=([^;]*)(;|$)");
			if (arr = document.cookie.match(reg))
				return (arr[2]);
			else
				return null;
		},
		// 添加DOM事件     
		addEvent: function(obj, EventName, callBack, options){
			if(obj.addEventListener){   //FF     
				obj.addEventListener(EventName, callBack, options);     
			}else if(obj.attachEvent){//IE        
				obj.attachEvent('on'+EventName, callBack);     
			}else{        
				obj["on"+EventName]=callBack;      
			}   
		},
		// 触发鼠标事件
		triggerMouseEvent: function( el, ev, bubbles, cancelable ){
			var e = document.createEvent("MouseEvents");
			e.initEvent(ev, bubbles, cancelable);
			el.dispatchEvent(e); 
		},
		// 文件分片
		cuteFile: function( file, chunk_size ){
			if( !file || !chunk_size ){
				return;
			}
			file._cute = file._cute?{
				_start: file._cute._end,
				_end: file._cute._end+chunk_size
			}:{
				_start: 0,
				_end: chunk_size
			};
			return file.slice(file._cute._start, file._cute._end);
		},
		
	});

	// HTTP请求&文件上传
	utils.extendAttrs(utils, {
		// Ajax请求, opt={uri, method, datas, headers, async}
		AjaxRequest: function( opt ){
			opt.xhrPayload = "";
			opt.async  = undefined==opt.async?true:opt.async;
			opt.datas  = opt.datas?opt.datas:{};
			opt.method = opt.method?opt.method.toUpperCase():"GET";
			opt.headers= opt.headers?opt.headers:{};

			{
				// 构建请求参数
				opt.payload = utils.buildPayload(opt.datas);
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
					return opt.onreadystatechange( opt.xhr, opt );
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
				if(opt.async && callback){
					opt.onreadystatechange = callback;
				}
				opt.xhr.send(opt.xhrPayload);
				return opt.async||!callback?opt.xhr:callback(opt.xhr, opt);
			};
			return opt;
		},
		// HTTP请求负载构建, 去除空字段并排序字段
		buildPayload: function (params) {
			if(!params){ return ""; }
			// 去除无效字段
			var paramsmap = new Map( );
			for(var key in params){
				if (params[key] == undefined || params[key] == null || params[key].length == 0) {
					continue;
				}
				paramsmap.set(key, params[key]);
			}
			// 构建请求负载
			var payload, payload_encode = "";
			if(paramsmap.size > 0 ){
				var keys = paramsmap.keys( ).sort( );
				var payloads = []; var payloads_encode = [];
				for (var i = 0; i < keys.length; i++) {
					var val = paramsmap.get(keys[i]);
					payloads.push(keys[i]+"="+val);
					payloads_encode.push(keys[i]+"="+encodeURIComponent(val));
				}
				payload = payloads.join("&");
				payload_encode = payloads_encode.join("&");
			}
			return {
				payload: payload,
				payload_encode: payload_encode,
			};
		},
		// obj to url
		json2query: function(obj, prefix) {
			prefix = prefix || '';
			var pairs = [];
			var has = Object.prototype.hasOwnProperty;
			if ('string' !== typeof prefix){
				prefix = '?';
			}
			for (var key in obj) {
				if (has.call(obj, key)) {
					pairs.push(key +'='+ encodeURIComponent(obj[key]));
				}
			}
			return pairs.length ? prefix + pairs.join('&') : '';
		},
		// From表单POST上传
		uploadByFormData: function( url, file, ctrl ){
			var uploader = {
				ctrl: {
					filekey: "file",
					method: "POST",
					header: { },
					form: {},
					loadstart: function( e ){ },
					load: function( e ){ },
					loadend:  function( e ){ },
					error: function( e ){ },
					abort: function( e ){ },
					progress: function( e ){ },
				},
				xhr: new XMLHttpRequest( ),
				formData: new FormData( ),
				abort: function( ){
					uploader.xhr.abort( );
				},
				start: function( ){
					uploader.xhr.open(uploader.ctrl.method, url);
					for( var key in uploader.ctrl.header ){
						uploader.xhr.setRequestHeader(key, uploader.ctrl.header[key])
					}
					uploader.xhr.overrideMimeType("application/octet-stream");
					uploader.xhr.send( uploader.formData );
				}
			};
			// 合并配置
			uploader.ctrl = utils.extendAttrs(uploader.ctrl, ctrl);
		    // 添加表单
			for(var key in uploader.ctrl.form ){
				uploader.formData.append(key, uploader.ctrl.form[key]);
			}
			uploader.formData.append(uploader.ctrl.filekey, file);
			// 绑定事件
			uploader.xhr.upload.addEventListener('loadstart', uploader.ctrl.loadstart);
		    uploader.xhr.upload.addEventListener('load', uploader.ctrl.load);
			uploader.xhr.upload.addEventListener("loadend", uploader.ctrl.loadend);
			uploader.xhr.upload.addEventListener("progress", uploader.ctrl.progress);
		    uploader.xhr.upload.addEventListener('error', uploader.ctrl.error);
		    uploader.xhr.upload.addEventListener('abort', uploader.ctrl.abort);
		    uploader.xhr.onreadystatechange = function(e){
		    	var xhr = e.target;
		    	if(xhr.readyState == 4 && xhr.status!=200){
		    		uploader.ctrl.error(new Error(xhr.responseText?xhr.responseText:xhr.statusText) );
		    	}
		    };
		    return uploader;
		},
		// FileReader 请求体上传
		uploadByFileReader: function ( url, blob, ctrl ){
			var uploader = {
				ctrl: {
					method: "POST",
					header: {
						"Content-Type": "text/plain"
					},
					progress: function( loaded, e ){ }
				},
				xhr: new XMLHttpRequest( ),
				reader: new FileReader( ),
				abort: function( ){
					uploader.reader.abort( );
					uploader.xhr.abort( );
				},
				start: function( ){
					uploader.xhr.open(uploader.ctrl.method, url);
					for( var key in uploader.ctrl.header ){
						uploader.xhr.setRequestHeader(key, uploader.ctrl.header[key])
					}
					uploader.xhr.overrideMimeType("application/octet-stream");
					self.reader.readAsArrayBuffer(blob);
				}
			};
			uploader.ctrl = utils.extendAttrs(uploader.ctrl, ctrl);
			uploader.xhr.upload.addEventListener("progress", function(e){
					if(e.lengthComputable){
						uploader.ctrl.progress(Math.round((e.loaded * 100) / e.total), e);
					}
			}, false);
			uploader.reader.onload = function(evt) {
				self.xhr.send(evt.target.result);
			};
			return uploader;
		},
	});
	return utils;
}));