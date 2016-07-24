package maintenance

import (
	"github.com/valyala/fasthttp"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/globs"
)

var underMaintenance bool = false

func GetMaintenance() bool {
	return underMaintenance
}

func SetMaintenance(tf bool) {
	underMaintenance = tf
}

func MaintenanceHandler(r *fasthttp.RequestCtx) {
	r.Write(globs.Output[underMaintenance])
}
