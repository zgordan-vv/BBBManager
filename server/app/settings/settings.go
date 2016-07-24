package settings

import (
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/tokens"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/users"
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
	"os/exec"
	"strconv"
)

type jsonArray struct {
	Params map[string]string `json: "params"`
}

type Param struct {
	tip string
	min int
	max int
}

const defaultsPath string = "../defaults/"

func evaluateParam(key string, value string, tmpl map[string]Param) bool {

	must, ok := tmpl[key]
	if !ok {fmt.Println("can't get type in a map"); return false}
	switch must.tip {
		case "int": {
			value2int, err := strconv.Atoi(value)
			if (err != nil) || (value2int < must.min) || ((must.max != 0) && (value2int > must.max)) {fmt.Println("int"); return false}
		}
		case "bool": {
			_, err := strconv.ParseBool(value)
			if err != nil {return false}
		}
	}
	return true
}

func getXMLParam(file, queryString string) string {
	xmlQuery := exec.Command("bash", "-c", "xmlstarlet sel -t -v \""+queryString+"\" "+file)
	output, err := xmlQuery.CombinedOutput()
	if err != nil {return ""} else {
		return string(output)
	}
}

func updateXMLParam(file, param, value string) ([]byte, error) {
	str := "xmlstarlet ed -u \""+param+"\" -v "+value+" "+file+" > /tmp/bbbfstmp.xml && mv -f /tmp/bbbfstmp.xml "+file
	updateParam := exec.Command("bash", "-c", str)
	output, err := updateParam.CombinedOutput()
	fmt.Printf("Output is %s, err is %s\n", output, err)
	return output, err
}

func resetDefaults(r *fasthttp.RequestCtx, path, filename string) error {
	
	username := users.GetUserName(r)
	user, ok := users.GetUser(username)

	secToken := r.FormValue("tokensec")
	answer := r.FormValue("answer")
	if (!ok) || (!user.IsAdmin) || (!tokens.CheckSec(secToken,username)) {return errors.New("403")}
	if string(answer) != "yes" {return errors.New("400")}

	reset := exec.Command("cp", "-f", defaultsPath+filename, path)
	resetOutput, err := reset.CombinedOutput()
	fmt.Println(string(resetOutput))
	if err != nil {
		fmt.Println(err)
		return err
	} else { return nil }
}
