package api

import (
	"fmt"
	"log"

	"golang.org/x/oauth2"

	"net/http"
	"tb-api-examples/golang/src/util"

	"github.com/DavidCai1993/request"
)

const htmlIndex = `<html><body>
<a href="/TBLogin">Log in with TB</a>
</body></html>
`

//OauthCfg ...
var OauthCfg = &oauth2.Config{
	ClientID:     "c9c44aa0-45f8-11e7-85e5-25300cc3a657",
	ClientSecret: "e297f011-ea56-4421-8be9-6477933e1591",
	RedirectURL:  "http://localhost:3000/tb/callback",
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://account.teambition.com/oauth2/authorize",
		TokenURL: "https://account.teambition.com/oauth2/access_token",
	},
}
var clientid = oauth2.SetAuthURLParam("client_id", OauthCfg.ClientID)
var clientsecret = oauth2.SetAuthURLParam("client_secret", OauthCfg.ClientSecret)

//HandleMain ...
func HandleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, htmlIndex)

}

//HandleTBLogin 跳转授权地址
func HandleTBLogin(w http.ResponseWriter, r *http.Request) {

	url := OauthCfg.AuthCodeURL("", clientid)
	log.Printf("URL: %v\n", url)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)

}

//HandleTBCallback 处理授权回调
func HandleTBCallback(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")

	res := new(util.Tokenresult)
	//获取accesstoken
	_, err := request.Post("https://account.teambition.com/oauth2/access_token").Send(map[string]string{
		"client_id":     OauthCfg.ClientID,
		"client_secret": OauthCfg.ClientSecret,
		"code":          code,
	}).JSON(res)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(res.AccessToken)

	url := "https://www.teambition.com/api/users/me" + "?" + "access_token=" + res.AccessToken
	res1 := new(util.MeResult)
	_, err1 := request.Get(url).JSON(res1)
	if err1 != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(res1.Name)

	//创建项目：
	pid, err := util.CreateProject(res.AccessToken)
	if err != nil {
		log.Println("项目创建失败:", err.Error())
	}
	log.Println("项目ID: ", pid, " 创建成功!")
	//添加用户test@teambition.com
	err = util.AddProjectMember(res.AccessToken, pid)
	if err != nil {
		log.Println("邀请成员失败:", err.Error())
	}
	//创建任务列表
	tlid, err := util.CreateTasklist(res.AccessToken, pid)
	if err != nil {
		log.Println("任务列表创建失败:", err.Error())
	}
	log.Println("任务列表ID: ", tlid, " 创建成功!")
	//创建任务
	tid, err := util.CreateTask(res.AccessToken, tlid)
	if err != nil {
		log.Println("任务创建失败:", err.Error())
	}
	log.Println("任务ID: ", tid, " 创建成功!")
	//创建日程
	eid, err := util.CreatEvent(res.AccessToken, pid)
	if err != nil {
		log.Println("日程创建失败:", err.Error())
	}
	log.Println("日程ID: ", eid, " 创建成功!")
	//更新日程标题
	err = util.UpdateEvent(res.AccessToken, eid)
	if err != nil {
		log.Println("日程标题更新失败:", err.Error())
	}
	log.Println("日程ID: ", eid, " 标题更新成功!")
	//完成任务
	err = util.DoneTask(res.AccessToken, tid)
	if err != nil {
		log.Println("任务完成失败:", err.Error())
	}
	log.Println("任务ID: ", eid, " 任务完成成功!")
	//删除任务
	err = util.DeleteTask(res.AccessToken, tid)
	if err != nil {
		log.Println("删除任务失败:", err.Error())
	}
	//删除日程
	err = util.DeleteEvent(res.AccessToken, eid)
	if err != nil {
		log.Println("删除日程失败:", err.Error())
	}
	//删除项目
	err = util.DeleteProject(res.AccessToken, pid)
	if err != nil {
		log.Println("删除项目失败:", err.Error())
	}

}
