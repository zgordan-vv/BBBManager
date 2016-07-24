package settings

import (
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/app/maintenance"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/out"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/check"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/tokens"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/users"
	"fmt"
	"encoding/json"
	"os/exec"
	"github.com/valyala/fasthttp"
)

const tomcatCfgPath string = "/var/lib/tomcat7/webapps/bigbluebutton/WEB-INF/classes/"
const tomcatCfgFile string = "bigbluebutton.properties"
var tomcatFullPath string = tomcatCfgPath + tomcatCfgFile

const sedstr string = check.AP+"s/^.[^=]*=\\(.*\\)$/\\1/"+check.AP

var tomcatTmpl = map[string]Param{
	"maxNumPages": Param{"int", 0, 0},
	"defaultMaxUsers": Param{"int", 0, 0},
	"defaultMeetingDuration": Param{"int", 0, 0},
	"defaultMeetingExpireDuration": Param{"int", 0, 0},
	"defaultMeetingCreateJoinDuration": Param{"int", 0, 0},
	"disableRecordingDefault": Param{"bool", 0, 0},
	"allowStartStopRecording": Param{"bool", 0, 0},
	"securitySalt": Param{"string", 0, 0},
}

func GetTomcatHandler(r *fasthttp.RequestCtx) {
	username := users.GetUserName(r)
	user, ok := users.GetUser(username)

	if (!ok) || (!user.IsAdmin) {out.Out403(r); return}

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

func SetTomcatHandler(r *fasthttp.RequestCtx) {

	if maintenance.GetMaintenance() {out.Out500(r); return}
	
	username := users.GetUserName(r)
	user, ok := users.GetUser(username)

	secToken := r.FormValue("tokensec")
	if (!ok) || (!user.IsAdmin) || (!tokens.CheckSec(secToken,username)) {out.Out403(r); return}

	jsonObj := r.FormValue("settings")
	tomcatSettings := jsonArray{}
	if err := json.Unmarshal(jsonObj, &tomcatSettings); err != nil {fmt.Println(err); out.Out500(r); return}

	params := tomcatSettings.Params;
	for key, value := range(params) {
		if !evaluateParam(key, value, tomcatTmpl) {out.Out500(r); return}
	}

	output, err := execPlain(r, params, tomcatFullPath)
	fmt.Println(string(output))
	if err != nil {out.Out500(r)} else {
		err = setSalt(params["securitySalt"])
		if err != nil {out.Out500(r)} else {restartTomcat(r); out.Outnil(r)}
	}
}

func ResetTomcatHandler(r *fasthttp.RequestCtx) {

	if maintenance.GetMaintenance() {out.Out500(r); return}
	
	username := users.GetUserName(r)
	user, ok := users.GetUser(username)

	secToken := r.FormValue("tokensec")
	if (!ok) || (!user.IsAdmin) || (!tokens.CheckSec(secToken,username)) {out.Out403(r); return}

	maintenance.SetMaintenance(true)
	defer maintenance.SetMaintenance(false)

	err := resetDefaults(r, tomcatCfgPath, tomcatCfgFile)
	if err == nil {restartTomcat(r)} else {out.Out500(r)}
}

func restartTomcat(r *fasthttp.RequestCtx) {
	
	restart := exec.Command("service", "tomcat7", "restart")
	restartOutput, err := restart.CombinedOutput()
	fmt.Print("restart output: "); fmt.Println(string(restartOutput))
	fmt.Print("restart error: "); fmt.Println(err)

	if err != nil {
		fmt.Print("Error is ")
		fmt.Println(err)
		out.Out500(r)
	} else {out.Outnil(r)}
}
