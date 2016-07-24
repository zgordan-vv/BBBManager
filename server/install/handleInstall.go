package install

import (
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/out"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/globs"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/users"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/sessions"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/auth"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/check"
	"encoding/json"
	"github.com/valyala/fasthttp"
	"fmt"
)

func RestInstalled(r *fasthttp.RequestCtx) { r.Write(globs.Output[Installed]) }

func NameUniqHandler(r *fasthttp.RequestCtx) {
	name := r.FormValue("name")
	if users.NameUniq(string(name)) {out.Outnil(r)} else {out.Out(r, "The name exists")}
}

func InstallHandler(r *fasthttp.RequestCtx) {
	jsonObj := r.FormValue("installdata")
	data := struct{
		Name string	`json:"name"`
		Fullname string	`json:"fullname"`
		Pwd string	`json:"pwd"`
	}{}
	if json.Unmarshal(jsonObj, &data) != nil {return}
	login := data.Name
	fullname := data.Fullname
	pwd := data.Pwd

	if !check.Validate(check.LOGIN, login) || !check.Validate(check.LOGIN, fullname) || (len(login) < 1) || (len(pwd) < 6) || !users.NameUniq(login) {out.Out500(r); return}

	user := &users.User{login, fullname, true, auth.PassEncrypt(pwd)}
	if err := users.SaveUser(*user); err != nil {fmt.Println("Users have not been saved", err); out.Out(r,"DontSaved"); return}
	Installed = true
	sessions.NewSessionValid(r,login)
}
