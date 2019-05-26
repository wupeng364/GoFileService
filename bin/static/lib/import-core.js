if(typeof __import_core__=="undefined"){   
	__import_core__=true;  
	var version = "2019.5.24"; 
	// JS
	document.write("<script type='text/javascript' src='/lib/common/extends.js?v="+version+"'></script>");
	document.write("<script type='text/javascript' src='/lib/common/tools.js?v="+version+"'></script>");
	document.write("<script type='text/javascript' src='/lib/3party_library/axios.min.js?v="+version+"'></script>");
	document.write("<script type='text/javascript' src='/lib/3party_library/promise.auto.min.js?v="+version+"'></script>");
	document.write("<script type='text/javascript' src='/lib/3party_library/vue/vue.min.js?v="+version+"'></script>");
	document.write("<script type='text/javascript' src='/lib/3party_library/vue/vue-router.js?v="+version+"'></script>");
	document.write("<script type='text/javascript' src='/lib/3party_library/iview/iview.min.js?v="+version+"'></script>");
	document.write("<script type='text/javascript' src='/lib/utils/futil.js?v="+version+"'></script>");
	document.write("<script type='text/javascript' src='/lib/utils/fhttp.js?v="+version+"'></script>");
	document.write("<script type='text/javascript' src='/lib/vue.components.js?v="+version+"'></script>");
	
	// CSS
	document.write("<link rel='stylesheet' type='text/css' href='/lib/3party_library/iview/styles/iview.css?v="+version+"'/>");
	document.write("<link rel='stylesheet' type='text/css' href='/css/common.css?v="+version+"'/>");

	{ // 停止所有拖拽动作
		try{
			document.ondrop		 = function( Even ){ Even.preventDefault( ); Even.stopPropagation( ); };
			document.ondragover  = function( Even ){ Even.preventDefault( ); Even.stopPropagation( ); };
			document.ondragleave = function( Even ){ Even.preventDefault( ); Even.stopPropagation( ); };
			document.ondragenter = function( Even ){ Even.preventDefault( ); Even.stopPropagation( ); };
		}catch( Err_Catch ){ }
	};
}//endif