// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// 附加在 $utils 上的属性
;(function (root, factory) {
	if (typeof exports === "object") {
		module.exports = exports = factory();
	}else {
		if( root.$utils && root.$utils.extendAttrs ){
			root.$utils.extendAttrs(root.$utils, factory( ));
		}else{
			root.$utils = factory( );
		}
	}
}(this, function ( ){
	"use strict";
	let utils = {
		// 获取图标URL
		iconUrl:function(path){
			let icons = ["ai","avi","bmp","catdrawing","catpart","catproduct","cdr","csv","doc","docx","dps","dpt","dwg","eio","eml","et","ett","exb","exe","file","flv","fold","gif","htm","html","jpeg","jpg","mht","mhtml","mid","mp3","mp4","mpeg","msg","odp","ods","odt","pdf","png","pps","ppt","pptx","prt","psd","rar","rm","rmvb","rtf","sldprt","swf","tif","txt","url","wav","wma","wmv","wps","wpt","xls","xlsx","zip"];
			if(path){
				let type = path.substring(path.lastIndexOf(".")+1).toLowerCase();				
				for(let i=0;i<icons.length;i++){
					if(icons[i] == type ){
						return "/img/file_icons/"+type+".png";
					}
				}
			}
			return "/img/file_icons/file.png";
				
		},
	};
	return utils;
}));