package oauth

import (
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/globs"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/utils"
	"git.bbbdemo.slava.zgordan.ru/BBBManager/server/out"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
	fboauth "golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/linkedin"
	"github.com/valyala/fasthttp"
	"os"
	"strconv"
)

type FBUser struct {
	Id	string `json:"id"`
	Email	string `json:"email"`
	Username	string `json:"name"`
}

type LinkedInUser struct {
	Id	string `json:"id"`
	FirstName	string `json:"firstName"`
	LastName	string `json:"lastName"`
	Headline	string `json:"headline"`
}

type OauthUser struct {
	Login string
	FullName string
	LastToken oauth2.Token
}

var (
	githubOauthConf *oauth2.Config
	fbOauthConf *oauth2.Config
	linkedinOauthConf *oauth2.Config
	EMPTY_TOKEN = oauth2.Token{}
	oauthURLs map[string]string = map[string]string{"FB":"https://graph.facebook.com/me?access_token=","IN":"https://api.linkedin.com/v1/people/~?format=json&oauth2_access_token="}
)

func InitOauth(){
	githubOauthConf = &oauth2.Config{
		ClientID: os.Getenv(globs.PREFIX + "_GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv(globs.PREFIX + "_GITHUB_CLIENT_SECRET"),
		Scopes: []string{"user:mail", "repo"},
		Endpoint: githuboauth.Endpoint,
		RedirectURL: os.Getenv(globs.PREFIX + "_GITHUB_REDIRECT_URL"),
	}
	fbOauthConf = &oauth2.Config{
		ClientID: os.Getenv(globs.PREFIX + "_FB_CLIENT_ID"),
		ClientSecret: os.Getenv(globs.PREFIX + "_FB_CLIENT_SECRET"),
		Scopes: []string{"email", "user_location", "user_about_me"},
		Endpoint: fboauth.Endpoint,
		RedirectURL: os.Getenv(globs.PREFIX + "_FB_REDIRECT_URL"),
	}
	linkedinOauthConf = &oauth2.Config{
		ClientID: os.Getenv(globs.PREFIX + "_IN_CLIENT_ID"),
		ClientSecret: os.Getenv(globs.PREFIX + "_IN_CLIENT_SECRET"),
		Scopes: []string{"r_basicprofile"},
		Endpoint: linkedin.Endpoint,
		RedirectURL: os.Getenv(globs.PREFIX + "_IN_REDIRECT_URL"),
	}
}

func GetGitHubLoginURLHandler(r *fasthttp.RequestCtx) {
	oauthStateString := string(utils.CW(32))
	url := githubOauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
	out.Out(r, url)
}

func GetFBLoginURLHandler(r *fasthttp.RequestCtx) {
	url := fbOauthConf.AuthCodeURL("")
	out.Out(r, url)
}

func GetLinkedInLoginURLHandler(r *fasthttp.RequestCtx) {
	oauthStateString := string(utils.CW(32))
	url := linkedinOauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
	out.Out(r, url)
}

func GithubCallback(r *fasthttp.RequestCtx) {
	code := r.FormValue("code")
	token, err := githubOauthConf.Exchange(oauth2.NoContext, string(code))
	if err != nil {
		fmt.Println(err)
		r.Redirect("/",302)
	}
	user,err := getGitHubUser(token)
	if err != nil {
		fmt.Println("Getting user error")
		r.Redirect("/",302)
	}

	oauthUser := OauthUser{
		Login: "<OAUTH>_GH_"+strconv.Itoa(*user.ID),
		FullName: *user.Name,
		LastToken: *token,
	}

	OauthLogin(r, &oauthUser)
}

func getGitHubUser(token *oauth2.Token) (*github.User, error) {
	oauthClient := githubOauthConf.Client(oauth2.NoContext, token)
	client := github.NewClient(oauthClient)
	user, _, err := client.Users.Get("")
	return user, err
}

func FbCallback(r *fasthttp.RequestCtx) {
	var userData FBUser
	code := r.FormValue("code")
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
	err = json.Unmarshal(response, &userData)
	if err != nil {
		fmt.Println("Getting unmarshal error")
		r.Redirect("/",302)
	}

	oauthUser := OauthUser{
		Login: "<OAUTH>_FB_"+userData.Id,
		FullName: userData.Username,
		LastToken: *token,
	}

	OauthLogin(r, &oauthUser)
}

func LinkedinCallback(r *fasthttp.RequestCtx) {
	var userData LinkedInUser
	code := r.FormValue("code")
	token, err := linkedinOauthConf.Exchange(oauth2.NoContext, string(code))
	if err != nil {
		fmt.Println(err)
		r.Redirect("/",302)
	}
	_, response, err := fasthttp.Get(nil, "https://api.linkedin.com/v1/people/~?format=json&oauth2_access_token=" + token.AccessToken)
	if err != nil {
		fmt.Println("Getting user error")
		r.Redirect("/",302)
	}
	err = json.Unmarshal(response, &userData)
	if err != nil {
		fmt.Println("Getting unmarshal error")
		r.Redirect("/",302)
	}
	fullname := userData.FirstName+" "+userData.LastName
	if fullname == "" {fullname="Anonymous user"}

	oauthUser := OauthUser{
		Login: "<OAUTH>_IN_"+userData.Id,
		FullName: fullname,
		LastToken: *token,
	}

	OauthLogin(r, &oauthUser)
}
