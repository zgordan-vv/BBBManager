package initialize

import(
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/auth/oauth"
)

func Init() {
	initPort()
	initGlobs()
	oauth.InitOauth()
}
