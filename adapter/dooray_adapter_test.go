package adapter

import (
	"fmt"
	"github.com/gwiyeomgo/adapters/config"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTask_Get_SUCCESS(t *testing.T) {
	task := Task{}
	key := "projectNo"
	val := "123"
	task.Set(key, val)
	projectNo := task.Get("projectNo")
	assert.Equal(t, val, projectNo)
}

func TestTask_Set_SUCCESS(t *testing.T) {
	task := Task{}
	tests := []struct {
		action string
		key    string
		val    string
	}{
		{action: "add title", key: "title", val: "제목"},
		{action: "add mineType", key: "mineType", val: "text/x-markdown"},
	}

	for _, tt := range tests {
		task.Set(tt.key, tt.val)
	}

	for _, tt := range tests {
		t.Run(tt.action, func(t *testing.T) {
			assert.Equal(t, tt.val, task[tt.key])
		})
	}

}
func TestTask_New_SUCCESS(t *testing.T) {
	task := Task{}
	task.Set("content", "내용")
	task.Set("mimeType", "text/x-markdown")
	task.Set("organizationMemberId", "123")

	expected := `{"body":{"content":"내용","mimeType":"text/x-markdown"},"dueDateFlag":true,"priority":"none","subject":"title","users":{"to":[{"member":{"organizationMemberId":"123"},"type":"member"}]}}`
	reqBody := task.New()
	assert.Equal(t, expected, reqBody)
}

func TestTask_Send_SUCCESS(t *testing.T) {
	// setUp WebServer Fixture
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/project/v1/projects/1/posts" {
			w.WriteHeader(200)
			fmt.Println("보냄?")
		} else {
			w.WriteHeader(400)
		}
	}))
	defer server.Close()

	config.Config.Dooray.Project.List.ErrorEvent.ProjectNo = "1"
	config.Config.Dooray.Project.List.ErrorEvent.ProjectMemberGroupId = "1"
	config.Config.Dooray.Project.Url = fmt.Sprintf("%v/project/v1/projects", server.URL)
	config.Config.Dooray.ApiKey = "TEST API Key"

	task := Task{}
	task.Set("title", "")
	task.Set("mineType", "text/x-markdown")
	task.Set("content", "내용")
	task.Set("projectNo", config.Config.Dooray.Project.List.ErrorEvent.ProjectNo)
	task.Set("organizationMemberId", config.Config.Dooray.Project.List.ErrorEvent.ProjectMemberGroupId)

	projectNo := task.Get("projectNo")
	doorayUrl := config.Config.Dooray.Project.Url + "/" + projectNo + "/posts"
	apiKey := config.Config.Dooray.ApiKey

	adapter := DoorayAdapter{Task: &task, DoorayUrl: doorayUrl, ApiKey: apiKey}
	adapter.Send()
}

func TestTask_ChangeContentType_SUCCESS(t *testing.T) {
	content := map[string]interface{}{}
	request := map[string]interface{}{}
	request["host"] = "host 입력"
	content["request"] = request
	content["error"] = fmt.Sprintf("[ERROR] %v \n", "stackTrace")

	result := ChangeContentType(content)
	fmt.Println(result)
}
