package oauth

import (
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/db"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/globs"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/users"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/sessions"
	"gopkg.in/mgo.v2/bson"
	"github.com/valyala/fasthttp"
)

func OauthLogin(r *fasthttp.RequestCtx, oauthUser *OauthUser) {
	login := oauthUser.Login
	client, err := db.InitRedis()
	if err == nil {
		defer client.Close()
		userExists,err := client.Do("SISMEMBER", globs.DBPREFIX+"users", login)
		if err == nil {
			if userExists.(int64) == 0 {
				sessions.NewSessionValid(r, login)
			}
		}
	}

	s, db := db.InitMongo()
	defer s.Close();

	usersC := db.C("userDetails")
	user := users.User{}
	usersC.Find(bson.M{"name":login}).One(&user)
	if (user == users.User{}) || (user.FullName != oauthUser.FullName) {
		oauthRegister(r, oauthUser)
	} else {
		users.SaveUser(user)
	}
	sessions.NewSessionValid(r, oauthUser.Login)
	r.Redirect("/",302)
}
