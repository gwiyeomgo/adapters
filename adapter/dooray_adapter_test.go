package adapter

import (
	"encoding/json"
	"fmt"
	"github.com/gwiyeomgo/adapters/config"
	"github.com/mssola/user_agent"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDooray_SendTask(t *testing.T) {
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

	errorStackTrace := "[ERROR] Key: test"
	token := "Bearer key"
	content := Content{
		User{
			"이름",
			"1",
		},
		Request{
			server.URL,
			"/api/admin/donations/visit/classifications",
			"POST",
			" application/json",
			token,
			`{

								"memberId" : "1",
							
							}`,
			""},
		Response{
			400,
			`{}`,
		},
		Error{
			errorStackTrace,
		},
	}

	// when
	err := DoorayAdapter{}.SendTask(NewDoorayTask(content))

	// then
	assert.Nil(t, err)
}
func GetHumanizeBrowserUserAgent(browserUserAgent string) string {
	if len(browserUserAgent) == 0 {
		return ""
	}

	ua := user_agent.New(browserUserAgent)
	humanizeUserAgent := map[string]interface{}{}
	humanizeUserAgent["mobile"] = ua.Mobile()
	humanizeUserAgent["platform"] = ua.Platform()
	humanizeUserAgent["os"] = ua.OS()

	name, version := ua.Engine()
	engine := map[string]interface{}{}
	engine["name"] = name
	engine["version"] = version
	humanizeUserAgent["engine"] = engine

	name, version = ua.Browser()
	browser := map[string]interface{}{}
	browser["name"] = name
	browser["version"] = version
	humanizeUserAgent["browser"] = browser

	jsonString, err := json.Marshal(humanizeUserAgent)
	if err != nil {
		log.Errorln(err)
		return ""
	}

	return string(jsonString)
}
