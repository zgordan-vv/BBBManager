package settings

import (
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/out"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/users"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/tokens"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/app/maintenance"
	"fmt"
	"encoding/json"
	"github.com/valyala/fasthttp"
	"os/exec"
	"strings"
)

const clientCfgPath string = "/var/www/bigbluebutton/client/conf/"
const clientCfgFile string = "config.xml"
var clientFullPath string = clientCfgPath + clientCfgFile

var clientTmpl = map[string]Param{
	"ChatModule/privateEnabled": Param{"bool", 0, 0},
	"UsersModule/enableSettingsButton": Param{"bool", 0, 0},
	"DeskShareModule/showButton": Param{"bool", 0, 0},
	"DeskShareModule/autoStart": Param{"bool", 0, 0},
	"DeskShareModule/autoFullScreen": Param{"bool", 0, 0},
	"VideoconfModule/autoStart": Param{"bool", 0, 0},
	"VideoconfModule/skipCamSettingsCheck": Param{"bool", 0, 0},
	"VideoconfModule/showButton": Param{"bool", 0, 0},
	"VideoconfModule/showCloseButton": Param{"bool", 0, 0},
	"VideoconfModule/smoothVideo": Param{"bool", 0, 0},
	"VideoconfModule/displayAvatar": Param{"bool", 0, 0},
	"VideoconfModule/videoQuality": Param{"int", 0, 100},
	"VideoconfModule/camQualityBandwidth": Param{"int", 0, 100},
	"VideoconfModule/camQualityPicture": Param{"int", 0, 100},
}

func GetClientHandler(r *fasthttp.RequestCtx) {
	username := users.GetUserName(r)
	user, ok := users.GetUser(username)

	if (!ok) || (!user.IsAdmin) {out.Out403(r); return}

	result := jsonArray{}
	params := map[string]string{}
	for paired := range(clientTmpl) {
		pair := strings.Split(paired, "/")
		module, param := pair[0], pair[1]
		params[paired] = getXMLParam(clientFullPath, "/config/modules/module[@name='"+module+"']/@"+param)
	}
	result.Params = params
	output, _ := json.Marshal(result)
	r.Write(output)
}

func SetClientHandler(r *fasthttp.RequestCtx) {
	
	if maintenance.GetMaintenance() {out.Out500(r); return}
	
	username := users.GetUserName(r)
	user, ok := users.GetUser(username)

	secToken := r.FormValue("tokensec")
	if (!ok) || (!user.IsAdmin) || (!tokens.CheckSec(secToken,username)) {out.Out403(r); return}

	jsonObj := r.FormValue("settings")
	clientSettings := jsonArray{}
	if err := json.Unmarshal([]byte(jsonObj), &clientSettings); err != nil {fmt.Println(err); out.Out500(r); return}

	params := clientSettings.Params;

	for paired := range(clientTmpl) {
		value := params[paired]
		if !evaluateParam(paired, value, clientTmpl) {return}
		pair := strings.Split(paired, "/")
		module, param := pair[0], pair[1]
		path := "/config/modules/module[@name='"+module+"']/@"+param
		output, err := updateXMLParam(clientFullPath, path, value)
		fmt.Println("output "+string(output))
		if err != nil {out.Out500(r); return}
	}

	restartClient(r)
}

func ResetClientHandler(r *fasthttp.RequestCtx) {

	if maintenance.GetMaintenance() {out.Out500(r); return}
	
	username := users.GetUserName(r)
	user, ok := users.GetUser(username)

	secToken := r.FormValue("tokensec")
	if (!ok) || (!user.IsAdmin) || (!tokens.CheckSec(secToken,username)) {out.Out403(r); return}

	maintenance.SetMaintenance(true)
	defer maintenance.SetMaintenance(false)

	err := resetDefaults(r, clientCfgPath, clientCfgFile)
	if err == nil {restartClient(r); fmt.Println("restarted")} else {out.Out500(r)}
}

func restartClient(r *fasthttp.RequestCtx){

	restart := exec.Command("bbb-conf", "--clean")
	err := restart.Start()
	fmt.Print("error is ")
	fmt.Print(err)

	if err != nil {
		fmt.Print("Error is ")
		fmt.Println(err)
		out.Out500(r)
	} else {out.Outnil(r)}

	result := restart.Wait()
	fmt.Println(result)

	out.Outnil(r)
}
