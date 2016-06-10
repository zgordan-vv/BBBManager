package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"regexp"
	"encoding/json"
	"os"
)

type SafetyState struct {
	SALT0 string
	SALT1 string
	SITEKEY string
	SECRETKEY string
	GITHUB_CLIENT_ID string
	GITHUB_CLIENT_SECRET string
	FB_CLIENT_ID string
	FB_CLIENT_SECRET string
}

var (
	SALT0, SALT1, SITEKEY, SECRETKEY, GITHUB_CLIENT_ID, GITHUB_CLIENT_SECRET, FB_CLIENT_ID, FB_CLIENT_SECRET string
	controlToken map[string]string = make(map[string]string)
	submitToken map[string]string = make(map[string]string)
	secToken map[string]string = make(map[string]string)
)

func initKeys() {
	var loadstate SafetyState
	file, err := os.Open("KEYS")
	defer file.Close()
	if err != nil {
		return
	}
	decoder := json.NewDecoder(file)
	check(decoder.Decode(&loadstate))
	SALT0 = loadstate.SALT0
	SALT1 = loadstate.SALT1
	SITEKEY = loadstate.SITEKEY
	SECRETKEY = loadstate.SECRETKEY
	GITHUB_CLIENT_ID = loadstate.GITHUB_CLIENT_ID
	GITHUB_CLIENT_SECRET = loadstate.GITHUB_CLIENT_SECRET
	FB_CLIENT_ID = loadstate.FB_CLIENT_ID
	FB_CLIENT_SECRET = loadstate.FB_CLIENT_SECRET
}

func validate(exp, input string) bool {
	result, err := regexp.MatchString(exp, input)
	if err != nil {result = false}
	return result
}

func passEncrypt(pwd string) string {
	h := md5.New()
	io.WriteString(h, pwd)
	pwdmd5 := fmt.Sprintf("%x", h.Sum(nil))
	io.WriteString(h,SALT0)
	io.WriteString(h,pwdmd5)
	io.WriteString(h,SALT1)
	return fmt.Sprintf("%x",h.Sum(nil))
}
