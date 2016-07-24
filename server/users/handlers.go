package users

import (
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/check"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/out"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/auth"
	"github.com/valyala/fasthttp"
	"encoding/json"
)

func GetUserName(r *fasthttp.RequestCtx) string {
	username := r.UserValue("username")
	if username != nil {return username.(string)} else {return ""}
}

func RestGetUserName(r *fasthttp.RequestCtx) {
	out.Out(r, GetUserName(r))
}

func RestGetUserFullName(r *fasthttp.RequestCtx) {
	fullName := ""
	username := GetUserName(r)
	user, ok := GetUser(username)
	if ok {fullName = user.FullName}
	out.Out(r, fullName)
}

func RestGetUser(r *fasthttp.RequestCtx) {
	username := GetUserName(r)
	user, _ := GetUser(username)
	output, err := json.Marshal(user)
	if err == nil {r.Write(output)}
}

func CheckAuthHandler(r *fasthttp.RequestCtx) {
	output := "guest"
	username := GetUserName(r)
	user, ok := GetUser(username)
	if ok {
		if username != "" {
			if user.IsAdmin {output = "admin"} else {output = "user"}
		}
	}	
	out.Out(r, output)
}

func ProfileSaveHandler(r *fasthttp.RequestCtx) {
	username := GetUserName(r)
	if username == "" { out.Out(r, "NotSaved"); return }
	user, ok := GetUser(username)
	if !ok { out.Out(r, "NotSaved"); return }
	fullname := string(r.FormValue("fullname"))
	if !check.CheckLogin(fullname) { out.Out(r, "WrongFullName") } else {
		oldpwd := string(r.FormValue("oldpwd"))
		pwd := string(r.FormValue("pwd"))
		if changed, correct := CheckPwdChange(username, oldpwd, pwd); !correct {out.Out403(r)} else {
			user.FullName = fullname
			if changed {user.Keyword = auth.PassEncrypt(pwd)}
			if err := SaveUser(user); err != nil {out.Out(r, "DontSave")}
			out.Out(r, "ok")
		}
	}
}
