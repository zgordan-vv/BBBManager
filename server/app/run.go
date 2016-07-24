package app

import (
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/app/settings"
)

func RunApp() {
	settings.InitSettings()
	serveHTTP()
}
