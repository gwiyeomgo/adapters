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

	/*err := CreateTask(c)
	if err != nil {
		log.Fatalln(err)
	}*/
	// then
	//	assert.Nil(t, err)
}
