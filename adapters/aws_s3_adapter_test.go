package adapters

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gwiyeomgo/adapters/config"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

//https://aws.amazon.com/ko/blogs/developer/mocking-out-then-aws-sdk-for-go-for-unit-testing/

// MockS3Client is a mock implementation of the S3 client for testing purposes.
type MockS3Client struct {
	mock.Mock
}

// PutObject mocks the PutObject method of the S3 client.
func (m *MockS3Client) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return nil, nil
}

// 내 코드의 행위만 검증하는 테스트로 유닛테스트
// 외부에 의존한 모듈들은 mockup
// 기존의 라이브러리를 호출할때 추상을 호출하도록 변경(테스트할때 mockup한 결과물을 추상에 주입)
func Test_awsS3Adapter_UploadFile_URL_creation(t *testing.T) {
	// setUp WebServer Fixture
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer server.Close()

	// Create an instance of the MockS3Client
	mockS3 := new(MockS3Client)
	config.Config.AwsS3.HttpEndPoint = server.URL
	// Create an instance of the AwsS3Adapter and set its S3 client to the mockS3
	adapter := NewAwsS3Adapter(mockS3)

	// Set up test data
	path := "test/path"
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	file, _ := os.Open("../testdata/users.png")
	fileContents, _ := ioutil.ReadAll(file)
	fi, _ := file.Stat()
	file.Close()
	part, _ := writer.CreateFormFile("files", fi.Name())
	writer.WriteField("path", "test")
	part.Write(fileContents)

	writer.Close()

	fileSize := int64(len(fileContents))
	fileHeader := &multipart.FileHeader{
		Filename: file.Name(),
		Size:     fileSize,
	}

	// Call the method under test
	accessUrl, err := adapter.UploadFile(path, file, fileHeader)
	if err != nil {
		t.Errorf("UploadFile returned an error: %v", err)
	}
	// 모킹된 객체에 대해 모든 기대한 호출이 실제로 발생했는지를 확인
	mockS3.AssertExpectations(t)
	fmt.Println(accessUrl)
}
