package main

import (
	"encoding/json"
	"errors"
	"fmt"
	adapter "github.com/gwiyeomgo/adapters/adapters"
	"github.com/gwiyeomgo/adapters/config"
	"github.com/labstack/echo"
	"net/http"
	"net/url"
	"os"
	"time"
)

const maxFileSize = int64(10 * 1024 * 1024)
const maxImageSize = int64(1 * 1024 * 1024)
const DateLayout8 = "20060102"

func main() {
	e := echo.New()
	config.ConfigureEnvironment("./")
	e.GET("/", GetDoorayMember)
	e.POST("/dooray/task", CreateTask)
	e.POST("/gmail/email", SendGmail)
	e.POST("/aws/sqs/message", CreateSQSMessage)
	e.POST("/api/files", CreateS3Uload)
	e.Logger.Fatal(e.Start(":1323"))
}

func CreateS3Uload(ctx echo.Context) error {
	var image bool

	q := ctx.QueryParam("image")
	if q == "1" || q == "true" {
		image = true
	}

	path := ctx.FormValue("path")

	form, err := ctx.MultipartForm()
	if err != nil {
		return err
	}

	files := form.File["files"]

	uploadPath := fmt.Sprintf("%s/%v", path, time.Now().Format(DateLayout8))
	fileUrls := make([]string, len(files))

	for i, file := range files {
		if image {
			if file.Size > maxImageSize {
				maxSize := fmt.Sprintf("최대: %dM", maxImageSize/(1024*1024))
				return errors.New("파일 사이즈가 너무 큽니다." + maxSize)
			}
		} else {
			if file.Size > maxFileSize {
				maxSize := fmt.Sprintf("최대: %dM", maxFileSize/(1024*1024))
				return errors.New("파일 사이즈가 너무 큽니다." + maxSize)
			}
		}

		src, err := file.Open()
		if err != nil {

			return err
		}
		defer src.Close()

		accessUrl, err := adapter.AwsS3Adapter().UploadFile(uploadPath, src, file)
		if err != nil {
			return err
		}

		_, err = url.ParseRequestURI(accessUrl)
		if err != nil {
			return errors.New(err.Error())
		}

		fileUrls[i] = accessUrl
	}

	resp := make(map[string][]string)
	resp["accessUrl"] = fileUrls

	return ctx.JSON(http.StatusOK, resp)
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
