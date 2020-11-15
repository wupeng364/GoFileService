// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// VUE组件, 无业务属性
; (function (factory) {
	try {
		factory();
	} catch (e) {
		console.error(e);
	}
})(function () {
	"use strict";
	// 自动适应高度, 自动减去某个值
	Vue.directive('minus-height', {
		// 绑定钩子函数
		bind: function (el, binding, vnode) { },
		// 绑定到节点函数
		inserted: function () { },
		// 组件更新钩子函数
		update: function (el, binding, vnode) { },
		// 组件更新完成
		componentUpdated: function (el, binding, vnode, vnodeold) {
			if (vnode.v_UnBindMinusHeight) {
				vnode.v_UnBindMinusHeight();
			}
			if (vnodeold.v_UnBindMinusHeight) {
				vnodeold.v_UnBindMinusHeight();
			}
			vnode.v_MinusHeight = function (listen) {
				try {
					if (el.parentNode.clientHeight == 0) {
						return;
					}
					vnode.componentInstance.height = el.parentNode.clientHeight - binding.value;
				} catch (e) { }

				if (listen === true) {
					window.addEventListener("resize", vnode.v_MinusHeight, false);
				}
			};
			vnode.v_UnBindMinusHeight = function () {
				window.removeEventListener("resize", vnode.v_MinusHeight, false);
				vnode.v_MinusHeight = undefined;
				vnode.v_UnBindMinusHeight = undefined;
			}
			vnode.v_MinusHeight(true);
		},
		// 解除指令 
		unbind: function (el, binding, vnode) {
			if (vnode.v_UnBindMinusHeight) {
				vnode.v_UnBindMinusHeight();
			}
		}
	});
	// 
	/**
	 * 右键菜单
	 * menus = {id: {name: 'xxx', icon:'', show:true, handler: func, divided:false, }}
	 */
	Vue.component("right-click-menu", {
		props: ['bindRef', 'menus'],
		data: function () {
			return {
				posX: 0,
				posY: 0,
				currentVisible: false
			}
		},
		template: `
		<Dropdown :style='locatorStyle' placement='right-start' trigger='custom' :visible='currentVisible' @on-click='onClick' @on-clickoutside='handleCancel'>
			<Dropdown-Menu slot='list'>
				<Dropdown-Item v-for='(val, key) in menus' v-show='val.show' :name='key' :divided='val.divided'>
				<i v-if='val.icon' :class='val.icon' style='padding-right: 3px;'></i> <span style='vertical-align: middle;'>{{val.name}}</span>
				</Dropdown-Item>
			</Dropdown-Menu>
		</Dropdown>
		`,
		computed: {
			locatorStyle: function () {
				return {
					position: 'fixed',
					left: this.posX + 'px',
					top: this.posY + 'px',
					maxHeight: 'unset'
				}
			}
		},
		methods: {
			onClick: function (name) {
				this.currentVisible = false
				if (this.menus[name] && this.menus[name].handler) {
					try {
						this.menus[name].handler(this.menus[name]);
					} catch (e) { console.error(e); }
				}
			},
			handleContextmenu: function (e) {
				e.returnValue = false;
				e.preventDefault();
				e.stopPropagation();
				this.currentVisible = false;
				if (e.button === 2) {
					if (this.posX !== e.clientX) { this.posX = e.clientX; }
					if (this.posY !== e.clientY) { this.posY = e.clientY; }
					this.$nextTick(() => {
						this.currentVisible = true;
					});
				}
			},
			handleCancel: function () {
				this.currentVisible = false
			},
			getRefNode: function () {
				if (!this.bindRef) {
					return document;
				}
				let node = this;
				while (true) {
					if (node.$refs[this.bindRef]) {
						break;
					}
					node = node.$parent;
					if (!node) { break; }
				}
				return node ? node.$refs[this.bindRef].$el : {};
			},
		},
		mounted: function () {
			let node = this.getRefNode();
			node.addEventListener('contextmenu', this.handleContextmenu, false);
			node.addEventListener('mouseup', this.handleContextmenu, false);
		},
		destroyed: function () {
			try {
				// let node = this.getRefNode();
				// node.removeEventListener('contextmenu', this.handleContextmenu, false);
				// node.removeEventListener('mouseup', this.handleContextmenu, false);
			} catch (e) { }
		},
		watch: {
		}
	});
});