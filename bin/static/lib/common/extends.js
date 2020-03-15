// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// String原型拓展&拓展对象

"use strict";

// 路径处理去掉最后一个 / 避免 [ /admin/ ]出现
String.prototype.getPath = function( ){
	var temp=this.toString( );
	if( this!=null && this!="" ){	
		var lstIndex = this.lastIndexOf("/"); 
		if( lstIndex == this.length-1 ){ 
			temp = this.substring(0,this.lastIndexOf("/"));
		}
	}
	return temp;
};

// 路径转换
String.prototype.parsePath = function( op ){
	if( op && op["windows"] ){
		return this.replaceStr("//", "/").replaceStr("/", "\\");
	}
	return this.replaceStr("\\", "/").replaceStr("//", "/");
};

// 替换字符
String.prototype.replaceStr = function( splitStr, replaceStr ){

	var strs = this.split( splitStr.toString( ) ); 
	var temp = ""; 
	for (var i=0;i<strs.length ;i++ ){ 
		temp += strs[i];
		if( i != strs.length-1 && replaceStr ){
			temp +=replaceStr;
		} 
	} 
	return temp;
};

// 获取绝对路径
String.prototype.getAbsolutePath = function( basePath ){
	try{
	  var js_dir = basePath.parsePath( ).getPath( );
	  var relative_path_arry = this.toString( ).split("/");
	  var upLevelCout = 0;
	  var relative_path_splite = "";
	  for( var i = 0; i<relative_path_arry.length; i++ ) {
	    if( relative_path_arry[i] != ".." && relative_path_arry[i] != "." ){
	        relative_path_splite+="/"+relative_path_arry[i];
	    }else if( relative_path_arry[i] == ".." ){ 
	      upLevelCout++; 
	    }
	  }
	  for( var i=0; i<upLevelCout;i++ ){
	    js_dir = js_dir.getParent( );
	  }
	  js_dir += relative_path_splite;
	  return js_dir.parsePath( );
	}catch(e){ 
		return basePath;
	}
};
// 得到名字
String.prototype.getName = function( B_GetSuffixed ){
	var temp = this.toString( ).parsePath( );
	temp = temp.substring( temp.lastIndexOf("/")+1 );
	if( B_GetSuffixed == false && temp.lastIndexOf(".") != -1 ){
	 	temp = temp.substring( 0,temp.lastIndexOf(".") );	
	}
	return temp;
};
// 得到后缀 
String.prototype.getSuffixed = function( B_HavePoint ){
	var temp = this.toString( ).getName( );
	if( temp.lastIndexOf(".") == -1 ){
		return "";
	}
	var subIndex = temp.lastIndexOf(".");
	if( B_HavePoint == false ){
		subIndex++;
	}
	temp = temp.substring( subIndex );
	return temp;
}

// 得到父级路径 
String.prototype.getParent = function( ){
	var temp = this.getPath(  ); 
	if( temp.indexOf("/") > -1 ){
		temp = temp.substring(0,temp.lastIndexOf("/"));
		temp = (temp == ""?"/":temp);
	}else{
		temp = temp.substring(0,temp.lastIndexOf("\\"));
		temp = (temp == ""?"\\":temp);
	}
	return temp;
};

// 转换成系统路径
String.prototype.getPathForKass = function( parentPath ){ 
	var temp = this.getPath( );
	if( parentPath!=null && parentPath.length > 0 ){
		parentPath = parentPath.getPath( ); 
	 	temp = temp.substring( parentPath.length ); 
	}
	return temp;
};

// 是否以XXX开始
String.prototype.startWith = function(str){  
    if(str==null||str==""||this.length==0||str.length>this.length)  
      return false;  
    if(this.substr(0,str.length)==str)  
      return true;  
    else  
      return false;  
    return true;  
};
// 是否以XXX结束
String.prototype.endWith = function(str){  
    if(str==null||str==""||this.length==0||str.length>this.length)  
      return false;  
    if(this.substring(this.length-str.length)==str)  
      return true;  
    else  
      return false;  
    return true;  
} 
// 获取ip地址端口
String.prototype.getPort =function(  ) {
  var lastIndex1 = this.lastIndexOf(":");
  var lastIndex2 = this.lastIndexOf("/");
  if ( lastIndex2 < lastIndex1 ) {
      if( lastIndex1 >0 && lastIndex1 < this.length && lastIndex1 > 3 && lastIndex2 > 5 ){
          return this.substring(lastIndex1+1);
      }
  }
   return null;
}
/*
* 添加协议头 www.baidu.com ==> https://www.baidu.com
* 添加协议默认端口 https://www.baidu.com ===> https://www.baidu.com:443
* return https://www.baidu.com:443
*/
String.prototype.getStandardUrl = function( Tx_Default ){
	var Tx_Temp = this.toString( );
	if( !Tx_Temp ){ return ""; }
	var B_AddProtocol = false;
	// www.baidu.com ==> https://www.baidu.com
	if( !Tx_Temp.startWith( "http://" ) && !Tx_Temp.startWith( "https://" ) ){
		Tx_Temp = (Tx_Default?Tx_Default:"http")+"://"+Tx_Temp;
	}

	// https://www.baidu.com ===> https://www.baidu.com:443
	if( Tx_Temp.startWith("http://") && Tx_Temp.lastIndexOf(":") < 6 ){
	    Tx_Temp +=":80";
	}else if( Tx_Temp.startWith("https://") && Tx_Temp.lastIndexOf(":") < 7 ){
	    Tx_Temp +=":443";
	}
	return Tx_Temp.toLowerCase( );
};
// 获取完整ip地址格式
String.prototype.getServer=function( ){
	var Tx_Temp = this.toString( );
	// 
	var I_StartIndex = Tx_Temp.startWith( "http://" )?7:(Tx_Temp.startWith( "https://" )?8:0); 
	var I_EndIndex   = Tx_Temp.lastIndexOf(":") > 7?Tx_Temp.lastIndexOf(":"):Tx_Temp.length;

    return Tx_Temp.substring(I_StartIndex, I_EndIndex);
};
// 获取完整地址后面的路径
String.prototype.getServerPath=function( ){
    return this.replace(/^.*?\:\/\/[^\/]+/, "");
};
/*
* 去掉协议头 https://www.baidu.com:443 ==> www.baidu.com:443
* 去掉协议默认端口 www.baidu.com:443 ===> www.baidu.com
* return www.baidu.com
*/
String.prototype.getHost = function( ){
	var Tx_Temp = this.toString( );
	if( !Tx_Temp ){ return ""; }
	if( Tx_Temp.startWith( "http://" ) ){
		Tx_Temp = Tx_Temp.substr( 7 );
		if( Tx_Temp.endWith( ":80" ) ){
			Tx_Temp = Tx_Temp.substr( 0, Tx_Temp.length-3 );
		}

	}else if( Tx_Temp.startWith( "https://" ) ){
		Tx_Temp = Tx_Temp.substr( 8 );
		if( Tx_Temp.endWith( ":443" ) ){
			Tx_Temp = Tx_Temp.substr( 0, Tx_Temp.length-4 );
		}
	}
	return Tx_Temp;
};

/*
* 获取协议头 默认 http
*/
String.prototype.getUrlProtocol = function( Tx_Default ){
	var Tx_ServerStr = this.toString( );
	if( Tx_ServerStr ){ 
		if( Tx_ServerStr.startWith( "http://" ) ){
			return "http";

		}else if( Tx_ServerStr.startWith( "https://" ) ){
			return "https";
		}
	}
	return Tx_Default?Tx_Default:"http";
};

Array.prototype.remove=function(dx){
　　if(isNaN(dx)||dx>this.length){return false;}
　　for(var i=0,n=0;i<this.length;i++)
　　{
　　　　if(this[i]!=this[dx])
　　　　{
　　　　　　this[n++]=this[i]
　　　　}
　　}
　　this.length-=1
};

// Map 兼容的map对象
;(function(root, handler){
	root.Map = handler;
})(this, function( ){
	var _ = this;
	this.entry = new Object( );
	this.size = 0;
	this.set = function (key , value){
		if(!this.has(key)){
			_.size ++ ;
		}
		_.entry[key] = value;
	};
	this.get = function (key){
		return this.has(key) ? _.entry[key] : null;
	};
	this.has = function ( key ){
		return (key in _.entry);
	};
	this.keys = function( ){
		var keys = new Array();
		for(var prop in _.entry){
			keys.push(prop);
		}
		return keys;
	};
	this.delete = function ( key ){
		if( this.has(key) && ( delete _.entry[key] ) ){
			_.size --;
		}
	};
	this.values = function(){
	    var values = new Array();
	    for(var prop in _.entry){
	      values.push(_.entry[prop]);
	    }
	    return values;
  	};
	this.clear = function( ){
		_.size = 0;
		_.entry = new Object( );
	};
});