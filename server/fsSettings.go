package main

import (
	"fmt"
	"encoding/json"
	"os/exec"
	"github.com/valyala/fasthttp"
)

const fsCfgFile string = "/opt/freeswitch/conf/autoload_configs/conference.conf.xml"
const xmlParam = "configuration/profiles/profile[@name='default']/param[@name='energy-level']/@value"


var fsTmpl = map[string]Param{
	"energyLevel": Param{"int", 0, 999},
}

func getFreeswitchHandler(r *fasthttp.RequestCtx) {
	username := getUserName(r)
	user, ok := getUser(username)

	if (!ok) || (!user.IsAdmin) {out403(r); return}

	result := jsonArray{}
	params := map[string]string{}
	params["energyLevel"] = getFSParam()
	result.Params = params
	output, _ := json.Marshal(result)
	r.Write(output)
}

func getFSParam() string {
	xmlQuery := exec.Command("bash", "-c", "xmlstarlet sel -t -v \""+xmlParam+"\" "+fsCfgFile)
	getFSOutput, err := xmlQuery.CombinedOutput()
	if err != nil {return ""} else {
		return string(getFSOutput)
	}
}

func setFreeswitchHandler(r *fasthttp.RequestCtx) {
	
	username := getUserName(r)
	user, ok := getUser(username)

	secToken := r.FormValue("tokensec")
	if (!ok) || (!user.IsAdmin) || (!checkSec(secToken,username)) {out403(r); return}

	jsonObj := r.FormValue("settings")
	fsSettings := jsonArray{}
	if err := json.Unmarshal([]byte(jsonObj), &fsSettings); err != nil {fmt.Println(err); out500(r); return}

	params := fsSettings.Params;
	if !evaluateParams(params, fsTmpl) {out500(r); return}

	path := "/opt/freeswitch/conf/autoload_configs/"
	file := "conference.conf.xml"
	energyLevelPath := "/configuration/profiles/profile/param[@name='energy-level']/@value"

	output, err := updateXMLParam(path+file, energyLevelPath, params["energyLevel"])
	fmt.Println("output "+string(output))
	if err != nil {out500(r); return}

	restart := exec.Command("service", "bbb-freeswitch", "restart")
	restartOutput, err := restart.CombinedOutput()
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
	fmt.Printf("String is %s\n", str)
	updateParam := exec.Command("bash", "-c", str)
	output, err := updateParam.CombinedOutput()
	fmt.Printf("Output is %s, err is %s\n", output, err)
	return output, err
}
