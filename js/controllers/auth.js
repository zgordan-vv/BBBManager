;(function(){
'use strict';

function registerCtrl($rootScope, $http, $location, REST) {

	var rc = this;

	rc.show=false;

	rc.updateCheckName = function() {
		$http.get('api/regexp?word='+rc.regname+'&type=name').then(
			function(response){
				rc.nameValid = response.data;
			}, function(){
				rc.nameValid = "Server error";
			}
		);
		$http.get('api/nameUniq?name='+rc.regname).then(
			function(response){
				rc.nameUniq = response.data;
			}, function(){
				rc.nameUniq = "Server error";
			}
		);
	};
	
	rc.updateCheckFName = function() {
		$http.get('api/regexp?word='+rc.fullname+'&type=desc').then(
			function(response){
				rc.fNameValid = response.data;
			}, function(){
				rc.fNameValid = "Server error";
			}
		);
	};
	
	rc.submitregister = function() {
		var post = $.param({
			registerdata: JSON.stringify({
				name: rc.regname,
				fullname: rc.fullname,
				pwd: rc.regpwd,
				pwdconf: rc.pwdconf,
			})
		});
		REST.post("api/register", post).then(function(response){
			if (response != "200") {rc.msg = "Server error"; rc.show=true;} else {
				REST.checkAuth().then(function(){
					$rootScope.authorized = 'user';
					$location.url("/");
				}, function(){
					$rootScope.authorized = 'guest';
					rc.msg="Server error";
					rc.show="true";
				});
			}
		}, function(){
			rc.msg = "Server error"; rc.show=true;
		})
	};
}

function authCtrl($http, $window, $location, $rootScope, REST) {

	var ac = this;

	ac.submitauth = function() {
		var post = $.param({
			login: JSON.stringify({
				login: ac.name,
				password: ac.pwd,
			})
		});
		REST.post("/api/login", post).then(function(response){
			if (response == "403") {
				ac.warn = true;
			} else {
				REST.checkAuth().then(function(response){
					$rootScope.authorized = response;
					$location.url("/");
				}, function(){
					$rootScope.authorized = 'guest';
					$location.url("/");
				});
			}
		}, function(response){
			$rootScope.authorized = 'guest';
			$location.url("/");
		});
	};

	ac.oauthLogin = function(provider) {
		REST.get("/api/get"+provider+"LoginURL").then(function(response){
			console.log(response);
//			$window.location = response;
		},function(error){
			alert(error);
		});
	};
}

function quitCtrl($http, $rootScope) {
	$rootScope.authorized = 'guest';
	$http.get("/api/quit");
}

angular.module("AuthCs", [])
.controller('registerCtrl', ['$rootScope', '$http', '$location', 'REST', registerCtrl])
.controller('authCtrl', ['$http', '$window', '$location', '$rootScope', 'REST', authCtrl])
.controller('quitCtrl', ['$http', '$rootScope', quitCtrl])
})();
