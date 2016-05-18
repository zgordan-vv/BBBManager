package main

import (
	"fmt"
	"encoding/json"
	"os/exec"
	"net/http"
	"strconv"
)

type jsonArray struct {
	Params [][]string `json: "params"`
}

type Param struct {
	tip string
	min int
	max int
}

const tomcatCfgFile string = "/var/lib/tomcat7/webapps/bigbluebutton/WEB-INF/classes/bigbluebutton.properties"
const sedstr string = AP+"s/^.*=\\(.*\\)$/\\1/"+AP

var tomcatTmpl = map[string]Param{
	"maxNumPages": Param{"int", 0, 0},
	"defaultWelcomeMessage": Param{"string", 0, 0},
	"defaultWelcomeMessageFooter": Param{"string", 0, 0},
	"defaultMaxUsers": Param{"int", 0, 0},
	"defaultMeetingDuration": Param{"int", 0, 0},
	"defaultMeetingExpireDuration": Param{"int", 0, 0},
	"defaultMeetingCreateJoinDuration": Param{"int", 0, 0},
	"disableRecordingDefault": Param{"bool", 0, 0},
	"allowStartStopRecording": Param{"bool", 0, 0},
	"bbb.web.logoutURL": Param{"url", 0, 0},
	"defaultAvatarURL": Param{"url", 0, 0},
}

func getTomcatHandler(w http.ResponseWriter, r *http.Request) {
	result := jsonArray{}
	params := [][]string{}
	for key := range(tomcatTmpl) {
		params = append(params, []string{key, getTomcatParam(key)})
	}
	fmt.Println(params)
	result.Params = params
	fmt.Println(result)
}

func getTomcatParam(param string) string {
	getTomcat := exec.Command("bash", "-c", "cat "+tomcatCfgFile+" | grep "+param+" | sed "+sedstr)
	getTomcatOutput, err := getTomcat.CombinedOutput()
	if err != nil {return ""} else {return string(getTomcatOutput)}
}

func execPlain (params [][]string, file string) ([]byte, error) {
	fmt.Println(params)
	for _, param := range(params) {
		name := param[0]
		value := param[1]
		str := "s/^"+name+".*/"+name+"="+value+"/"
		command := exec.Command("sed", "-i", str, file)
		output, err := command.CombinedOutput()
			fmt.Print("output: "); fmt.Println(string(output))
			fmt.Print("error: "); fmt.Println(err)
		if (output != nil) || (err != nil) {
			return output, err
		}
	}
	restart := exec.Command("service", "tomcat7", "restart")
	restartOutput, err := restart.CombinedOutput()
	fmt.Print("restart output: "); fmt.Println(string(restartOutput))
	fmt.Print("restart error: "); fmt.Println(err)
	return restartOutput, err
}

func evaluateParams(params [][]string, tmpl map[string]Param) bool {

	for _, param := range(params) {
		key := param[0]
		must, ok := tmpl[key]
		if !ok {return false}
		value:= param[1]
		switch must.tip {
			case "int": {
				value2int, err := strconv.Atoi(value)
				if (err != nil) || (value2int < must.min) || ((must.max != 0) && (value2int > must.max)) {return false}
			}
			case "bool": {
				_, err := strconv.ParseBool(value)
				if err != nil {return false}
			}
			case "url": {
				if !checkDomainName(value) {return false}
			}
		}
	}
	return true
}

func setTomcatHandler(w http.ResponseWriter, r *http.Request) {
	
	username := getUserName(r)
	user, ok := getUser(username)

	secToken := r.FormValue("tokensec")
	fmt.Println("stage 1, tokensec is "+secToken); fmt.Println(user.IsAdmin)
	if (!ok) || (!user.IsAdmin) || (!checkSec(secToken,username)) {out403(w); return}
	fmt.Println("stage 2")

	jsonObj := r.FormValue("settings")
	fmt.Println(jsonObj)
	tomcatSettings := jsonArray{}
	if json.Unmarshal([]byte(jsonObj), &tomcatSettings) != nil {out500(w); return}
	fmt.Println(tomcatSettings)

	fmt.Println("stage 3")
	params := tomcatSettings.Params;
	fmt.Println(tomcatSettings.Params)
	if !evaluateParams(params, tomcatTmpl) {out500(w); return}

	fmt.Println("stage 4")
	
	output, err := execPlain(params, tomcatCfgFile)
	fmt.Println("stage 5")
	if err != nil || output != nil {out500(w)} else {outnil(w)}
}

/*func setFreeswitchSettings() {
	path := "/opt/freeswitch/conf/autoload_configs"
	file := "conference.conf.xml"
	energyLevelPath := "/configuration/profiles/profile/param[@name='energy-level']/@value"
}

func setClientSettings() {
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
