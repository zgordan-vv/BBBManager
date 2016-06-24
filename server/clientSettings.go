package main

import (
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

func getClientHandler(r *fasthttp.RequestCtx) {
	username := getUserName(r)
	user, ok := getUser(username)

	if (!ok) || (!user.IsAdmin) {out403(r); return}

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

func setClientHandler(r *fasthttp.RequestCtx) {
	
	if getMaintenance() {out500(r); return}
	
	username := getUserName(r)
	user, ok := getUser(username)

	secToken := r.FormValue("tokensec")
	if (!ok) || (!user.IsAdmin) || (!checkSec(secToken,username)) {out403(r); return}

	jsonObj := r.FormValue("settings")
	clientSettings := jsonArray{}
	if err := json.Unmarshal([]byte(jsonObj), &clientSettings); err != nil {fmt.Println(err); out500(r); return}

	params := clientSettings.Params;

	for paired := range(clientTmpl) {
		value := params[paired]
		if !evaluateParam(paired, value, clientTmpl) {return}
		pair := strings.Split(paired, "/")
		module, param := pair[0], pair[1]
		path := "/config/modules/module[@name='"+module+"']/@"+param
		output, err := updateXMLParam(clientFullPath, path, value)
		fmt.Println("output "+string(output))
		if err != nil {out500(r); return}
	}

	restartClient(r)
}

func resetClientHandler(r *fasthttp.RequestCtx) {

	if getMaintenance() {out500(r); return}
	
	username := getUserName(r)
	user, ok := getUser(username)

	secToken := r.FormValue("tokensec")
	if (!ok) || (!user.IsAdmin) || (!checkSec(secToken,username)) {out403(r); return}

	setMaintenance(true)
	defer setMaintenance(false)

	err := resetDefaults(r, clientCfgPath, clientCfgFile)
	if err == nil {restartClient(r); fmt.Println("restarted")} else {out500(r)}
}

func restartClient(r *fasthttp.RequestCtx){

	restart := exec.Command("bbb-conf", "--clean")
	err := restart.Start()
	fmt.Print("error is ")
	fmt.Print(err)

	if err != nil {
		fmt.Print("Error is ")
		fmt.Println(err)
		out500(r)
	} else {outnil(r)}

	result := restart.Wait()
	fmt.Println(result)

	outnil(r)
}
