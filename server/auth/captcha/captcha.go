package captcha

import (
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/globs"
	"os"
)

var SITEKEY, SECRETKEY string

func initCaptcha() {
	SITEKEY = os.Getenv(globs.PREFIX + "_SITEKEY")
	SECRETKEY = os.Getenv(globs.PREFIX + "_SECRETKEY")
}
