package sessions

import (
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/globs"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/utils"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/valyala/fasthttp"
	"strings"
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

func NewSessionValid(r *fasthttp.RequestCtx, user string) {
	sessionID := utils.CW(32)
	c1 := fasthttp.Cookie{}
	c1.SetKey(VALIDCOOKIE)
	c1.SetValueBytes(sessionID)
	c1.SetPath("/")
	c1.SetExpire(time.Now().Add(time.Second*36000000))
	c1.SetHTTPOnly(true)
	r.Response.Header.SetCookie(&c1)
	SaveUserSession(string(sessionID), user)
}

func DeleteUserSession(r *fasthttp.RequestCtx) {
	cookie := r.Request.Header.Cookie(VALIDCOOKIE)
	if cookie == nil { return } else {
		client, err := redis.Dial("tcp", ":6379")
		if err == nil {client.Do("DEL", globs.DBPREFIX+"usersession:"+string(cookie))}
		r.SetUserValue("username", "")
	}
}

func SaveUserSession(sessionID, username string) {
	var duration time.Duration
	client, err := redis.Dial("tcp", ":6379")
	if err == nil {
		defer client.Close()
		if strings.HasPrefix(username, "<OAUTH>") {duration = time.Hour*24} else {duration = time.Hour*2400}
		_, err = client.Do("HMSET", globs.DBPREFIX+"usersession:"+sessionID,"Username",username,"Expires",time.Now().Add(duration).Unix())
	}
	if err != nil {
		fmt.Println("Session save has failed: ", err)
	}
}

func LoadUserSession(sessionID string) UserSession {
	var userSession UserSession
	client, err := redis.Dial("tcp", ":6379")
	if err == nil {
		defer client.Close()
		dataRaw, err := client.Do("HGETALL", globs.DBPREFIX+"usersession:"+sessionID)
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
