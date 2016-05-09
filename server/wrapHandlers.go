package main

import (
	"fmt"
	"github.com/gorilla/context"
	"net/http"
	"time"
	"strings"
)

const STATCOOKIE="4BBBStats"

var bm, cr float32

func mw(fn http.HandlerFunc) http.HandlerFunc {
	return generalWrapper(checkHandler(fn))
}

func mwGuestOk(fn http.HandlerFunc) http.HandlerFunc {
	return generalWrapper(guestOkHandler(fn))
}

func mwInstall(fn http.HandlerFunc) http.HandlerFunc {
	return timeMeasure(recoverHandler(fn))
}

func generalWrapper(fn http.HandlerFunc) http.HandlerFunc {
	return timeMeasure(installCheck(recoverHandler(referrerHandler(fn))))
}

func timeMeasure(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		fn(w, r)
		t2 := time.Now()
		bm += float32(t2.Sub(t1))
		cr++
		//fmt.Println(bm / cr)
	}
}

func installCheck(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if installed {fn(w,r)} else {out(w,"notInstalled")}
	}
}

func recoverHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
				out(w, "recovered")
			}
		}()
		fn(w, r)
	}
}

func referrerHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		host := ""
		hostslice := strings.Split(r.Referer(),"/")
		if len(hostslice) >= 3 {host = hostslice[2]}
		ip := r.Header.Get("X-Real-Ip")
		clientID := ""
		if c, err := r.Cookie(STATCOOKIE); (err != nil) || (c.Value == "") {
			clientID = CW(8)
			c := http.Cookie{Name: STATCOOKIE, Value: clientID, Path: "/", MaxAge: 31104000, HttpOnly: true}
			http.SetCookie(w, &c)
		} else {
			clientID = c.Value
		}
		addStats(host,ip,clientID)
		fn(w,r)
	}
}

func checkHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(VALIDCOOKIE)
		if err != nil {
			out403(w)
		} else {
			sessionID := c.Value
			if sessionID == "" {
				out403(w)
			} else {
				defer context.Clear(r)
				session := loadUserSession(sessionID)
				username := session.Username
				if (username == "") || (time.Now().Unix()>=session.Expires) {
					out403(w)
				} else {
					context.Set(r, "username", username)
					fn(w, r)
				}
			}
		}
	}
}

func guestOkHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(VALIDCOOKIE)
		if err != nil {
			context.Set(r, "username", "")
			fn(w, r)
		} else {
			sessionID := c.Value
			if sessionID == "" {
				context.Set(r, "username", "")
				fn(w, r)
			} else {
				defer context.Clear(r)
				session := loadUserSession(sessionID)
				username := session.Username
				if (username == "") || (time.Now().Unix()>=session.Expires) {
					context.Set(r, "username", "")
					fn(w, r)
				} else {
					context.Set(r, "username", username)
					fn(w, r)
				}
			}
		}
	}
}
