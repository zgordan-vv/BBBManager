;(function(){
'use strict';

function REST($http, $q) {
	var instance = this;

	instance.checkInst = function(){
		var defer = $q.defer();
		$http.get('api/checkInstall').then(function(response){
			defer.resolve(response.data);
		}, function(){defer.reject('false');});
		return defer.promise;
	};

	instance.checkAuth = function() {
		var defer = $q.defer();
		$http.get('api/checkAuth').then(function(response){
			defer.resolve(response.data);
		}, function(){defer.reject('guest');});
		return defer.promise;
	};

	instance.get = function(url) {
		var defer = $q.defer();
		$http.get(url).then(function(response){
			defer.resolve(response.data);
		}, function(response){defer.reject(response.data);});
		return defer.promise;
	};

	instance.post = function(url, data){
		var defer = $q.defer();
		$http.post(url, data, { headers: { 'Content-Type': 'application/x-www-form-urlencoded; charset=UTF-8'} }).then(function(response) {
			defer.resolve(response.data);
		}, function(response){defer.reject(response.data);});
		return defer.promise;
	};

	instance.apiCall = function(command, qobj, respFunc, errFunc) {

		var qstr = BBBglob.getstring(qobj);

		var obj = $.param({
			string: command+qstr,
		});
		instance.post('/api/getsha', obj).then(function(checksum){
			console.log(BBBglob.BBBURL+command+'?'+qstr+'&checksum='+checksum);
			instance.get(BBBglob.BBBURL+command+'?'+qstr+'&checksum='+checksum).then(
				respFunc, errFunc
			);
		});
	}

	instance.toggleRunning = function(meetingID, on) {
		$http.get("/api/toggleRunning?meetingID="+meetingID+"&on="+on);
	}

	return instance;
};

function Userdata($resource){
	return $resource('api/getUser', {}, {});
};

function Meeting($resource){
	return $resource('api/meetings?meetingID=:meetingID', {}, {
		query: {method:'GET', params:{meetingID:'all'}, isArray:true}
	});
};

function Passwords($resource){
	return $resource('api/passwords?meetingID=:meetingID', {}, {});
};

angular.module('srv', ['ngResource'])
.factory('REST', ['$http', '$q', REST])
.factory('Userdata', ['$resource', Userdata])
.factory('Meeting', ['$resource', Meeting])
.factory('Passwords', ['$resource', Passwords]);
})();
