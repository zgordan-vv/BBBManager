package db

import (
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/globs"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/utils"
	"github.com/garyburd/redigo/redis"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func InitMongo() (*mgo.Session, *mgo.Database) {
	session,err := mgo.Dial("127.0.0.1")
	utils.Check(err)
	session.SetMode(mgo.Monotonic, true)
	return session, session.DB(globs.DBPREFIX[:len(globs.DBPREFIX)-1])
}

func InitRedis() (redis.Conn, error) {
	client, err := redis.Dial("tcp", ":6379")
	return client, err
}

func NewID() bson.ObjectId {
	return bson.NewObjectId()
}
