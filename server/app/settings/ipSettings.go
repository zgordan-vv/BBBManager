package settings

import (
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/app/maintenance"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/out"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/check"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/tokens"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/users"
	"fmt"
	"os/exec"
	"github.com/valyala/fasthttp"
)

func GetIPHandler(r *fasthttp.RequestCtx) {

	if maintenance.GetMaintenance() {out.Outnil(r); return}
	
	username := users.GetUserName(r)
	user, ok := users.GetUser(username)

	if (!ok) || (!user.IsAdmin) {out.Outnil(r); return}

	out.Out(r, ConnSettings.IP)
	
}

func SetIPHandler(r *fasthttp.RequestCtx) {

	if maintenance.GetMaintenance() {out.Out500(r); return}
	
	username := users.GetUserName(r)
	user, ok := users.GetUser(username)

	secToken := r.FormValue("tokensec")
	if (!ok) || (!user.IsAdmin) || (!tokens.CheckSec(secToken,username)) {out.Out403(r); return}

	ip := string(r.FormValue("ip"))
	if !check.Validate(check.IP, ip) { out.Out500(r);return } else {
		restart := exec.Command("bbb-conf", "--setip", ip)
		maintenance.SetMaintenance(true)
		output, err := restart.CombinedOutput()
		fmt.Println(string(output))
		maintenance.SetMaintenance(false)
		if err != nil { fmt.Println(err); out.Out500(r) } else { ConnSettings.IP = ip; out.Outnil(r) }
	}
}
