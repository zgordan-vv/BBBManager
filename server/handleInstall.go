package main

import (
	"encoding/json"
	"net/http"
	"fmt"
)

func restInstalled(w http.ResponseWriter, r *http.Request) { w.Write(output[installed]) }

func nameUniqHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if nameUniq(name) {w.Write(nil)} else {out(w, "The name exists")}
}

func installHandler(w http.ResponseWriter, r *http.Request) {
	jsonObj := r.FormValue("installdata")
	data := struct{
		Name string	`json:"name"`
		Fullname string	`json:"fullname"`
		Pwd string	`json:"pwd"`
		Dbprefix string	`json:"dbprefix"`
		Domainname string	`json:"domainname"`
	}{}
	if json.Unmarshal([]byte(jsonObj), &data) != nil {return}
	login := data.Name
	fullname := data.Fullname
	dbprefix := data.Dbprefix
	dname := data.Domainname
	pwd := data.Pwd

	if !validate(LOGIN, login) || !validate(LOGIN, fullname) || (len(login) < 1) || (len(pwd) < 6) || !validate(CHARNUM, dbprefix) || !validate(DOMAIN, dname) || !nameUniq(login) {out500(w); return}

	DBPREFIX = dbprefix+":"
	DOMAINNAME = dname
	saveGlobs(DBPREFIX,DOMAINNAME)
	user := &User{login, fullname, true, passEncrypt(pwd)}
	if err := saveUser(*user); err != nil {fmt.Println("Users have not been saved", err); out(w,"DontSaved"); return}
	installed = true
	newSessionValid(w,r,login)
}
