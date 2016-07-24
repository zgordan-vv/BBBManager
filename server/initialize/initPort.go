package initialize

import (
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/globs"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/utils"
	"encoding/json"
	"os"
)

type PortState struct {
	Port string
}

func initPort() {
	var loadstate PortState
	file, err := os.Open("conf/PORT")
	defer file.Close()
	if err != nil {
		return
	}
	decoder := json.NewDecoder(file)
	utils.Check(decoder.Decode(&loadstate))
	globs.PORT = loadstate.Port
	return
}
