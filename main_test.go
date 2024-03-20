package main

import (
	"bytes"
	"fmt"
	"github.com/gwiyeomgo/adapters/config"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func init() {
	fmt.Println("______________________________")
	config.Config.AwsS3.HttpEndPoint = "http://120.0.2.4"
	config.Config.AwsS3.Bucket = "test2"
	config.Config.AwsS3.Region = "ap-northeast-2"
	config.Config.AwsS3.AccessKeyId = "zz"
	config.Config.AwsS3.SecretAccessKey = "zz"
}
func TestCreateS3Uload(t *testing.T) {
	//상태 테스트 - 직접 연결
	//direct AWS S3 connection in UploadFile
	t.Skip()
	// given
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	file, _ := os.Open("./testdata/users.png")
	fileContents, _ := ioutil.ReadAll(file)
	fi, _ := file.Stat()
	file.Close()
	part, _ := writer.CreateFormFile("files", fi.Name())
	writer.WriteField("path", "test")
	part.Write(fileContents)

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/files", body)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	rec := httptest.NewRecorder()

	e := echo.New()
	ctx := e.NewContext(req, rec)

	err := CreateS3Uload(ctx)
	if err != nil {
		t.Errorf("CreateS3Uload 실패: %v", err)
	}

	assert.Equal(t, http.StatusOK, rec.Code)

}
