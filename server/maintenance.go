package main

import "github.com/valyala/fasthttp"

var underMaintenance bool = false

func getMaintenance() bool {
	return underMaintenance
}

func setMaintenance(tf bool) {
	underMaintenance = tf
}

func maintenanceHandler(r *fasthttp.RequestCtx) {
	r.Write(output[underMaintenance])
}
