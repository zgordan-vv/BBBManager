package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
)

type Meeting struct{
	ID bson.ObjectId `json:"id" bson:"_id"`
	Title string	`json:"title"`
}

type MeetingDetails struct{
	ID  bson.ObjectId `json:"id" bson:"_id"`
	Title string	`json:"title"`
	Desc string	`json:"desc"`
	Author string	`json:"author"`
	Welcome string	`json:"welcome"`
	Duration int	`json:"duration"`
	Isrec bool	`json:"isrec"`
	Autorec bool	`json:"autorec"`
	AllowStartStopRec bool	`json:"allowstartstoprec"`
}

type Passwords struct {
	ID  bson.ObjectId `json:"id" bson:"_id"`
	Admpwd string	`json:"admpwd"`
	Pwd string	`json:"pwd"`
}

func meetingsHandler(r *fasthttp.RequestCtx) {

	s, db := initMongo()
	defer s.Close()

	meetings := db.C("meetings")

	index := string(r.FormValue("meetingID"))

	if index == "all" {
		meetingsList := []Meeting{}
		check(meetings.Find(nil).All(&meetingsList))
		
		output, err := json.Marshal(meetingsList)
		if err == nil {r.Write(output)}
		return
	}

	details := db.C("meetingDetails")

	meetingDetails := MeetingDetails{}

	check(details.FindId(bson.ObjectIdHex(index)).One(&meetingDetails))

	output, err := json.Marshal(meetingDetails)
	if err == nil {r.Write(output)}
}

func passwordsHandler(r *fasthttp.RequestCtx) {

	passwords := Passwords{}

	username := getUserName(r)
	user, ok := getUser(username)

	if ok && user.IsAdmin {

		s, db := initMongo()
		defer s.Close()

		pwds := db.C("passwords")

		index := string(r.FormValue("meetingID"))

		check(pwds.FindId(bson.ObjectIdHex(index)).One(&passwords))
	}

	output, err := json.Marshal(passwords)
	if err == nil {r.Write(output)}
}

func createMeetingHandler(r *fasthttp.RequestCtx) {
	username := getUserName(r)
	user, ok := getUser(username)

	secToken := r.FormValue("tokensec")
//	if (!ok) || (!user.IsAdmin) || (!checkSec(secToken,username)) || (!dupControl(username)) {return}
	if (!ok) {fmt.Println("user doesn't exist"); return}
	if (!user.IsAdmin) {fmt.Println("User isn't admin"); return}
	if (!checkSec(secToken,username)) {fmt.Println("wrong sectoken"); return}
	if !dupControl(username) {fmt.Println("Dup detected"); return}

	meetingJson := r.FormValue("meeting")
	newMeeting := MeetingDetails{}
	err := json.Unmarshal(meetingJson, &newMeeting); if err != nil {fmt.Println(err); return}
	title := newMeeting.Title
	desc := newMeeting.Desc

	newMeeting.Author = user.FullName;

	passwords := Passwords{}

	pwdsJson := r.FormValue("passwords")
	err = json.Unmarshal(pwdsJson, &passwords); if err != nil {fmt.Println(err); return}

	admpwd := passwords.Admpwd
	pwd := passwords.Pwd

	if !validate(LOGIN, title) || !validate(DESC, desc) || (len(title) < 1) || (len(admpwd) < 6) || (len(pwd) < 6) {out500(r); return}

	s, db := initMongo()
	defer s.Close()

	meetings := db.C("meetings")
	details := db.C("meetingDetails")
	pwds := db.C("passwords")

	if !meetingUniq(db, title) {out500(r); return}

	id := newID()
	newMeeting.ID = id
	passwords.ID = id

	err = meetings.Insert(Meeting{
		ID: id,
		Title: title,
	})
	if err != nil {out500(r);  return}

	err = details.Insert(newMeeting)
	if err != nil {out500(r);  return}

	err = pwds.Insert(passwords)
	if err != nil {out500(r);  return}

	outnil(r)
}

func editMeetingHandler(r *fasthttp.RequestCtx) {
	username := getUserName(r)
	user, ok := getUser(username)

	secToken := r.FormValue("tokensec")
//	if (!ok) || (!user.IsAdmin) || (!checkSec(secToken,username)) || (!dupControl(username)) {return}
	if (!ok) {fmt.Println("user doesn't exist"); return}
	if (!user.IsAdmin) {fmt.Println("User isn't admin"); return}
	if (!checkSec(secToken,username)) {fmt.Println("wrong sectoken"); return}

	jsonObj := r.FormValue("meeting")
	newMeeting := MeetingDetails{}
	if json.Unmarshal(jsonObj, &newMeeting) != nil {return}
	title := newMeeting.Title
	desc := newMeeting.Desc

	newMeeting.Author = user.FullName;

	passwords := Passwords{}
	pwdsJson := r.FormValue("passwords")
	err := json.Unmarshal(pwdsJson, &passwords); if err != nil {fmt.Println(err); return}

	admpwd := passwords.Admpwd
	pwd := passwords.Pwd

	if !validate(LOGIN, title) || !validate(DESC, desc) || (len(title) < 1) || (len(admpwd) < 6) || (len(pwd) < 6) {out500(r); return}

	s, db := initMongo()
	defer s.Close()

	meetings := db.C("meetings")
	details := db.C("meetingDetails")
	pwds := db.C("passwords")

	id := newMeeting.ID

	err = meetings.Update(bson.M{"_id":id}, Meeting{
		ID: id,
		Title: title,
	})
	if err != nil {out500(r);  return}

	err = details.Update(bson.M{"_id":id}, newMeeting)
	if err != nil {out500(r);  return}

	err = pwds.Update(bson.M{"_id":id}, passwords)
	if err != nil {out500(r);  return}

	outnil(r)
}

func deleteMeetingHandler(r *fasthttp.RequestCtx) {
	username := getUserName(r)
	user, ok := getUser(username)

	secToken := r.FormValue("tokensec")
	if (!ok) {fmt.Println("user doesn't exist"); return}
	if (!user.IsAdmin) {fmt.Println("User isn't admin"); return}
	if (!checkSec(secToken,username)) {fmt.Println("wrong sectoken"); return}

	s, db := initMongo()
	defer s.Close()

	id := string(r.FormValue("meetingID"))

	meetings := db.C("meetings")
	details := db.C("meetingDetails")
	passwords := db.C("passwords")

	err := meetings.RemoveId(bson.ObjectIdHex(id))
	if err != nil {out500(r);  return}
	err = details.RemoveId(bson.ObjectIdHex(id))
	if err != nil {out500(r);  return}
	err = passwords.RemoveId(bson.ObjectIdHex(id))
	if err != nil {out500(r);  return}

	outnil(r)
}

func meetingUniqHandler(r *fasthttp.RequestCtx) {
	title := string(r.FormValue("title"))
	s, db := initMongo()
	defer s.Close()
	if !meetingUniq(db, title) {out(r, "Login already exists"); return}
	outnil(r)
}

func meetingUniq(db *mgo.Database, title string) bool {
	
	meetings := db.C("meetings")
	details := db.C("meetingDetails")

	ctrl:= Meeting{}
	ctrlDetails := MeetingDetails{}
	meetings.Find(bson.M{"title":title}).One(&ctrl)
	details.Find(bson.M{"title":title}).One(&ctrlDetails)
	return (ctrl.Title == "") && (ctrlDetails.Title == "")
}

func checkPwdHandler(r *fasthttp.RequestCtx) {
	pwd := string(r.FormValue("pwd"))
	id := string(r.FormValue("id"))
	passwords := Passwords{}
	s, db := initMongo()
	defer s.Close()
	pwds := db.C("passwords")
	pwds.FindId(bson.ObjectIdHex(id)).One(&passwords)
	r.Write(output[(passwords.Admpwd == pwd) || (passwords.Pwd == pwd)])
}

func isRunningHandler(r *fasthttp.RequestCtx) {
	id := r.FormValue("meetingID")
	str := "isMeetingRunningmeetingID="
	shaword := sha(appendAll([][]byte{[]byte(str), id, []byte(settings.Secret)}))
	full := "http://"+settings.IP+"/bigbluebutton/api/isMeetingRunning?meetingID="+string(id)+"&checksum="+string(shaword)
	_, result, _ := fasthttp.Get(nil, full)
	r.Write(output[bytes.Contains(result, []byte("true"))])
}
