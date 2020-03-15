// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// VUE组件 - 常用,带有业务属性的

// 文件拷贝
Vue.component("fs-copyfile", {
	props: ["show-dailog", "src-paths", "dest-path", "copy-settings"],
	data: function( ){
		var self = this;
		return {
			operations:{
				stop: "discontinue",
				ignore: "ignore",
				ignoreall: "ignoreall",
				replace: "replace",
				replaceall: "replaceall",
			},
			copySettings: {
				ignore: false,
				replace: false,
				width: 450,
			},
			operationData: {
				multiCount: 0,
				opCount: 0,
				token: "",
				srcPaths: [],
				destPath: "",
				nowSrcPath: "",
				nowDstPath: "",
			},
			copyError: {
				IsError: false,
				Error: "",
				IsExist: false,
			},
		};
	},
	template: ""+
		"<Modal v-model=\"showDailog\" :closable=\"false\" :mask-closable=\"false\" :width=\"copySettings.width\">"+
		    "<p slot=\"header\" style=\"color:#f60;text-align:center\">"+
		        "<span>复制文件</span>"+
		    "</p>"+
			"<p style=\"text-align:left;\">"+
				"<div>"+
					"<table>"+
						"<tr>"+
							"<td>正在处理:</td>"+
							"<td>{{operationData.multiCount+operationData.opCount}}</td>"+
						"</tr>"+
						"<tr>"+
							"<td width=\"30%\">复制路径:</td>"+
							"<td>{{operationData.nowSrcPath}}</td>"+
						"</tr>"+
						"<tr>"+
							"<td>目标位置:</td>"+
							"<td>{{operationData.nowDstPath}}</td>"+
						"</tr>"+
						"<tr v-show=\"copyError.IsError\">"+
							"<td>出现错误:</td>"+
							"<td>{{copyError.Error}}</td>"+
						"</tr>"+
					"</table>"+
				"</div>"+
		    "</p>"+
		    "<div slot=\"footer\">"+
		       	"<div v-show=\"copyError.IsError\">"+
					"<i-button long :style=\"{margin: '5px 2px'}\" @click=\"doIgnore\">跳过</i-button>"+
					"<i-button long :style=\"{margin: '5px 2px'}\" v-show=\"copyError.IsExist\" @click=\"doReplace\">覆盖</i-button>"+
					"<div style=\"margin:8px 0px;color:#57a3f3; \">"+
						"<Checkbox v-model=\"copySettings.ignore\">自动跳过出错文件</Checkbox>"+
						"<Checkbox v-model=\"copySettings.replace\" v-show=\"copyError.IsExist\">自动覆盖重复文件</Checkbox>"+
					"</div>"+
		       	"</div>"+
				"<i-button type=\"error\" long :style=\"{margin: '5px 2px'}\" @click=\"doStop\">终止</i-button>"+
		    "</div>"+
		"</Modal>",
	created: function( ){

	},
	methods: {
		doCopy: function( ){
			var self = this;
			if( this.operationData.srcPaths && this.operationData.srcPaths.length > 0){
				var tempSrc = this.operationData.srcPaths[this.operationData.srcPaths.length-1];
				var tempDst = this.operationData.destPath+"/"+tempSrc.Path.getName( ).parsePath( );
				$fsApi.CopyAsync(tempSrc.Path, tempDst, this.copySettings.replace, this.copySettings.ignore).then(function( data ){
					if( self.operationData.srcPaths.length > 1 ){
						self.operationData.srcPaths = self.operationData.srcPaths.slice(0, self.operationData.srcPaths.length-1);
					}else{
						self.operationData.srcPaths = [];
					}
					self.operationData.token = data;
				}).catch(function( err ){
					self.$Message.error(err.toString( ));
					self.$emit("on-error");
				});
			}
		},
		doRefreshPs: function( ){
			if( !this.operationData.token || this.operationData.token == "" ){
				return;
			}
			var self = this;
			$fsApi.BatchOperationTokenInfo(this.operationData.token).then(function(data){
				/*
					{
						"CountIndex":7,
						"ErrorString":"",
						"Src":"/files/Mount01/glibc-ports-2.15.tar.gz",
						"Dst":"/files/.cache/Mount01/glibc-ports-2.15.tar.gz",
						"IsSrcExist":false,
						"IsDstExist":false,
						"IsReplace":false,
						"IsReplaceAll":false,
						"IsIgnore":false,
						"IsIgnoreAll":false,
						"IsComplete":false,
						"IsDiscontinue":false
					}
				 */
				data = JSON.parse(data);
				// console.log( data )
				if( data.CountIndex > 0 ){
					self.operationData.opCount = data.CountIndex;
				}
				if(data.IsComplete){
					self.operationData.token = "";
					if( data.IsDiscontinue ){
						self.showDailog = false;
						self.$Message.error("复制已终止")
						self.$emit("on-stop");
					}else{
						if( !self.operationData.srcPaths || self.operationData.srcPaths.length == 0 ){
							self.showDailog = false;
							if( data.ErrorString && data.ErrorString.length > 0 ){
								self.$Message.error(data.ErrorString);
							}else{
								self.$Message.success("复制完成");
							}
							self.$emit("on-end");
						}else{
							self.operationData.multiCount += self.operationData.opCount;
							self.doCopy( );
						}
					}
					return;
				}
				self.operationData.nowSrcPath = data.Src;
				self.operationData.nowDstPath = data.Dst;
				// 
				self.copyError.IsError = (data.ErrorString&&data.ErrorString.length>0)?true:false;
				self.copyError.Error = self.parseError(data);
				self.copyError.IsExist = data.IsDstExist;
				setTimeout(function( ){
					self.doRefreshPs( );
				}, 100);

			}).catch(function( err ){
				self.$Message.error(err.toString( ));
			});
		},
		doStop: function( ){
			if( this.operationData.token && this.operationData.token.length > 0 ){
				var self = this;
				$fsApi.SetBatchOperationToken(this.operationData.token, this.operations.stop).then(function(data){
					self.operationData.opCount  = 0;
					self.operationData.multiCount  = 0;
					self.operationData.srcPaths = [];
					self.operationData.destPath = [];
				}).catch(function( err ){
					self.$Message.error(err.toString( ));
				});
			}
		},
		doIgnore: function( ){
			if( !this.operationData.token || this.operationData.token == "" ){
				return;
			}
			var self = this;
			$fsApi.SetBatchOperationToken(this.operationData.token, this.copySettings.ignore?this.operations.ignoreall:this.operations.ignore).then(function(data){
				// console.log(data);
			}).catch(function( err ){
				self.$Message.error(err.toString( ));
			});
		},
		doReplace: function( ){
			if( !this.operationData.token || this.operationData.token == "" ){
				return;
			}
			var self = this;
			$fsApi.SetBatchOperationToken(this.operationData.token, this.copySettings.replace?this.operations.replaceall:this.operations.replace).then(function(data){
				// console.log(data);
			}).catch(function( err ){
				self.$Message.error(err.toString( ));
			});
		},
		parseError: function( data ){
			if (data && data.ErrorString) {
				if( data.IsDstExist ){
					return "目标位置已存在: "+data.Dst
				}else if( !data.IsSrcExist ){
					return "源目录不存在: "+data.Src
				}else{
					return data.ErrorString;
				}
			}
			return "";
		}

	},
	watch: {
		"operationData.token": function(n, o){
			if( n && n!= "" ){
				this.doRefreshPs( );
			}
		},
		"showDailog": function(n, o){
			if( n ){
				this.copyError.IsError = false;
				this.copyError.Error   = "";
				this.copyError.IsExist = false;
				this.copySettings.ignore = false;
				this.copySettings.replace = false;
				this.operationData.opCount  = 0;
				this.operationData.multiCount  = 0;
				this.operationData.showDailog = this.showDailog;
				this.operationData.srcPaths = this.srcPaths;
				this.operationData.destPath = this.destPath;
				this.doCopy( );
			}
		},
	}

});
// 文件拷贝
Vue.component("fs-movefile", {
	props: ["show-dailog", "src-paths", "dest-path", "move-settings"],
	data: function( ){
		var self = this;
		return {
			operations:{
				stop: "discontinue",
				ignore: "ignore",
				ignoreall: "ignoreall",
				replace: "replace",
				replaceall: "replaceall",
			},
			moveSettings: {
				ignore: false,
				replace: false,
				width: 450,
			},
			operationData: {
				multiCount: 0,
				opCount: 0,
				token: "",
				srcPaths: [],
				destPath: "",
				nowSrcPath: "",
				nowDstPath: "",
			},
			moveError: {
				IsError: false,
				Error: "",
				IsExist: false,
			},
		};
	},
	template: ""+
		"<Modal v-model=\"showDailog\" :closable=\"false\" :mask-closable=\"false\" :width=\"moveSettings.width\">"+
		    "<p slot=\"header\" style=\"color:#f60;text-align:center\">"+
		        "<span>移动文件</span>"+
		    "</p>"+
			"<p style=\"text-align:left;\">"+
				"<div>"+
					"<table>"+
						"<tr>"+
							"<td>正在处理:</td>"+
							"<td>{{operationData.multiCount+operationData.opCount}}</td>"+
						"</tr>"+
						"<tr>"+
							"<td width=\"30%\">移动路径:</td>"+
							"<td>{{operationData.nowSrcPath}}</td>"+
						"</tr>"+
						"<tr>"+
							"<td>目标位置:</td>"+
							"<td>{{operationData.nowDstPath}}</td>"+
						"</tr>"+
						"<tr v-show=\"moveError.IsError\">"+
							"<td>出现错误:</td>"+
							"<td>{{moveError.Error}}</td>"+
						"</tr>"+
					"</table>"+
				"</div>"+
		    "</p>"+
		    "<div slot=\"footer\">"+
		       	"<div v-show=\"moveError.IsError\">"+
					"<i-button long :style=\"{margin: '5px 2px'}\" @click=\"doIgnore\">跳过</i-button>"+
					"<i-button long :style=\"{margin: '5px 2px'}\" v-show=\"moveError.IsExist\" @click=\"doReplace\">覆盖</i-button>"+
					"<div style=\"margin:8px 0px;color:#57a3f3; \">"+
						"<Checkbox v-model=\"moveSettings.ignore\">自动跳过出错文件</Checkbox>"+
						"<Checkbox v-model=\"moveSettings.replace\" v-show=\"moveError.IsExist\">自动覆盖重复文件</Checkbox>"+
					"</div>"+
		       	"</div>"+
				"<i-button type=\"error\" long :style=\"{margin: '5px 2px'}\" @click=\"doStop\">终止</i-button>"+
		    "</div>"+
		"</Modal>",
	created: function( ){

	},
	methods: {
		doMove: function( ){
			var self = this;
			if( this.operationData.srcPaths && this.operationData.srcPaths.length > 0){
				var tempSrc = this.operationData.srcPaths[this.operationData.srcPaths.length-1];
				var tempDst = this.operationData.destPath+"/"+tempSrc.Path.getName( ).parsePath( );
				$fsApi.MoveAsync(tempSrc.Path, tempDst, this.moveSettings.replace, this.moveSettings.ignore).then(function( data ){
					if( self.operationData.srcPaths.length > 1 ){
						self.operationData.srcPaths = self.operationData.srcPaths.slice(0, self.operationData.srcPaths.length-1);
					}else{
						self.operationData.srcPaths = [];
					}
					self.operationData.token = data;
				}).catch(function( err ){
					self.$Message.error(err.toString( ));
					self.$emit("on-error");
				});
			}
		},
		doRefreshPs: function( ){
			if( !this.operationData.token || this.operationData.token == "" ){
				return;
			}
			var self = this;
			$fsApi.BatchOperationTokenInfo(this.operationData.token).then(function(data){
				/*
					{
						"CountIndex":7,
						"ErrorString":"",
						"Src":"/files/Mount01/glibc-ports-2.15.tar.gz",
						"Dst":"/files/.cache/Mount01/glibc-ports-2.15.tar.gz",
						"IsSrcExist":false,
						"IsDstExist":false,
						"IsReplace":false,
						"IsReplaceAll":false,
						"IsIgnore":false,
						"IsIgnoreAll":false,
						"IsComplete":false,
						"IsDiscontinue":false
					}
				 */
				data = JSON.parse(data);
				// console.log( data )
				if( data.CountIndex > 0 ){
					self.operationData.opCount = data.CountIndex;
				}
				if(data.IsComplete){
					self.operationData.token = "";
					if( data.IsDiscontinue ){
						self.showDailog = false;
						self.$Message.error("移动已终止")
						self.$emit("on-stop");
					}else{
						if( !self.operationData.srcPaths || self.operationData.srcPaths.length == 0 ){
							self.showDailog = false;
							if( data.ErrorString && data.ErrorString.length > 0 ){
								self.$Message.error(data.ErrorString);
							}else{
								self.$Message.success("移动完成");
							}
							self.$emit("on-end");
						}else{
							self.operationData.multiCount += self.operationData.opCount;
							self.doMove( );
						}
					}
					return;
				}
				self.operationData.nowSrcPath = data.Src;
				self.operationData.nowDstPath = data.Dst;
				// 
				self.moveError.IsError = (data.ErrorString&&data.ErrorString.length>0)?true:false;
				self.moveError.Error = self.parseError(data);
				self.moveError.IsExist = data.IsDstExist;
				setTimeout(function( ){
					self.doRefreshPs( );
				}, 100);

			}).catch(function( err ){
				self.$Message.error(err.toString( ));
			});
		},
		doStop: function( ){
			if( this.operationData.token && this.operationData.token.length > 0 ){
				var self = this;
				$fsApi.SetBatchOperationToken(this.operationData.token, this.operations.stop).then(function(data){
					self.operationData.opCount  = 0;
					self.operationData.multiCount  = 0;
					self.operationData.srcPaths = [];
					self.operationData.destPath = [];
				}).catch(function( err ){
					self.$Message.error(err.toString( ));
				});
			}
		},
		doIgnore: function( ){
			if( !this.operationData.token || this.operationData.token == "" ){
				return;
			}
			var self = this;
			$fsApi.SetBatchOperationToken(this.operationData.token, this.moveSettings.ignore?this.operations.ignoreall:this.operations.ignore).then(function(data){
				// console.log(data);
			}).catch(function( err ){
				self.$Message.error(err.toString( ));
			});
		},
		doReplace: function( ){
			if( !this.operationData.token || this.operationData.token == "" ){
				return;
			}
			var self = this;
			$fsApi.SetBatchOperationToken(this.operationData.token, this.moveSettings.replace?this.operations.replaceall:this.operations.replace).then(function(data){
				// console.log(data);
			}).catch(function( err ){
				self.$Message.error(err.toString( ));
			});
		},
		parseError: function( data ){
			if (data && data.ErrorString) {
				if( data.IsDstExist ){
					return "目标位置已存在: "+data.Dst
				}else if( !data.IsSrcExist ){
					return "源目录不存在: "+data.Src
				}else{
					return data.ErrorString;
				}
			}
			return "";
		}

	},
	watch: {
		"operationData.token": function(n, o){
			if( n && n!= "" ){
				this.doRefreshPs( );
			}
		},
		"showDailog": function(n, o){
			if( n ){
				this.moveError.IsError = false;
				this.moveError.Error   = "";
				this.moveError.IsExist = false;
				this.moveSettings.ignore = false;
				this.moveSettings.replace = false;
				this.operationData.opCount  = 0;
				this.operationData.multiCount  = 0;
				this.operationData.showDailog = this.showDailog;
				this.operationData.srcPaths = this.srcPaths;
				this.operationData.destPath = this.destPath;
				this.doMove( );
			}
		},
	}

});
// 文件|文件夹选择
Vue.component("fs-selector", {
	props:["show-dailog", "settings", "start-path", "select-muti", "select-file", "select-dir"],
	data: function( ){
		var self = this;
		return{
			isSelectFile: false,
			isSelectDir: false,
			isSelectMuti: false,
			fsStatus:{
				loadedPath:"",
				fsLoading: true,
			},
			fsSettings: {
				rootname: "/",
				width: 680,
				height: 450
			},
			fsColumns: [
				{
					title: "#",	
					width: 60,
					render: function(h, params){
						return h('checkbox', {
							props: {
								value: params.row._checked==undefined?false:params.row._checked
							},
							on: {
								'on-change': function(val){
									self.fsStatus.fsLoading = true;
									if( !self.isSelectMuti ){
										for(var i=0; i<self.fsData.length; i++ ){
											if( i == params.index ){ continue; }
											self.fsData[i]['_checked'] = false;
										}
									}
									self.$set(self.fsData[params.index], '_checked', val);
									if(val){
										self.putSelect(params.row);
									}else{
										self.removeSelect(params.row);
									}
									self.fsStatus.fsLoading = false;
								}                                    
							}
						});
					}
                },
				{
                    title: '文件名称',
                    key: 'Path',
                    render:function (h, params){
                    	return h("fs-fileicon", {
                    		props:{
                    			node: params.row,
                    			isEditor:false
                    		},
                    		on:{
                    			click: self.doOpenDir
                    		}
                    	});
                    }
                },
                {
                    title: '修改时间',
                    key: 'CtTime',
                    width:180,
                    render:function (h, params){
                    	return h("span", $utils.long2Time(params.row.CtTime));
                    }
                }
			],
			fsData: [],
			selectedDates: [],
		}
	},
	template:"<Modal v-model=\"showDailog\" :closable=\"false\" :mask-closable=\"false\" :width=\"fsSettings.width\"  @on-ok=\"even_ok\" @on-cancel=\"even_cancel\">"+
			"<div style=\"height: 30px;line-height: 30px;padding-left:5px;\">"+
				"<fs-address v-if=\"fsStatus.loadedPath&&fsStatus.loadedPath.length>0\" :depth=\"4\" :rootname=\"fsSettings.rootname\" :path=\"fsStatus.loadedPath\" @click=\"goToPath\"></fs-address>"+
			"</div>"+
		    "<i-table :loading=\"fsStatus.fsLoading\" :columns=\"fsColumns\" :data=\"fsData\" :height=\"fsSettings.height\" @on-row-click=\"onRowClick\" @on-selection-change=\"onSelectionChange\" ></i-table>"+
	    "</Modal>",
	created: function( ){
		
	},
	methods: {
		init: function( ){
			if( this.setting ){
				this.fsSettings.width = this.settings.width?this.settings.width:this.fsSettings.width;
				this.fsSettings.height = this.settings.height?this.settings.height:this.fsSettings.height;
			}
			// 
			this.isSelectFile = this.selectFile?this.selectFile:false;
			this.isSelectDir  = this.selectDir?this.selectDir:false;
			this.isSelectMuti = this.selectMuti?eval(this.selectMuti):false;
			this.selectedDates = [];
			// 
			if( this.startPath&&this.startPath.length>0 ){
				this.goToPath(this.startPath);
			}
		},
		goToPath: function( path ){
			if( this.fsStatus.loadedPath != path ){
				this.fsStatus.fsLoading = true;
				this.fsStatus.loadedPath = path;
			}
		},
		doOpenDir: function( node ){
			if( !node.IsFile ){
				this.goToPath( node.Path )
			}
		},
		even_ok: function( ){
			this.fsStatus.fsLoading = true;
			if( this.fsData ){
				for(var i=0; i<this.fsData.length; i++){
					this.fsData[i]["_checked"] = false;
				}
			}
			this.fsStatus.fsLoading = false;
			this.$emit("on-select", (this.selectedDates&&this.selectedDates.length>0)?this.selectedDates:(this.isSelectDir?[{
				"Path": this.fsStatus.loadedPath,
				"IsFile": false
			}]:[]));
		},
		even_cancel: function( ){
			this.$emit("on-cancel");
		},
		onRowClick: function(row, index){
          for (var i = 0; i < this.fsData.length; i++) {
            if( this.fsData[i]._checked && index != i ){
              this.$set( this.fsData[i], "_checked", false);
            }
          }
          this.selectedDates = [row];
          this.$set( this.fsData[index], "_checked", true);
        },
		putSelect: function( row ){
			if(!this.isSelectMuti){ 
				this.selectedDates = [row];
			}else{
				for(var i=0; i<this.selectedDates.length; i++){
					if(this.selectedDates[i].Path == row.Path){
						return;
					}
				}
				this.selectedDates.push( row );
			}
		},
		removeSelect: function( row ){
			for(var i = this.selectedDates.length - 1; i >= 0; i--){
				if( this.selectedDates[i].Path == row.Path ){
					this.selectedDates.remove( i );
				}
			}
		},
		onSelectionChange: function( selection, row ){

		},
	},
	watch: {
		showDailog: function(n, o){
			if(n){
				this.init( );
			}else{
				this.fsStatus.loadedPath = "";
			}
		},
		'fsStatus.loadedPath': function(n, o){
			if( !n ){ return; }
			var self = this;
			$fsApi.List(n).then(function( data ){
				self.fsData = [];
				self.selectedDates = [];
				for(var i=0; i<data.length; i++){
					if( (self.isSelectFile && data[i].IsFile) || (self.isSelectDir && !data[i].IsFile) ){
						data[i]["_checked"] = false;
						self.fsData.push( data[i] );
					}
				}
				self.fsStatus.fsLoading = false;
			}).catch(function(err){
				self.fsStatus.fsLoading = false;
				self.$Message.error(err.toString( ));
			});
		},
	}
});
// 地址栏
Vue.component("fs-address",{
	props:["path", "root", "rootname", "depth"],
	data:function(){
		return{
			paths:[],
			max:6,
			showrootname:""
		}
	},
	template:"<breadcrumb separator=\">\">"
				+"<breadcrumb-item>"
					+"<a href=\"###\" @click=\"address_GoToRoot()\">{{showrootname}}</a>"
				+"</breadcrumb-item>"
				+"<breadcrumb-item v-for=\"(item,index) in paths\" v-show='item&&(index>=paths.length-2 || index<=max)'>"
					// +"<a href=\"###\" s-show='index>=paths.length-2 || index<=max'>...</a>"
					+"<a href=\"###\" @click=\"address_GoToPath(item,index)\">{{item}}</a>"
				+"</breadcrumb-item>"
			+"</breadcrumb>",
	created:function(){	
		this.max = this.depth?(this.depth-2>0?this.depth-2:2):this.max;
		this.buildPaths( );
    },
	methods:{
		buildPaths:function(){
			this.showrootname = this.rootname?this.rootname:"/"
			this.paths = this.path.split("/");		
		},
		address_GoToRoot:function( ){
			this.$emit("click", this.root?this.root:"/");
		},
		address_GoToPath:function(item, index){
			var path = "";
			for( i=0; i<=index; i++ ){
				if( this.paths[i] ){
					path += "/"+this.paths[i];
				}
			}
			this.$emit("click", path);
		}
	},
	watch:{
		path:function(v1,v2){			
			this.buildPaths();
		}
	}
});
// 文件图标
Vue.component("fs-fileicon",{
	props:["node", "isEditor"],
	data:function(){
		return{
			iconTypes:["ai","avi","bmp","catdrawing","catpart","catproduct","cdr","csv","doc","docx","dps","dpt","dwg","eio","eml","et","ett","exb","exe","file","flv","fold","gif","htm","html","jpeg","jpg","mht","mhtml","mid","mp3","mp4","mpeg","msg","odp","ods","odt","pdf","png","pps","ppt","pptx","prt","psd","rar","rm","rmvb","rtf","sldprt","swf","tif","txt","url","wav","wma","wmv","wps","wpt","xls","xlsx","zip"],
			icon:"",
			filename:""
		}
	},
	template:"<div style='width:100%;'>"+
				"<div style='overflow: hidden;text-overflow:ellipsis;white-space:nowrap;'>"+
					"<img :src=\"icon\" style='vertical-align:middle;margin-right:10px;width:32px;height:32px'>"+
					"<i-input v-if=\"isEditor\" v-model='filename' style='width:200px' @on-enter=\"doSave\"></i-input>"+
					"<span v-else style='cursor:pointer' @click=\"doClick\" :title='filename'>{{filename}}</span>"+
				"</div>"+
			"</div>",		
	methods:{		
		doClick:function(){
			this.$emit("click",this.node);
		},		
		doSave:function(){	
			if( !this.filename ){
				return;
			}			
			if( this.isEditor === true ){ 
				this.isEditor = false;						
				this.$emit("doRename", this.node.Path, this.filename);	
			}					
		},
		getFileIcon:function(path){
			if(!this.node.IsFile){
				return "/img/file_icons/folder.png";				
			}
			return $utils.iconUrl(path);			
		},
		initvalue:function(){
			this.filename = this.node.Path.getName( );
			this.icon = this.getFileIcon( this.filename );
		}
	},	
	created:function(){
		var that = this;
		this.initvalue( );
		this.$nextTick(function (){
			window.addEventListener("keydown",function(e){
				if(!e){
					e = window.event;
				 }
				 if((e.keyCode || e.which) == 13){
					that.doSave();
				 }
			});

			window.addEventListener("mousedown",function(e){
				var dom = e.target;
				if(dom == that.$el.querySelector("input") ){					
				}else{
					that.doSave();
				}
			});
		});

    },
    watch:{
		node:{
			deep:true,
			handler:function(newval,oldval){									
				this.initvalue();
			}
		}
	}
});

