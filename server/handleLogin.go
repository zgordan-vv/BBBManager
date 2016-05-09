package main

import (
	"net/http"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"github.com/garyburd/redigo/redis"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	jsonObj := r.FormValue("login")
	logindata := struct{
		Login string	`json:"login"`
		Password string	`json:"password"`
	}{}
	if json.Unmarshal([]byte(jsonObj), &logindata) != nil {return}
	login := logindata.Login
	password := logindata.Password

	client, err := redis.Dial("tcp", ":6379")
	if err == nil {
		defer client.Close()
		userExists,err := client.Do("SISMEMBER", DBPREFIX+"users", login)
		if err == nil {
			if userExists.(int64) == 0 {
				out403(w)
				return
			}
		}
	}

	s, db := initMongo()
	defer s.Close();

	usersC := db.C("userDetails")
	user := User{}
	usersC.Find(bson.M{"name":login}).One(&user)
	if (user == User{}) || (passEncrypt(password) != user.Keyword) {out403(w); return} else {
		newSessionValid(w, r, login)
	}
}

func quitHandler(w http.ResponseWriter, r *http.Request) {
	deleteUserSession(w, r)
}
