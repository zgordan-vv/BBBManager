package oauth

import (
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/users"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/out"
	"fmt"
	"github.com/valyala/fasthttp"
)

func oauthRegister(r *fasthttp.RequestCtx, user *OauthUser) {
	oauthUser := users.User{user.Login, user.FullName, false, ""}
	if err := users.SaveUser(oauthUser); err != nil {fmt.Println("Users have not been saved", err); out.Out(r, "DontSaved"); return}
}
