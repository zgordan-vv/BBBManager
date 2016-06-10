package main

import (
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
	fboauth "golang.org/x/oauth2/facebook"
	"github.com/valyala/fasthttp"
	"os"
	"strings"
)

var (
	githubOauthConf = &oauth2.Config{
		ClientID: GITHUB_CLIENT_ID,
		ClientSecret: GITHUB_CLIENT_SECRET,
		Scopes: []string{"user:mail", "repo"},
		Endpoint: githuboauth.Endpoint,
		RedirectURL: "http://bbbdemo.slava.zgordan.ru/api/GHCB",
	}
	fbOauthConf = &oauth2.Config{
		ClientID: FB_CLIENT_ID,
		ClientSecret: FB_CLIENT_SECRET,
		Scopes: []string{"email", "user_birthday", "user_location", "user_about_me"},
		Endpoint: fboauth.Endpoint,
		RedirectURL: "http://bbbdemo.slava.zgordan.ru/api/FBCB",
	}
)

func getGitHubLoginURLHandler(r *fasthttp.RequestCtx) {
	oauthStateString := string(CW(32))
	url := githubOauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
	url = strings.Replace(url, "client_id=", "client_id="+GITHUB_CLIENT_ID, -1)
	out(r, url)
}

func getFBLoginURLHandler(r *fasthttp.RequestCtx) {
	url := fbOauthConf.AuthCodeURL("")
	url = strings.Replace(url, "client_id=", "client_id="+FB_CLIENT_ID, -1)
	out(r, url)
}

func githubCallback(r *fasthttp.RequestCtx) {
	code := r.FormValue("code")
	fmt.Println(string(code))
	token, err := githubOauthConf.Exchange(oauth2.NoContext, string(code))
	if err != nil {
		fmt.Fprintf(r, err.Error())
		r.Redirect("/",302)
	}
	oauthClient := githubOauthConf.Client(oauth2.NoContext, token)
	client := github.NewClient(oauthClient)
	user, _, err := client.Users.Get("")
	if err != nil {
		fmt.Println("Getting user error")
		r.Redirect("/",302)
	}
	fmt.Println(user.Login)
	r.Redirect("/",302)
}

func fbCallback(r *fasthttp.RequestCtx) {
	code := r.FormValue("code")
	fmt.Println(string(code))
	fmt.Println("Client ID is "+fbOauthConf.ClientID)
	fmt.Println("Real Client ID is "+FB_CLIENT_ID)
	token, err := fbOauthConf.Exchange(oauth2.NoContext, string(code))
	if err != nil {
		fmt.Println(err)
		r.Redirect("/",302)
	}
	_, response, err := fasthttp.Get(nil, "https://graph.facebook.com/me?access_token=" + token.AccessToken)
	if err != nil {
		fmt.Println("Getting user error")
		r.Redirect("/",302)
	}
	fmt.Println("Responded")
	fmt.Println(string(response))
	r.Redirect("/",302)
}
