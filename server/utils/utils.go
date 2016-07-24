package utils

import (
	"crypto/rand"
	"encoding/base64"
	"os/exec"
	"io"
)

func CW(n int) []byte {
	b := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return nil
	}
	result := []byte(base64.URLEncoding.EncodeToString(b))
	return result
}

func Xmlstarlet(path, file, param, value string) string {
	comstr := "cd '+path+' && xmlstarlet ed -u '+param+' -v '+value+' tmp.xml && mv -f tmp.xml '+file"
	command := exec.Command("bash", "-c", comstr)
	err := command.Start()
	if err != nil {return "500"}
	err = command.Wait()
	if err != nil {return "500"}
	return ""
}

func AppendAll(slices [][]byte) []byte {
	result := []byte{}
	for _, slice := range(slices) {
		result = append(result, slice...)
	}
	return result
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}
