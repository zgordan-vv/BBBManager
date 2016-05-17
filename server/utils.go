package main

import (
	"crypto/rand"
	"encoding/base64"
	"os/exec"
	"io"
)

func CW(n int) string {
	b := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func xmlstarlet(path, file, param, value string) string {
	comstr := "cd '+path+' && xmlstarlet ed -u '+param+' -v '+value+' tmp.xml && mv -f tmp.xml '+file"
	command := exec.Command("bash", "-c", comstr)
	err := command.Start()
	if err != nil {return "500"}
	err = command.Wait()
	if err != nil {return "500"}
	return ""
}
