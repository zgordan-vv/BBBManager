package settings

import (
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/check"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/out"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/users"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"errors"
	"regexp"
	"os/exec"
)

type Settings struct {
	IP string	`json:"ip"`
	Secret string	`json:"secret"`
}

var ConnSettings Settings

func InitSettings() {
	cmd := exec.Command("bbb-conf", "--secret")
	output, err := cmd.CombinedOutput()
	if err != nil {fmt.Println(err); return}

	url_exp := regexp.MustCompile("\\/\\/(.*)\\/bigbluebutton")
	url := url_exp.FindStringSubmatch(string(output))[1]

	salt_exp := regexp.MustCompile("Secret: (.*)")
	salt := salt_exp.FindStringSubmatch(string(output))[1]

	ConnSettings = Settings{url, salt}
}

func GetSettingsHandler(r *fasthttp.RequestCtx) {
	username := users.GetUserName(r)
	user, ok := users.GetUser(username)

	if (!ok) || (!user.IsAdmin) {out.Out403(r); return}

	output, err := json.Marshal(ConnSettings)
	if err == nil {r.Write(output)}
}

func setSalt(salt string) error {

	if !check.Validate(check.DIGITS, salt) {return errors.New("Not validated "+ salt +" ...")}

	ConnSettings.Secret = salt
	return nil
}
