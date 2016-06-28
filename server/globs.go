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
	dbpr := DBPREFIX[:len(DBPREFIX)-1]
	if !checkDBPrefix(dbpr) || !checkDomainName(DOMAINNAME) {return}

	s, db := initMongo()
	defer s.Close();

	usersC := db.C("userDetails")
	usersArr := []User{}
	err = usersC.Find(bson.M{"isadmin":true}).All(&usersArr)

	if (err != nil) || (len(usersArr) == 0) {return}
	installed = true
}

func saveGlobs(dbprefix, domainname string) {
	var savestate State
	savestate.DB = dbprefix
	savestate.DN = domainname
	file, err := os.Create("GLOBS")
	check(err)
	defer file.Close()
	encoder := json.NewEncoder(file)
	check(encoder.Encode(savestate))
}
