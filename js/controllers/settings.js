;(function(){
'use strict';

function settingsCtrl($http, REST){

	var sc = this;

	sc.show = false;

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
			console.log(qstr);
			REST.get(qstr).then(function(response){
				var getMeetings = BBBglob.x2j(response);
				if (getMeetings.returncode != "SUCCESS") {
					sc.msg = "Wrong IP or secret";
					sc.show = true;
					return;
				} else {
					REST.post("/api/setSettings", post).then(function(response){
						if (response == "403") {
							sc.msg = "You are not authorized to change settings";
							sc.show = true;
						} else {
							sc.msg = "New settings are saved!"
							sc.show = true;
						}
					}, function(response){
						sc.msg = "Server error. Please, try later.";
						sc.show = true;
					});
				}
			}, function(){
				sc.msg = "No response from server";
				sc.show = true;
				return;
			});
		});
	};
};

angular.module("SettingsCs", []).controller('settingsCtrl', ['$http', 'REST', settingsCtrl])
})();
