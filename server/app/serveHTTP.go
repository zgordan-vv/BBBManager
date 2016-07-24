package app

import (
	"github.com/valyala/fasthttp"
	"strings"

	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/mw"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/globs"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/check"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/auth/oauth"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/auth/login"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/install"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/users"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/tokens"

	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/app/settings"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/app/maintenance"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/app/meetings"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/app/bbbapi"
)

func requestHandler(ctx *fasthttp.RequestCtx) {
	url := string(ctx.Path())
	switch url {

//common handlers

		case "/api/install": mw.MwInstall(install.InstallHandler)(ctx)
		case "/api/checkInstall": mw.MwInstall(install.RestInstalled)(ctx)
		case "/api/nameUniq": mw.MwGuestOk(install.NameUniqHandler)(ctx)

		case "/api/register": mw.MwGuestOk(login.BrowserRegisterHandler)(ctx)
		case "/api/login": mw.MwGuestOk(login.BrowserLoginHandler)(ctx)
		case "/api/quit": mw.MwGuestOk(login.QuitHandler)(ctx)

		case "/api/getGitHubLoginURL": mw.MwGuestOk(oauth.GetGitHubLoginURLHandler)(ctx)
		case "/api/getFBLoginURL": mw.MwGuestOk(oauth.GetFBLoginURLHandler)(ctx)
		case "/api/getLinkedInLoginURL": mw.MwGuestOk(oauth.GetLinkedInLoginURLHandler)(ctx)

		case "/api/GHCB": mw.MwGuestOk(oauth.GithubCallback)(ctx)
		case "/api/FBCB": mw.MwGuestOk(oauth.FbCallback)(ctx)
		case "/api/INCB": mw.MwGuestOk(oauth.LinkedinCallback)(ctx)

		case "/api/getName": mw.Mw(users.RestGetUserName)(ctx)
		case "/api/getFullName": mw.Mw(users.RestGetUserFullName)(ctx)
		case "/api/getUser": mw.Mw(users.RestGetUser)(ctx)
		case "/api/checkAuth": mw.MwGuestOk(users.CheckAuthHandler)(ctx)
		case "/api/profilesave": mw.Mw(users.ProfileSaveHandler)(ctx)

		case "/api/getSecToken": mw.Mw(tokens.SecTokenHandler)(ctx)
		case "/api/getDupToken": mw.Mw(tokens.DupTokenHandler)(ctx)

		case "/api/regexp": mw.MwInstall(check.RegexpHandler)(ctx)

//app handlers

		case "/api/getsha": mw.Mw(bbbapi.ShaHandler)(ctx)
		case "/api/getsecrsha": mw.Mw(bbbapi.ShaSecHandler)(ctx)

		case "/api/meetings": mw.MwGuestOk(meetings.MeetingsHandler)(ctx)
		case "/api/isRunning": mw.MwGuestOk(meetings.IsRunningHandler)(ctx)
		case "/api/passwords": mw.Mw(meetings.PasswordsHandler)(ctx)
		case "/api/checkPwd": mw.Mw(meetings.CheckPwdHandler)(ctx)
		case "/api/join": mw.Mw(meetings.JoinHandler)(ctx)
		case "/api/meetingUniq": mw.Mw(meetings.MeetingUniqHandler)(ctx)
		case "/api/create": mw.Mw(meetings.CreateMeetingHandler)(ctx)
		case "/api/edit": mw.Mw(meetings.EditMeetingHandler)(ctx)
		case "/api/delete": mw.Mw(meetings.DeleteMeetingHandler)(ctx)

		case "/api/getIP": mw.Mw(settings.GetIPHandler)(ctx)
		case "/api/setIP": mw.Mw(settings.SetIPHandler)(ctx)
		case "/api/getTomcat": mw.Mw(settings.GetTomcatHandler)(ctx)
		case "/api/setTomcat": mw.Mw(settings.SetTomcatHandler)(ctx)
		case "/api/resetTomcat": mw.Mw(settings.ResetTomcatHandler)(ctx)
		case "/api/getFreeswitch": mw.Mw(settings.GetFreeswitchHandler)(ctx)
		case "/api/setFreeswitch": mw.Mw(settings.SetFreeswitchHandler)(ctx)
		case "/api/resetFreeswitch": mw.Mw(settings.ResetFreeswitchHandler)(ctx)
		case "/api/getClient": mw.Mw(settings.GetClientHandler)(ctx)
		case "/api/setClient": mw.Mw(settings.SetClientHandler)(ctx)
		case "/api/resetClient": mw.Mw(settings.ResetClientHandler)(ctx)

		case "/api/getMaintenance": mw.Mw(maintenance.MaintenanceHandler)(ctx)

		default: fasthttp.FSHandler(".."+url, strings.Count(url,"/"))(ctx)

	}
}

func serveHTTP() {
	fasthttp.ListenAndServe("127.0.0.1:"+globs.PORT, requestHandler)
}
