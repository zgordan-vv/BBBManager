;(function(){
'use strict';

function CheckInst($rootScope, REST) {
	REST.checkInst().then(function(data){$rootScope.installed = data;}, function(){$rootScope.installed = 'false';});
	REST.checkAuth().then(function(data){$rootScope.authorized = data;}, function(){$rootScope.authorized = 'guest';});
};

function installCtrl($http, $window, REST) {
	var ic = this;

	ic.name = "";
	ic.fullname = "";
	ic.dbprefix = "";
	ic.dname = "";
	
	ic.updateCheckName = function() {
		$http.get('api/regexp?word='+ic.name+'&type=name').then(
			function(response){
				ic.nameValid = response.data;
			}, function(){
				ic.nameValid = "Server error";
			}
		);
		$http.get('api/nameUniq?name='+ic.name).then(
			function(response){
				ic.nameUniq = response.data;
			}, function(){
				ic.nameUniq = "Server error";
			}
		);
	};
	
	ic.updateCheckFName = function() {
		$http.get('api/regexp?word='+ic.fullname+'&type=desc').then(
			function(response){
				ic.fNameValid = response.data;
			}, function(){
				ic.fNameValid = "Server error";
			}
		);
	};
	
	ic.updateCheckDBPrefix = function() {
		$http.get('api/regexp?word='+ic.dbprefix+'&type=charnum').then(
			function(response){
				ic.dbPrefixValid = response.data;
			}, function(){
				ic.dbPrefixValid = "Server error";
			}
		);
	};
	
	ic.updateCheckDName = function() {
		$http.get('api/regexp?word='+ic.dname+'&type=domain').then(
			function(response){
				ic.dNameValid = response.data;
			}, function(){
				ic.dNameValid = "Server error";
			}
		);
	};

	ic.submitinstall = function() {
		var post = $.param({
			installdata: JSON.stringify({
				name: ic.name,
				fullname: ic.fullname,
				pwd: ic.pwd,
				dbprefix: ic.dbprefix,
				domainname: ic.dname,
			})
		});
		REST.post("/api/install", post).then(function(){
			$window.location.reload();
		}, function(){
			ic.feedback = "Server error";
		});
	};
}

angular.module("InstallCs", [])
.controller('CheckInst', ['$rootScope', 'REST', CheckInst])
.controller('installCtrl', ['$http', '$window', 'REST', installCtrl])
})();
