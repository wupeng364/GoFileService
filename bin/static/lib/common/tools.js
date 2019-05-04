;(function (root, factory) {
	/**
	 * @description 通用工具集
	 * @author	wupeng364@outlook.com
	 * @license  MIT
	*/
	if (typeof exports === "object") {
		module.exports = exports = factory();
	}else {
		// Global (browser)
		root.$tools = factory();
	}
}(this, function ( ){

	/**
	* 导出对象
	* 初始对象操作类型的函数, 其他函数分类别通过 extendAttrs 来加载
	* 可以根据需要来加载部分类别函数
	*/
	var _Commons = {
		/**
		 * 拓展|覆盖对象属性 
		 * B_IsCopy 指定是否使用拷贝的方式给原有对象添加属性, 默认使用拷贝方式
		*/ 
		extendAttrs: function( OB_Src, OB_Add, B_IsCopy ){
			if( !OB_Src ){
				return OB_Add;
			}
			if(!_Commons.isType( OB_Src, "object" )){
				OB_Src = {};
			}
			for( var Tx_Key in OB_Add){
				if(_Commons.isType( OB_Add[ Tx_Key ], "object" ) && B_IsCopy !== false ){
					OB_Src[ Tx_Key ] = _Commons.extendAttrs(OB_Src[Tx_Key], OB_Add[ Tx_Key ]);
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
			return _Commons.isType(o, "Array");
		},
		// 是否是字符
		isString: function(o){
			if(o==undefined){
				return false;
			}
			return _Commons.isType(o, "string");;
		},
	};

	/**
	 * 字符转换|生成 类型函数
	*/ 
	_Commons.extendAttrs(_Commons, {
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

	/**
	 * 时间转换 类型函数
	*/ 
	_Commons.extendAttrs(_Commons, {
		// 时间戳转 yyyy-MM-dd HH:mm 格式时间
		long2Time: function( time ){
			if( !time || time == 0 ){
				return "";
			}
			return _Commons.dateFormat(new Date(time), 'yyyy-MM-dd HH:mm')
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
		// 获取今天凌晨
		getTerdayStart: function( ){
			var date = new Date( );
			date.setHours( 0, 0, 0, 0 );
			return date;
		},
		// 获取昨天日期
		getYesterday: function( ){
			var date = new Date( );
			date.setHours( 0, 0, 0, 0 );
			date.setDate( date.getDate( ) - 1 );
			return date;
		},
		// 获取本周的第一天
		getCurrentWeekFirst: function( ){
			var date = new Date( );
			date.setHours( 0, 0, 0, 0 );
			date.setDate( date.getDate( ) - (date.getDay( )||7)+1 );
			return date;
		},
		// 获取本周的最后一天
		getCurrentWeekLast: function( ){
			var date = new Date( );
			date.setHours( 0, 0, 0, 0 );
			date.setDate( date.getDate( ) - (date.getDay( )||7)+1 );
			date.setDate( date.getDate( )+6 );
			return date;
		},
		// 获取当前月的第一天
		getCurrentMonthFirst: function( ){
			var date = new Date( );
			date.setDate(1);
			date.setHours(0, 0, 0, 0);
			return date;
		},
		// 获取当前月的最后一天
		getCurrentMonthLast: function( ){
			var date = new Date( );
			date.setMonth( date.getMonth( ) + 1 );
			date.setDate(0);
			date.setHours(23, 59, 59, 0);
			return date;
		},
		// 获取上个月的第一天
		getAfterMonthFirst: function( ){
			var date = new Date( );
			date.setMonth( date.getMonth( ) - 1 );
			date.setDate(1);
			date.setHours(0, 0, 0, 0);
			return date;
		},
		// 获取上个月的最后一天
		getAfterMonthLast: function( ){
			var date = new Date( );
			date.setDate(0);
			date.setHours(23, 59, 59, 0);
			return date;
		},
	});

	/**
	 * HTML | DOM操作 类型函数
	*/ 
	_Commons.extendAttrs(_Commons, {
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
		// From表单POST上传
		uploadByFormData: function( url, file, ctrl ){
			return new (function( _url, _file, _ctrl){
				var self = this;
				// 
				this.ctrl = {
					fileFormName: "file",
					method: "POST",
					header: { },
					form: {},
					loadstart: function( e ){ },
					load: function( e ){ },
					loadend:  function( e ){ },
					error: function( e ){ },
					abort: function( e ){ },
					progress: function( e ){ },
				};
				this.ctrl = _Commons.extendAttrs(this.ctrl, _ctrl);

				// HttpRequest
				this.xhr = new XMLHttpRequest( );
				this.xhr.upload.addEventListener('loadstart', this.ctrl.loadstart);
			    this.xhr.upload.addEventListener('load', this.ctrl.load);
				this.xhr.upload.addEventListener("loadend", this.ctrl.loadend);
				this.xhr.upload.addEventListener("progress", this.ctrl.progress);
			    this.xhr.upload.addEventListener('error', this.ctrl.error);
			    this.xhr.upload.addEventListener('abort', this.ctrl.abort);

				// FormData
				this.formData = new FormData( );
				for(var key in this.ctrl.form ){
					this.formData.append(key, this.ctrl.form[key]);
				}
				this.formData.append(this.ctrl.fileFormName, _file);

				// Method
				this.abort = function( ){
					this.xhr.abort( );
				};
				// Post
				this.start = function( ){
					this.xhr.open(this.ctrl.method, _url);
					for( var key in this.ctrl.header ){
						this.xhr.setRequestHeader(key, this.ctrl.header[key])
					}
					this.xhr.overrideMimeType("application/octet-stream");
					self.xhr.send( self.formData );
				};

			})( url, file, ctrl );
		},
		// FileReader 请求体上传
		uploadByFileReader: function ( url, blob, ctrl ){
			return new (function( _url, _blob, _ctrl){
				var self = this;
				// 
				this.ctrl = {
					method: "POST",
					header: {
						"Content-Type": "text/plain"
					},
					progress: function( loaded, e ){ }
				};
				this.ctrl = _Commons.extendAttrs(this.ctrl, _ctrl);

				// HttpRequest
				this.xhr = new XMLHttpRequest( );
				this.xhr.upload.addEventListener("progress", function(e){
					if(e.lengthComputable){
						self.ctrl.progress(Math.round((e.loaded * 100) / e.total), e);
					}
				}, false);

				// FileReader
				this.reader = new FileReader( );
				this.reader.onload = function(evt) {
					Uploader._Consolelog( evt )
					self.xhr.send(evt.target.result);
				};
				// Method
				this.abort = function( ){
					this.reader.abort( );
					this.xhr.abort( );
				};
				// Post
				this.start = function( ){
					this.xhr.open(this.ctrl.method, _url);
					for( var key in this.ctrl.header ){
						this.xhr.setRequestHeader(key, this.ctrl.header[key])
					}
					this.xhr.overrideMimeType("application/octet-stream");
					self.reader.readAsArrayBuffer(_blob);
				};

			})( url, blob, ctrl );
		},
	});

	return _Commons;
}));