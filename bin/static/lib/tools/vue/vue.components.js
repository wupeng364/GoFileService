// Copyright (C) 2020 WuPeng <wupeng364@outlook.com>.
// Use of this source code is governed by an MIT-style.
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
// IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// VUE组件 - 常用,带有业务属性的
;(function(factory){
	try {
		factory( );
	} catch (e) {
		console.error(e);
	}
})(function( ){
	// 文件拷贝
	Vue.component("fs-copyfile", {
		props: ["show-dailog", "src-paths", "dest-path", "copy-settings"],
		data: function( ){
			let _ = this;
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
				let _ = this;
				if( this.operationData.srcPaths && this.operationData.srcPaths.length > 0){
					let tempSrc = this.operationData.srcPaths[this.operationData.srcPaths.length-1];
					let tempDst = this.operationData.destPath+"/"+tempSrc.Path.getName( ).parsePath( );
					$fsApi.CopyAsync(tempSrc.Path, tempDst, this.copySettings.replace, this.copySettings.ignore).then(function( data ){
						if( _.operationData.srcPaths.length > 1 ){
							_.operationData.srcPaths = _.operationData.srcPaths.slice(0, _.operationData.srcPaths.length-1);
						}else{
							_.operationData.srcPaths = [];
						}
						_.operationData.token = data;
					}).catch(function( err ){
						_.$Message.error(err.toString( ));
						_.$emit("on-error");
					});
				}
			},
			doRefreshPs: function( ){
				if( !this.operationData.token || this.operationData.token == "" ){
					return;
				}
				let _ = this;
				$fsApi.AsyncExecToken("CopyFile", this.operationData.token).then(function(data){
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
						_.operationData.opCount = data.CountIndex;
					}
					if(data.IsComplete){
						_.operationData.token = "";
						if( data.IsDiscontinue ){
							_.showDailog = false;
							_.$Message.error("复制已终止")
							_.$emit("on-stop");
						}else{
							if( !_.operationData.srcPaths || _.operationData.srcPaths.length == 0 ){
								_.showDailog = false;
								if( data.ErrorString && data.ErrorString.length > 0 ){
									_.$Message.error(data.ErrorString);
								}else{
									_.$Message.success("复制完成");
								}
								_.$emit("on-end");
							}else{
								_.operationData.multiCount += _.operationData.opCount;
								_.doCopy( );
							}
						}
						return;
					}
					_.operationData.nowSrcPath = data.Src;
					_.operationData.nowDstPath = data.Dst;
					// 
					_.copyError.IsError = (data.ErrorString&&data.ErrorString.length>0)?true:false;
					_.copyError.Error = _.parseError(data);
					_.copyError.IsExist = data.IsDstExist;
					setTimeout(function( ){
						_.doRefreshPs( );
					}, 100);

				}).catch(function( err ){
					_.$Message.error(err.toString( ));
				});
			},
			doStop: function( ){
				if( this.operationData.token && this.operationData.token.length > 0 ){
					let _ = this;
					$fsApi.AsyncExecToken("CopyFile", this.operationData.token, {operation: this.operations.stop}).then(function(data){
						_.operationData.opCount  = 0;
						_.operationData.multiCount  = 0;
						_.operationData.srcPaths = [];
						_.operationData.destPath = [];
					}).catch(function( err ){
						_.$Message.error(err.toString( ));
					});
				}
			},
			doIgnore: function( ){
				if( !this.operationData.token || this.operationData.token == "" ){
					return;
				}
				let _ = this;
				$fsApi.AsyncExecToken("CopyFile", this.operationData.token, {operation: this.copySettings.ignore?this.operations.ignoreall:this.operations.ignore}).then(function(data){
					// console.log(data);
				}).catch(function( err ){
					_.$Message.error(err.toString( ));
				});
			},
			doReplace: function( ){
				if( !this.operationData.token || this.operationData.token == "" ){
					return;
				}
				let _ = this;
				$fsApi.AsyncExecToken("CopyFile", this.operationData.token, {operation: this.copySettings.replace?this.operations.replaceall:this.operations.replace}).then(function(data){
					// console.log(data);
				}).catch(function( err ){
					_.$Message.error(err.toString( ));
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
			let _ = this;
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
				let _ = this;
				if( this.operationData.srcPaths && this.operationData.srcPaths.length > 0){
					let tempSrc = this.operationData.srcPaths[this.operationData.srcPaths.length-1];
					let tempDst = this.operationData.destPath+"/"+tempSrc.Path.getName( ).parsePath( );
					$fsApi.MoveAsync(tempSrc.Path, tempDst, this.moveSettings.replace, this.moveSettings.ignore).then(function( data ){
						if( _.operationData.srcPaths.length > 1 ){
							_.operationData.srcPaths = _.operationData.srcPaths.slice(0, _.operationData.srcPaths.length-1);
						}else{
							_.operationData.srcPaths = [];
						}
						_.operationData.token = data;
					}).catch(function( err ){
						_.$Message.error(err.toString( ));
						_.$emit("on-error");
					});
				}
			},
			doRefreshPs: function( ){
				if( !this.operationData.token || this.operationData.token == "" ){
					return;
				}
				let _ = this;
				$fsApi.AsyncExecToken("MoveFile",this.operationData.token).then(function(data){
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
						_.operationData.opCount = data.CountIndex;
					}
					if(data.IsComplete){
						_.operationData.token = "";
						if( data.IsDiscontinue ){
							_.showDailog = false;
							_.$Message.error("移动已终止")
							_.$emit("on-stop");
						}else{
							if( !_.operationData.srcPaths || _.operationData.srcPaths.length == 0 ){
								_.showDailog = false;
								if( data.ErrorString && data.ErrorString.length > 0 ){
									_.$Message.error(data.ErrorString);
								}else{
									_.$Message.success("移动完成");
								}
								_.$emit("on-end");
							}else{
								_.operationData.multiCount += _.operationData.opCount;
								_.doMove( );
							}
						}
						return;
					}
					_.operationData.nowSrcPath = data.Src;
					_.operationData.nowDstPath = data.Dst;
					// 
					_.moveError.IsError = (data.ErrorString&&data.ErrorString.length>0)?true:false;
					_.moveError.Error = _.parseError(data);
					_.moveError.IsExist = data.IsDstExist;
					setTimeout(function( ){
						_.doRefreshPs( );
					}, 100);

				}).catch(function( err ){
					_.$Message.error(err.toString( ));
				});
			},
			doStop: function( ){
				if( this.operationData.token && this.operationData.token.length > 0 ){
					let _ = this;
					$fsApi.AsyncExecToken("MoveFile", this.operationData.token, {operation: this.operations.stop}).then(function(data){
						_.operationData.opCount  = 0;
						_.operationData.multiCount  = 0;
						_.operationData.srcPaths = [];
						_.operationData.destPath = [];
					}).catch(function( err ){
						_.$Message.error(err.toString( ));
					});
				}
			},
			doIgnore: function( ){
				if( !this.operationData.token || this.operationData.token == "" ){
					return;
				}
				let _ = this;
				$fsApi.AsyncExecToken("MoveFile", this.operationData.token, {operation: this.moveSettings.ignore?this.operations.ignoreall:this.operations.ignore}).then(function(data){
					// console.log(data);
				}).catch(function( err ){
					_.$Message.error(err.toString( ));
				});
			},
			doReplace: function( ){
				if( !this.operationData.token || this.operationData.token == "" ){
					return;
				}
				let _ = this;
				$fsApi.AsyncExecToken("MoveFile", this.operationData.token, {operation: this.moveSettings.replace?this.operations.replaceall:this.operations.replace}).then(function(data){
					// console.log(data);
				}).catch(function( err ){
					_.$Message.error(err.toString( ));
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
			let _ = this;
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
										_.fsStatus.fsLoading = true;
										if( !_.isSelectMuti ){
											for(let i=0; i<_.fsData.length; i++ ){
												if( i == params.index ){ continue; }
												_.fsData[i]['_checked'] = false;
											}
										}
										_.$set(_.fsData[params.index], '_checked', val);
										if(val){
											_.putSelect(params.row);
										}else{
											_.removeSelect(params.row);
										}
										_.fsStatus.fsLoading = false;
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
									click: _.doOpenDir
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
		template:"<Modal v-model=\"showDailog\" :title=\"selectFile?'选择文件':'选择目录'\" :width=\"fsSettings.width\" @on-ok=\"onOk\" @on-cancel=\"onCancel\">"+
				"	<div style=\"height: 30px;line-height: 30px;padding-left:5px;\">"+
				"		<fs-address v-if=\"fsStatus.loadedPath&&fsStatus.loadedPath.length>0\" :depth=\"4\" :rootname=\"fsSettings.rootname\" :path=\"fsStatus.loadedPath\" @click=\"goToPath\"></fs-address>"+
				"	</div>"+
				"	<i-table :loading=\"fsStatus.fsLoading\" :columns=\"fsColumns\" :data=\"fsData\" :height=\"fsSettings.height\" @on-row-click=\"onRowClick\" @on-selection-change=\"onSelectionChange\" ></i-table>"+
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
			onOk: function( ){
				this.fsStatus.fsLoading = true;
				if( this.fsData ){
					for(let i=0; i<this.fsData.length; i++){
						this.fsData[i]["_checked"] = false;
					}
				}
				this.fsStatus.fsLoading = false;
				this.$emit("on-select", (this.selectedDates&&this.selectedDates.length>0)?this.selectedDates:(this.isSelectDir?[{
					"Path": this.fsStatus.loadedPath,
					"IsFile": false
				}]:[]));
			},
			onCancel: function( ){
				this.$emit("on-cancel");
			},
			onRowClick: function(row, index){
			for (let i = 0; i < this.fsData.length; i++) {
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
					for(let i=0; i<this.selectedDates.length; i++){
						if(this.selectedDates[i].Path == row.Path){
							return;
						}
					}
					this.selectedDates.push( row );
				}
			},
			removeSelect: function( row ){
				for(let i = this.selectedDates.length - 1; i >= 0; i--){
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
				let _ = this;
				$fsApi.List(n).then(function( data ){
					data = JSON.parse(data);
					_.fsData = [];
					_.selectedDates = [];
					for(let i=0; i<data.length; i++){
						if( (_.isSelectFile && data[i].IsFile) || (_.isSelectDir && !data[i].IsFile) ){
							data[i]["_checked"] = false;
							_.fsData.push( data[i] );
						}
					}
					_.fsStatus.fsLoading = false;
				}).catch(function(err){
					_.fsStatus.fsLoading = false;
					_.$Message.error(err.toString( ));
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
				let path = "";
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
		template:"<div style='width:100%;' @click='onBoxClick'>"+
					"<div style='overflow: hidden;text-overflow:ellipsis;white-space:nowrap;'>"+
						"<img :src=\"icon\" style='vertical-align:middle;margin-right:10px;width:32px;height:32px'>"+
						"<i-input v-if=\"isEditor\" v-model='filename' style='width: calc(100% - 45px);' @on-blur='doSave' @on-enter='doSave'></i-input>"+
						"<span v-else style='cursor:pointer' @click=\"onNameClick\" :title='filename'>{{filename}}</span>"+
					"</div>"+
				"</div>",		
		methods:{	
			onBoxClick: function(e){
				if( e.srcElement.tagName.toUpperCase() == 'INPUT'){
					e.stopPropagation();
				}
			},	
			onNameClick: function(e){
				e.stopPropagation();
				this.$emit("click",this.node, e);
			},		
			doSave:function(){	
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
			let that = this;
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
					let dom = e.target;
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
	// 文件上传
	Vue.component("fs-upload",{
		props:["show-drawer", "parent", "drag-ref"],
		data:function(){
			return{
				maxuploading: 5,    // 最大正在上传的个数
				countuploading: 0,  // 正在上传的个数
				dindex: 0,          // 当前数据下标
				files: [],          // 文件
				queueend: true,     // 队列可用为空
			}
		},
		template:"<Drawer title=\"上传文件\" width=\"450px\" v-model=\"showDrawer\" @on-close=\"$emit('on-close')\">" + 
		"		  		<div class=\"ivu-upload\">" + 
		"		  			<div class=\"ivu-upload ivu-upload-drag\" ref=\"uploadDrag\" @click=\"doSelectFiels\">" + 
		"		  				<input ref=\"upload_selector_file\" type=\"file\" multiple=\"multiple\" class=\"ivu-upload-input\">" + 
		"		  				<div style=\"padding: 20px 0px;\">" + 
		"		  				<i class=\"ivu-icon ivu-icon-ios-cloud-upload\" style=\"font-size: 52px; color: rgb(51, 153, 255);\"></i>" + 
		"		  				<p>点击或拖拽到此处</p>" + 
		"		  				</div>" + 
		"		  			</div> " + 
		"		  			<ul class=\"ivu-upload-list\">" + 
		"		  				<li v-for=\"temp in files\" v-if=\"!temp._upload.removed\" class=\"ivu-upload-list-file\">" + 
		"		  				<span>{{temp.name}}</span>" + 
		"		  				<i class=\"ivu-icon ivu-icon-ios-close ivu-upload-list-remove\" @click=\"removeTask(temp._upload.index)\"></i>" + 
		"		  				<div class=\"ivu-progress ivu-progress-normal ivu-progress-show-info\">" + 
		"		  					<div class=\"ivu-progress-outer\">" + 
		"		  					<div class=\"ivu-progress-inner\">" + 
		"		  						<div v-if=\"!temp._upload.ps||temp._upload.ps<100\" class=\"ivu-progress-bg\" :style=\"{width: (temp._upload.ps?temp._upload.ps:0)+'%', height: '2px'}\"></div>" + 
		"		  						<div v-else class=\"ivu-progress-success-bg\" style=\"width: 100%; height: 2px;\"></div>" + 
		"		  					</div>" + 
		"		  					</div> " + 
		"		  					<span class=\"ivu-progress-text\">" + 
		"		  					<span :style=\"{color: temp._upload.err?'#b42525':'#515a6e'}\" class=\"ivu-progress-text-inner\">{{temp._upload.err?temp._upload.err:temp._upload.ps+'%'}}</span>" + 
		"		  					</span>" + 
		"		  				</div>" + 
		"		  				</li>" + 
		"		  			</ul>" + 
		"		  		</div>" + 
		"		  </Drawer>",		
		methods:{		
			// 上传-触发选择文件
			doSelectFiels: function( ev ){
				let _ = this; let emited = false;
				$utils.addEvent(this.$refs.upload_selector_file, "change", function( ev_data ){
					if( ev_data.target.files ){
						for(let i = 0; i < ev_data.target.files.length; i++){
							let fs = ev_data.target.files[i];
							fs._upload = {
								base: _.parent,
								index: _.files.length,
								updater: false,
								ps: 0,
							};
							if(!emited){ emited = true; _.$emit('on-start'); }
							_.files.push( fs );
							_.doStartUpload( );
						}
					}
					_.$refs.upload_selector_file.value = "";
				}, {once: true});
				$utils.triggerMouseEvent(this.$refs.upload_selector_file, "click");
			},
			// 上传-触发上传动作
			doStartUpload: function( ){
				if( this.files && this.files.length > 0 ){
					if( this.countuploading >= this.maxuploading){
						return;
					}
					if(this.dindex >= this.files.length){
						this.queueend = true; return;
					}else if(this.queueend){
						this.queueend = false;
					}
					this.countuploading ++;
					let file = this.files[this.dindex++];
					if( file._upload.removed ){
						this.countuploading --;
						this.$nextTick(this.doStartUpload);
						return;
					}
					// file._upload.index = this.dindex-1;
					let _ = this;
					let opts = {
						form: { },
						header: { },
						progress: function( e ){
							// 数据源为数组, 需要直接设置数组
							file._upload.ps = Math.round((e.loaded/e.total)*1000)/10;
							_.$set(_.files, file._upload.index, file);
						},
						error: function( e ){
							file._upload.err = e?e.toString():"上传失败";
							_.$set(_.files, file._upload.index, file);
						},
						abort: function( e ){
							file._upload.err = "上传取消";
							_.$set(_.files, file._upload.index, file);
						},
						loadstart: function( e ){ },
						loadend: function( e ){
							_.countuploading--;
							file._upload.ended = true;
							_.$nextTick(function( ){
								_.doStartUpload( );
								//_.removeTask(file._upload.index);
							});
						}
					};
					// 预备开始
					file._upload.started = true;
					$fsApi.GetUploadUrl(file._upload.base+"/"+file.name).then(function(url){
						file._upload.updater = $utils.uploadByFormData(url, file, opts);
						file._upload.updater.start( );
					}).catch(function(err){
						opts.error(err);
						opts.loadend();
					});
				}
			},
			// 上传 - 移除任务
			removeTask: function( index ){
				let file = this.files[index];
				if( file ){
					// 正在传输
					if( !file._upload.ended && file._upload.started ){
						file._upload.updater.abort( );
					}
					file._upload.removed = true;
					this.$set(this.files, index, file);
				}
			},
			// 上传 - 监听拖拽
			doListenDrag: function( key ){
				let _ = this;
				this.$nextTick(function( ){
					let obj = undefined;
					if( key ){
						if(_.$refs[key]){
							obj = _.$refs[key];				
						}else if(_.$parent && _.$parent.$refs){
							obj = _.$parent.$refs[key];
						}
					}
					if(!obj){ return; }
					obj.ondrop = function( evn ){
						evn.preventDefault( );
						let emited = false;
						let fileList = evn.dataTransfer.files;
						for(let i = 0; i < fileList.length; i++){
							let fs = fileList[i];
							if(!fs.type && fs.size == 0){
								continue
							}
							fs._upload = {
								base: _.parent,
								index: _.files.length,
								updater: false,
								ps: 0,
							};
							if(!emited){ emited = true; _.$emit('on-start'); }
							_.files.push( fs );
							_.doStartUpload( );
						}
					}
				});
			},
		},	
		created:function(){
			this.doListenDrag('uploadDrag');
			this.doListenDrag(this.dragRef);
		},
		watch:{
			countuploading: function(n, o){
				if(n <= 0 && this.queueend){
					this.$emit('on-end');
				}
			}
		}
	});
});
