// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// VUE组件, 无业务属性

// 自动适应高度, 自动减去某个值
Vue.directive('minus-height',{
	// 绑定钩子函数
    bind: function (el,binding,vnode){	
    },
	// 绑定到节点函数
	inserted:function(){
	},
	// 组件更新钩子函数
    update: function (el, binding, vnode) {		
    },
	// 组件更新完成
    componentUpdated: function (el, binding, vnode, vnodeold){
		if(vnode.v_UnBindMinusHeight){
			vnode.v_UnBindMinusHeight( );
		}
		if(vnodeold.v_UnBindMinusHeight){
			vnodeold.v_UnBindMinusHeight( );
		}
		vnode.v_MinusHeight = function( listen ){
			try{
				if( el.parentNode.clientHeight == 0 ){
					return;
				}
				vnode.componentInstance.height = el.parentNode.clientHeight-binding.value;	
			}catch(e){ }
	
			if(listen === true){
				window.addEventListener("resize", vnode.v_MinusHeight, false);
			}
		};
		vnode.v_UnBindMinusHeight = function(){
			window.removeEventListener("resize", vnode.v_MinusHeight, false);
			vnode.v_MinusHeight = undefined;
			vnode.v_UnBindMinusHeight = undefined;
		}
		vnode.v_MinusHeight( true );
	},
	// 解除指令 
    unbind: function( el, binding, vnode ){
		vnode.v_UnBindMinusHeight();
    }
});