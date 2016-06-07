package main

import (
	"fmt"
	"encoding/json"
	"os/exec"
	"github.com/valyala/fasthttp"
	"strconv"
	"strings"
)

type jsonArray struct {
	Params map[string]string `json: "params"`
}

type Param struct {
	tip string
	min int
	max int
}

const tomcatCfgFile string = "/var/lib/tomcat7/webapps/bigbluebutton/WEB-INF/classes/bigbluebutton.properties"
const sedstr string = AP+"s/^.[^=]*=\\(.*\\)$/\\1/"+AP
const url_prefix = "${bigbluebutton.web.serverURL}/"

var tomcatTmpl = map[string]Param{
	"maxNumPages": Param{"int", 0, 0},
//	"defaultWelcomeMessage": Param{"string", 0, 0},
//	"defaultWelcomeMessageFooter": Param{"string", 0, 0},
	"defaultMaxUsers": Param{"int", 0, 0},
	"defaultMeetingDuration": Param{"int", 0, 0},
	"defaultMeetingExpireDuration": Param{"int", 0, 0},
	"defaultMeetingCreateJoinDuration": Param{"int", 0, 0},
	"disableRecordingDefault": Param{"bool", 0, 0},
	"allowStartStopRecording": Param{"bool", 0, 0},
//	"bbb.web.logoutURL": Param{"url", 0, 0},
//	"defaultAvatarURL": Param{"url", 0, 0},
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
	getTomcat := exec.Command("bash", "-c", "cat "+tomcatCfgFile+" | grep "+param+" | sed "+sedstr)
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

func execPlain (params map[string]string, file string) ([]byte, error) {
	fmt.Println(params)
	for name, value := range(params) {
		str := "s/^"+name+".*/"+name+"="+value+"/"
		command := exec.Command("sed", "-i", str, file)
		output, err := command.CombinedOutput()
			fmt.Println("str is "+str)
			fmt.Print("output: "); fmt.Println(string(output))
			fmt.Print("error: "); fmt.Println(err)
		if (len(output) >0 ) || (err != nil) {
			fmt.Println("execPlain error")
			return output, err
		}
	}
	restart := exec.Command("service", "tomcat7", "restart")
	restartOutput, err := restart.CombinedOutput()
	fmt.Print("restart output: "); fmt.Println(string(restartOutput))
	fmt.Print("restart error: "); fmt.Println(err)
	return restartOutput, err
}

func evaluateParams(params map[string]string, tmpl map[string]Param) bool {

	for key, value := range(params) {
		must, ok := tmpl[key]
		if !ok {fmt.Println("can't get type in a map"); return false}
		switch must.tip {
			case "int": {
				value2int, err := strconv.Atoi(value)
				if (err != nil) || (value2int < must.min) || ((must.max != 0) && (value2int > must.max)) {fmt.Println("int"); return false}
			}
			case "bool": {
				_, err := strconv.ParseBool(value)
				if err != nil {fmt.Println("bool case"); return false}
			}
			case "url": {
				if !checkDomainName(value) {
					if strings.HasPrefix(value, url_prefix) {
						valueWOPrefix := strings.TrimPrefix(value, url_prefix)
						if !checkDomainName(valueWOPrefix) {return false}
					} else if value != "" {fmt.Printf("url case %s\n", value); return false}
				}
			}
		}
	}
	return true
}

func setTomcatHandler(r *fasthttp.RequestCtx) {
	
	username := getUserName(r)
	user, ok := getUser(username)

	secToken := r.FormValue("tokensec")
	if (!ok) || (!user.IsAdmin) || (!checkSec(secToken,username)) {out403(r); return}

	jsonObj := r.FormValue("settings")
	tomcatSettings := jsonArray{}
	if err := json.Unmarshal(jsonObj, &tomcatSettings); err != nil {fmt.Println(err); out500(r); return}

	params := tomcatSettings.Params;
	if !evaluateParams(params, tomcatTmpl) {out500(r); return}

	output, err := execPlain(params, tomcatCfgFile)
	fmt.Println("output "+string(output))
	if err != nil {
		fmt.Print("Error is ")
		fmt.Println(err)
		out500(r)
	} else {outnil(r)}
}

/*func setClientSettings() {
	path := "/var/www/bigbluebutton/client/conf"
	file := "config.xml"
	muteOnStart := "/config/modules/meeting/"
	modules := "/config/modules/module"
	chatModule := modules+"[@name=ChatModule]/"
	usersModule := modules+"[@name=UsersModule]/"
	DeskShareModule := modules+"[@name=DeskShareModule]/"
	VideoconfModule := modules+"[@name=VideoconfModule]/"
	PresentModule := modules+"[@name=PresentModule]/"
}

func setConnectionSettings() {
//	bbb-conf --setip ip
//	bbb-conf --setcesret secret
}
*/
