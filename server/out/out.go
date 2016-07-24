package out

import (
	"github.com/valyala/fasthttp"
)

func Out(w *fasthttp.RequestCtx, s string) {
	w.Write([]byte(s))
}

func Outnil(w *fasthttp.RequestCtx) {
	w.Write(nil)
}

func Out403(w *fasthttp.RequestCtx) {
	w.Write([]byte{'4','0','3'})
}

func Out200(w *fasthttp.RequestCtx) {
	w.Write([]byte{'2','0','0'})
}

func Out500(w *fasthttp.RequestCtx) {
	w.Write([]byte{'5','0','0'})
}
