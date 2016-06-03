package main

import (
	"encoding/json"
	"net/http"
	"gopkg.in/mgo.v2/bson"
)

type Settings struct {
	ID string	`json:"id"`
	IP string	`json:"ip"`
	Secret string	`json:"secret"`
	Webrtc bool	`json:"webrtc"`
}

var settings = Settings{ID:"settings"}

func initSettings() {
	s, db := initMongo()
	defer s.Close()

	settings0 := Settings{}

	settingsC := db.C("settings")
	err := settingsC.Find(bson.M{"id":"settings"}).One(&settings0)
	if err == nil {settings = settings0}
}

func getSettingsHandler(w http.ResponseWriter, r *http.Request) {
	username := getUserName(r)
	user, ok := getUser(username)

	if (!ok) || (!user.IsAdmin) {out403(w); return}

	out, err := json.Marshal(settings)
	if err == nil {w.Write(out)}
}

func setSettingsHandler(w http.ResponseWriter, r *http.Request) {

	username := getUserName(r)
	user, ok := getUser(username)

	secToken := r.FormValue("tokensec")
	if (!ok) || (!user.IsAdmin) || (!checkSec(secToken,username)) {return}

	jsonObj := r.FormValue("settings")
	s, db := initMongo()
	defer s.Close()
	settingsC := db.C("settings")

	newSettings := struct{
		ID string	`json:"id"`
		IP string	`json:"ip"`
		Secret string	`json:"secret"`
		Webrtc bool	`json:"webrtc"`
	}{ID:"settings"}

	if json.Unmarshal([]byte(jsonObj), &newSettings) != nil {return}

	if !validate(IP, newSettings.IP) || !validate(DIGITS, newSettings.Secret) {out(w, "Not validated"+newSettings.IP+" "+newSettings.Secret+" ..."); return}

	_, err := settingsC.Upsert(bson.M{"id":"settings"}, newSettings)
	if err != nil {out500(w);  return}
	settings = newSettings
	outnil(w)
}
