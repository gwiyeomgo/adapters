package main

import (
	"encoding/json"
	"fmt"
	"github.com/gwiyeomgo/adapters/config"
	"github.com/labstack/echo"
	"net/http"
)

func main() {
	e := echo.New()
	config.ConfigureEnvironment("./")
	e.GET("/", GetDoorayMember)
	e.POST("/dooray/task", CreateTask)
	e.Logger.Fatal(e.Start(":1323"))
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

	task := Task{}
	task.Set("title", "")
	task.Set("mimType", mineType)
	task.Set("content", string(bytes))
	task.Set("projectNo", config.Config.Dooray.Project.List.ErrorEvent.ProjectNo)
	task.Set("organizationMemberId", config.Config.Dooray.Project.List.ErrorEvent.ProjectMemberGroupId)

	projectNo := task.Get("projectNo")
	doorayUrl := config.Config.Dooray.Project.Url + "/" + projectNo + "/posts"
	apiKey := config.Config.Dooray.ApiKey

	adapter := DoorayAdapter{task: &task, doorayUrl: doorayUrl, apiKey: apiKey}
	return adapter.Send()
}
