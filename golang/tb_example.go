package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/DavidCai1993/request"
	"golang.org/x/oauth2"
)

const htmlIndex = `<html><body>
<a href="/TBLogin">Log in with TB</a>
</body></html>
`

// Tokenresult ...
type Tokenresult struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

//MeResult ...
type MeResult struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Avatar string `json:"avatarUrl"`
}

//OrgResult ...
type OrgResult struct {
	ID        string `json:"_id"`
	Name      string `json:"name"`
	CreatorID string `json:"_creatorId"`
}

//ProjectResult ...
type ProjectResult struct {
	ID          string `json:"_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

//TasklistResult ...
type TasklistResult struct {
	ID          string `json:"_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

//TaskResult ...
type TaskResult struct {
	ID      string `json:"_id"`
	Content string `json:"content"`
	Note    string `json:"note"`
}

// EventResult ...
type EventResult struct {
	ID        string `json:"_id"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Title     string `json:"title"`
}

var ctx = context.Background()

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

//进行router定义
func main() {

	http.HandleFunc("/auth", handleMain)
	http.HandleFunc("/TBLogin", handleTBLogin)
	http.HandleFunc("/tb/callback", handleTBCallback)
	fmt.Println(http.ListenAndServe(":3000", nil))
}

//
func handleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, htmlIndex)

}

//跳转授权地址
func handleTBLogin(w http.ResponseWriter, r *http.Request) {

	url := OauthCfg.AuthCodeURL("", clientid)
	log.Printf("URL: %v\n", url)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)

}

//处理授权回调
func handleTBCallback(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	res := new(Tokenresult)
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
	res1 := new(MeResult)
	_, err1 := request.Get(url).JSON(res1)
	if err1 != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(res1.Name)

	//创建项目：
	pid, err := createProject(res.AccessToken)
	if err != nil {
		log.Println("项目创建失败:", err.Error())
	}
	log.Println("项目ID: ", pid, " 创建成功!")
	//添加用户test@qq.com
	err = addProjectMember(res.AccessToken, pid)
	if err != nil {
		log.Println("邀请成员失败:", err.Error())
	}
	//创建任务列表
	tlid, err := createTasklist(res.AccessToken, pid)
	if err != nil {
		log.Println("任务列表创建失败:", err.Error())
	}
	log.Println("任务列表ID: ", tlid, " 创建成功!")
	//创建任务
	tid, err := createTask(res.AccessToken, tlid)
	if err != nil {
		log.Println("任务创建失败:", err.Error())
	}
	log.Println("任务ID: ", tid, " 创建成功!")
	//创建日程
	eid, err := creatEvent(res.AccessToken, pid)
	if err != nil {
		log.Println("日程创建失败:", err.Error())
	}
	log.Println("日程ID: ", eid, " 创建成功!")
	//更新日程标题
	err = updateEvent(res.AccessToken, eid)
	if err != nil {
		log.Println("日程标题更新失败:", err.Error())
	}
	log.Println("日程ID: ", eid, " 标题更新成功!")
	//完成任务
	err = doneTask(res.AccessToken, tid)
	if err != nil {
		log.Println("任务完成失败:", err.Error())
	}
	log.Println("任务ID: ", eid, " 任务完成成功!")
	//删除任务
	err = deleteTask(res.AccessToken, tid)
	if err != nil {
		log.Println("删除任务失败:", err.Error())
	}
	//删除日程
	err = deleteEvent(res.AccessToken, eid)
	if err != nil {
		log.Println("删除日程失败:", err.Error())
	}
	//删除项目
	err = deleteProject(res.AccessToken, pid)
	if err != nil {
		log.Println("删除项目失败:", err.Error())
	}

}

//创建项目（不带_organizationId则默认创建在个人项目中）
func createProject(token string) (string, error) {
	url := "https://www.teambition.com/api/projects" + "?" + "access_token=" + token
	res := new(ProjectResult)
	_, err := request.Post(url).Send(map[string]string{
		"name":        "我的项目2",
		"description": "Teambition的项目",
	}).JSON(res)
	if err != nil {
		fmt.Println(err.Error())
		return "", err

	}
	fmt.Println("项目信息名称: ", res.Name, " 项目ID: ", res.ID, " 项目描述: ", res.Description)
	return res.ID, nil
}

//创建任务分组,将之前创建的project的projectid传入
func createTasklist(token string, pid string) (string, error) {
	url := "https://www.teambition.com/api/tasklists" + "?" + "access_token=" + token
	res := new(TasklistResult)
	_, err := request.Post(url).Send(map[string]string{
		"title":       "我的任务分组1",
		"description": "Teambition的任务分组",
		"_projectId":  pid,
	}).JSON(res)
	if err != nil {

		return "", err
	}
	log.Println("任务列表名称: ", res.Title, " 任务列表ID: ", res.ID, " 任务列表描述: ", res.Description)
	return res.ID, nil
}

//创建任务，将之前创建的tasklist的tasklistid传入
func createTask(token string, tlid string) (string, error) {
	url := "https://www.teambition.com/api/tasks" + "?" + "access_token=" + token
	res := new(TaskResult)
	_, err := request.Post(url).Send(map[string]string{
		"content":     "我的任务1",
		"note":        "Teambition的分组",
		"_tasklistId": tlid,
	}).JSON(res)
	if err != nil {

		return "", err
	}
	log.Println("任务名称: ", res.Content, " 任务ID: ", res.ID, " 任务备注: ", res.Note)
	return res.ID, nil
}

//完成任务，将之前创建的task的taskid传入
func doneTask(token string, tid string) error {
	url := "https://www.teambition.com/api/tasks/" + tid + "/isDone" + "?" + "access_token=" + token
	_, err := request.Put(url).Send(map[string]interface{}{
		"isDone": true,
	}).JSON()
	if err != nil {
		fmt.Println(err.Error())
		return err

	}
	return nil
}

//创建日程，将之前创建的project的projectid传入
func creatEvent(token string, pid string) (string, error) {
	url := "https://www.teambition.com/api/events" + "?" + "access_token=" + token
	res := new(EventResult)
	_, err := request.Post(url).Send(map[string]interface{}{
		"_projectId": pid,
		"title":      "我的日程1",
		"startDate":  "2020-06-01 12:00:00",
		"endDate":    "2020-06-03 12:00:00",
	}).JSON(res)

	if err != nil {

		return "", err
	}
	log.Println("日程名称: ", res.Title, " 日程ID: ", res.ID, " 日程开始时间: ", res.StartDate, " 日程截止时间: ", res.EndDate)
	return res.ID, err

}

//更新日程标题,将之前创建的event的eventid传入
func updateEvent(token string, eid string) error {
	url := "https://www.teambition.com/api/events/" + eid + "?" + "access_token=" + token

	_, err := request.Put(url).Send(map[string]interface{}{
		"title": "我的日程2",
	}).JSON()

	if err != nil {
		return err
	}
	return nil
}

//添加项目成员,将之前创建的project的projectid传入
func addProjectMember(token string, pid string) error {
	url := "https://www.teambition.com/api/v2/projects/" + pid + "/members" + "?" + "access_token=" + token
	_, err := request.Post(url).Send(map[string]interface{}{
		"email": "test@qq.com",
	}).End()
	if err != nil {
		return err

	}
	log.Println("添加用户:378013446@qq.com 完成")
	return nil

}

//删除任务
func deleteTask(token string, tid string) error {
	url := "https://www.teambition.com/api/tasks/" + tid + "?" + "access_token=" + token
	_, err := request.Delete(url).End()
	if err != nil {

		return err

	}
	log.Println("删除任务成功")
	return nil
}

//删除日程
func deleteEvent(token string, eid string) error {
	url := "https://www.teambition.com/api/events/" + eid + "?" + "access_token=" + token
	_, err := request.Delete(url).End()
	if err != nil {
		return err
	}
	log.Println("删除日程成功")
	return nil
}

//删除项目
func deleteProject(token string, pid string) error {
	url := "https://www.teambition.com/api/projects/" + pid + "?" + "access_token=" + token
	_, err := request.Delete(url).End()
	if err != nil {
		return err
	}
	log.Println("删除项目成功")
	return nil
}
