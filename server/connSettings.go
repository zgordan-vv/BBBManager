package main

import (
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

var settings Settings

func initSettings() {
	cmd := exec.Command("bbb-conf", "--secret")
	output, err := cmd.CombinedOutput()
	if err != nil {fmt.Println(err); return}

	url_exp := regexp.MustCompile("\\/\\/(.*)\\/bigbluebutton")
	url := url_exp.FindStringSubmatch(string(output))[1]

	salt_exp := regexp.MustCompile("Secret: (.*)")
	salt := salt_exp.FindStringSubmatch(string(output))[1]

	settings = Settings{url, salt}
}

func getSettingsHandler(r *fasthttp.RequestCtx) {
	username := getUserName(r)
	user, ok := getUser(username)

	if (!ok) || (!user.IsAdmin) {out403(r); return}

	out, err := json.Marshal(settings)
	if err == nil {r.Write(out)}
}

func setSalt(salt string) error {

	if !validate(DIGITS, salt) {return errors.New("Not validated "+ salt +" ...")}

	settings.Secret = salt
	return nil
}
