package main

import (
	"github.com/valyala/fasthttp"
	"gopkg.in/mgo.v2/bson"
	"github.com/garyburd/redigo/redis"
	"encoding/json"
)

type User struct {
	Name string	`json:"name"`
	FullName string	`json:"fullname"`
	IsAdmin bool	`json:"isadmin"`
	Keyword string
}

func getUserName(r *fasthttp.RequestCtx) string {
	username := r.UserValue("username")
	if username != nil {return username.(string)} else {return ""}
}

func restGetUserName(r *fasthttp.RequestCtx) {
	out(r, getUserName(r))
}

func restGetUserFullName(r *fasthttp.RequestCtx) {
	fullName := ""
	username := getUserName(r)
	user, ok := getUser(username)
	if ok {fullName = user.FullName}
	out(r, fullName)
}

func restGetUser(r *fasthttp.RequestCtx) {
	username := getUserName(r)
	user, _ := getUser(username)
	out, err := json.Marshal(user)
	if err == nil {r.Write(out)}
}

func checkAuthHandler(r *fasthttp.RequestCtx) {
	output := "guest"
	username := getUserName(r)
	user, ok := getUser(username)
	if ok {
		if username != "" {
			if user.IsAdmin {output = "admin"} else {output = "user"}
		}
	}	
	out(r, output)
}

func profileSaveHandler(r *fasthttp.RequestCtx) {
	username := getUserName(r)
	if username == "" { out(r, "NotSaved"); return }
	user, ok := getUser(username)
	if !ok { out(r, "NotSaved"); return }
	fullname := string(r.FormValue("fullname"))
	if !checkLogin(fullname) { out(r, "WrongFullName") } else {
		oldpwd := string(r.FormValue("oldpwd"))
		pwd := string(r.FormValue("pwd"))
		if changed, correct := checkPwdChange(username, oldpwd, pwd); !correct {out403(r)} else {
			user.FullName = fullname
			if changed {user.Keyword = passEncrypt(pwd)}
			if err := saveUser(user); err != nil {out(r, "DontSave")}
			out(r, "ok")
		}
	}
}

func checkPwdChange(username,old,new string) (bool, bool) {
	if (old == "") && (new == "") {return false, true} else {
		if !checkPwd(username,old) {return true, false} else {return true,true}
	}
}

func nameUniq(name string) bool {

	client, err := redis.Dial("tcp", ":6379")
	if err == nil {
		defer client.Close()
		userExists,err := client.Do("SISMEMBER", DBPREFIX+"users", name)
		if err == nil {
			if userExists.(int64) != 0 {
				return false
			}
		}
	}

	s, db := initMongo()
	defer s.Close();

	usersC := db.C("userDetails")
	userDetails := User{}
	usersC.Find(bson.M{"Name":name}).One(&userDetails)
	return userDetails == User{}

}

func saveUser(user User) error {

	s, db := initMongo()
	defer s.Close();

	userDetailsC := db.C("userDetails")

	username := user.Name

	client, err := redis.Dial("tcp", ":6379")
	if err == nil {
		defer client.Close()
		_, err = client.Do("SADD", DBPREFIX+"users",username)
		if err != nil {return err}
	} else {return err}

	_, err = userDetailsC.Upsert(bson.M{"name":username}, user)
	return err
}

func getUser(username string) (User, bool) {

	client, err := redis.Dial("tcp", ":6379")
	if err == nil {
		defer client.Close()
		userExists,err := client.Do("SISMEMBER", DBPREFIX+"users", username)
		if err == nil {
			if userExists.(int64) == 1 {
				s, db := initMongo()
				defer s.Close();

				userDetailsC := db.C("userDetails")

				user := User{}
				err = userDetailsC.Find(bson.M{"name":username}).One(&user)
				if err == nil {return user, true}
			}
		}
	}
	return User{}, false
}
