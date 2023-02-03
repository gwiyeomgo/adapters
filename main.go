package main

import (
	"encoding/json"
	"fmt"
	adapter "github.com/gwiyeomgo/adapters/adapters"
	"github.com/gwiyeomgo/adapters/config"
	"github.com/labstack/echo"
	"net/http"
	"os"
)

func main() {
	e := echo.New()
	config.ConfigureEnvironment("./")
	e.GET("/", GetDoorayMember)
	e.POST("/dooray/task", CreateTask)
	e.POST("/gmail/email", SendGmail)
	e.POST("/aws/sqs/message", CreateSQSMessage)
	e.Logger.Fatal(e.Start(":1323"))
}

func SendGmail(c echo.Context) error {
	sender := adapter.MailAdapter{}.NewEmailSender()
	body := "Subject: Hello World \r\n TEST 내용"
	return sender.Send([]string{"gwiyeomgo@gmail.com"}, []byte(body))
}

func GetDoorayMember(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func CreateTask(c echo.Context) error {
	mineType := "text/x-markdown"
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

	bytes, err := json.Marshal(content)
	if err != nil {
		return err
	}

	task := adapter.Task{}
	task.Set("title", "")
	task.Set("mimType", mineType)
	task.Set("content", string(bytes))
	task.Set("projectNo", config.Config.Dooray.Project.List.ErrorEvent.ProjectNo)
	task.Set("organizationMemberId", config.Config.Dooray.Project.List.ErrorEvent.ProjectMemberGroupId)

	projectNo := task.Get("projectNo")
	doorayUrl := config.Config.Dooray.Project.Url + "/" + projectNo + "/posts"
	apiKey := config.Config.Dooray.ApiKey

	adapter := adapter.DoorayAdapter{Task: &task, DoorayUrl: doorayUrl, ApiKey: apiKey}
	return adapter.Send()
}

func CreateSQSMessage(c echo.Context) error {
	svc := adapter.NewSQS()
	result, err := svc.ListQueues(nil)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var queueURL string
	for _, t := range result.QueueUrls {
		queueURL = *t
	}
	output, err := adapter.SendMessage(svc, "`{test:2}`", queueURL)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println(output)
	return nil
}
