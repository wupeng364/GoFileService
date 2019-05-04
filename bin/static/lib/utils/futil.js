;(function (root, factory) {
	/**
	 * @description 附加在 $tools 上的属性
	 * @author	wupeng364@outlook.com
	*/
	if( $tools && $tools.extendAttrs ){
		$tools.extendAttrs($tools, factory( ));
	}else{
		root.$tools = factory( );
	}
}(this, function ( ){
	var _Utils = {
		// 获取图标URL
		iconUrl:function(path){
			var icons = ["ai","avi","bmp","catdrawing","catpart","catproduct","cdr","csv","doc","docx","dps","dpt","dwg","eio","eml","et","ett","exb","exe","file","flv","fold","gif","htm","html","jpeg","jpg","mht","mhtml","mid","mp3","mp4","mpeg","msg","odp","ods","odt","pdf","png","pps","ppt","pptx","prt","psd","rar","rm","rmvb","rtf","sldprt","swf","tif","txt","url","wav","wma","wmv","wps","wpt","xls","xlsx","zip"];
			if(path){
				var type = path.substring(path.lastIndexOf(".")+1).toLowerCase();				
				for(var i=0;i<icons.length;i++){
					if(icons[i] == type ){
						return "/img/file_icons/"+type+".png";
					}
				}
			}
			return "/img/file_icons/file.png";
				
		},
	};
	return _Utils;
}));