;(function(){
'use strict';

function settingsCtrl($http, $location, REST){

	var sc = this;

	REST.get('/api/getMaintenance').then(function(response){
		if (response == 'true') {
			sc.maintenance = true;
			sc.serverWarning = "Server is busy or restarting, wait 2 minutes and reload the page";
			settings();
		} else {
			sc.maintenance = false;
			sc.serverWarning = "";
			settings();
		}
	}, function(error){
		sc.maintenance = true;
		sc.serverWarning = "Server is busy or restarting, wait 2 minutes and reload the page";
		settings();
	});

	function settings(){

	sc.show = false;
	sc.waiting = false;
	sc.settingsTab = "tomcat";

	sc.profilename = {
		"default": "Default",
		"wideband": "Wide band",
		"ultrawideband": "Ultrawide band",
		"cdquality": "CD quality",
		"sla": "Shared line appearance",
	}

	sc.settingsSwitch = function(option){
		sc.settingsTab = option;
	}

	function saved(){
		sc.waiting = false;
		sc.msg = "New settings are saved";
		sc.show = true;
	}

	function notSaved(){
		sc.waiting = false;
		sc.msg = "Server error, please try later";
		sc.show = true;
	}

	function notAuth(){
		sc.waiting = false;
		sc.msg = "You are not authorized to change settings";
		sc.show = true;
	}

	function wait(){
		$location.url('/wait').replace();
	}

	sc.resetDefaults = function(setting) {
		$location.url("/resetDefaults?setting="+setting);
	}

	$http.get('/api/getSecToken').then(function(response){
		sc.tokensec = response.data;
	}, function(){ sc.tokensec = ""; });

	REST.get("/api/getIP").then(function(ip){
		if (ip != "") {	sc.ip = ip; }
	});

	REST.get("/api/getTomcat").then(function(tomcatData){
		if (tomcatData.Params) {
			sc.tomcat = tomcatData.Params;
			sc.tomcat.maxNumPages = +sc.tomcat.maxNumPages;
			sc.tomcat.defaultMaxUsers = +sc.tomcat.defaultMaxUsers;
			sc.tomcat.defaultMeetingCreateJoinDuration = +sc.tomcat.defaultMeetingCreateJoinDuration;
			sc.tomcat.defaultMeetingExpireDuration = +sc.tomcat.defaultMeetingExpireDuration;
			sc.tomcat.defaultMeetingDuration = +sc.tomcat.defaultMeetingDuration;
			sc.tomcat.allowStartStopRecording = sc.tomcat.allowStartStopRecording == "true";
			sc.tomcat.disableRecordingDefault = sc.tomcat.disableRecordingDefault == "true";
		} else { wait() }
	});

	REST.get("/api/getFreeswitch").then(function(fsData){
		if (fsData.Params) { sc.freeswitch = fsData.Params; } else { wait() }
	})

	REST.get("/api/getClient").then(function(clientData){
		if (clientData.Params) {
			sc.client = clientData.Params;
			for (var param in sc.client) {
				if (param == 'VideoconfModule/videoQuality' || param == 'VideoconfModule/camQualityBandwidth' || param == 'VideoconfModule/camQualityPicture') {
					sc.client[param] = +sc.client[param]
				} else {
					sc.client[param] = sc.client[param] == 'true';
				}
			}
		} else { wait() }
	})

	sc.updateCheckSecret = function() {
		$http.get('api/regexp?word='+sc.settings.secret+'&type=num').then(function(response){
			sc.secretValid = response.data;
		}, function(){
			sc.secretValid = "Server error";
		});
	};

	sc.submitIP = function(){
		REST.get('/api/getMaintenance').then(function(response){
			if (response != 'false') {
				wait();
			} else {
				var ip = sc.ip;
		
				var post = $.param({
					tokensec: sc.tokensec,
					ip: sc.ip,
					})
				};
				sc.waiting = true;
				REST.post("/api/setIP", post).then(function(response){
					sc.waiting = false;
					if (response == "403") {notAuth()} else {
					if (response == "500") {notSaved()} else {saved()}}
				}, function(error){
					wait();
				});
		}, function(error){
			wait();
		});
	};

	sc.submittomcat = function() {
		REST.get('/api/getMaintenance').then(function(response){
			if (response != 'false') {
				wait();
			} else {
				var params = JSON.parse(JSON.stringify(sc.tomcat));
		
				for (var param in params) {
					params[param] = params[param].toString();
				}
		
				var post = $.param({
					tokensec: sc.tokensec,
					settings: JSON.stringify({
						params: params,
					})
				});
				sc.waiting = true;
				REST.post("/api/setTomcat", post).then(function(response){
					sc.waiting = false;
					if (response == "403") {notAuth()} else {
						if (response == "500") {notSaved()} else {saved()}
					}
				}, function(){notSaved()});
			}
		}, function(error){
			wait();
		});
	};

	sc.audioprofiles = ["default", "wideband", "ultrawideband", "cdquality", "sla"];

	sc.submitfs = function(){
		REST.get('/api/getMaintenance').then(function(response){
			if (response != 'false') {
				wait();
			} else {
				var post = $.param({
					tokensec: sc.tokensec,
					settings: JSON.stringify({
						params: sc.freeswitch,
					})
				});
				sc.waiting = true;
				REST.post("/api/setFreeswitch", post).then(function(response){
					sc.waiting = false;
					if (response == "403") {notAuth()} else {
					if (response == "500") {notSaved()} else {saved()}}
				}, function(){notSaved()});
			}
		}, function(error){
			wait();
		});
	};

	sc.submitclient = function(){
		REST.get('/api/getMaintenance').then(function(response){
			if (response != 'false') {
				wait();
			} else {
				var params = JSON.parse(JSON.stringify(sc.client));
				for (var param in params) {
					params[param] = params[param].toString();
				}
		
				var post = $.param({
					tokensec: sc.tokensec,
					settings: JSON.stringify({
						params: params,
					})
				});
				sc.waiting = true;
				REST.post("/api/setClient", post).then(function(response){
					sc.waiting = false;
					if (response == "403") {notAuth()} else {
					if (response == "500") {notSaved()} else {saved()}}
				}, function(){
					notSaved();
				});
			}
		}, function(error){
			wait();
		});
	};
	}
}

function resetDefaults($http, $location, $routeParams, REST){

	function wait(){
		$location.url('/wait').replace();
	}

	var rd = this;
	rd.waiting = false;

	$http.get('/api/getSecToken').then(function(response){
		rd.tokensec = response.data;
	}, function(){ rd.tokensec = ""; });

	rd.setting = $routeParams.setting;

	rd.yes = function() {
		REST.get('/api/getMaintenance').then(function(response){
			if (response != 'false') {
				wait();
			} else {
				if (!rd.setting) { wait(); } else {
					var post = $.param({
						tokensec: rd.tokensec,
						setting: rd.setting,
						answer: "yes",
					});
					rd.waiting = true;
					REST.post("/api/reset"+rd.setting, post).then(function(response){
						rd.waiting = false;
						if (response == "403") {
							$location.url("/").replace();
						} else {
							if (response == "500") {
								wait();
							} else {
								$location.url("/settings").replace();
							}
						}
					}, function(response){
						wait();
					});
				}
			}
		}, function(error){
			wait();
		})
	};

	rd.no = function() {
		$location.url("/settings").replace();
	}
}

angular.module("SettingsCs", [])
.controller('settingsCtrl', ['$http', '$location', 'REST', settingsCtrl])
.controller('resetDefaults', ['$http', '$location', '$routeParams', 'REST', resetDefaults])
})();
