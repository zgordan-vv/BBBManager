package login

import (
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/auth"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/globs"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/db"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/out"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/sessions"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/users"
	"fmt"
	"github.com/valyala/fasthttp"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
)

func BrowserLoginHandler(r *fasthttp.RequestCtx) {
	jsonObj := r.FormValue("login")
	logindata := struct{
		Login string	`json:"login"`
		Password string	`json:"password"`
	}{}
	if json.Unmarshal(jsonObj, &logindata) != nil {return}
	login := logindata.Login
	password := logindata.Password

	client, err := db.InitRedis()
	if err == nil {
		defer client.Close()
		userExists,err := client.Do("SISMEMBER", globs.DBPREFIX+"users", login)
		if err == nil {
			if userExists.(int64) == 0 {
				out.Out403(r)
				return
			}
		}
	}

	s, db := db.InitMongo()
	defer s.Close();

	usersC := db.C("userDetails")
	user := users.User{}
	usersC.Find(bson.M{"name":login}).One(&user)
	if (user == users.User{}) || (auth.PassEncrypt(password) != user.Keyword) {
				fmt.Println(auth.PassEncrypt(password))
				fmt.Println(user.Keyword)
		out.Out403(r); return} else {
		sessions.NewSessionValid(r, login)
	}
}

func QuitHandler(r *fasthttp.RequestCtx) {
	sessions.DeleteUserSession(r)
}
