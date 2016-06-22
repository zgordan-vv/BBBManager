;(function(){
'use strict';

function settingsCtrl($http, REST){

	var sc = this;

	sc.show = false;
	sc.waiting = false;
	sc.settingsTab = "conn";

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
		sc.msg = "New settings are saved";
		sc.show = true;
	}

	function notSaved(){
		sc.msg = "Server error, please try later";
		sc.show = true;
	}

	function notAuth(){
		sc.msg = "You are not authorized to change settings";
		sc.show = true;
	}

	function wait(){
		sc.waiting = false;
		sc.msg = "Please wait about 2 minutes until server restarts"
		sc.show = true;
	}

	$http.get('/api/getSecToken').then(function(response){
		sc.tokensec = response.data;
	}, function(){ sc.tokensec = ""; });

	REST.get("/api/getSettings").then(function(connData){sc.settings = connData;});

	REST.get("/api/getTomcat").then(function(tomcatData){
		sc.tomcat = tomcatData.Params;
		sc.tomcat.maxNumPages = +sc.tomcat.maxNumPages;
		sc.tomcat.defaultMaxUsers = +sc.tomcat.defaultMaxUsers;
		sc.tomcat.defaultMeetingCreateJoinDuration = +sc.tomcat.defaultMeetingCreateJoinDuration;
		sc.tomcat.defaultMeetingExpireDuration = +sc.tomcat.defaultMeetingExpireDuration;
		sc.tomcat.defaultMeetingDuration = +sc.tomcat.defaultMeetingDuration;
		sc.tomcat.allowStartStopRecording = sc.tomcat.allowStartStopRecording == "true";
		sc.tomcat.disableRecordingDefault = sc.tomcat.disableRecordingDefault == "true";
	});

	REST.get("/api/getFreeswitch").then(function(fsData){
		sc.freeswitch = fsData.Params;
	})

	REST.get("/api/getClient").then(function(clientData){
		sc.client = clientData.Params;
		for (var param in sc.client) {
			if (param == 'VideoconfModule/videoQuality' || param == 'VideoconfModule/camQualityBandwidth' || param == 'VideoconfModule/camQualityPicture') {
				sc.client[param] = +sc.client[param]
			} else {
				sc.client[param] = sc.client[param] == 'true';
			}
		}
	})

	sc.updateCheckIP = function() {
		$http.get('api/regexp?word='+sc.settings.ip+'&type=ip').then(function(response){
			sc.ipValid = response.data;
		}, function(){
			sc.ipValid = "Server error";
		});
	};

	sc.updateCheckSecret = function() {
		$http.get('api/regexp?word='+sc.settings.secret+'&type=num').then(function(response){
			sc.secretValid = response.data;
		}, function(){
			sc.secretValid = "Server error";
		});
	};

	sc.submitsettings = function() {
		var post = $.param({
			tokensec: sc.tokensec,
			settings: JSON.stringify({
				ip: sc.settings.ip,
				secret: sc.settings.secret,
			})
		});
		var qstr = 'getMeetings';
		sc.waiting = true;
		REST.get('/api/getsecrsha?string='+qstr+'&secret='+sc.settings.secret).then(function(checksum){
			qstr='http://'+sc.settings.ip+'/bigbluebutton/api/'+qstr+'?checksum='+checksum;
			REST.get(qstr).then(function(response){
				var getMeetings = BBBglob.x2j(response);
				sc.waiting = false;
				if (getMeetings.returncode != "SUCCESS") {
					sc.msg = "Wrong IP or secret";
					sc.show = true;
					return;
				} else {
					REST.post("/api/setSettings", post).then(function(response){
						if (response == "403") {notAuth()} else if (response == "500") {wait()} else {saved()}
					}, wait());
				}
			}, function(){
				wait();
			});
		}, function(){
			wait();
		});
	};

	sc.submittomcat = function() {
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
				if (response == "500") {wait()} else {saved()}
			}
		}, function(){wait()});
	};

	sc.audioprofiles = ["default", "wideband", "ultrawideband", "cdquality", "sla"];

	sc.submitfs = function(){
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
			if (response == "500") {wait()} else {saved()}}
		}, function(){wait()});
	};

	sc.submitclient = function(){
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
			if (response == "500") {wait()} else {saved()}}
		}, function(){
			wait();
		});
	};
};

angular.module("SettingsCs", []).controller('settingsCtrl', ['$http', 'REST', settingsCtrl])
})();
