package main

import (
	"fmt"
	"encoding/json"
	"os/exec"
	"github.com/valyala/fasthttp"
)

const tomcatCfgPath string = "/var/lib/tomcat7/webapps/bigbluebutton/WEB-INF/classes/"
const tomcatCfgFile string = "bigbluebutton.properties"
var tomcatFullPath string = tomcatCfgPath + tomcatCfgFile

const sedstr string = AP+"s/^.[^=]*=\\(.*\\)$/\\1/"+AP
const url_prefix = "${bigbluebutton.web.serverURL}/"

var tomcatTmpl = map[string]Param{
	"maxNumPages": Param{"int", 0, 0},
	"defaultMaxUsers": Param{"int", 0, 0},
	"defaultMeetingDuration": Param{"int", 0, 0},
	"defaultMeetingExpireDuration": Param{"int", 0, 0},
	"defaultMeetingCreateJoinDuration": Param{"int", 0, 0},
	"disableRecordingDefault": Param{"bool", 0, 0},
	"allowStartStopRecording": Param{"bool", 0, 0},
}

func getTomcatHandler(r *fasthttp.RequestCtx) {
	username := getUserName(r)
	user, ok := getUser(username)

	if (!ok) || (!user.IsAdmin) {out403(r); return}

	result := jsonArray{}
	params := map[string]string{}
	for key := range(tomcatTmpl) {
		params[key] = getTomcatParam(key)
	}

	result.Params = params
	output, _ := json.Marshal(result)
	r.Write(output)
}

func getTomcatParam(param string) string {
	getTomcat := exec.Command("bash", "-c", "cat "+tomcatFullPath+" | grep "+param+" | sed "+sedstr)
	getTomcatOutput, err := getTomcat.CombinedOutput()
	if err != nil {return ""} else {
		var str string
		if len(getTomcatOutput) > 0 {
			str = string(getTomcatOutput[0:len(getTomcatOutput)-1])
		} else {
			str = string(getTomcatOutput[0:len(getTomcatOutput)])
		}
		return str[0:len(str)]
	}
}

func execPlain (r *fasthttp.RequestCtx, params map[string]string, file string) ([]byte, error) {

	for name, value := range(params) {
		str := "s/^"+name+".*/"+name+"="+value+"/"
		command := exec.Command("sed", "-i", str, file)
		output, err := command.CombinedOutput()
		if (len(output) >0 ) || (err != nil) {
			return output, err
		}
	}
	return nil, nil
}

func setTomcatHandler(r *fasthttp.RequestCtx) {

	if getMaintenance() {out500(r); return}
	
	username := getUserName(r)
	user, ok := getUser(username)

	secToken := r.FormValue("tokensec")
	if (!ok) || (!user.IsAdmin) || (!checkSec(secToken,username)) {out403(r); return}

	jsonObj := r.FormValue("settings")
	tomcatSettings := jsonArray{}
	if err := json.Unmarshal(jsonObj, &tomcatSettings); err != nil {fmt.Println(err); out500(r); return}

	params := tomcatSettings.Params;
	for key, value := range(params) {
		if !evaluateParam(key, value, tomcatTmpl) {out500(r); return}
	}

	output, err := execPlain(r, params, tomcatFullPath)
	fmt.Println(string(output))
	if err != nil {out500(r)} else { restartTomcat(r) }
}

func resetTomcatHandler(r *fasthttp.RequestCtx) {

	setMaintenance(true)
	defer setMaintenance(false)

	err := resetDefaults(r, tomcatCfgPath, tomcatCfgFile)
	if err == nil {restartTomcat(r)} else {out500(r)}
}

func restartTomcat(r *fasthttp.RequestCtx) {
	
	restart := exec.Command("service", "tomcat7", "restart")
	restartOutput, err := restart.CombinedOutput()
	fmt.Print("restart output: "); fmt.Println(string(restartOutput))
	fmt.Print("restart error: "); fmt.Println(err)

	if err != nil {
		fmt.Print("Error is ")
		fmt.Println(err)
		out500(r)
	} else {outnil(r)}
}
