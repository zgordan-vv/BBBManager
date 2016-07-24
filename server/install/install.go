package install

import (
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/users"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/db"
	"gopkg.in/mgo.v2/bson"
)

var Installed bool

func CheckInstall () bool {
	s, mongodb := db.InitMongo()
	defer s.Close();

	usersC := mongodb.C("userDetails")
	usersArr := []users.User{}
	err := usersC.Find(bson.M{"isadmin":true}).All(&usersArr)

	Installed = (err == nil) && (len(usersArr) != 0)
	return Installed
}
