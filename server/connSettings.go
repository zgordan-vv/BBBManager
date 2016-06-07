package main

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
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

func getSettingsHandler(r *fasthttp.RequestCtx) {
	username := getUserName(r)
	user, ok := getUser(username)

	if (!ok) || (!user.IsAdmin) {out403(r); return}

	out, err := json.Marshal(settings)
	if err == nil {r.Write(out)}
}

func setSettingsHandler(r *fasthttp.RequestCtx) {

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

	if json.Unmarshal(jsonObj, &newSettings) != nil {return}

	if !validate(IP, newSettings.IP) || !validate(DIGITS, newSettings.Secret) {out(r, "Not validated"+newSettings.IP+" "+newSettings.Secret+" ..."); return}

	_, err := settingsC.Upsert(bson.M{"id":"settings"}, newSettings)
	if err != nil {out500(r);  return}
	settings = newSettings
	outnil(r)
}
