'use strict';

var BBBglob = {};

BBBglob.BBBURL = 'http://85.143.223.74/bigbluebutton/api/';

BBBglob.x2j = function(xmlstr) {
	var x2js = new X2JS();
	var jsonResp = x2js.xml_str2json(xmlstr);
	return jsonResp.response;
}

BBBglob.getstring = function(obj) {
	var result = "";
	for (var key in obj) {
		result += '&'+key+'='+encodeURIComponent(obj[key]);
	};
	result = result.substring(1);
	return result;
}
