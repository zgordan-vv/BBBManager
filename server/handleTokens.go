package main

import (
	"net/http"
)

var secTokens map[string][]string = make(map[string][]string)
var dupTokens map[string]string = make(map[string]string)
var dupUsedTokens map[string]string = make(map[string]string)

func secTokenHandler(w http.ResponseWriter, r *http.Request) {
	username := getUserName(r)
	secToken := CW(32)
	secTokens[username] = append(secTokens[username],secToken)
	out(w, secToken)
}

func checkSec(sectoken string, username string) bool {
	tokens, ok := secTokens[username]
	if !ok {return false} else {
		return true
		for i, token := range(tokens) {
			if token == sectoken {
				tokens = append(tokens[:i], tokens[i+1:]...)
				secTokens[username] = tokens
				return true
			}
		}
		return false
	}
}

func dupTokenHandler(w http.ResponseWriter, r *http.Request) {
	username := getUserName(r)
	if username != "" {dupTokens[username] = CW(8)}
}

func dupControl(username string) bool {
	if dupTokens[username] == dupUsedTokens[username] {return false} else {
		dupUsedTokens[username] = dupTokens[username]
		return true
	}
}
