package main

import (
	"github.com/valyala/fasthttp"
)

var secTokens map[string][][]byte = make(map[string][][]byte)
var dupTokens map[string][]byte = make(map[string][]byte)
var dupUsedTokens map[string][]byte = make(map[string][]byte)

func secTokenHandler(r *fasthttp.RequestCtx) {
	username := getUserName(r)
	secToken := CW(32)
	secTokens[username] = append(secTokens[username],secToken)
	r.Write(secToken)
}

func checkSec(sectoken []byte, username string) bool {
	tokens, ok := secTokens[username]
	if !ok {return false} else {
		return true
		for i, token := range(tokens) {
			if string(token) == string(sectoken) {
				tokens = append(tokens[:i], tokens[i+1:]...)
				secTokens[username] = tokens
				return true
			}
		}
		return false
	}
}

func dupTokenHandler(r *fasthttp.RequestCtx) {
	username := getUserName(r)
	if username != "" {dupTokens[username] = CW(8)}
}

func dupControl(username string) bool {
	if string(dupTokens[username]) == string(dupUsedTokens[username]) {return false} else {
		dupUsedTokens[username] = dupTokens[username]
		return true
	}
}
