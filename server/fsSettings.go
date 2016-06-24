package main

import (
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

func getFreeswitchHandler(r *fasthttp.RequestCtx) {
	username := getUserName(r)
	user, ok := getUser(username)

	if (!ok) || (!user.IsAdmin) {out403(r); return}

	result := jsonArray{}
	params := map[string]string{}
	for profile := range(fsTmpl) {
		params[profile] = getXMLParam(fsFullPath, "configuration/profiles/profile[@name='"+profile+"']/param[@name='energy-level']/@value")
	}
	result.Params = params
	output, _ := json.Marshal(result)
	r.Write(output)
}

func setFreeswitchHandler(r *fasthttp.RequestCtx) {
	
	if getMaintenance() {out500(r); return}
	
	username := getUserName(r)
	user, ok := getUser(username)

	secToken := r.FormValue("tokensec")
	if (!ok) || (!user.IsAdmin) || (!checkSec(secToken,username)) {out403(r); return}

	jsonObj := r.FormValue("settings")
	fsSettings := jsonArray{}
	if err := json.Unmarshal([]byte(jsonObj), &fsSettings); err != nil {fmt.Println(err); out500(r); return}

	params := fsSettings.Params;

	for profile := range(fsTmpl) {
		param := params[profile]
		if !evaluateParam(profile, param, fsTmpl) {return}
		energyLevelPath := "/configuration/profiles/profile[@name='"+profile+"']/param[@name='energy-level']/@value"
		output, err := updateXMLParam(fsFullPath, energyLevelPath, param)
		fmt.Println("output "+string(output))
		if err != nil {out500(r); return}
	}

	restartFreeswitch(r)
}

func resetFreeswitchHandler(r *fasthttp.RequestCtx) {

	setMaintenance(true)
	defer setMaintenance(false)

	err := resetDefaults(r, fsCfgPath, fsCfgFile)
	if err == nil {restartFreeswitch(r)} else {out500(r)}
}

func restartFreeswitch(r *fasthttp.RequestCtx) {

	restart := exec.Command("service", "bbb-freeswitch", "restart")
	setMaintenance(true)
	restartOutput, err := restart.CombinedOutput()
	setMaintenance(false)
	fmt.Print("restart output: "); fmt.Println(string(restartOutput))
	fmt.Print("restart error: "); fmt.Println(err)
	if err != nil {
		fmt.Print("Error is ")
		fmt.Println(err)
		out500(r)
	} else {outnil(r)}
	
}
