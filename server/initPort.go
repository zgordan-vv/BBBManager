package main

import (
	"encoding/json"
	"os"
)

type PortState struct {
	Port string
}

func initPort() {
	var loadstate PortState
	file, err := os.Open("PORT")
	defer file.Close()
	if err != nil {
		return
	}
	decoder := json.NewDecoder(file)
	check(decoder.Decode(&loadstate))
	PORT = loadstate.Port
	return
}
