package main

import (
	"github.com/valyala/fasthttp"
)

func out(w *fasthttp.RequestCtx, s string) {
	w.Write([]byte(s))
}

func outnil(w *fasthttp.RequestCtx) {
	w.Write(nil)
}

func out403(w *fasthttp.RequestCtx) {
	w.Write([]byte{'4','0','3'})
}

func out200(w *fasthttp.RequestCtx) {
	w.Write([]byte{'2','0','0'})
}

func out500(w *fasthttp.RequestCtx) {
	w.Write([]byte{'5','0','0'})
}
