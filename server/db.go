package main

import (
	"github.com/garyburd/redigo/redis"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func initMongo() (*mgo.Session, *mgo.Database) {
	session,err := mgo.Dial("127.0.0.1")
	check(err)
	session.SetMode(mgo.Monotonic, true)
	return session, session.DB(DBPREFIX[:len(DBPREFIX)-1])
}

func newID() bson.ObjectId {
	return bson.NewObjectId()
}

func addStats(host, ip, clientID string) {
	client, err := redis.Dial("tcp", ":6379")
	if err == nil {
		defer client.Close()
		record := DBPREFIX+"ipstats:"+clientID+":"+ip
			client.Do("HINCRBY", record, "internal", 1)
		if host != DOMAINNAME {
			client.Do("HINCRBY", record, "external", 1)
		}
	}
}
