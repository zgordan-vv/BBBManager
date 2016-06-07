package main

import (
	"github.com/valyala/fasthttp"
)

func regexpHandler(r *fasthttp.RequestCtx) {

	word := string(r.FormValue("word"))
	typ := string(r.FormValue("type"))
	exp := ""
	errors := map[string]string{}
	errors["name"] = "Name letter error"
	errors["desc"] = "Description letter error"
	errors["url"] = "Wrong URL"
	errors["num"] = "Must contain only numbers"
	errors["ip"] = "Must contain IP address"
	errors["charnum"] = "Must containt only letters and numbers"
	errors["domain"] = "Must be the domain name in the FQDN format without http://"

	switch typ {
		case "name":
			exp = LOGIN
		case "desc":
			exp = DESC
		case "url":
			exp = URL
		case "num":
			exp = DIGITS
		case "ip":
			exp = IP
		case "charnum":
			exp = CHARNUM
		case "domain":
			exp = DOMAIN
	}
	if validate(exp, word) {outnil(r)} else {out(r, errors[typ])}
}
