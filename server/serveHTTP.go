package main

import (
"strings"
"github.com/valyala/fasthttp"
)

func requestHandler(ctx *fasthttp.RequestCtx) {
	url := string(ctx.Path())
	switch url {

		case "/api/install": mwInstall(installHandler)(ctx)
		case "/api/checkInstall": mwInstall(restInstalled)(ctx)

		case "/api/register": mwGuestOk(registerHandler)(ctx)
		case "/api/nameUniq": mwGuestOk(nameUniqHandler)(ctx)
		case "/api/profilesave": mw(profileSaveHandler)(ctx)

		case "/api/login": mwGuestOk(loginHandler)(ctx)
		case "/api/quit": mwGuestOk(quitHandler)(ctx)
		case "/api/checkAuth": mwGuestOk(checkAuthHandler)(ctx)

		case "/api/getName": mw(restGetUserName)(ctx)
		case "/api/getFullName": mw(restGetUserFullName)(ctx)
		case "/api/getUser": mw(restGetUser)(ctx)

		case "/api/getSecToken": mw(secTokenHandler)(ctx)
		case "/api/getDupToken": mw(dupTokenHandler)(ctx)

		case "/api/getsha": mw(shaHandler)(ctx)
		case "/api/getsecrsha": mw(shaSecHandler)(ctx)

		case "/api/meetings": mwGuestOk(meetingsHandler)(ctx)
		case "/api/isRunning": mwGuestOk(isRunningHandler)(ctx)
		case "/api/passwords": mw(passwordsHandler)(ctx)
		case "/api/checkPwd": mw(checkPwdHandler)(ctx)
		case "/api/join": mw(joinHandler)(ctx)
		case "/api/meetingUniq": mw(meetingUniqHandler)(ctx)
		case "/api/create": mw(createMeetingHandler)(ctx)
		case "/api/edit": mw(editMeetingHandler)(ctx)
		case "/api/delete": mw(deleteMeetingHandler)(ctx)

		case "/api/regexp": mwInstall(regexpHandler)(ctx)

		case "/api/getSettings": mw(getSettingsHandler)(ctx)
		case "/api/setSettings": mw(setSettingsHandler)(ctx)
		case "/api/getTomcat": mw(getTomcatHandler)(ctx)
		case "/api/setTomcat": mw(setTomcatHandler)(ctx)
		case "/api/getFreeswitch": mw(getFreeswitchHandler)(ctx)
		case "/api/setFreeswitch": mw(setFreeswitchHandler)(ctx)

		default: fasthttp.FSHandler(".."+url, strings.Count(url,"/"))(ctx)

	}
}

func serveHTTP() {
	fasthttp.ListenAndServe("127.0.0.1:"+PORT, requestHandler)
}
