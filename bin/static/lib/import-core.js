// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

if(typeof __import_core__=="undefined"){   
	__import_core__=true;  
	var version = "2020.07.04"; 
	// JS
	document.write("<script type='text/javascript' src='/lib/common/extends.js?v="+version+"'></script>");
	document.write("<script type='text/javascript' src='/lib/common/utils.js?v="+version+"'></script>");
	document.write("<script type='text/javascript' src='/lib/3party_library/vue/vue.min.js?v="+version+"'></script>");
	document.write("<script type='text/javascript' src='/lib/3party_library/iview/iview.min.js?v="+version+"'></script>");
	document.write("<script type='text/javascript' src='/lib/3party_library/md5.js?v="+version+"'></script>");
	document.write("<script type='text/javascript' src='/lib/tools/futil.js?v="+version+"'></script>");
	document.write("<script type='text/javascript' src='/lib/tools/apitools.js?v="+version+"'></script>");
	document.write("<script type='text/javascript' src='/lib/common/vue/vue.components.js?v="+version+"'></script>");
	document.write("<script type='text/javascript' src='/lib/tools/vue/vue.components.js?v="+version+"'></script>");
	
	// CSS
	document.write("<link rel='stylesheet' type='text/css' href='/lib/3party_library/iview/styles/iview.css?v="+version+"'/>");
	document.write("<link rel='stylesheet' type='text/css' href='/css/common.css?v="+version+"'/>");

	{ // 停止所有拖拽动作
		try{
			document.ondrop		 = function( Even ){ Even.preventDefault( ); Even.stopPropagation( ); };
			document.ondragover  = function( Even ){ Even.preventDefault( ); Even.stopPropagation( ); };
			document.ondragleave = function( Even ){ Even.preventDefault( ); Even.stopPropagation( ); };
			document.ondragenter = function( Even ){ Even.preventDefault( ); Even.stopPropagation( ); };
			document.oncontextmenu = function( Even ){ Even.preventDefault( ); Even.stopPropagation( ); };
		}catch( Err_Catch ){ }
	};
}//endif
