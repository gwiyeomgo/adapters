package adapters

import (
	"net/http"
	"fmt"
)

type DoorayAdapter struct {
}

func NewDoorayTask(content Content) (string, string, string, string) {

	title := fmt.Sprintf("%v ...", content.Error.ErrorStackTrace[:40])
	contents := fmt.Sprintf("# 로그인 사용자\n"+
		"* memberId:<span style=\"color:#9933ff\"> %v</span>\n"+
		"* Name:<span style=\"color:#9933ff\"> %v</span>\n\n"+
		"# HTTP Request\n"+
		"* HOST :<span style=\"color:#9933ff\"> %v</span>\n"+
		"* URL :<span style=\"color:#9933ff\"> %v</span>\n"+
		"* Method :<span style=\"color:#9933ff\"> %v</span>\n* Header\n    "+
		"* Content-Type:<span style=\"color:#9933ff\"> %v</span>\n    "+
		"* Authorization\n```\n%v\n```\n"+
		"* BrowserUserAgent\n```\n%v\n```\n"+
		"* Body\n```\n%v\n```\n\n"+
		"# HTTP Response\n"+
		"* status code : <span style=\"color:#9933ff\">%v</span>\n"+
		"* response body :\n```\n%v\n```\n\n"+
		"# Error 원인\n```\n%v\n```",
		content.User.MemberId,
		content.User.Name,
		content.Request.HOST,
		content.Request.URL,
		content.Request.Method,
		content.Request.ContentType,
		content.Request.Authorization,
		content.Request.BrowserUserAgent,
		content.Request.RequestBody,
		content.Response.statusCode,
		content.Response.ResponseBody,
		content.Error.ErrorStackTrace)

	projectNo := config.Config.Dooray.Project.List.ErrorEvent.ProjectNo
	projectMemberGroupId := config.Config.Dooray.Project.List.ErrorEvent.ProjectMemberGroupId
	if content.Response.statusCode == http.StatusNotFound {
		projectNo = config.Config.Dooray.Project.List.ErrorNotFoundEvent.ProjectNo
		projectMemberGroupId = config.Config.Dooray.Project.List.ErrorNotFoundEvent.ProjectMemberGroupId
	}
	return title, contents, projectNo, projectMemberGroupId
}

type Content struct {
	User     User
	Request  Request
	Response Response
	Error    Error
}

type User struct {
	Name     string
	MemberId string
}

type Request struct {
	HOST             string
	URL              string
	Method           string
	ContentType      string
	Authorization    string
	RequestBody      string
	BrowserUserAgent string
}

type Response struct {
	statusCode   int64
	ResponseBody string
}

type Error struct {
	ErrorStackTrace string
}


func (d DoorayAdapter) SendTask(title string, contents string, projectNo string, projectMemberGroupId string) error {
	doorayUrl := config.Config.Dooray.Project.Url + "/" + projectNo + "/posts"
	apiKey := config.Config.Dooray.ApiKey
	requestBody := map[string]interface{}{
		"users": map[string]interface{}{
			"cc": []interface{}{
				map[string]interface{}{
					"type": "group",
					"group": map[string]interface{}{
						"projectMemberGroupId": projectMemberGroupId,
					},
				},
			},
		},
		"subject": title,
		"body": map[string]interface{}{
			"mimeType": "text/x-markdown",
			"content":  contents,
		},
		"dueDateFlag": false,
		"priority":    "none",
	}

	client := rest.Client{}

	return client.
		Request().
		SetBody(requestBody).
		SetHeader("Authorization", apiKey).
		Post(doorayUrl)
}
