package main

import (
	"crypto/sha1"
	"encoding/hex"
	"net/http"
)

func shaHandler(w http.ResponseWriter, r *http.Request) {
	username := getUserName(r)
	user, ok := getUser(username)
	if !ok || !user.IsAdmin {return}
	str := r.FormValue("string")
	w.Write(sha(str + settings.Secret))
}

func shaSecHandler(w http.ResponseWriter, r *http.Request) {
	username := getUserName(r)
	user, ok := getUser(username)
	if !ok || !user.IsAdmin {return}
	str := r.FormValue("string")
	secret := r.FormValue("secret")
	w.Write(sha(str + secret))
}

func sha(str string) []byte {
	x := sha1.Sum([]byte(str))
	return []byte(hex.EncodeToString(x[0:20]))
}

func joinHandler(w http.ResponseWriter, r *http.Request) {
	username := getUserName(r)
	_, ok := getUser(username)
	if !ok {return}
	qstr := r.FormValue("string")
	checksum := sha("join"+qstr + settings.Secret)
	w.Write(append([]byte( "join"+"?"+qstr+"&checksum="),checksum...))
}
