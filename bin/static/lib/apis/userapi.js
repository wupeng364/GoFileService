"use strict";
/**
 *@description UserApi
 *@author	wupeng364@outlook.com
*/
;(function (root, factory) {
	if (typeof exports === "object") {
		module.exports = exports = factory();
	}else {
		// Global (browser)
		root.$userApi = factory();
	}
}(this, function ( ){
	var _Api = { };

	/**
	 * 登陆
	 * return {UserId, AccessKey, SecretKey}
	 */
	_Api.login = function(user, pwd){
		return new Promise(function(resolve, reject){
			$fhttp.AjaxRequest({
				uri: "/userapi/checkpwd",
				datas: {
					"userId": user,
					"pwd": pwd,
				},				
			}).do(function(xhr, opt){
				if(xhr.readyState === 4){
					var res = $fhttp.apiResponseFormat(xhr);
					if( res.Code === 200 ){
						resolve( res.Data );
					}else{
						reject( res.Data );
					}
				}
			});
		});
	};
	return _Api;
}));
