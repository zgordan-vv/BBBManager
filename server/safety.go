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
}

var (
	SALT0, SALT1, SITEKEY, SECRETKEY string
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
