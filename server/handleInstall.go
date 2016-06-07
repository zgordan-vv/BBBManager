package main

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"fmt"
)

func restInstalled(r *fasthttp.RequestCtx) { r.Write(output[installed]) }

func nameUniqHandler(r *fasthttp.RequestCtx) {
	name := r.FormValue("name")
	if nameUniq(string(name)) {outnil(r)} else {out(r, "The name exists")}
}

func installHandler(r *fasthttp.RequestCtx) {
	jsonObj := r.FormValue("installdata")
	data := struct{
		Name string	`json:"name"`
		Fullname string	`json:"fullname"`
		Pwd string	`json:"pwd"`
		Dbprefix string	`json:"dbprefix"`
		Domainname string	`json:"domainname"`
	}{}
	if json.Unmarshal(jsonObj, &data) != nil {return}
	login := data.Name
	fullname := data.Fullname
	dbprefix := data.Dbprefix
	dname := data.Domainname
	pwd := data.Pwd

	if !validate(LOGIN, login) || !validate(LOGIN, fullname) || (len(login) < 1) || (len(pwd) < 6) || !validate(CHARNUM, dbprefix) || !validate(DOMAIN, dname) || !nameUniq(login) {out500(r); return}

	DBPREFIX = dbprefix+":"
	DOMAINNAME = dname
	saveGlobs(DBPREFIX,DOMAINNAME)
	user := &User{login, fullname, true, passEncrypt(pwd)}
	if err := saveUser(*user); err != nil {fmt.Println("Users have not been saved", err); out(r,"DontSaved"); return}
	installed = true
	newSessionValid(r,login)
}
