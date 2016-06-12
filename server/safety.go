package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"regexp"
	"os"
)

var (
	SALT0, SALT1, SITEKEY, SECRETKEY string
	controlToken map[string]string = make(map[string]string)
	submitToken map[string]string = make(map[string]string)
	secToken map[string]string = make(map[string]string)
)

func initKeys() {
	SALT0 = os.Getenv("BBB_SALT0")
	SALT1 = os.Getenv("BBB_SALT1")
	SITEKEY = os.Getenv("BBB_SITEKEY")
	SECRETKEY = os.Getenv("BBB_SECRETKEY")
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
