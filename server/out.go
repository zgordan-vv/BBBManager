package main

import (
	"net/http"
)

func out(w http.ResponseWriter, s string) {
	w.Write([]byte(s))
}

func outnil(w http.ResponseWriter) {
	w.Write(nil)
}

func out403(w http.ResponseWriter) {
	w.Write([]byte{'4','0','3'})
}

func out200(w http.ResponseWriter) {
	w.Write([]byte{'2','0','0'})
}

func out500(w http.ResponseWriter) {
	w.Write([]byte{'5','0','0'})
}
