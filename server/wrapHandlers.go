package main

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"time"
	"strings"
)

const STATCOOKIE="4BBBStats"

var bm, cr float32

func mw(fn fasthttp.RequestHandler) fasthttp.RequestHandler {
	return generalWrapper(checkHandler(fn))
}

func mwGuestOk(fn fasthttp.RequestHandler) fasthttp.RequestHandler {
	return generalWrapper(guestOkHandler(fn))
}

func mwInstall(fn fasthttp.RequestHandler) fasthttp.RequestHandler {
	return timeMeasure(recoverHandler(fn))
}

func generalWrapper(fn fasthttp.RequestHandler) fasthttp.RequestHandler {
	return timeMeasure(installCheck(recoverHandler(referrerHandler(fn))))
}

func timeMeasure(fn fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(r *fasthttp.RequestCtx) {
		t1 := time.Now()
		fn(r)
		t2 := time.Now()
		bm += float32(t2.Sub(t1))
		cr++
//		fmt.Println(bm / cr)
	}
}

func installCheck(fn fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(r *fasthttp.RequestCtx) {
		if installed {fn(r)} else {out(r,"notInstalled")}
	}
}

func recoverHandler(fn fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(r *fasthttp.RequestCtx) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
				out(r, "recovered")
			}
		}()
		fn(r)
	}
}

func referrerHandler(fn fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(r *fasthttp.RequestCtx) {
		host := ""
		hostslice := strings.Split(string(r.Referer()),"/")
		if len(hostslice) >= 3 {host = hostslice[2]}
		ip := r.Request.Header.Peek("X-Real-Ip")
		clientID := ""
		if statCookie := r.Request.Header.Cookie(STATCOOKIE); statCookie == nil {
			clientID = string(CW(8))
			c := fasthttp.Cookie{}
			c.SetKey(STATCOOKIE)
			c.SetValue(clientID)
			c.SetPath("/")
			c.SetExpire(time.Now().Add(time.Second*31104000))
			c.SetHTTPOnly(true)
			r.Response.Header.SetCookie(&c)
		} else {
			clientID = string(statCookie)
		}
		addStats(host,string(ip),clientID)
		fn(r)
	}
}

func checkHandler(fn fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(r *fasthttp.RequestCtx) {
		sessionID := r.Request.Header.Cookie(VALIDCOOKIE)
		if sessionID == nil {
			out403(r)
		} else {
			session := loadUserSession(string(sessionID))
			username := session.Username
			if (username == "") || (time.Now().Unix()>=session.Expires) {
				out403(r)
			} else {
				checkIfLogged(r, username)
				r.SetUserValue("username", username)
				fn(r)
			}
		}
	}
}

func checkIfLogged(r *fasthttp.RequestCtx, username string) {
	if strings.HasPrefix(username, "FB_") {fmt.Println("Facebook guy!")}
	if strings.HasPrefix(username, "GH_") {fmt.Println("GitHub guy!")}
}

func guestOkHandler(fn fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(r *fasthttp.RequestCtx) {
		sessionID := r.Request.Header.Cookie(VALIDCOOKIE)
		if sessionID == nil {
			r.SetUserValue("username", "")
			fn(r)
		} else {
			session := loadUserSession(string(sessionID))
			username := session.Username
			if (username == "") || (time.Now().Unix()>=session.Expires) {
				r.SetUserValue("username", "")
				fn(r)
			} else {
				r.SetUserValue("username", username)
				fn(r)
			}
		}
	}
}
