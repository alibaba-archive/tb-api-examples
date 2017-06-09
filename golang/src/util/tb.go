package util

import (
	"fmt"
	"log"

	"github.com/DavidCai1993/request"
)

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

//CreateProject 创建项目（不带_organizationId则默认创建在个人项目中）
func CreateProject(token string) (string, error) {
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

//CreateTasklist 创建任务分组,将之前创建的project的projectid传入
func CreateTasklist(token string, pid string) (string, error) {
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

//CreateTask 创建任务，将之前创建的tasklist的tasklistid传入
func CreateTask(token string, tlid string) (string, error) {
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

//DoneTask 完成任务，将之前创建的task的taskid传入
func DoneTask(token string, tid string) error {
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

//CreatEvent 创建日程，将之前创建的project的projectid传入
func CreatEvent(token string, pid string) (string, error) {
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

//UpdateEvent 更新日程标题,将之前创建的event的eventid传入
func UpdateEvent(token string, eid string) error {
	url := "https://www.teambition.com/api/events/" + eid + "?" + "access_token=" + token

	_, err := request.Put(url).Send(map[string]interface{}{
		"title": "我的日程2",
	}).JSON()

	if err != nil {
		return err
	}
	return nil
}

//AddProjectMember 添加项目成员,将之前创建的project的projectid传入
func AddProjectMember(token string, pid string) error {
	url := "https://www.teambition.com/api/v2/projects/" + pid + "/members" + "?" + "access_token=" + token
	_, err := request.Post(url).Send(map[string]interface{}{
		"email": "test@teambition.com",
	}).End()
	if err != nil {
		return err

	}
	log.Println("添加用户:test@teambition.com 完成")
	return nil

}

//DeleteTask 删除任务
func DeleteTask(token string, tid string) error {
	url := "https://www.teambition.com/api/tasks/" + tid + "?" + "access_token=" + token
	_, err := request.Delete(url).End()
	if err != nil {

		return err

	}
	log.Println("删除任务成功")
	return nil
}

//DeleteEvent 删除日程
func DeleteEvent(token string, eid string) error {
	url := "https://www.teambition.com/api/events/" + eid + "?" + "access_token=" + token
	_, err := request.Delete(url).End()
	if err != nil {
		return err
	}
	log.Println("删除日程成功")
	return nil
}

//DeleteProject 删除项目
func DeleteProject(token string, pid string) error {
	url := "https://www.teambition.com/api/projects/" + pid + "?" + "access_token=" + token
	_, err := request.Delete(url).End()
	if err != nil {
		return err
	}
	log.Println("删除项目成功")
	return nil
}
