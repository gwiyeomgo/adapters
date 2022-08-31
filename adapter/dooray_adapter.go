package adapter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gwiyeomgo/adapters/config"
	"net/http"
)

type DoorayAdapter struct {
}

func NewDoorayTask(content Content) (string, string, string, string) {

	title := fmt.Sprintf("%v ...", content.Error.ErrorStackTrace)
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
		content.Response.StatusCode,
		content.Response.ResponseBody,
		content.Error.ErrorStackTrace)

	projectNo := config.Config.Dooray.Project.List.ErrorEvent.ProjectNo
	projectMemberGroupId := config.Config.Dooray.Project.List.ErrorEvent.ProjectMemberGroupId

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
	StatusCode   int64
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
			"to": []interface{}{
				map[string]interface{}{
					"type": "member",
					"member": map[string]interface{}{
						"organizationMemberId": projectMemberGroupId,
					},
				},
			},
		},
		"subject": "title",
		"body": map[string]interface{}{
			"mimeType": "text/html", //"text/x-markdown",
			"content":  "contents",
		},
		"dueDateFlag": true,
		"priority":    "none",
	}
	b, err := json.Marshal(requestBody)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	reqBody := bytes.NewBufferString(string(b))
	/*	reqBody := bytes.NewBufferString(`{
	    "users": {
	        "to": [{
	            "type": "member",
	            "member": {
	                "organizationMemberId": "3352267848439658321"
	            }
	        }, {
	            "type": "emailUser",
	            "emailUser": {
	                "emailAddress": "",
	                "name": ""
	            }
	        }],
	        "cc": [{
	            "type": "member",
	            "member": {
	                "organizationMemberId": "3352267848439658321"
	            }
	        }]
	    },
	    "subject": "장애 대응 테스트 파일",
	    "body": {
	        "mimeType": "text/html",
	        "content": "장애 대응 테스트 입니다."
	    },
	    "dueDate": "2022-01-08T18:00:00+09:00",
	    "dueDateFlag": true,
	    "milestoneId" :"1",
	    "tagIds": ["1", "2"],
	    "priority": "none"
	}`)*/
	req, _ := http.NewRequest("POST", doorayUrl, reqBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", apiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	/*	r, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var result = map[string]interface{}{}
		if err := json.Unmarshal(r, &result); err != nil {
			return err
		}
		fmt.Println(result)*/
	return nil
}
