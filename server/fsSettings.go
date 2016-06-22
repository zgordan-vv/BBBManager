package main

import (
	"fmt"
	"encoding/json"
	"os/exec"
	"github.com/valyala/fasthttp"
)

const fsCfgFile string = "/opt/freeswitch/conf/autoload_configs/conference.conf.xml"

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
		params[profile] = getXMLParam(fsCfgFile, "configuration/profiles/profile[@name='"+profile+"']/param[@name='energy-level']/@value")
	}
	result.Params = params
	output, _ := json.Marshal(result)
	r.Write(output)
}

func getXMLParam(file, queryString string) string {
	xmlQuery := exec.Command("bash", "-c", "xmlstarlet sel -t -v \""+queryString+"\" "+file)
	output, err := xmlQuery.CombinedOutput()
	if err != nil {return ""} else {
		return string(output)
	}
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
		output, err := updateXMLParam(fsCfgFile, energyLevelPath, param)
		fmt.Println("output "+string(output))
		if err != nil {out500(r); return}
	}

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

func updateXMLParam(file, param, value string) ([]byte, error) {
	str := "xmlstarlet ed -u \""+param+"\" -v "+value+" "+file+" > /tmp/bbbfstmp.xml && mv -f /tmp/bbbfstmp.xml "+file
	updateParam := exec.Command("bash", "-c", str)
	output, err := updateParam.CombinedOutput()
	fmt.Printf("Output is %s, err is %s\n", output, err)
	return output, err
}
