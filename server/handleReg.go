package main

import (
	"github.com/valyala/fasthttp"
//	"github.com/haisum/recaptcha"
	"encoding/json"
	"fmt"
)

//var recap recaptcha.R = recaptcha.R{Secret: SECRETKEY}

func registerHandler(r *fasthttp.RequestCtx) {
	jsonObj := r.FormValue("registerdata")
	data := struct{
		Name string	`json:"name"`
		Fullname string	`json:"fullname"`
		Pwd string	`json:"pwd"`
		Pwdconf string	`json:"pwdconf"`
	}{}
	if json.Unmarshal(jsonObj, &data) != nil {return}
	login := data.Name
	fullname := data.Fullname
	pwd := data.Pwd

	if !validate(LOGIN, login) || !validate(LOGIN, fullname) || (len(login) < 1) || (len(pwd) < 6) || !nameUniq(login) || (pwd != data.Pwdconf) {out500(r); return}

	user := &User{login, fullname, false, passEncrypt(pwd)}
	if err := saveUser(*user); err != nil {fmt.Println("Users have not been saved", err); out(r, "DontSaved"); return}
	out200(r)
	newSessionValid(r,login)
}
