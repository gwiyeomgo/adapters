package main

import (
	"fmt"
	"github.com/gwiyeomgo/adapters/adapter"
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
	content := adapter.Content{
		Request: adapter.Request{
			HOST:             "req.Host",
			URL:              "req.RequestURI",
			Method:           "req.Method",
			ContentType:      "req.Header.Get(echo.HeaderContentType)",
			Authorization:    "token",
			RequestBody:      "string(body)",
			BrowserUserAgent: "GetHumanizeBrowserUserAgent(req.UserAgent())",
		},
		Response: adapter.Response{
			StatusCode:   int64(500),
			ResponseBody: "",
		},
		Error: adapter.Error{
			ErrorStackTrace: fmt.Sprintf("[ERROR] %v \n", "stackTrace"),
		},
	}
	content.User = adapter.User{}
	return adapter.DoorayAdapter{}.SendTask(adapter.NewDoorayTask(content))
}
