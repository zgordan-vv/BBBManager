package bbbapi

import (
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/app/settings"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/users"
	"crypto/sha1"
	"encoding/hex"
	"github.com/valyala/fasthttp"
)

func ShaHandler(r *fasthttp.RequestCtx) {
	username := users.GetUserName(r)
	user, ok := users.GetUser(username)
	if !ok || !user.IsAdmin {return}
	str := r.FormValue("string")
	r.Write(Sha(append(str, []byte(settings.ConnSettings.Secret)...)))
}

func ShaSecHandler(r *fasthttp.RequestCtx) {
	username := users.GetUserName(r)
	user, ok := users.GetUser(username)
	if !ok || !user.IsAdmin {return}
	str := r.FormValue("string")
	secret := r.FormValue("secret")
	r.Write(Sha(append(str, secret...)))
}

func Sha(str []byte) []byte {
	x := sha1.Sum(str)
	return []byte(hex.EncodeToString(x[0:20]))
}
