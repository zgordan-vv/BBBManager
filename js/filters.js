;(function(){
'use strict';

angular.module('filters', [])
.filter('yn', function(){
	return function(input) {
		return (!input || input=='false') ? 'no':'yes';
	}
});
})();
