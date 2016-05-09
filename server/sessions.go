package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/context"
	"net/http"
	"strconv"
	"time"
)

const (
	VALIDCOOKIE string = "4BBB_e8h64k9a"
	LASTURLCOOKIE string = "4BBB_g49wp0a6"
)

type UserSession struct {
	Username string
	Expires	int64
}

func newSessionValid(w http.ResponseWriter, r *http.Request, user string) {
	sessionID := CW(32)
	c1 := http.Cookie{Name: VALIDCOOKIE, Value: sessionID, Path: "/", MaxAge: 36000000, HttpOnly: true}
	http.SetCookie(w, &c1)
	saveUserSession(sessionID, user)
}

func deleteUserSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(VALIDCOOKIE)
	if err != nil {	return } else {
		if sessionID := cookie.Value; sessionID == "" { return } else {
			client, err := redis.Dial("tcp", ":6379")
			if err == nil {client.Do("DEL", DBPREFIX+"usersession:"+sessionID)}
			context.Delete(r, "username")
		}
	}
}

func saveUserSession(sessionID, username string) {
	client, err := redis.Dial("tcp", ":6379")
	if err == nil {
		defer client.Close()
		_, err = client.Do("HMSET", DBPREFIX+"usersession:"+sessionID,"Username",username,"Expires",time.Now().Add(time.Hour*2400).Unix())
	}
	if err != nil {
		fmt.Println("Session save has failed: ", err)
	}
}

func loadUserSession(sessionID string) UserSession {
	var userSession UserSession
	client, err := redis.Dial("tcp", ":6379")
	if err == nil {
		defer client.Close()
		dataRaw, err := client.Do("HGETALL", DBPREFIX+"usersession:"+sessionID)
		if err == nil {
			data := dataRaw.([]interface{})
			if len(data) != 0 {
				userSession.Username = string(data[1].([]byte))
				userSession.Expires, err = strconv.ParseInt((string(data[3].([]byte))),10,64)
				if err != nil {userSession.Expires = time.Now().Unix()}
			}
		}
	}
	return userSession
}
