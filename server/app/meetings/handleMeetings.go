package meetings

import (
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/app/settings"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/app/bbbapi"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/check"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/globs"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/db"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/tokens"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/utils"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/users"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/out"
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

func MeetingsHandler(r *fasthttp.RequestCtx) {

	s, db := db.InitMongo()
	defer s.Close()

	meetings := db.C("meetings")

	index := string(r.FormValue("meetingID"))

	if index == "all" {
		meetingsList := []Meeting{}
		utils.Check(meetings.Find(nil).All(&meetingsList))
		
		output, err := json.Marshal(meetingsList)
		if err == nil {r.Write(output)}
		return
	}

	details := db.C("meetingDetails")

	meetingDetails := MeetingDetails{}

	utils.Check(details.FindId(bson.ObjectIdHex(index)).One(&meetingDetails))

	output, err := json.Marshal(meetingDetails)
	if err == nil {r.Write(output)}
}

func PasswordsHandler(r *fasthttp.RequestCtx) {

	passwords := Passwords{}

	username := users.GetUserName(r)
	user, ok := users.GetUser(username)

	if ok && user.IsAdmin {

		s, db := db.InitMongo()
		defer s.Close()

		pwds := db.C("passwords")

		index := string(r.FormValue("meetingID"))

		utils.Check(pwds.FindId(bson.ObjectIdHex(index)).One(&passwords))
	}

	output, err := json.Marshal(passwords)
	if err == nil {r.Write(output)}
}

func CreateMeetingHandler(r *fasthttp.RequestCtx) {

	username := users.GetUserName(r)
	user, ok := users.GetUser(username)

	secToken := r.FormValue("tokensec")
//	if (!ok) || (!user.IsAdmin) || (!tokens.CheckSec(secToken,username)) || (!tokens.DupControl(username)) {return}
	if (!ok) {fmt.Println("user doesn't exist"); return}
	if (!user.IsAdmin) {fmt.Println("User isn't admin"); return}
	if (!tokens.CheckSec(secToken,username)) {fmt.Println("wrong sectoken"); return}
	if !tokens.DupControl(username) {fmt.Println("Dup detected"); return}

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

	if !check.Validate(check.LOGIN, title) || !check.Validate(check.DESC, desc) || (len(title) < 1) || (len(admpwd) < 6) || (len(pwd) < 6) {out.Out500(r); return}

	s, mongodb := db.InitMongo()
	defer s.Close()

	meetings := mongodb.C("meetings")
	details := mongodb.C("meetingDetails")
	pwds := mongodb.C("passwords")

	if !meetingUniq(mongodb, title) {out.Out500(r); return}

	id := db.NewID()
	newMeeting.ID = id
	passwords.ID = id

	err = meetings.Insert(Meeting{
		ID: id,
		Title: title,
	})
	if err != nil {out.Out500(r);  return}

	err = details.Insert(newMeeting)
	if err != nil {out.Out500(r);  return}

	err = pwds.Insert(passwords)
	if err != nil {out.Out500(r);  return}

	out.Outnil(r)
}

func EditMeetingHandler(r *fasthttp.RequestCtx) {
	username := users.GetUserName(r)
	user, ok := users.GetUser(username)

	secToken := r.FormValue("tokensec")
//	if (!ok) || (!user.IsAdmin) || (!checkSec(secToken,username)) || (!dupControl(username)) {return}
	if (!ok) {fmt.Println("user doesn't exist"); return}
	if (!user.IsAdmin) {fmt.Println("User isn't admin"); return}
	if (!tokens.CheckSec(secToken,username)) {fmt.Println("wrong sectoken"); return}

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

	if !check.Validate(check.LOGIN, title) || !check.Validate(check.DESC, desc) || (len(title) < 1) || (len(admpwd) < 6) || (len(pwd) < 6) {out.Out500(r); return}

	s, db := db.InitMongo()
	defer s.Close()

	meetings := db.C("meetings")
	details := db.C("meetingDetails")
	pwds := db.C("passwords")

	id := newMeeting.ID

	err = meetings.Update(bson.M{"_id":id}, Meeting{
		ID: id,
		Title: title,
	})
	if err != nil {out.Out500(r);  return}

	err = details.Update(bson.M{"_id":id}, newMeeting)
	if err != nil {out.Out500(r);  return}

	err = pwds.Update(bson.M{"_id":id}, passwords)
	if err != nil {out.Out500(r);  return}

	out.Outnil(r)
}

func DeleteMeetingHandler(r *fasthttp.RequestCtx) {
	username := users.GetUserName(r)
	user, ok := users.GetUser(username)

	secToken := r.FormValue("tokensec")
	if (!ok) {fmt.Println("user doesn't exist"); return}
	if (!user.IsAdmin) {fmt.Println("User isn't admin"); return}
	if (!tokens.CheckSec(secToken,username)) {fmt.Println("wrong sectoken"); return}

	s, db := db.InitMongo()
	defer s.Close()

	id := string(r.FormValue("meetingID"))

	meetings := db.C("meetings")
	details := db.C("meetingDetails")
	passwords := db.C("passwords")

	err := meetings.RemoveId(bson.ObjectIdHex(id))
	if err != nil {out.Out500(r);  return}
	err = details.RemoveId(bson.ObjectIdHex(id))
	if err != nil {out.Out500(r);  return}
	err = passwords.RemoveId(bson.ObjectIdHex(id))
	if err != nil {out.Out500(r);  return}

	out.Outnil(r)
}

func MeetingUniqHandler(r *fasthttp.RequestCtx) {
	title := string(r.FormValue("title"))
	s, db := db.InitMongo()
	defer s.Close()
	if !meetingUniq(db, title) {out.Out(r, "Login already exists"); return}
	out.Outnil(r)
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

func CheckPwdHandler(r *fasthttp.RequestCtx) {
	pwd := string(r.FormValue("pwd"))
	id := string(r.FormValue("id"))
	passwords := Passwords{}
	s, db := db.InitMongo()
	defer s.Close()
	pwds := db.C("passwords")
	pwds.FindId(bson.ObjectIdHex(id)).One(&passwords)
	r.Write(globs.Output[(passwords.Admpwd == pwd) || (passwords.Pwd == pwd)])
}

func IsRunningHandler(r *fasthttp.RequestCtx) {
	id := r.FormValue("meetingID")
	str := "isMeetingRunningmeetingID="
	shaword := bbbapi.Sha(utils.AppendAll([][]byte{[]byte(str), id, []byte(settings.ConnSettings.Secret)}))
	full := "http://"+settings.ConnSettings.IP+"/bigbluebutton/api/isMeetingRunning?meetingID="+string(id)+"&checksum="+string(shaword)
	_, result, _ := fasthttp.Get(nil, full)
	r.Write(globs.Output[bytes.Contains(result, []byte("true"))])
}

func JoinHandler(r *fasthttp.RequestCtx) {
	username := users.GetUserName(r)
	_, ok := users.GetUser(username)
	if !ok {return}
	qstr := r.FormValue("string")
	checksum := bbbapi.Sha(utils.AppendAll([][]byte{[]byte("join"),qstr,[]byte(settings.ConnSettings.Secret)}))
	r.Write(utils.AppendAll([][]byte{[]byte("join?"),qstr,[]byte("&checksum="),checksum}))
}
