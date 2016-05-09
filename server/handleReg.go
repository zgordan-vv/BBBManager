package main

import (
	"net/http"
//	"github.com/haisum/recaptcha"
	"encoding/json"
	"fmt"
)

//var recap recaptcha.R = recaptcha.R{Secret: SECRETKEY}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	jsonObj := r.FormValue("registerdata")
	data := struct{
		Name string	`json:"name"`
		Fullname string	`json:"fullname"`
		Pwd string	`json:"pwd"`
		Pwdconf string	`json:"pwdconf"`
	}{}
	if json.Unmarshal([]byte(jsonObj), &data) != nil {return}
	login := data.Name
	fullname := data.Fullname
	pwd := data.Pwd

	if !validate(LOGIN, login) || !validate(LOGIN, fullname) || (len(login) < 1) || (len(pwd) < 6) || !nameUniq(login) || (pwd != data.Pwdconf) {out500(w); return}

	user := &User{login, fullname, false, passEncrypt(pwd)}
	if err := saveUser(*user); err != nil {fmt.Println("Users have not been saved", err); out(w, "DontSaved"); return}
	out200(w)
	newSessionValid(w,r,login)
}
