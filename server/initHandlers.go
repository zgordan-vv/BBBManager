package main

import "net/http"

func initHandlers() {

	http.HandleFunc("/api/quit", mwGuestOk(quitHandler))
	http.HandleFunc("/api/register", mwGuestOk(registerHandler))
	http.HandleFunc("/api/profilesave", mw(profileSaveHandler))

	http.HandleFunc("/api/getSecToken", mw(secTokenHandler))
	http.HandleFunc("/api/getDupToken", mw(dupTokenHandler))
	http.HandleFunc("/api/getsha", mw(shaHandler))
	http.HandleFunc("/api/getsecrsha", mw(shaSecHandler))

	http.HandleFunc("/api/install", mwInstall(installHandler))
	http.HandleFunc("/api/checkInstall", mwInstall(restInstalled))
	http.HandleFunc("/api/nameUniq", mwGuestOk(nameUniqHandler))

	http.HandleFunc("/api/login", mwGuestOk(loginHandler))

	http.HandleFunc("/api/checkAuth", mwGuestOk(checkAuthHandler))
	http.HandleFunc("/api/getName", mw(restGetUserName))
	http.HandleFunc("/api/getFullName", mw(restGetUserFullName))
	http.HandleFunc("/api/getUser", mw(restGetUser))

	http.HandleFunc("/api/meetings", mwGuestOk(meetingsHandler))
	http.HandleFunc("/api/toggleRunning", mw(toggleRunningHandler))
	http.HandleFunc("/api/passwords", mw(passwordsHandler))
	http.HandleFunc("/api/meetingUniq", mw(meetingUniqHandler))
	http.HandleFunc("/api/isRunning", mwGuestOk(isRunningHandler))
	http.HandleFunc("/api/create", mw(createMeetingHandler))
	http.HandleFunc("/api/edit", mw(editMeetingHandler))
	http.HandleFunc("/api/delete", mw(deleteMeetingHandler))
	http.HandleFunc("/api/checkPwd", mw(checkPwdHandler))
	http.HandleFunc("/api/join", mw(joinHandler))

	http.HandleFunc("/api/regexp", mwInstall(regexpHandler))
	http.HandleFunc("/api/getSettings", mw(getSettingsHandler))
	http.HandleFunc("/api/setSettings", mw(setSettingsHandler))
	http.HandleFunc("/api/getTomcat", mw(getTomcatHandler))
	http.HandleFunc("/api/setTomcat", mw(setTomcatHandler))
	http.HandleFunc("/api/getFreeswitch", mw(getFreeswitchHandler))
	http.HandleFunc("/api/setFreeswitch", mw(setFreeswitchHandler))

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("../"))))
	http.ListenAndServe(":"+PORT, nil)
}
