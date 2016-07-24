package mw

import (
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/install"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/out"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/sessions"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/stats"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/utils"
	"fmt"
	"github.com/valyala/fasthttp"
	"time"
	"strings"
)

const STATCOOKIE="4BBBStats"

var bm, cr float32

func Mw(fn fasthttp.RequestHandler) fasthttp.RequestHandler {
	return generalWrapper(checkHandler(fn))
}

func MwGuestOk(fn fasthttp.RequestHandler) fasthttp.RequestHandler {
	return generalWrapper(guestOkHandler(fn))
}

func MwInstall(fn fasthttp.RequestHandler) fasthttp.RequestHandler {
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
		if install.Installed {fn(r)} else {out.Out(r,"notInstalled")}
	}
}

func recoverHandler(fn fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(r *fasthttp.RequestCtx) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
				out.Out(r, "recovered")
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
			clientID = string(utils.CW(8))
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
		stats.AddStats(host,string(ip),clientID)
		fn(r)
	}
}

func checkHandler(fn fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(r *fasthttp.RequestCtx) {
		sessionID := r.Request.Header.Cookie(sessions.VALIDCOOKIE)
		if sessionID == nil {
			out.Out403(r)
		} else {
			session := sessions.LoadUserSession(string(sessionID))
			username := session.Username
			if (username == "") || (time.Now().Unix()>=session.Expires) {
				r.SetUserValue("username", "")
				out.Out403(r)
			} else {
				r.SetUserValue("username", username)
				fn(r)
			}
		}
	}
}

func guestOkHandler(fn fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(r *fasthttp.RequestCtx) {
		sessionID := r.Request.Header.Cookie(sessions.VALIDCOOKIE)
		if sessionID == nil {
			r.SetUserValue("username", "")
			fn(r)
		} else {
			session := sessions.LoadUserSession(string(sessionID))
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
