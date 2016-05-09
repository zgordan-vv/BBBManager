;(function(){
'use strict';

function NavLogged($http, $rootScope, CheckAuth) {
	var nlc = this;
	var rs = $rootScope;
	rs.$watch('authorized', function(){updatenav(rs.authorized||'guest');})
	function updatenav(str) {
		$http.get('resources/navs/'+str+'.json').then(function(response){
			nlc.items = response.data;
		}, function(){
			nlc.items = [];
		});
	};
};

angular.module("CommonCs", []).controller('NavLogged', ['$http', '$rootScope', NavLogged]);
})();
