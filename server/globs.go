package main

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"os"
)

type State struct {
	DB string
	DN string
}

func initGlobs() {
	var loadstate State
	file, err := os.Open("conf/GLOBS")
	defer file.Close()

	if err != nil {
		return
	}

	decoder := json.NewDecoder(file)
	check(decoder.Decode(&loadstate))
	DBPREFIX = loadstate.DB
	DOMAINNAME = loadstate.DN
	if !checkDBPrefix(DBPREFIX) || !checkDomainName(DOMAINNAME) {return}
	DBPREFIX += ":"
	s, db := initMongo()
	defer s.Close();

	usersC := db.C("userDetails")
	usersArr := []User{}
	err = usersC.Find(bson.M{"isadmin":true}).All(&usersArr)

	if (err != nil) || (len(usersArr) == 0) {return}
	installed = true
}
