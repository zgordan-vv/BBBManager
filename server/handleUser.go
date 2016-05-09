package main

import (
	"github.com/gorilla/context"
	"net/http"
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

func getUserName(r *http.Request) string {
	username, ok := context.GetOk(r,"username")
	if ok { return username.(string) } else { return ""}
}

func restGetUserName(w http.ResponseWriter, r *http.Request) {
	out(w, getUserName(r))
}

func restGetUserFullName(w http.ResponseWriter, r *http.Request) {
	fullName := ""
	username := getUserName(r)
	user, ok := getUser(username)
	if ok {fullName = user.FullName}
	out(w, fullName)
}

func restGetUser(w http.ResponseWriter, r *http.Request) {
	username := getUserName(r)
	user, _ := getUser(username)
	out, err := json.Marshal(user)
	if err == nil {w.Write(out)}
}

func checkAuthHandler(w http.ResponseWriter, r *http.Request) {
	output := "guest"
	username := getUserName(r)
	user, ok := getUser(username)
	if ok {
		if username != "" {
			if user.IsAdmin {output = "admin"} else {output = "user"}
		}
	}	
	out(w, output)
}

func profileSaveHandler(w http.ResponseWriter, r *http.Request) {
	username := getUserName(r)
	if username == "" { out(w, "NotSaved"); return }
	user, ok := getUser(username)
	if !ok { out(w, "NotSaved"); return }
	fullname := r.FormValue("fullname")
	if !checkLogin(fullname) { out(w, "WrongFullName") } else {
		oldpwd := r.FormValue("oldpwd")
		pwd := r.FormValue("pwd")
		if changed, correct := checkPwdChange(username, oldpwd, pwd); !correct {out403(w)} else {
			user.FullName = fullname
			if changed {user.Keyword = passEncrypt(pwd)}
			if err := saveUser(user); err != nil {out(w, "DontSave")}
			out(w, "ok")
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
