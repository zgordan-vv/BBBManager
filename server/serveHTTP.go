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

		case "/api/register": mwGuestOk(browserRegisterHandler)(ctx)
		case "/api/nameUniq": mwGuestOk(nameUniqHandler)(ctx)
		case "/api/profilesave": mw(profileSaveHandler)(ctx)

		case "/api/login": mwGuestOk(browserLoginHandler)(ctx)
		case "/api/quit": mwGuestOk(quitHandler)(ctx)
		case "/api/checkAuth": mwGuestOk(checkAuthHandler)(ctx)

		case "/api/getGitHubLoginURL": mwGuestOk(getGitHubLoginURLHandler)(ctx)
		case "/api/getFBLoginURL": mwGuestOk(getFBLoginURLHandler)(ctx)
		case "/api/getLinkedInLoginURL": mwGuestOk(getLinkedInLoginURLHandler)(ctx)

		case "/api/GHCB": mwGuestOk(githubCallback)(ctx)
		case "/api/FBCB": mwGuestOk(fbCallback)(ctx)
		case "/api/INCB": mwGuestOk(linkedinCallback)(ctx)

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
		case "/api/resetTomcat": mw(resetTomcatHandler)(ctx)
		case "/api/getFreeswitch": mw(getFreeswitchHandler)(ctx)
		case "/api/setFreeswitch": mw(setFreeswitchHandler)(ctx)
		case "/api/resetFreeswitch": mw(resetFreeswitchHandler)(ctx)
		case "/api/getClient": mw(getClientHandler)(ctx)
		case "/api/setClient": mw(setClientHandler)(ctx)
		case "/api/resetClient": mw(resetClientHandler)(ctx)

		case "/api/getMaintenance": mw(maintenanceHandler)(ctx)

		default: fasthttp.FSHandler(".."+url, strings.Count(url,"/"))(ctx)

	}
}

func serveHTTP() {
	fasthttp.ListenAndServe("127.0.0.1:"+PORT, requestHandler)
}
