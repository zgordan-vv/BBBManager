package stats

import (
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/db"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/globs"
)

func AddStats(host, ip, clientID string) {
	client, err := db.InitRedis()
	if err == nil {
		defer client.Close()
		record := globs.DBPREFIX+"ipstats:"+clientID+":"+ip
			client.Do("HINCRBY", record, "internal", 1)
		if host != globs.DOMAINNAME {
			client.Do("HINCRBY", record, "external", 1)
		}
	}
}
