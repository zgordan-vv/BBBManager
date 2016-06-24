package main

import (
	"fmt"
	"os/exec"
	"github.com/valyala/fasthttp"
)

func getIPHandler(r *fasthttp.RequestCtx) {

	if getMaintenance() {outnil(r); return}
	
	username := getUserName(r)
	user, ok := getUser(username)

	if (!ok) || (!user.IsAdmin) {outnil(r); return}

	out(r, settings.IP)
	
}

func setIPHandler(r *fasthttp.RequestCtx) {

	if getMaintenance() {out500(r); return}
	
	username := getUserName(r)
	user, ok := getUser(username)

	secToken := r.FormValue("tokensec")
	if (!ok) || (!user.IsAdmin) || (!checkSec(secToken,username)) {out403(r); return}

	ip := string(r.FormValue("ip"))
	if !validate(IP, ip) { out500(r);return } else {
		restart := exec.Command("bbb-conf", "--setip", ip)
		setMaintenance(true)
		output, err := restart.CombinedOutput()
		fmt.Println(string(output))
		setMaintenance(false)
		if err != nil { fmt.Println(err); out500(r) } else { settings.IP = ip; outnil(r) }
	}
}
