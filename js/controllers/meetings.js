;(function(){

'use strict';

function MeetingsView($rootScope, $http, $location, REST, Meeting, Passwords, Userdata) {

	var mv = this;

	mv.show=false;
	mv.showyn = false;
	mv.link = "false";
	mv.showdetails = false;

	mv.meetings = Meeting.query();

	mv.change = function(meetingID){
		mv.link = "false";
		Meeting.get({meetingID: meetingID}, function(meeting){}).$promise.then(function(result){
			mv.meeting = result;
			mv.isrunning = mv.meeting.isrunning;
			Passwords.get({meetingID: meetingID}, function(){}).$promise.then(function(result){
				mv.passwords = result;
				if ($rootScope.authorized == 'admin') {
					REST.apiCall("getMeetingInfo", {meetingID: meetingID, password: mv.passwords.admpwd}, function(response){
						var meetingJson = BBBglob.x2j(response);
						var attObj = meetingJson.attendees;
						if (!attObj || !attObj.attendee) {mv.attendees = [];} else if (attObj.attendee.constructor !== Array) {mv.attendees=new Array(attObj.attendee);} else {mv.attendees = attObj.attendee;};
					}, function(error){
						mv.attendees=[];
					});
					REST.apiCall("getRecordings", {meetingID: meetingID}, function(response){
						var recordingsJson = BBBglob.x2j(response);
						var recordObj = recordingsJson.recordings;
						if (!recordObj || !recordObj.recording) {mv.recordings = [];} else if (recordObj.recording.constructor !== Array) {mv.recordings=new Array(recordObj.recording);} else {mv.recordings = recordObj.recording;};
					}, function(error){
						mv.recordings=[];
					});
				} else {
					mv.attendees=[];
					mv.recordings = [];
				}
			});
		});
		mv.showdetails = true;
	};

	mv.create = function(){
		$location.url("/create");
	}

	mv.createonserver = function(meetingID){

		REST.apiCall("create", {
			meetingID: meetingID,
			name: mv.meeting.desc,
			welcome: mv.meeting.welcome,
			duration: mv.meeting.duration,
			moderatorPW: mv.passwords.admpwd,
			attendeePW: mv.passwords.pwd,
			record: mv.meeting.isrec || false,
			autoStartRecording: mv.meeting.autorec || false,
			allowStartStopRecording: mv.meeting.allowstartstoprec || false,
			logoutURL: 'http://slava.zgordan.ru',
		}, function(response){
			mv.join(meetingID, mv.passwords.admpwd);
			REST.toggleRunning(meetingID, true);
		}, function(error){
			mv.msg = "The meeting is not created.\nServer response: "+error;
			mv.show=true;
		});
	};

	mv.join = function(meetingID, meetingpwd){
		Userdata.get({}, function(){}).$promise.then(function(user){
			if (!user) {return;}
			var name = user.name;
			var fullname = user.fullname;
			REST.get('api/checkPwd?pwd='+meetingpwd+'&id='+meetingID).then(function(response){
				if (response == 'true'){
					var command="join";
					var qstr = BBBglob.getstring({
						fullName: fullname,
						meetingID: meetingID,
						userID: name,
						password: meetingpwd,
					});
					var obj = $.param({
						string: qstr,
					});

					REST.post('/api/join',obj).then(function(result){
						mv.link = BBBglob.BBBURL+result;
					});
				} else {
					mv.msg = "Wrong user password!";
					mv.show = true;
				};
			}, function(){
					mv.msg = "Server error";
					mv.show = true;
			});
		});
	};

	mv.endalert = function(meetingID){
		mv.msgyn = "Are you sure to end meeting "+mv.meeting.title+" ?";
		mv.showyn = true;
	};

	mv.end = function(meetingID){

		REST.toggleRunning(meetingID, false);
		REST.apiCall("end", {
			meetingID: meetingID,
			password: mv.passwords.admpwd,
		}, function(){
			mv.showyn = false;
			mv.change(meetingID);
		}, function(){
			mv.showyn = false;
			mv.change(meetingID);
		});
	};

	mv.cancel = function(){
		mv.showyn = false;
	};
};

function createValid($http, $location, REST){

	var cv = this;

	cv.show = false;

	$http.get('/api/getDupToken');

	$http.get('/api/getSecToken').then(function(response){
		cv.tokensec = response.data;
	}, function(){ cv.tokensec = ""; });

	cv.titleValid="";
	cv.titleUniq="";
	cv.descValid="";
	cv.title = "";
	cv.desc = "";
	cv.admpwd = "";
	cv.admpwdconf = "";
	cv.pwd = "";
	cv.pwdconf = "";

	cv.updateCheckTitle = function() {
		$http.get('api/regexp?word='+cv.title+'&type=name').then(
			function(response){
				cv.titleValid = response.data;
			}, function(){
				cv.titleValid = "Server error";
			}
		);
		$http.get('api/meetingUniq?title='+cv.title).then(function(response){
			cv.titleUniq = response.data;
		}, function(){
			cv.titleUniq = "Server error";
		});
	};

	cv.updateCheckDesc = function() {
		$http.get('api/regexp?word='+cv.desc+'&type=desc').then(function(response){
			cv.descValid = response.data;
		}, function(){
			cv.titleValid = "Server error";
		});
	};

	cv.submitcreate = function() {
		if (cv.title == "") {cv.msg = "Title is empty!"; cv.show = true; return;}
		if (cv.titleValid != "") {cv.msg = "Title contains wrong characters!"; cv.show = true; return;}
		if (cv.titleUniq != "") {cv.msg = "There is another meeting with \'"+cv.title+"\' title!"; cv.show = true; return;}
		if (cv.descValid != "") {cv.msg = "Description contains wrong characters!"; cv.show = true; return;}
		if (!cv.admpwd || cv.admpwd.length < 6) {cv.msg = "Admin password must be 6 symbols or more!"; cv.show = true; return;}
		if (cv.admpwd != cv.admpwdconf) {cv.msg = "Admin passwords must match!"; cv.show = true; return;}
		if (!cv.pwd || cv.pwd.length < 6) {cv.msg = "User password must be 6 symbols or more!"; cv.show = true; return;}
		if (cv.pwd != cv.pwdconf) {cv.msg = "User passwords must match!"; cv.show = true; return;}
		var post = $.param({
			tokensec: cv.tokensec,
			meeting: JSON.stringify({
				id: "",
				title: cv.title,
				desc: cv.desc,
				author: cv.author,
				welcome: cv.welcome,
				duration: cv.duration,
				isrec: cv.isrec,
				autorec: cv.autorec,
				allowstartstoprec: cv.allowstartstoprec,
			}),
			passwords: JSON.stringify({
				admpwd: cv.admpwd,
				pwd: cv.pwd,
			}),
		});
		REST.post("/api/create", post).then(function(response){
			if (response == "403") {
				cv.warn = "You are not authorized to create meetings";
			} else {
				$location.url("/").replace;
			}
		}, function(response){
			cv.warn = "Server error. Please, try later.";
		});
	};
};

function editValid($http, $location, $routeParams, Meeting, Passwords, REST){

	var ev = this;

	$http.get('/api/getSecToken').then(function(response){
		ev.tokensec = response.data;
	}, function(){ ev.tokensec = ""; });

	ev.id = $routeParams.id;
	ev.meeting = Meeting.get({meetingID: ev.id}, function(meeting){});
	ev.passwords = Passwords.get({meetingID: ev.id}, function(meeting){});
	ev.admpwd = "";
	ev.admpwdconf = "";
	ev.pwd = "";
	ev.pwdconf = "";

	ev.titleValid="";
	ev.titleUniq="";
	ev.descValid="";

	ev.updateCheckTitle = function() {
		$http.get('api/regexp?word='+ev.meeting.title+'&type=name').then(
			function(response){
				ev.titleValid = response.data;
			}, function(){
				ev.titleValid = "Server error";
			}
		);
	};

	ev.updateCheckDesc = function() {
		$http.get('api/regexp?word='+ev.meeting.desc+'&type=desc').then(function(response){
			ev.descValid = response.data;
		}, function(){
			ev.titleValid = "Server error";
		});
	};

	ev.submitedit = function() {
		if (ev.meeting.title == "") {ev.msg = "Title is empty!"; ev.show = true; return;}
		if (ev.titleValid != "") {ev.msg = "Title contains wrong characters!"; ev.show = true; return;}
		if (ev.titleUniq != "") {ev.msg = "There is another meeting with the same title!"; ev.show = true; return;}
		if (ev.descValid != "") {ev.msg = "Description contains wrong characters!"; ev.show = true; return;}
		if (!ev.passwords.admpwd || ev.passwords.admpwd.length < 6) {ev.msg = "Admin password must be 6 symbols or more!"; ev.show = true; return;}
		if (ev.passwords.admpwd != ev.admpwdconf) {ev.msg = "Admin passwords must match!"; ev.show = true; return;}
		if (!ev.passwords.pwd || ev.passwords.pwd.length < 6) {ev.msg = "User password must be 6 symbols or more!"; ev.show = true; return;}
		if (ev.passwords.pwd != ev.pwdconf) {ev.msg = "User passwords must match!"; ev.show = true; return;}
		var post = $.param({
			tokensec: ev.tokensec,
			meeting: JSON.stringify({
				id: ev.meeting.id,
				title: ev.meeting.title,
				desc: ev.meeting.desc,
				author: ev.meeting.author,
				welcome: ev.meeting.welcome,
				duration: ev.meeting.duration,
				isrec: ev.meeting.isrec,
				autorec: ev.meeting.autorec,
				allowstartstoprec: ev.meeting.allowstartstoprec,
			}),
			passwords: JSON.stringify({
				admpwd: ev.passwords.admpwd,
				pwd: ev.passwords.pwd,
			}),
		});
		REST.post("/api/edit", post).then(function(response){
			if (response == "403") {
				ev.warn = "You are not authorized to edit meetings";
			} else {
				$location.url("/").replace();
			}
		}, function(response){
			ev.warn = "Server error. Please, try later.";
		});
	};
}

function deleteMeeting($http, $location, $routeParams, REST){

	var dm = this;

	$http.get('/api/getSecToken').then(function(response){
		dm.tokensec = response.data;
	}, function(){ dm.tokensec = ""; });

	dm.id = $routeParams.id;
	dm.title = $routeParams.title;

	dm.yes = function() {
		var post = $.param({
			tokensec: dm.tokensec,
			meetingID: dm.id,
		});
		REST.post("/api/delete", post).then(function(response){
			if (response == "403") {
				dm.warn = "You are not authorized to delete meetings";
			} else {
				$location.url("/").replace();
			}
		}, function(response){
			dm.warn = "Server error. Please, try later.";
		});
	};

	dm.no = function() {
		$location.url("/").replace();
	}
};

angular.module("MeetingCs", [])
.controller('MeetingsView', ['$rootScope', '$http', '$location', 'REST', 'Meeting', 'Passwords', 'Userdata', MeetingsView])
.controller('createValid', ['$http', '$location', 'REST', createValid])
.controller('editValid', ['$http', '$location', '$routeParams', 'Meeting', 'Passwords', 'REST', editValid])
.controller('deleteMeeting', ['$http', '$location', '$routeParams', 'REST', deleteMeeting])
})();
