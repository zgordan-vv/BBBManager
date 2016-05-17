;(function(){
'use strict';

function settingsCtrl($http, REST){

	var sc = this;

	sc.show = false;

	function saved(){
		sc.msg = "New settings are saved";
		sc.show = true;
	}

	function notSaved(){
		sc.msg = "Server error, please try later";
		sc.show = true;
	}

	function notAuth(){
		sc.msg = "Server error, please try later";
		sc.show = true;
	}

	$http.get('/api/getSecToken').then(function(response){
		sc.tokensec = response.data;
	}, function(){ sc.tokensec = ""; });

	REST.get("/api/getSettings").then(function(data){sc.settings = data;});

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
		REST.get('/api/getsecrsha?string='+qstr+'&secret='+sc.settings.secret).then(function(checksum){
			qstr='http://'+sc.settings.ip+'/bigbluebutton/api/'+qstr+'?checksum='+checksum;
			REST.get(qstr).then(function(response){
				var getMeetings = BBBglob.x2j(response);
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
		var post = $.param({
			tokensec: sc.tokensec,
			settings: JSON.stringify({
				params: [
					["maxNumPages", sc.tomcat.maxNumPages.toString()],
					["defaultWelcomeMessage", sc.tomcat.defaultWelcomeMessage],
					["defaultWelcomeMessageFooter", sc.tomcat.defaultWelcomeMessageFooter],
					["defaultMaxUsers", sc.tomcat.defaultMaxUsers.toString()],
					["defaultMeetingDuration", sc.tomcat.defaultMeetingDuration.toString()],
					["defaultMeetingExpireDuration", sc.tomcat.defaultMeetingExpireDuration.toString()],
					["defaultMeetingCreateJoinDuration", sc.tomcat.defaultMeetingCreateJoinDuration.toString()],
					["disableRecordingDefault", sc.tomcat.disableRecordingDefault.toString()],
					["allowStartStopRecording", sc.tomcat.allowStartStopRecording.toString()],
					["bbb.web.logoutURL", sc.tomcat.bbbWebLogoutURL],
					["defaultAvatarURL", sc.tomcat.defaultAvatarURL],
				],
			})
		});
		REST.post("/api/setTomcat", post).then(function(response){
			if (response == "403") {notAuth()} else {saved()}
		}, notSaved());
	};
};

angular.module("SettingsCs", []).controller('settingsCtrl', ['$http', 'REST', settingsCtrl])
})();
