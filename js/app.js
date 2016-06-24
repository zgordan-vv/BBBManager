;(function(){
'use strict'

angular.module('BBBapp', ['ngRoute', 'AuthCs', 'SettingsCs', 'InstallCs', 'MeetingCs', 'CommonCs', 'filters', 'srv', 'dir'])
.config(['$routeProvider',
	function($routeProvider) {
		$routeProvider.
		 when('/create', {
			templateUrl: 'tmpl/meetingCreate.tmpl'
		 }).
		 when('/delete', {
			templateUrl: 'tmpl/meetingDelete.tmpl'
		 }).
		 when('/edit', {
			templateUrl: 'tmpl/meetingEdit.tmpl'
		 }).
		 when('/login', {
			templateUrl: 'tmpl/auth.tmpl'
		 }).
		 when('/quit', {
			templateUrl: 'tmpl/quit.tmpl'
		 }).
		 when('/settings', {
			templateUrl: 'tmpl/settings.tmpl'
		 }).
		 when('/resetDefaults', {
			templateUrl: 'tmpl/resetDefaults.tmpl'
		 }).
		 when('/wait', {
			templateUrl: 'tmpl/wait.tmpl'
		 }).
		 when('/', {
			templateUrl: 'tmpl/meetingSelect.tmpl'
		 }).
		 otherwise({
			redirectTo: '/'
		 })
	}
]);
})();
