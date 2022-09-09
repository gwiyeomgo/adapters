package main

import (
	"fmt"
	"github.com/gwiyeomgo/adapters/config"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateTask(t *testing.T) {
	// setUp WebServer Fixture
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/project/v1/projects/1/posts" {
			w.WriteHeader(200)
		} else {
			w.WriteHeader(400)
		}
	}))
	defer server.Close()

	// given
	config.Config.Dooray.Project.List.ErrorEvent.ProjectNo = "1"
	config.Config.Dooray.Project.List.ErrorEvent.ProjectMemberGroupId = "1"
	config.Config.Dooray.Project.Url = fmt.Sprintf("%v/project/v1/projects", server.URL)
	config.Config.Dooray.ApiKey = "TEST API Key"

	content := map[string]interface{}{}
	request := map[string]interface{}{}
	request["host"] = ""
	request["url"] = ""
	request["method"] = ""
	request["contentType"] = ""
	request["authorization"] = ""
	request["requestBody"] = ""
	request["browserUserAgent"] = ""
	content["request"] = request
	response := map[string]interface{}{}
	response["statusCode"] = "int64(500)"
	content["error"] = fmt.Sprintf("[ERROR] %v \n", "stackTrace")

	task := Task{}
	task.Set("title", "")
	task.Set("mineType", "text/x-markdown")
	task.Set("content", ChangeContentType(content))
	task.Set("projectNo", config.Config.Dooray.Project.List.ErrorEvent.ProjectNo)
	task.Set("organizationMemberId", config.Config.Dooray.Project.List.ErrorEvent.ProjectMemberGroupId)

	projectNo := task.Get("projectNo")
	doorayUrl := config.Config.Dooray.Project.Url + "/" + projectNo + "/posts"
	apiKey := config.Config.Dooray.ApiKey

	adapter := DoorayAdapter{task: &task, doorayUrl: doorayUrl, apiKey: apiKey}
	adapter.Send()
	// then
	//	assert.Nil(t, err)
}
