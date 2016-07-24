package main

import (
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/app"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/initialize"
)

func main() {
	initialize.Init()
	app.RunApp()
}
