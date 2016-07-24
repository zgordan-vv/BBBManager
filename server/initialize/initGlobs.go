package initialize

import (
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/check"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/install"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/utils"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/globs"
	"encoding/json"
	"os"
)

type State struct {
	DB string
	DN string
	PREFIX string
}

func initGlobs() {
	var loadstate State
	file, err := os.Open("conf/GLOBS")
	defer file.Close()

	if err != nil {
		return
	}

	decoder := json.NewDecoder(file)
	utils.Check(decoder.Decode(&loadstate))
	globs.DBPREFIX = loadstate.DB
	globs.DOMAINNAME = loadstate.DN
	globs.PREFIX = loadstate.PREFIX
	if !check.CheckDBPrefix(globs.DBPREFIX) || !check.CheckDomainName(globs.DOMAINNAME) {return}
	globs.DBPREFIX += ":"
	install.CheckInstall()
}
