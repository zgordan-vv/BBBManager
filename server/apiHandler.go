package main

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/valyala/fasthttp"
)

func shaHandler(r *fasthttp.RequestCtx) {
	username := getUserName(r)
	user, ok := getUser(username)
	if !ok || !user.IsAdmin {return}
	str := r.FormValue("string")
	r.Write(sha(append(str, []byte(settings.Secret)...)))
}

func shaSecHandler(r *fasthttp.RequestCtx) {
	username := getUserName(r)
	user, ok := getUser(username)
	if !ok || !user.IsAdmin {return}
	str := r.FormValue("string")
	secret := r.FormValue("secret")
	r.Write(sha(append(str, secret...)))
}

func sha(str []byte) []byte {
	x := sha1.Sum(str)
	return []byte(hex.EncodeToString(x[0:20]))
}

func joinHandler(r *fasthttp.RequestCtx) {
	username := getUserName(r)
	_, ok := getUser(username)
	if !ok {return}
	qstr := r.FormValue("string")
	checksum := sha(appendAll([][]byte{[]byte("join"),qstr,[]byte(settings.Secret)}))
	r.Write(appendAll([][]byte{[]byte("join?"),qstr,[]byte("&checksum="),checksum}))
}
