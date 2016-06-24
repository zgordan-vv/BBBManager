;(function(){
'use strict';

function mainNav(){
	return {
  	restrict: 'E',
	templateUrl: 'tmpl/mainnav.tmpl'
	};
};

function meetingsList(){
	return {
  	restrict: 'E',
	templateUrl: 'tmpl/mlist.tmpl'
	};
};

function meetingDetails(){
	return {
  	restrict: 'E',
	templateUrl: 'tmpl/details.tmpl'
	};
};

function install(){
	return {
	restrict: 'E',
	templateUrl: 'tmpl/install.tmpl'
	};
};

function auth(){
	return {
	restrict: 'E',
	templateUrl: 'tmpl/auth.tmpl'
	};
};

function warnAuth (){
	return {
	restrict: 'E',
	templateUrl: 'tmpl/warnauth.tmpl'
	};
};

function welcome(){
	return {
	restrict: 'E',
	templateUrl: 'tmpl/welcome.tmpl'
	};
};

function connSettingsForm(){
	return {
		restrict: 'E',
		templateUrl: 'tmpl/settings/connSettingsForm.tmpl',
	};
};

function tomcatSettingsForm(){
	return {
	restrict: 'E',
	templateUrl: 'tmpl/settings/tomcatSettingsForm.tmpl'
	};
};

function fsSettingsForm(){
	return {
	restrict: 'E',
	templateUrl: 'tmpl/settings/fsSettingsForm.tmpl'
	};
};

function clientSettingsForm(){
	return {
	restrict: 'E',
	templateUrl: 'tmpl/settings/clientSettingsForm.tmpl'
	};
};

function showtab() {
	return {
		link: function (scope, element, attrs) {
			element.click(function(e) {
				e.preventDefault();
				$(element).tab('show');
			});
		}
	};
};

function alert() {
	return {
		restrict: 'E',
		scope: { show: '=' },
		replace: true, // Replace with the template below
		transclude: true, // we want to insert custom content inside the directive
		link: function(scope, element, attrs) {
			scope.removeAlert = function() { scope.show = false; };
		},
		templateUrl: 'tmpl/alert.tmpl' // See below
	};
};

function waiting() {
	return {
		restrict: 'E',
		templateUrl: 'tmpl/waiting.tmpl' // See below
	};
};

angular.module('dir', [])
.directive('auth', auth)
.directive('warnAuth', warnAuth)
.directive('install', install)
.directive('mainNav', mainNav)
.directive('meetingsList', meetingsList)
.directive('meetingDetails', meetingDetails)
.directive('connSettingsForm', connSettingsForm)
.directive('tomcatSettingsForm', tomcatSettingsForm)
.directive('fsSettingsForm', fsSettingsForm)
.directive('clientSettingsForm', clientSettingsForm)
.directive('showtab', showtab)
.directive('welcome', welcome)
.directive('alert', alert)
.directive('waiting', waiting)
})();
