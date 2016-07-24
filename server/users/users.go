package users

import (
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/db"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/globs"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/auth"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Name string	`json:"name"`
	FullName string	`json:"fullname"`
	IsAdmin bool	`json:"isadmin"`
	Keyword string
}

func CheckPwdChange(username,old,new string) (bool, bool) {
	if (old == "") && (new == "") {return false, true} else {
		if !CheckPwd(username,old) {return true, false} else {return true,true}
	}
}

func NameUniq(name string) bool {

	client, err := db.InitRedis()
	if err == nil {
		defer client.Close()
		userExists,err := client.Do("SISMEMBER", globs.DBPREFIX+"users", name)
		if err == nil {
			if userExists.(int64) != 0 {
				return false
			}
		}
	}

	s, db := db.InitMongo()
	defer s.Close();

	usersC := db.C("userDetails")
	userDetails := User{}
	usersC.Find(bson.M{"Name":name}).One(&userDetails)
	return userDetails == User{}

}

func SaveUser(user User) error {

	s, mongodb := db.InitMongo()
	defer s.Close();

	userDetailsC := mongodb.C("userDetails")

	username := user.Name

	client, err := db.InitRedis()
	if err == nil {
		defer client.Close()
		_, err = client.Do("SADD", globs.DBPREFIX+"users",username)
		if err != nil {return err}
	} else {return err}

	_, err = userDetailsC.Upsert(bson.M{"name":username}, user)
	return err
}

func GetUser(username string) (User, bool) {

	client, err := db.InitRedis()
	if err == nil {
		defer client.Close()
		userExists,err := client.Do("SISMEMBER", globs.DBPREFIX+"users", username)
		if err == nil {
			if userExists.(int64) == 1 {
				s, db := db.InitMongo()
				defer s.Close();

				userDetailsC := db.C("userDetails")

				user := User{}
				err = userDetailsC.Find(bson.M{"name":username}).One(&user)
				if err == nil {return user, true}
			}
		}
	}
	return User{}, false
}

func CheckPwd(username,pwd string) bool {
	user, ok := GetUser(username)
	if ok {return user.Keyword == auth.PassEncrypt(pwd)} else {return false}
}
