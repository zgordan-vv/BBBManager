;(function(){
'use strict';

function settingsCtrl($http, REST){

	var sc = this;

	sc.show = false;
	sc.waiting = false;
	sc.settingsTab = "conn";

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

	$http.get('/api/getSecToken').then(function(response){
		sc.tokensec = response.data;
	}, function(){ sc.tokensec = ""; });

	REST.get("/api/getSettings").then(function(connData){sc.settings = connData;});

	REST.get("/api/getTomcat").then(function(tomcatData){
		sc.tomcat = tomcatData.Params;
		sc.tomcat.maxNumPages = Number(sc.tomcat.maxNumPages);
		sc.tomcat.defaultMaxUsers = Number(sc.tomcat.defaultMaxUsers);
		sc.tomcat.defaultMeetingCreateJoinDuration = Number(sc.tomcat.defaultMeetingCreateJoinDuration);
		sc.tomcat.defaultMeetingExpireDuration = Number(sc.tomcat.defaultMeetingExpireDuration);
		sc.tomcat.defaultMeetingDuration = Number(sc.tomcat.defaultMeetingDuration);
		sc.tomcat.allowStartStopRecording = sc.tomcat.allowStartStopRecording == "true";
		sc.tomcat.disableRecordingDefault = sc.tomcat.disableRecordingDefault == "true";
	});

	REST.get("/api/getFreeswitch").then(function(fsData){
		sc.freeswitch = fsData.Params;
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
						if (response == "403") {notAuth()} else {saved()}
					}, notSaved());
				}
			}, function(){
				sc.msg = "No response from server";
				sc.show = true;
				return;
			});
		});
	};

	sc.submittomcat = function() {
		var params = JSON.parse(JSON.stringify(sc.tomcat));
/*		params.maxNumPages = params.maxNumPages.toString();
		params.defaultMaxUsers = params.defaultMaxUsers.toString();
		params.defaultMeetingCreateJoinDuration = params.defaultMeetingCreateJoinDuration.toString();
		params.defaultMeetingExpireDuration = params.defaultMeetingExpireDuration.toString();
		params.defaultMeetingDuration = params.defaultMeetingDuration.toString();
		params.allowStartStopRecording = params.allowStartStopRecording.toString()
		params.disableRecordingDefault = params.disableRecordingDefault.toString();*/

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
		}, function(){notSaved();});
	};

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
			if (response == "500") {notSaved()} else {saved()}}
		}, function(){notSaved();});
	};
};

angular.module("SettingsCs", []).controller('settingsCtrl', ['$http', 'REST', settingsCtrl])
})();
