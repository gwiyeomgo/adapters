package adapters

import (
	"encoding/json"
	"fmt"
	"github.com/mssola/user_agent"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestDooray_SendTask(t *testing.T) {
	// setUp WebServer Fixture
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/project/v1/projects/12/posts" {
			w.WriteHeader(200)
		} else {
			w.WriteHeader(400)
		}
	}))
	defer server.Close()

	// given
	config.Config.Dooray.Project.List.ErrorEvent.ProjectNo = ""
	config.Config.Dooray.Project.Url = fmt.Sprintf("%v/project/v1/projects", server.URL)
	config.Config.Dooray.ApiKey = "TEST API Key"

	errorStackTrace := "[ERROR] Key: 'DonationVisitTxCreate.Mobile' Error:Field validation for 'Mobile' failed on the 'required' tag goroutine 113014 [running]: sharing-platform-service/config/handler.CustomHTTPErrorHandler({0xe52cc0, 0xc0008acdc8}, {0xe840b8, 0xc000290d90})  /var/app/staging/config/handler/error_handler.go:77 +0x94"
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
								"takeOverMethodType" : "VISIT",
								"siteCode" : "200020",
								"memberId" : 1,
								"mobile" : "",
								"name" : "gwiyeomgo",
								"agreed" : true,
								"campaignId" : 4,
								"items" : [
								{
								"itemType" : "100",
								"quantity" : 1
								}
								],
								"note" : ""
							}`,
			GetHumanizeBrowserUserAgent("Mozilla/5.0 (Linux; U; Android 2.3.7; en-us; Nexus One Build/FRF91) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1"),
		},
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