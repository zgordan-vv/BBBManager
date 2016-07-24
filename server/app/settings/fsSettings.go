package settings

import (
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/app/maintenance"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/out"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/tokens"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/users"
	"fmt"
	"encoding/json"
	"os/exec"
	"github.com/valyala/fasthttp"
)

const fsCfgPath string = "/opt/freeswitch/conf/autoload_configs/"
const fsCfgFile string = "conference.conf.xml"
var fsFullPath string = fsCfgPath + fsCfgFile

var fsTmpl = map[string]Param{
	"default": Param{"int", 0, 999},
	"wideband": Param{"int", 0, 999},
	"ultrawideband": Param{"int", 0, 999},
	"cdquality": Param{"int", 0, 999},
	"sla": Param{"int", 0, 999},
}

func GetFreeswitchHandler(r *fasthttp.RequestCtx) {
	username := users.GetUserName(r)
	user, ok := users.GetUser(username)

	if (!ok) || (!user.IsAdmin) {out.Out403(r); return}

	result := jsonArray{}
	params := map[string]string{}
	for profile := range(fsTmpl) {
		params[profile] = getXMLParam(fsFullPath, "configuration/profiles/profile[@name='"+profile+"']/param[@name='energy-level']/@value")
	}
	result.Params = params
	output, _ := json.Marshal(result)
	r.Write(output)
}

func SetFreeswitchHandler(r *fasthttp.RequestCtx) {
	
	if maintenance.GetMaintenance() {out.Out500(r); return}
	
	username := users.GetUserName(r)
	user, ok := users.GetUser(username)

	secToken := r.FormValue("tokensec")
	if (!ok) || (!user.IsAdmin) || (!tokens.CheckSec(secToken,username)) {out.Out403(r); return}

	jsonObj := r.FormValue("settings")
	fsSettings := jsonArray{}
	if err := json.Unmarshal([]byte(jsonObj), &fsSettings); err != nil {fmt.Println(err); out.Out500(r); return}

	params := fsSettings.Params;

	for profile := range(fsTmpl) {
		param := params[profile]
		if !evaluateParam(profile, param, fsTmpl) {return}
		energyLevelPath := "/configuration/profiles/profile[@name='"+profile+"']/param[@name='energy-level']/@value"
		output, err := updateXMLParam(fsFullPath, energyLevelPath, param)
		fmt.Println("output "+string(output))
		if err != nil {out.Out500(r); return}
	}

	restartFreeswitch(r)
}

func ResetFreeswitchHandler(r *fasthttp.RequestCtx) {

	if maintenance.GetMaintenance() {out.Out500(r); return}
	
	username := users.GetUserName(r)
	user, ok := users.GetUser(username)

	secToken := r.FormValue("tokensec")
	if (!ok) || (!user.IsAdmin) || (!tokens.CheckSec(secToken,username)) {out.Out403(r); return}

	maintenance.SetMaintenance(true)
	defer maintenance.SetMaintenance(false)

	err := resetDefaults(r, fsCfgPath, fsCfgFile)
	if err == nil {restartFreeswitch(r)} else {out.Out500(r)}
}

func restartFreeswitch(r *fasthttp.RequestCtx) {

	restart := exec.Command("service", "bbb-freeswitch", "restart")
	restartOutput, err := restart.CombinedOutput()
	fmt.Print("restart output: "); fmt.Println(string(restartOutput))
	fmt.Print("restart error: "); fmt.Println(err)
	if err != nil {
		fmt.Print("Error is ")
		fmt.Println(err)
		out.Out500(r)
	} else {
		fmt.Println("No error")
		out.Outnil(r)
	}
	
}
