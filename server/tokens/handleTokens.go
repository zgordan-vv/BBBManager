package tokens

import (
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/users"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/utils"
	"github.com/valyala/fasthttp"
)

func SecTokenHandler(r *fasthttp.RequestCtx) {
	username := users.GetUserName(r)
	secToken := utils.CW(32)
	secTokens[username] = append(secTokens[username],secToken)
	r.Write(secToken)
}

func DupTokenHandler(r *fasthttp.RequestCtx) {
	username := users.GetUserName(r)
	if username != "" {dupTokens[username] = utils.CW(8)}
}
